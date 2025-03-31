package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vishenosik/CherryWatch/internal/app"
	appctx "github.com/vishenosik/CherryWatch/internal/app/context"
)

// @title           sso
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
	flag.Parse()
	runServer()
}

func runServer() {

	ctx := context.Background()

	// Инициализация приложения
	application := app.MustInitApp()

	application.MustRun()

	// Graceful shut down
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	stopctx, cancel := context.WithTimeout(appctx.WithSignalCtx(ctx, <-stop), time.Second*5)
	defer cancel()

	application.Stop(stopctx)
}
