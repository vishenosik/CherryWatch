package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vishenosik/CherryWatch/internal/app"
	"github.com/vishenosik/gocherry"
	webctx "github.com/vishenosik/gocherry/pkg/context"
)

// @title           CherryWatch
// @version         0.0.1
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
//
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host      localhost:8080
// @BasePath  /
//
// @securityDefinitions.basic  BasicAuth
//
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	gocherry.Flags()
	flag.Parse()

	// Logo()

	ctx := context.Background()

	// App init
	application, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	application.Start(ctx)

	// Graceful shut down
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	stopctx, cancel := context.WithTimeout(
		webctx.WithStopCtx(ctx, <-stop),
		time.Second*5,
	)
	defer cancel()

	application.Stop(stopctx)
}
