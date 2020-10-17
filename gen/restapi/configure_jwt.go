// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/pytorchtw/go-jwt-server/gen/restapi/operations/token"
	"github.com/pytorchtw/go-jwt-server/services"
	"github.com/pytorchtw/go-jwt-server/utils"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"os"

	"github.com/volatiletech/sqlboiler/boil"
	"log"

	"net/http"
	"strings"

	"github.com/go-openapi/runtime/middleware"

	apiModels "github.com/pytorchtw/go-jwt-server/gen/models"
	"github.com/pytorchtw/go-jwt-server/gen/restapi/operations"
	"github.com/pytorchtw/go-jwt-server/gen/restapi/operations/user"
	dbModels "github.com/pytorchtw/go-jwt-server/repo/models"
)

//go:generate swagger generate server --target ../../gen --name Jwt --spec ../../swagger.yml

func configureFlags(api *operations.JwtAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

const (
	SecretKey = "this is a secretkey"
)

func resetDB(db *sql.DB, schemaName string) error {
	/*
		log.Println("resetting db")
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return err
		}
	*/

	log.Println("dropping test schema")
	_, err := db.Exec("DROP SCHEMA IF EXISTS " + schemaName + " CASCADE")
	if err != nil {
		return err
	}

	log.Println("creating test schema")
	_, err = db.Exec("CREATE SCHEMA " + schemaName)
	if err != nil {
		return err
	}

	log.Println("creating schema migrations table")
	sql := `CREATE TABLE IF NOT EXISTS schema_migrations
(
    version bigint NOT NULL,
    dirty boolean NOT NULL,
    CONSTRAINT schema_migrations_pkey PRIMARY KEY (version)
)`
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	//m, err := migrate.NewWithDatabaseInstance("file://"+utils.Basepath+"/db/migrations", "postgres", driver)
	m, err := migrate.NewWithDatabaseInstance("file://./db/migrations", "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		return err
	}
	log.Println("completed db migration steps")

	return nil
}

func setupDB() (*sql.DB, error) {
	basePath := utils.Basepath
	configFile := os.Getenv("CONFIG_FILE")
	viper.SetConfigName(configFile)
	viper.AddConfigPath(basePath)
	viper.SetEnvPrefix("TEST")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic("viper config not found")
	}
	viper.AutomaticEnv()

	dbHost := viper.GetString("db.host")
	dbPort := viper.GetInt("db.port")
	dbUser := viper.GetString("db.user")
	dbPassword := viper.GetString("db.pass")
	dbName := viper.GetString("db.dbname")
	sslMode := viper.GetString("db.sslmode")
	schemaName := viper.GetString("db.schema")

	serverEnv := viper.GetString("env")
	if serverEnv == "integration_test" {
		log.Println("running in integration_test env")
		log.Println("setting db schema to " + schemaName)

		connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s", dbHost, dbPort, dbUser, dbPassword, dbName, sslMode, schemaName)
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			panic(err.Error())
		}
		err = resetDB(db, schemaName)
		if err != nil {
			panic(err)
		}
		boil.SetDB(db)
		return db, nil

	} else if serverEnv == "production" {
		log.Println("running in production env")
		log.Println("setting db schema to " + schemaName)

		connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s", dbHost, dbPort, dbUser, dbPassword, dbName, sslMode, schemaName)
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			panic(err.Error())
		}
		boil.SetDB(db)
		return db, nil

	} else if serverEnv == "unit_test" {
		log.Println("running in unit test env")

		schemaName := viper.GetString("db.schema")
		if schemaName == "" {
			panic("error table name")
		}
		log.Println("setting db schema to " + schemaName)

		connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s", dbHost, dbPort, dbUser, dbPassword, dbName, sslMode, schemaName)
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			panic(err.Error())
		}
		boil.SetDB(db)
		return db, nil

	} else {
		panic("missing environment config file to execute")
	}
	return nil, nil
}

func configureAPI(api *operations.JwtAPI) http.Handler {

	db, err := setupDB()
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("error db setup")
	}

	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()
	api.TxtProducer = runtime.TextProducer()

	/*
		api.KeyAuth = func(token string) (*db_models.Principal, error) {
			if token == "abcdefuvwxyz" {
				prin := db_models.Principal(token)
				return &prin, nil
			}
			api.Logger("Access attempt with incorrect api key auth: %s", token)
			return nil, errors.New(401, "incorrect api key auth")
		}
	*/

	api.UserCreateHandler = user.CreateHandlerFunc(
		func(params user.CreateParams) middleware.Responder {
			ws := services.NewWebService(db)
			var dbUser = dbModels.User{}
			dbUser.Email = *params.Body.Email
			dbUser.Password = *params.Body.Password
			err := ws.CreateUser(&dbUser)
			if err != nil {
				log.Println(err)
				return user.NewCreateDefault(500)
			}
			//payloadUser := apiModels.User{}
			//payloadUser.ID = int64(dbUser.ID)
			//payloadUser.Email = &dbUser.Email
			//payloadUser.Password = &dbUser.Password
			//return user.NewCreateCreated().WithPayload(&payloadUser)
			params.Body.ID = int64(dbUser.ID)
			return user.NewCreateCreated().WithPayload(params.Body)
		})
	if api.UserCreateHandler == nil {
		api.UserCreateHandler = user.CreateHandlerFunc(func(params user.CreateParams) middleware.Responder {
			return middleware.NotImplemented("operation user.Create has not yet been implemented")
		})
	}

	api.TokenCreateTokenHandler = token.CreateTokenHandlerFunc(func(params token.CreateTokenParams) middleware.Responder {
		ws := services.NewWebService(db)
		ok, err := ws.VerifyUser(*params.Body.Email, *params.Body.Password)
		if err != nil {
			log.Println(err)
			return token.NewCreateTokenDefault(500)
		}
		if !ok {
			log.Println("error verifying user")
			return token.NewCreateTokenDefault(500)
		}
		createdToken, err := ws.CreateToken()
		if err != nil {
			log.Println(err)
			return token.NewCreateTokenDefault(500)
		}
		user := apiModels.User{}
		user.Email = params.Body.Email
		user.Token = createdToken.Token
		return token.NewCreateTokenCreated().WithPayload(&user)
	})

	api.GetGreetingHandler = operations.GetGreetingHandlerFunc(func(params operations.GetGreetingParams) middleware.Responder {
		name := swag.StringValue(params.Name)
		if name == "" {
			name = "World "
		}
		user := apiModels.User{}
		var email = "hello"
		user.Email = &email
		user.Password = &name
		return operations.NewGetGreetingOK().WithPayload(&user)
	})
	if api.GetGreetingHandler == nil {
		api.GetGreetingHandler = operations.GetGreetingHandlerFunc(func(params operations.GetGreetingParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetGreeting has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	/*
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.URL.String())
			//log.Println(r.Header)

			requestDump, err := httputil.DumpRequest(r, true)
			if err != nil {
				log.Println(err)
			}
			log.Println(string(requestDump))

			handler.ServeHTTP(w, r)
		})
	*/

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		if strings.HasSuffix(r.URL.String(), "/user") || strings.HasSuffix(r.URL.String(), "/token") || strings.HasSuffix(r.URL.String(), "/hello1") {
			handler.ServeHTTP(w, r)
			return
		}

		token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "unauthorized access to this resource")
			log.Println("unauthorized access to this resource")
			return
		}

		if token.Valid {
			//log.Println("token is valid")
			handler.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "token is not valid")
			log.Println("token is not valid")
		}
	})

}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedHeaders:   []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	return c.Handler(handler)
}
