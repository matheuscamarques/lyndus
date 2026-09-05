package main

import (
	"bitbucket.org/lyndus/backend/domain/bonder/services"
	"bitbucket.org/lyndus/backend/internal/bm/financial"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"

	"bitbucket.org/lyndus/backend/infra/auth"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"bitbucket.org/lyndus/backend/infra/handlers/bonder"
	"bitbucket.org/lyndus/backend/infra/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func NewPGXPool(host, user, pass, database string, port int, ctx context.Context, logLevel pgx.LogLevel) (*pgxpool.Pool, error) {

	// user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=10
	dbInfo := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", user, pass, host, port, database)
	conf, err := pgxpool.ParseConfig(dbInfo)
	if err != nil {
		return nil, err
	}

	//conf.ConnConfig.Logger = logger

	// Set the log level for pgx, if set.
	if logLevel != 0 {
		conf.ConnConfig.LogLevel = logLevel
	}

	// pgx, by default, does some I/O operation on initialization of a pool to check if the database is reachable.
	// Comment the following line if you don't want pgx to try to connect pool once the Connect function is called,
	//
	// If comment it, and your application seems stuck, you probably forgot to set up PGCONNECT_TIMEOUT,
	// and your code is hanging waiting for a connection to be established.
	conf.LazyConnect = true

	// pgxpool default max number of connections is the number of CPUs on your machine returned by runtime.NumCPU().
	// This number is very conservative, and you might be able to improve performance for highly concurrent applications
	// by increasing it.
	// conf.MaxConns = runtime.NumCPU() * 5

	pool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		//todo zap log
		return nil, fmt.Errorf("pgx connection error: %w", err)
	}
	return pool, nil
}

func main() {

	logger.New("Bonder")
	logger.Info("System UP")

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(logger.APILogger())

	err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	postgres.DB, err = postgres.InitDB(
		config.Config.PostgresHost,
		config.Config.PostgresUser,
		config.Config.PostgresPass,
		config.Config.PostgresDatabase,
		config.Config.PostgresPort,
	)

	if err != nil {
		log.Fatal(err)
	}
	err = auth.InitKeysPem()
	if err != nil {
		log.Fatal(err)
	}

	pgPool, err := NewPGXPool(config.Config.PostgresHost, config.Config.PostgresUser, config.Config.PostgresPass, config.Config.PostgresDatabase, config.Config.PostgresPort, context.Background(), pgx.LogLevelWarn)

	if err != nil {
		log.Fatal(err)
	}

	services.BonderCompany = services.NewBonderCompanyService()
	services.BonderClient = services.NewBonderClientService()
	services.BenefitService = services.NewBonderBenefitService()
	services.ClientUser = services.NewBonderClientUserService()
	services.BonderBS = services.NewBonderBsService()
	services.BsUser = services.NewBonderBSUserService()
	services.BonderBSCategory = services.NewBonderBSCategoryService()
	services.AppUser = services.NewAppUserService()
	services.ClientCategory = services.NewClientCategoryService()
	financial.Service = financial.NewService(pgPool)

	auth.Service = auth.NewAuthenticationService()
	corsOptions := cors.New(cors.Options{
		//AllowOriginFunc:  AllowOriginFunc,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(corsOptions.Handler)

	r.Post("/login", bonder.LoginHandler)
	r.Route("/", func(r chi.Router) {
		r.Use(auth.ValidateTokenMiddleware)
		r.Route("/bs", bonder.BSRouter)
		r.Route("/bm", bonder.BMRouter)
		r.Route("/client", bonder.ClientRouter)
		r.Route("/appuser", bonder.AppUserRouter)
		r.Route("/master", bonder.MasterRouter)
	})
	log.Println("Running... port:", config.Config.Port)
	// Start server
	err = http.ListenAndServe(":"+config.Config.Port, r)
	if err != nil {
		log.Fatal(err)
	}
	defer pgPool.Close()

}
