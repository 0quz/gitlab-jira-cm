package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/0quz/gitlab-jira-integration/pkg/api"
	"github.com/0quz/gitlab-jira-integration/pkg/middleware"
	"github.com/0quz/gitlab-jira-integration/pkg/model"
	"github.com/0quz/gitlab-jira-integration/pkg/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type App struct {
	Router *http.ServeMux
	DB     *gorm.DB
}

// App initialization
func main() {
	app := App{}
	app.initialize(
		os.Getenv("db_user"),
		os.Getenv("db_pass"),
		os.Getenv("postgres_ip"),
		os.Getenv("db_name"),
	)
	app.routers()
	app.run(":8000")
}

// DB connection
func (a *App) initialize(username, password, host, dbname string) {
	var err error
	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s", username, password, host, dbname)
	a.DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatal(err)
	}
	a.Router = http.NewServeMux()
}

// run app setting
func (a *App) run(addr string) {
	fmt.Printf("Server started at %s\n", addr)
	s := &http.Server{
		Addr:         addr,
		Handler:      a.Router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	fmt.Println("Recieved terminate, graceful shuwdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}

// routers setting
func (a *App) routers() {
	marcoAPI := initMarcoAPI(a.DB)
	a.Router.Handle("/gitlab/merge_request", middleware.MergeRequestMiddleware(marcoAPI.HandleMergeRequest()))
	a.Router.Handle("/jira/close", middleware.CloseMiddleware(marcoAPI.HandleClose()))
	a.Router.Handle("/check", marcoAPI.HandleHealthCheck())
}

// dependency outside to inside
func initMarcoAPI(db *gorm.DB) api.MarcoAPI {
	marcoModel := model.NewMarcoModel(db)
	marcoService := service.NewMarcoService(marcoModel)
	marcoAPI := api.NewMarcoAPI(marcoService)
	return marcoAPI
}
