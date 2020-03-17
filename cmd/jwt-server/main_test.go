package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	"github.com/pytorchtw/go-jwt-server/gen/models"
	"github.com/pytorchtw/go-jwt-server/gen/restapi"
	"github.com/pytorchtw/go-jwt-server/gen/restapi/operations"
	"github.com/pytorchtw/go-jwt-server/utils"
	"io/ioutil"

	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"

	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"

	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

const (
	dbName     = "docker"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "postgres123"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		//fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		log.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

var testManager struct {
	db               *sql.DB
	schemaName       string
	stopDB           func()
	stopHttpListener func()

	client   *http.Client
	api      *operations.JwtAPI
	server   *restapi.Server
	addr     string
	testUser *models.User
}

func createTestDatabase(t *testing.T) (*sql.DB, string, func()) {
	connectionString := fmt.Sprintf("port=%d user=%s password=%s dbname=%s sslmode=disable", dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		t.Fatalf("Fail to create database %s", err.Error())
	}

	rand.Seed(time.Now().UnixNano())
	schemaName := "test" + strconv.FormatInt(rand.Int63(), 10)

	_, err = db.Exec("CREATE SCHEMA " + schemaName)
	if err != nil {
		t.Fatalf("Fail to create schema. %s", err.Error())
	}

	// close db and reconnect with the new test schema set
	db.Close()
	connectionString = fmt.Sprintf("port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s", dbPort, dbUser, dbPassword, dbName, schemaName)
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		t.Fatalf("Fail to create database %s", err.Error())
	}

	return db, schemaName, func() {
		_, err := db.Exec("DROP SCHEMA " + schemaName + " CASCADE")
		if err != nil {
			t.Fatalf("Fail to drop schema. %s", err.Error())
		}
	}
}

func Test_setup(t *testing.T) {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	// replace the db with test db
	db, schemaName, dbStopFunc := createTestDatabase(t)
	testManager.db = db
	testManager.schemaName = schemaName
	testManager.stopDB = dbStopFunc

	os.Setenv("CONFIG_FILE", "config.unit_test")
	os.Setenv("TEST_DB_SCHEMA", schemaName)
	testManager.api = operations.NewJwtAPI(swaggerSpec)
	testManager.server = restapi.NewServer(testManager.api)
	testManager.server.ConfigureAPI()

	log.Println("created test database " + schemaName)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		t.Fatalf("error getting db instance data, %s", err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+utils.Basepath+"/db/migrations", "postgres", driver)
	if err != nil {
		t.Fatalf("error getting db instance data, %s", err.Error())
	}

	err = m.Up()
	if err != nil {
		t.Fatalf("error executing migrations steps, #{err.Error()}")
	}
	log.Println("completed db migration steps")

	s := &http.Server{
		Handler: testManager.server.GetHandler(),
	}

	l, err := net.Listen("tcp", ":0")
	ok(t, err)

	testManager.addr = l.Addr().String()
	testManager.stopHttpListener = func() {
		l.Close()
	}

	go s.Serve(l)
	log.Println("test api server configured and running")
}

func GetTestUser(email string, pass string) models.User {
	user := models.User{}
	user.Email = &email
	user.Password = &pass
	return user
}

func Test_createUser(t *testing.T) {
	user := GetTestUser("demo@demo.com", "testpassword")
	bytes, err := json.Marshal(user)
	ok(t, err)

	res, err := http.Post("http://"+testManager.addr+"/api/user",
		"application/json",
		strings.NewReader(string(bytes)))
	ok(t, err)
	equals(t, res.StatusCode, 201)
	log.Println(string(bytes))

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	ok(t, err)
	log.Println(string(result))
}

func Test_createExistingUser(t *testing.T) {
	user := GetTestUser("demo@demo.com", "testpassword")
	bytes, err := json.Marshal(user)
	ok(t, err)

	res, err := http.Post("http://"+testManager.addr+"/api/user",
		"application/json",
		strings.NewReader(string(bytes)))
	ok(t, err)
	equals(t, res.StatusCode, 500)

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	ok(t, err)
	log.Println(string(result))
}

func Test_createToken(t *testing.T) {
	user := GetTestUser("demo@demo.com", "testpassword")
	bytes, err := json.Marshal(user)
	ok(t, err)

	res, err := http.Post("http://"+testManager.addr+"/api/token",
		"application/json",
		strings.NewReader(string(bytes)))
	ok(t, err)
	equals(t, res.StatusCode, 201)

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	ok(t, err)
	log.Println(string(result))

	err = json.Unmarshal(result, &user)
	ok(t, err)
	log.Println(user.Token)
	testManager.testUser = &user
}

func Test_helloWithToken(t *testing.T) {
	testManager.client = &http.Client{}
	request, err := http.NewRequest("GET", "http://"+testManager.addr+"/api/hello", nil)
	request.Header.Add("Authorization", "Bearer "+testManager.testUser.Token)
	res, err := testManager.client.Do(request)
	ok(t, err)
	equals(t, res.StatusCode, 200)

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	ok(t, err)
	log.Println(string(result))
}

func Test_shutdown(t *testing.T) {
	testManager.stopDB()
	log.Println("dropped test database " + testManager.schemaName)

	testManager.stopHttpListener()
	log.Println("stopped http listener")

	testManager.server.Shutdown()
	log.Println("api server shutdowned")
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
