package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/willena/super-go-mod-proxy/config"
	"github.com/willena/super-go-mod-proxy/plugins"
	"github.com/willena/super-go-mod-proxy/types"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var logger, _ = zap.NewDevelopment()
var started = time.Now()

var mainConfig *config.Config
var pluginsInstances *types.PhasesPluginsInstance

var configFile = flag.String("config", "config.json", "Configuration file for super-go-proxy")

func main() {
	defer logger.Sync()
	flag.Parse()
	logger.Info("Starting Super-Go-Pro")

	mainConfig, err := config.LoadConfig(*configFile)
	if err != nil {
		return
	}
	pluginsInstances = plugins.CreateFromConfig(mainConfig)

	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/", StatusHandler).Methods("GET")
	r.HandleFunc("/{module:[A-Za-z.0-9/]+}/@v/list", ListVersionHandler).Methods("GET")
	r.HandleFunc("/{module:[A-Za-z.0-9/]+}/@v/{moduleVersion}.info", InfoVersionHandler).Methods("GET")
	r.HandleFunc("/{module:[A-Za-z.0-9/]+}/@v/{moduleVersion}.mod", ModVersionHandler).Methods("GET")
	r.HandleFunc("/{module:[A-Za-z.0-9/]+}/@v/{moduleVersion}.zip", ZipVersionHandler).Methods("GET")
	r.HandleFunc("/{module:[A-Za-z.0-9/]+}/@v/@latest", LatestVersionHandler).Methods("GET")
	http.Handle("/", r)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Minute * 5,
		IdleTimeout:  time.Minute * 5,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// RunPhase our server in a goroutine so that it doesn't block.
	go func() {
		logger.Info("Starting server on ", zap.String("bind", srv.Addr))
		logger.Info("Usage Metrics are available on /metrics")
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("Could not start http server", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func StatusHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(map[string]interface{}{
		"status":        "OK",
		"moduleVersion": "0.0.1",
		"uptime":        time.Now().Sub(started).Seconds(),
	})
	writer.Write(res)
}
