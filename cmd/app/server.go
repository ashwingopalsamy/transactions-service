package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/ashwingopalsamy/transactions-service/internal/handler"
	"github.com/ashwingopalsamy/transactions-service/internal/middleware"
	"github.com/ashwingopalsamy/transactions-service/internal/writer"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

const DefaultWebPort = 8080

type (
	// httpServer defines the HTTP server and its settings
	httpServer struct {
		server   *http.Server
		tracesCh chan string
		errorLog *log.Logger
	}

	// HTTPServer is the interface for starting and stopping the HTTP server
	HTTPServer interface {
		Start() error
		Stop(ctx context.Context) error
	}

	// Options are for configuring the server's behavior (e.g., port, timeouts)
	options struct {
		port           int
		errorLogger    *log.Logger
		readTimeout    time.Duration
		writeTimeout   time.Duration
		disableRecover bool
	}
	Option func(*options)
)

// NewServer creates a new instance of HTTPServer with options
func NewServer(h http.Handler, opts ...Option) HTTPServer {
	setup := options{
		port:         DefaultWebPort,
		readTimeout:  15 * time.Second,
		writeTimeout: 15 * time.Second,
	}

	for _, o := range opts {
		o(&setup)
	}

	tracesCh := make(chan string, runtime.GOMAXPROCS(0))
	if !setup.disableRecover {
		h = withRecovery(h, tracesCh)
	}

	return &httpServer{
		tracesCh: tracesCh,
		errorLog: setup.errorLogger,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", setup.port),
			Handler:      h,
			ErrorLog:     setup.errorLogger,
			ReadTimeout:  setup.readTimeout,
			WriteTimeout: setup.writeTimeout,
		},
	}
}

// Start begins the HTTP server to listen and serve requests
func (h *httpServer) Start() error {
	go h.outputStackTraces()
	return h.server.ListenAndServe()
}

// Stop gracefully shuts down the server
func (h *httpServer) Stop(ctx context.Context) error {
	close(h.tracesCh)
	return h.server.Shutdown(ctx)
}

// outputStackTraces listens for panic recovery and logs stack traces
func (h *httpServer) outputStackTraces() {
	for trace := range h.tracesCh {
		if h.errorLog != nil {
			h.errorLog.Printf("recovered from panic: \n %s", trace)
		}
	}
}

// withRecovery is middleware that recovers from panics and sends them to tracesCh
func withRecovery(next http.Handler, tracesCh chan string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stackTrace := fmt.Sprintf("panic recovered: %v", err)
				tracesCh <- stackTrace
				writer.WriteError(
					w, r.Context(),
					500,
					"internal_server_error",
					"Internal Server Error",
					"Something went wrong. Please try again later",
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// withPort allows you to override the default port.
func withPort(port int) Option {
	return func(o *options) {
		o.port = port
	}
}

// NewRouter creates a new router with all the routes registered
func NewRouter(accHandler *handler.AccountsHandler, trxHandler *handler.TransactionsHandler) http.Handler {
	router := chi.NewRouter()

	// Middlewares
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(middleware.SetRequestIDToContext)

	// Healthcheck route
	router.Get("/health", healthCheckHandler)

	// Account Routes
	router.Route("/v1/accounts", func(r chi.Router) {
		r.Post("/", accHandler.CreateAccount)
		r.Get("/{id}", accHandler.GetAccount)
	})

	// Transaction Routes
	router.Route("/v1/transactions", func(r chi.Router) {
		r.Post("/", trxHandler.CreateTransaction)
	})

	return router
}

func healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	writer.WriteJSON(w, 200, "OK")
}
