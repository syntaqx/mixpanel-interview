package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/syntaqx/chix"
	"github.com/unrolled/render"
	"go.uber.org/zap"
)

var (
	debug = kingpin.Flag("debug", "Enable debug mode").Bool()
	host  = kingpin.Flag("host", "HTTP listening network host").Envar("HOST").String()
	port  = kingpin.Flag("port", "HTTP listening network port").Default("8080").Envar("PORT").String()
)

func main() {
	kingpin.Parse()

	// Initialize our logger before anything.
	log, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("[mixpanel] %v", err))
	}

	// If we are in debug mode, we should use the more readable dev logger.
	if *debug == true {
		devLogger, err := zap.NewDevelopment()
		if err != nil {
			panic(fmt.Sprintf("[mixpanel] %v", err))
		}
		log = devLogger
	}

	// Initialize termination channel to be provided signal notifications.
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Rendering handler.
	renderer := render.New(render.Options{
		Directory:     "templates",
		Layout:        "layout",
		Extensions:    []string{".html"},
		IsDevelopment: *debug,
	})

	// Routing muxer
	r := chi.NewRouter()

	// Global middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(chix.NewZapLogger(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Heartbeat("/heartbeat"))

	// When a closes their connection midway through a request, the
	// http.CloseNotifier will cancel the request context (ctx).
	r.Use(middleware.CloseNotify)

	// Sets a timeout value on the request context that will signal via ctx.Done()
	// that the request has timed out and further processing should be haulted.
	r.Use(middleware.Timeout(60 * time.Second))

	// Simple Health route.
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		renderer.Text(w, http.StatusOK, "Hello World")
	})

	// Initialize our HTTP server with the configured address.
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(*host, *port),
		Handler: r,
	}

	go func() {
		log.Info("starting http server", zap.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Error("http server closed unexpectedly", zap.Error(err))
			}
		}
	}()

	// Block the process from exiting until a signal is provided.
	<-term
}
