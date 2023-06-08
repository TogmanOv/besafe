package main

import (
	"flag"
	"os"

	"github.com/forstes/besafe-go/customer/pkg/hash"
	"github.com/forstes/besafe-go/customer/pkg/logger"
	"github.com/forstes/besafe-go/customer/pkg/store/postgres"
	"github.com/forstes/besafe-go/customer/pkg/token"
	http "github.com/forstes/besafe-go/customer/services/customer/internal/delivery/http"
	v1 "github.com/forstes/besafe-go/customer/services/customer/internal/delivery/http/v1"
	"github.com/forstes/besafe-go/customer/services/customer/internal/repository"
	"github.com/forstes/besafe-go/customer/services/customer/internal/service"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	dbConnCfg := postgres.ConnectionConfig{}
	httpServerCfg := http.ServerConfig{}
	logger.New(os.Stderr, true)

	err := godotenv.Load(".env")
	if err != nil {
		log.Err(err).Msg("Failed to load .env file")
		os.Exit(1)
	}

	flag.IntVar(&httpServerCfg.Port, "http-port", 8000, "HTTP server port")
	flag.StringVar(&httpServerCfg.ReadTimeout, "http-read-timeout", "10s", "HTTP read timeout")
	flag.StringVar(&httpServerCfg.WriteTimeout, "http-write-timeout", "30s", "HTTP write timeout")
	flag.StringVar(&httpServerCfg.IdleTimeout, "http-idle-timeout", "1m", "HTTP idle timeout")

	flag.IntVar(&dbConnCfg.Port, "pg-port", 5432, "Postgres port")
	flag.StringVar(&dbConnCfg.Host, "pg-host", "localhost", "Postgres host")
	flag.StringVar(&dbConnCfg.User, "pg-user", os.Getenv("PG_USER"), "Postgres user")
	flag.StringVar(&dbConnCfg.Password, "pg-password", os.Getenv("PG_PASSWORD"), "Postgres password")
	flag.StringVar(&dbConnCfg.DbName, "pg-db-name", os.Getenv("PG_DB_NAME"), "Postgres DB name")
	flag.IntVar(&dbConnCfg.MaxOpenConnections, "pg-max-open-conns", 15, "Postgres max open connections")
	flag.StringVar(&dbConnCfg.MaxIdleTime, "pg-max-idle-time", "15m", "Postgres max connection idle time")
	flag.Parse()

	db, err := postgres.OpenDB(dbConnCfg)
	if err != nil {
		log.Err(err)
		os.Exit(1)
	}
	defer db.Close()

	log.Print("Connected to Postgres DB")

	passwordHasher := hash.NewSHA256Hasher("Bruhable")
	tokenManager, err := token.NewManager(os.Getenv("TOKEN_KEY"))
	userRepository := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepository, passwordHasher, tokenManager)

	httpServerV1 := http.NewHttpServer(v1.NewRouter(userService).GetRoutes(), httpServerCfg)
	err = httpServerV1.Serve()
	if err != nil {
		log.Err(err).Msg("Failed to start HTTP server")
	}
}
