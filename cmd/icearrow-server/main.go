// This module is derived from the Yopass project, specifically the yopass-server module:
// https://github.com/jhaals/yopass/blob/7e50bef6aacc5b401149914fd3472404f1b65e5c/cmd/yopass-server/main.go
//
// Yopass is an open-source project for sharing secrets in a quick and secure manner,
// and is distributed under the Apache License Version 2.0.
//
// Changes from the original module:
// - Removed memcached support.
// - Added support for Walrus.
//
// This module continues to be distributed under the terms of the Apache License Version 2.0.
// For more information, see the original project's license:
//
//	https://github.com/jhaals/yopass/blob/7e50bef6aacc5b401149914fd3472404f1b65e5c/LICENSE
//
package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jhaals/yopass/pkg/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"hippos/hippo"
)

var logLevel zapcore.Level

func init() {
	pflag.String("address", "", "listen address (default 0.0.0.0)")
	pflag.Int("port", 1337, "listen port")
	pflag.Int("max-length", 10000, "max length of encrypted secret")
	pflag.Int("metrics-port", -1, "metrics server listen port")
	pflag.String("redis", "redis://localhost:6379/0", "Redis URL")
	pflag.String("walrus", "http://localhost:31415", "Walrus URL")
	pflag.String("tls-cert", "", "path to TLS certificate")
	pflag.String("tls-key", "", "path to TLS key")
	pflag.Bool("force-onetime-secrets", false, "reject non onetime secrets from being created")
	pflag.CommandLine.AddGoFlag(&flag.Flag{Name: "log-level", Usage: "Log level", Value: &logLevel})

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	_ = viper.BindPFlags(pflag.CommandLine)

	pflag.Parse()
}

func main() {
	logger := configureZapLogger()

	var db server.Database
	var err error
	redisURL := viper.GetString("redis")
	walrusURL := viper.GetString("walrus")
	db, err = hippo.NewHippo(walrusURL, redisURL)
	if err != nil {
		logger.Fatal("Hippo configuration error", zap.Error(err))
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	registry.MustRegister(prometheus.NewGoCollector())

	cert := viper.GetString("tls-cert")
	key := viper.GetString("tls-key")
	quit := make(chan os.Signal, 1)

	y := server.New(db, viper.GetInt("max-length"), registry, viper.GetBool("force-onetime-secrets"), logger)
	yopassSrv := &http.Server{
		Addr:      fmt.Sprintf("%s:%d", viper.GetString("address"), viper.GetInt("port")),
		Handler:   y.HTTPHandler(),
		TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	}
	go func() {
		logger.Info("Starting yopass server", zap.String("address", yopassSrv.Addr))
		err := listenAndServe(yopassSrv, cert, key)
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("yopass stopped unexpectedly", zap.Error(err))
		}
	}()

	metricsServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", viper.GetString("address"), viper.GetInt("metrics-port")),
		Handler: metricsHandler(registry),
	}
	if port := viper.GetInt("metrics-port"); port > 0 {
		go func() {
			logger.Info("Starting yopass metrics server", zap.String("address", metricsServer.Addr))
			err := listenAndServe(metricsServer, cert, key)
			if !errors.Is(err, http.ErrServerClosed) {
				logger.Fatal("metrics server stopped unexpectedly", zap.Error(err))
			}
		}()
	}

	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	signal := <-quit
	log.Printf("Shutting down HTTP server %s", signal)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := yopassSrv.Shutdown(ctx); err != nil {
		logger.Fatal("shutdown error: %s", zap.Error(err))
	}
	if port := viper.GetInt("metrics-port"); port > 0 {
		if err := metricsServer.Shutdown(ctx); err != nil {
			logger.Fatal("shutdown error: %s", zap.Error(err))
		}
	}
	logger.Info("Server shut down")
}

// listenAndServe starts a HTTP server on the given addr. It uses TLS if both
// certFile and keyFile are not empty.
func listenAndServe(srv *http.Server, certFile string, keyFile string) error {
	if certFile == "" || keyFile == "" {
		return srv.ListenAndServe()
	}
	return srv.ListenAndServeTLS(certFile, keyFile)
}

// metricsHandler builds a handler to serve Prometheus metrics
func metricsHandler(r *prometheus.Registry) http.Handler {
	mx := http.NewServeMux()
	mx.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{EnableOpenMetrics: true}))
	return mx
}

// configureZapLogger uses the `log-level` command line argument to set and replace the zap global logger.
func configureZapLogger() *zap.Logger {
	loggerCfg := zap.NewProductionConfig()
	loggerCfg.Level.SetLevel(logLevel)

	logger, err := loggerCfg.Build()
	if err != nil {
		log.Fatalf("Unable to build logger %v", err)
	}
	zap.ReplaceGlobals(logger)
	return logger
}
