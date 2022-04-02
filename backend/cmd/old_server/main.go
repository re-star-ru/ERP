package main

import (
	"backend/internal/app/apiserver/api"
	"backend/internal/app/apiserver/api/catalog"
	"backend/internal/app/apiserver/db"
	"backend/internal/app/apiserver/img"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/asdine/storm"
	"github.com/spf13/viper"
)

type server struct {
	restSrv     *api.Rest
	AdminPasswd string
	Port        int
	terminated  chan struct{}
}

func initApplication() {
	readConfig()
	img.Init()
}

func main() {
	StartDeprecated()
}

func StartDeprecated() {

	var err error
	initApplication()

	db.Offers.S3path = viper.GetString("s3path")
	catalog.InitDB()

	if db.Offers.DB, err = storm.Open("storage/my.db"); err != nil {
		log.Fatalln(err)
	}

	db.Offers.Host = viper.GetString("host")

	if err := execute(); err != nil {
		log.Println("errror: ", err)
	}
}

func (s *server) run(ctx context.Context) error {
	if s.AdminPasswd != "" {
		log.Println("admin basic auth enabled")
	}

	go func() {
		// shutdown on context cancellation
		<-ctx.Done()
		log.Println("shutdown initiated")
		s.restSrv.Shutdown()
		log.Println("shutdown completed")
	}()

	s.restSrv.Run(s.Port)
	close(s.terminated)

	return nil
}

func execute() error {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Println("interrup signal")
		cancel()
	}()

	app := newServerApp()
	log.Println("start server on port", app.Port)

	//if err != nil {
	//	log.Println("failed to setup application", err)
	//	return err
	//}
	if err := app.run(ctx); err != nil {
		log.Println("app terminated with error", err)
		return err
	}

	log.Println("app terminated")

	return nil
}

func newServerApp() *server {
	svr := &api.Rest{
		Version: "0.1",
	}

	return &server{
		restSrv:     svr,
		Port:        viper.GetInt("port"),
		AdminPasswd: viper.GetString("AdminPasswd"),
		terminated:  make(chan struct{}),
	}
}

func readConfig() {
	viper.SetConfigFile("configs/server_config.yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}
