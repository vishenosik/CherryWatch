package app

import (
	// std

	// pkg
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"

	// internal
	"github.com/vishenosik/gocherry/pkg/collections"
	"github.com/vishenosik/gocherry/pkg/operation"
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
	Testing    []string `env:"TESTING"`
}

type RestServer struct {
	Port uint16 `env:"REST_PORT" default:"8080" desc:"REST server port"`
}

type GrpcServer struct {
	Port uint16 `env:"GRPC_PORT" default:"44844" desc:"gRPC server port"`
}

func mustLoadEnvConfig() Config {
	conf, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	return conf
}

func loadEnvConfig() (Config, error) {

	var conf Config

	fail := operation.FailWrapError(Config{}, "loadEnvConfig")

	if err := cleanenv.ReadConfig(".env", &conf); err != nil {
		return fail(errors.Wrap(err, "failed to read config"))
	}

	if collections.HasDuplicates(conf.GrpcConfig.Port, conf.RestConfig.Port) {
		return fail(ErrServerPortMustBeUnique)
	}

	return conf, nil
}
