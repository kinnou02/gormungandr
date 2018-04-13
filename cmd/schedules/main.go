package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"database/sql"

	_ "net/http/pprof"

	"github.com/CanalTP/gormungandr"
	"github.com/CanalTP/gormungandr/auth"
	"github.com/CanalTP/gormungandr/internal/schedules"
	_ "github.com/lib/pq"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func setupRouter() *gin.Engine {
	r := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	r.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, false))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://*"},
		AllowHeaders:     []string{"Access-Control-Request-Headers", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}

func initLog(jsonLog bool) {
	if jsonLog {
		// Log as JSON instead of the default ASCII formatter.
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	showVersion := pflag.Bool("version", false, "show version")
	pflag.Parse()
	if *showVersion {
		fmt.Printf("gormungandr %s built with %s", gormungandr.Version, runtime.Version())
		os.Exit(0)
	}

	logger := logrus.WithFields(logrus.Fields{
		"version": gormungandr.Version,
		"runtime": runtime.Version(),
	})
	config, err := schedules.GetConfig()
	if err != nil {
		logger.Fatalf("failure to load configuration: %+v", err)
	}
	initLog(config.JSONLog)
	logger = logger.WithFields(logrus.Fields{
		"config": config,
	})
	logger.Info("starting schedules")

	kraken := gormungandr.NewKraken("default", config.Kraken, config.Timeout)

	db, err := sql.Open("postgres", config.ConnectionString)
	if err != nil {
		logrus.Fatal("connection to postgres failed: ", err)
	}
	err = db.Ping()
	if err != nil {
		logrus.Fatal("connection to postgres failed: ", err)
	}

	if len(config.PprofListen) != 0 {
		go func() {
			logrus.Infof("pprof listening on %s", config.PprofListen)
			logrus.Error(http.ListenAndServe(config.PprofListen, nil))
		}()
	}

	r := setupRouter()
	cov := r.Group("/v1/coverage/:coverage")

	if !config.SkipAuth {
		cov.Use(auth.AuthenticationMiddleware(db))
	}

	cov.GET("/*filter", schedules.NoRouteHandler(kraken))
	err = r.Run(config.Listen)
	if err != nil {
		logrus.Errorf("failure to start: %+v", err)
	}
}
