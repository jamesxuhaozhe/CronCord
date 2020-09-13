package worker

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	EtcdEndpoints         []string
	EtcdDialTimeout       int
	MongodbUri            string
	MongodbConnectTimeout int
	JobLogBatchSize       int
	JobLogCommitTimeout   int
}

var WorkerConfig *Config

func InitConfig(filePath string) error {
	if filePath != "" {
		viper.SetConfigFile(filePath)
	} else {
		return errors.New("worker config file is not defined")
	}

	viper.SetConfigType("yaml")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()              // read the match env var
	viper.SetEnvPrefix("CRON_WORKER") // read the env var whose prefix starts with APISERVER
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	WorkerConfig = &Config{
		EtcdEndpoints:         []string{viper.GetString("etcd.endpoint")},
		EtcdDialTimeout:       viper.GetInt("etcd.dial_timeout"),
		MongodbUri:            viper.GetString("mongodb.endpoint"),
		MongodbConnectTimeout: viper.GetInt("mongodb.connection_timeout"),
		JobLogBatchSize:       viper.GetInt("logsave.batch_size"),
		JobLogCommitTimeout:   viper.GetInt("logsave.commit_timeout"),
	}

	watchConfig()

	return nil
}

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s", e.Name)
	})
}
