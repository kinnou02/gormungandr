package schedules

import (
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.String("listen", ":8080", "[IP]:PORT to listen")
	pflag.Duration("timeout", time.Second, "timeout for call to kraken")
	pflag.String("kraken", "tcp://localhost:3000", "zmq addr for kraken")
	pflag.String("pprof-listen", "", "address to listen for pprof. format: \"IP:PORT\"")
	pflag.Lookup("pprof-listen").NoOptDefVal = "localhost:6060"
	pflag.Bool("json-log", false, "enable json logging")
	pflag.StringP("connection-string", "c",
		"host=localhost user=navitia password=navitia dbname=jormungandr sslmode=disable",
		"connection string to the jormungandr database",
	)
	pflag.Bool("skip-auth", false, "disable authentication")
	pflag.String("newrelic-license", "", "license key new relic")
	pflag.String("newrelic-appname", "gormungandr", "application name in new relic")
	pflag.StringP("rabbitmq-dsn", "r", "amqp://guest:guest@localhost:5672/", "connection uri for rabbitmq")
	pflag.String("stats-exchange", "stat_persistor_exchange_topic", "exchange where to send stats")
	pflag.Bool("skip-stats", false, "disable statistics")
}

type Config struct {
	Listen           string
	Timeout          time.Duration
	Kraken           string
	PprofListen      string `mapstructure:"pprof-listen"`
	JSONLog          bool   `mapstructure:"json-log"`
	ConnectionString string `mapstructure:"connection-string"`
	SkipAuth         bool   `mapstructure:"skip-auth"`
	NewRelicLicense  string `mapstructure:"newrelic-license"`
	NewRelicAppName  string `mapstructure:"newrelic-appname"`
	RabbitmqDsn      string `mapstructure:"rabbitmq-dsn"`
	StatsExchange    string `mapstructure:"stats-exchange"`
	SkipStats        bool   `mapstructure:"skip-stats"`
}

func GetConfig() (Config, error) {
	var config Config
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return config, err
	}
	viper.SetEnvPrefix("SCHEDULES")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	err = viper.Unmarshal(&config)
	return config, err
}
