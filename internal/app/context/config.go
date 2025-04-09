package context

import (
	// std
	"flag"
	"log"
	"os"

	// pkg
	"github.com/joho/godotenv"
	"github.com/pkg/errors"

	// internal
	"github.com/vishenosik/web-tools/collections"
	"github.com/vishenosik/web-tools/env"
)

var (
	//
	ErrServerPortMustBeUnique = errors.New("port numbers must be unique")
)

type Config struct {
	Env        string `env:"ENV" default:"dev" validate:"oneof=dev prod test" desc:"The environment in which the application is running"`
	StorePath  string `env:"STORE_PATH" default:"./storage/CherryWatch.db" validate:"required" desc:"Path to sqlite store"`
	GrpcConfig GrpcServer
	RestConfig RestServer
}

type RestServer struct {
	Port uint16 `env:"REST_PORT" default:"8080" desc:"REST server port"`
}

type GrpcServer struct {
	Port uint16 `env:"GRPC_PORT" default:"44844" desc:"gRPC server port"`
}

func init() {

	flag.BoolFunc("config.info", "Show config schema information", env.ConfigInfo[Config](os.Stdout))
	flag.Func("config.doc", "Update config example in docs", env.ConfigDoc[Config]())

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

}

func mustLoadEnvConfig() Config {
	conf, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	return conf
}

func loadEnvConfig() (Config, error) {
	conf := env.ReadEnv[Config]()
	if collections.HasDuplicates(conf.GrpcConfig.Port, conf.RestConfig.Port) {
		return Config{}, ErrServerPortMustBeUnique
	}
	return conf, nil
}
