package main

import (
	"bitbucket.org/lyndus/backend/internal/pkg/asaas"
	"context"
	"io/ioutil"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4"

	"bitbucket.org/lyndus/backend/internal/bs/api"
	"bitbucket.org/lyndus/backend/internal/config"
	"bitbucket.org/lyndus/backend/internal/platform/postgres"
)

func initKeys() error {
	signBytes, err := ioutil.ReadFile(config.Config.PrivKeyPath)
	if err != nil {
		return err
	}

	config.SignKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return err
	}
	verifyBytes, err := ioutil.ReadFile(config.Config.PubKeyPath)
	if err != nil {
		return err
	}

	config.VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = initKeys()
	if err != nil {
		log.Fatal(err)
	}

	pgxLogLevel, err := postgres.LogLevelFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	pgxLogLevel = pgx.LogLevelWarn

	pgPool, err := postgres.NewPGXPool(config.Config.PostgresHost, config.Config.PostgresUser, config.Config.PostgresPass, config.Config.PostgresDatabase, config.Config.PostgresPort, context.Background(), &postgres.PGXStdLogger{}, pgxLogLevel)
	if err != nil {
		log.Fatal(err)
	}
	defer pgPool.Close()

	asaas.Asaas = asaas.NewAsaas(config.Config.AsaasAPIProd, config.Config.AsaasAPIToken)

	s, err := api.NewServer(":"+config.Config.Port, pgPool)

	s.Start()

}
