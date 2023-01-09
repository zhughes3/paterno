package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/file"
	"github.com/zhughes3/elliot/pkg/log"
	"github.com/zhughes3/elliot/pkg/persistence"
)

func main() {
	logger := mustCreateLogger()
	config := mustLoadConfig(logger)
	db := mustCreateDatabase(logger, config)
	srv := mustCreateService(logger, db, config)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatalf("Problem starting service: %v", err)
	}
}

func mustCreateLogger() log.Logger {
	return log.NewZeroLogger(os.Stderr)
}

func mustLoadConfig(logger log.Logger) *koanf.Koanf {
	k := koanf.New(".")
	if err := k.Load(file.Provider(".env"), dotenv.Parser()); err != nil {
		logger.Fatalf("Problem loading configuration from .env file: %v", err)
	}
	return k
}

func mustCreateDatabase(logger log.Logger, cfg *koanf.Koanf) persistence.DB {
	database, err := persistence.NewDB(logger, cfg)
	if err != nil {
		logger.Fatalf("Problem creating database connection: %v", err)
	}

	return database
}

func mustCreateService(logger log.Logger, db persistence.DB, cfg *koanf.Koanf) http.Server {
	port := cfg.String("HTTP_PORT")
	if len(port) == 0 {
		port = "5000"
	}
	r := chi.NewRouter()
	r.Get("/db/version", func(w http.ResponseWriter, r *http.Request) {
		var pgVersion string
		query := "SELECT version();"
		err := db.DB.QueryRow(query).Scan(&pgVersion)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(pgVersion))
	})
	logger.Infof("starting HTTP server at port: %s", port)
	return http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
}
