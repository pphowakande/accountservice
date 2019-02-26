// Poonam Phowakande
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

//Server Local HTTP server to server GraphQL queries
type Server struct {
	logger *log.Logger
	mux    *http.ServeMux
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// App - Application context
type App struct {
	Server   *Server
	Hostname string
	IP       string
	Module   string
	Port     int
	Router   *mux.Router
	logger   *log.Entry
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	w.Write(page)
}

// Initialize the environment for our app
func (a *App) Initialize() {
	a.Hostname, _ = os.Hostname()
	a.Module = "Account"
	a.Port = *port

	// All Log enteries will contain these two fields
	a.logger = log.WithFields(log.Fields{
		"Host":             a.Hostname,
		"Module":           a.Module,
		"Port":             a.Port,
		"X-Correlation-Id": uuid.Must(uuid.NewV4()).String(),
	})
	a.logger.Info(a.Module, "_Start")
	a.logger.Info("Rock n Roll")

	a.Router = mux.NewRouter()

	a.initializeBoltClient()

	a.logger.Info("Loaded Test Data")

	// Always set the latest routes
	a.initializeRoutes(a.Router)

}

// Run starts http service on given addr and responding to PACKAGE queries
func (a *App) Run(addr string) {
	a.Server = &Server{
		logger: log.New(),
		mux:    http.NewServeMux(),
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	h := &http.Server{Addr: addr, Handler: a.Server}

	go func() {
		log.Fatal(http.ListenAndServe(addr, a.Router))
		if err := h.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop
	// a.DB.Close()
	log.Println("Shutting down the Account server...")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	log.Info(a.Module, "_Stop")

	h.Shutdown(ctx)

}

func (a *App) initializeRoutes(r *mux.Router) {
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotImplemented)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Method not implemented",
		})
	})

	// endpoint to authenticate user/client
	r.HandleFunc("/authenticate", headerSetter(a.requestLogger(a.CreateToken))).Methods("POST")

	// public:create account
	r.HandleFunc("/account", headerSetter(a.requestLogger(a.ValidateMiddleware(a.NewAccount)))).Methods("POST")

	// public:return account details using id
	r.HandleFunc("/account/{id}", headerSetter(a.requestLogger(a.ValidateMiddleware(a.GetAccount)))).Methods("GET")

}

func (a *App) initializeBoltClient() {

	DBClient = &BoltClient{}

	DBClient.OpenBoltDb()

	DBClient.Seed()

}

var page = []byte(`
<!DOCTYPE html>
<html>
	<head>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<h1> Welcome to Golang Rest API</h1>
	</body>
</html>
`)
