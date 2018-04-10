package main

import (
	"net/http"
	"os"
	"time"

	_ "net/http/pprof"

	"github.com/CanalTP/gormungandr"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
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

func init_log(jsonLog bool) {
	if jsonLog {
		// Log as JSON instead of the default ASCII formatter.
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {

	config, err := GetConfig()
	if err != nil {
		logrus.Fatalf("failure to load configuration: %+v", err)
	}
	init_log(config.JsonLog)
	logrus.WithFields(logrus.Fields{
		"config": config,
	}).Debug("configuration loaded")

	kraken := gormungandr.NewKraken("default", config.Kraken, config.Timeout)

	if len(config.PprofListen) != 0 {
		go func() {
			logrus.Infof("pprof listening on %s", config.PprofListen)
			logrus.Error(http.ListenAndServe(config.PprofListen, nil))
		}()
	}

	r := setupRouter()
	r.GET("/v1/coverage/:coverage/*filter", NoRouteHandler(kraken))

	err = r.Run(config.Listen)
	if err != nil {
		logrus.Errorf("failure to start: %+v", err)
	}
}
