package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/zhughes3/elliot/pkg/log"
	"github.com/zhughes3/elliot/pkg/secret"
)

func main() {
	logger := mustCreateLogger()
	kv := mustCreateKeyVault(logger)
	srv := mustCreateService(logger, kv)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal("error starting service")
	}
}

func mustCreateLogger() log.Logger {
	return log.NewZeroLogger(os.Stderr)
}

func mustCreateKeyVault(logger log.Logger) secret.KeyVault {
	uri := "get_from_config"
	logger.Infof("using Azure KeyVault at %s", uri)
	keyVaultClient, err := secret.NewAzureKeyVaultClient(uri)
	if err != nil {
		logger.Fatalf("error setting up credentials for Azure KeyVault: %v", err)
	}
	return secret.NewAzureKeyVault(keyVaultClient)
}

func mustCreateService(logger log.Logger, kv secret.KeyVault) http.Server {
	r := chi.NewRouter()
	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		logger.Fatal("panic")
	})
	return http.Server{
		Addr:    ":30003",
		Handler: r,
	}
}
