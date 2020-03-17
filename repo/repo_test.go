package repo

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"os"
	"strconv"

	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/pytorchtw/go-jwt-server/repo/models"
	"github.com/pytorchtw/go-jwt-server/utils"
)

const (
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "postgres123"
	dbName     = "docker"
)

var testManager struct {
	db         *sql.DB
	schemaName string
	repo       Repo
	stop       func()
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
	db, schemaName, dbStopFunc := createTestDatabase(t)
	testManager.db = db
	testManager.schemaName = schemaName
	testManager.stop = dbStopFunc
	testManager.repo = NewDBRepo(testManager.db)
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
		t.Fatalf("error doing migration steps, %s", err.Error())
	}
	log.Println("completed db migration steps")
}

func Test_CreateUser(t *testing.T) {
	user := models.User{}
	user.Email = "test@test.com"
	err := testManager.repo.CreateUser(&user)
	if err != nil {
		t.Fatalf("error creating user, %s", err.Error())
	}
	err = testManager.repo.CreateUser(&user)
	if err == nil {
		t.Fatalf("error should not be nil, should not be able to create the same user again")
	}
	if user.ID != 1 {
		t.Fatal("user id should be 1")
	}
}

func Test_GetUser(t *testing.T) {
	email := "test@test.com"
	user, err := testManager.repo.GetUser(email)
	if err != nil {
		t.Fatalf("error getting user, %s", err.Error())
	}
	if user.Email != email {
		t.Fatalf("error email")
	}
}

func Test_UpdateUser(t *testing.T) {
	email := "test@test.com"
	user, err := testManager.repo.GetUser(email)
	if err != nil {
		t.Fatalf("error getting user, %s", err.Error())
	}
	if user.Email != email {
		t.Fatalf("error email")
	}

	username := "testusername"
	user.Username = username
	rowsAff, err := testManager.repo.UpdateUser(user)
	if err != nil {
		t.Fatalf("error updating user, %s", err.Error())
	}
	if rowsAff != 1 {
		t.Fatalf("error rows count affected")
	}

	user, err = testManager.repo.GetUser(email)
	if err != nil {
		t.Fatalf("error getting user, %s", err.Error())
	}
	if user.Username != username {
		t.Fatalf("error username")
	}
}

func Test_DeleteUser(t *testing.T) {
	email := "test@test.com"
	user, err := testManager.repo.GetUser(email)
	if err != nil {
		t.Fatalf("error getting user, %s", err.Error())
	}
	if user.Email != email {
		t.Fatalf("error email")
	}
	rowsAff, err := testManager.repo.DeleteUser(email)
	if err != nil {
		t.Fatalf("error deleting user, %s", err.Error())
	}
	if rowsAff != 1 {
		t.Fatalf("affected rows should be 1")
	}

	_, err = testManager.repo.GetUser(email)
	if err == nil {
		t.Fatalf("should return error when getting deleted user")
	}
}

func Test_shutdown(t *testing.T) {
	testManager.stop()
	log.Println("dropped test database " + testManager.schemaName)
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
