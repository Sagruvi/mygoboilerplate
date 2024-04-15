package run

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log"
	_ "mygoboilerplate/docs" // docs is generated by Swag CLI, you have to import it.
	"mygoboilerplate/internal/auth/controller"
	controller2 "mygoboilerplate/internal/geolocation/controller"
	"mygoboilerplate/internal/geolocation/repository"
	"mygoboilerplate/internal/geolocation/service"
	"mygoboilerplate/internal/metrics"
	"net/http"
	"net/http/httputil"
	"net/http/pprof"
	_ "net/http/pprof"
	"net/url"
	"os"
	"os/signal"
	"time"
)

type server struct {
	port string
	s    *http.Server
}

// method serve
func (s server) Serve() {
	log.Printf("Starting server on port %s\n", s.port)
	err := s.s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
func (s server) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}

var users = make(map[string]string)

//	@title			Dadata API Proxy
//	@version		1.0
//	@description	This is a sample server geolocation service.

// @host	localhost:8080/
// @BasePath	/api/address
func Run() {
	protocol := os.Getenv("PROTOCOL")
	newAuthController := controller.NewController("secret")
	repo := repository.NewRepository()
	repo.Migrate(context.Background())
	service := service.NewService(repo)
	newGeolocationController := controller2.NewController(service)
	hugoURL, err := url.Parse("http://hugo:1313")
	if err != nil {
		panic(err)
	}
	hugoProxy := httputil.NewSingleHostReverseProxy(hugoURL)

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Get("/register", func(writer http.ResponseWriter, request *http.Request) {
			newAuthController.Register(writer, request)
		})
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		))
		r.Get("/api/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello from API"))
		})
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			hugoProxy.ServeHTTP(w, r)
		})

	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(newAuthController.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/login", func(writer http.ResponseWriter, request *http.Request) {
			newAuthController.Login(writer, request)
		})
		r.Post("/api/address/geocode", func(writer http.ResponseWriter, request *http.Request) {
			newGeolocationController.Geocode(writer, request)
		})
		r.Post("/api/address/search", func(writer http.ResponseWriter, request *http.Request) {
			newGeolocationController.Search(writer, request)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(newAuthController.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)
		r.Get("/debug/pprof/goroutine", pprof.Index)
		r.Get("/debug/pprof/heap", pprof.Index)
		r.Get("/debug/pprof/threadcreate", pprof.Index)
		r.Get("/debug/pprof/block", pprof.Index)
		r.Get("/debug/pprof/mutex", pprof.Index)
	})
	var srv server = server{
		port: "8080",
		s: &http.Server{
			Addr:    ":8080",
			Handler: r,
		},
	}

	go func() {
		r, err := controller2.RPCGW{}.GetFactory(protocol, &service)
		if err != nil {
			log.Fatal(err)
		}
		r.CreateGateway().Run("1234")
	}()
	go srv.Serve()
	go metrics.PrometheusMiddleware()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
	// Установка тайм-аута для завершения работы
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")

}

type ReverseProxy struct {
	host string
	port string
}

func NewReverseProxy(host, port string) *ReverseProxy {
	return &ReverseProxy{
		host: host,
		port: port,
	}
}

func (rp *ReverseProxy) ReverseProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://"+rp.host+":"+rp.port+r.RequestURI, http.StatusFound)
}

// @contact.name	API Support
// @contact.url	https://github.com/go-chi/chi/issues
// @contact.email	6z6o8@example.com

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
