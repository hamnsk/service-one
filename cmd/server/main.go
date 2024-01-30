package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func newApp(tp AppTracer) *app {
	a := &app{
		httpServer: newHTTPServer(),
		tp:         tp,
	}
	return a
}

func newHTTPServer() *fiber.App {
	return fiber.New()
}

func (a *app) run() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-c
		a.shutdown()
	}()

	return a.httpServer.Listen(":9090")
}

func (a *app) routes() {
	a.httpServer.Get("/", a.audit)
}

func (a *app) shutdown() {

	ctx, serverCancel := context.WithTimeout(context.Background(), 15*time.Second)

	err := a.httpServer.ShutdownWithContext(ctx)
	if err != nil {
		a.fatalServer(err)
	}
	serverCancel()
}
func (a *app) fatalServer(err error) {
	log.Fatal(err.Error())
}

func (a *app) audit(c *fiber.Ctx) error {
	//time.Sleep(1 * time.Second)
	return c.SendString("accepted")
}

func main() {
	log.Println("Application start")
	tp, err := InitTracing()
	if err != nil {
		log.Fatalln("Unable to create a global trace provider", err)
	}

	defer func() {
		_ = tp.Shutdown(context.Background())
	}()

	app := newApp(tp)
	app.httpServer.Use(trace)
	//app.httpServer.Use(otelfiber.Middleware())

	app.routes()

	err = app.run()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Application successful terminated")
}
