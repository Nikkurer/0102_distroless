package main

import (
	"app/pkg/models"
	"app/pkg/svc/events"
	"context"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	host  = flag.String("host", "0.0.0.0", "server host")
	port  = flag.String("port", "9000", "server port")
	dbUrl = flag.String("dbUrl", "postgres://user:pass@localhost:5432/app", "db url")
)

func main() {
	var err error

	log.Print("starting: parse flags")
	flag.Parse()

	log.Print("starting: create connection pool")
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, *dbUrl)
	if err != nil {
		log.Println(err)
		return
	}
	defer pool.Close()

	log.Print("starting: create service")
	svc := events.NewSvc(pool)

	log.Print("starting: create router")
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Post("/api/events", func(writer http.ResponseWriter, request *http.Request) {
		model := &models.Event{
			Action:      request.PostFormValue("action"),
			Product:     request.PostFormValue("product"),
			Fingerprint: request.PostFormValue("fingerprint"),
		}

		if err := svc.Register(request.Context(), model); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		writer.WriteHeader(http.StatusOK)
	})

	addr := net.JoinHostPort(*host, *port)
	log.Printf("starting: listen and serve on: %s", addr)
	server := &http.Server{Addr: addr, Handler: router}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-sigChannel
		log.Printf("stopping: got %s signal", sig)

		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("stopping: graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}
