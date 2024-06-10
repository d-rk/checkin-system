package server

import (
	"context"
	"github.com/d-rk/checkin-system/pkg/api"
	"github.com/d-rk/checkin-system/pkg/checkin"
	"github.com/d-rk/checkin-system/pkg/database"
	"github.com/d-rk/checkin-system/pkg/user"
	"github.com/d-rk/checkin-system/pkg/websocket"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	netHttpMiddleware "github.com/oapi-codegen/nethttp-middleware"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func NewRouter(ctx context.Context) chi.Router {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db := database.Connect()

	ws := &websocket.Server{}

	userRepo := user.NewRepo(db)
	checkinRepo := checkin.NewRepo(db)

	userService := user.NewService(userRepo, ws)
	checkinService := checkin.NewService(checkinRepo, userService, ws)

	if err := checkinService.DeleteOldCheckIns(ctx); err != nil {
		panic(err)
	}

	return setupRouter(userService, checkinService, ws)
}

func Run() {

	router := NewRouter(context.Background())

	srv := &http.Server{
		Handler: router,
		Addr:    net.JoinHostPort("0.0.0.0", "8080"),
	}

	// And we serve HTTP until the world ends.
	log.Fatal(srv.ListenAndServe())
}

func setupRouter(userService user.Service, checkinService checkin.Service, ws *websocket.Server) chi.Router {

	router := chi.NewRouter()

	swagger, err := api.GetSwagger()
	if err != nil {
		slog.Error("error loading swagger spec", "err", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	apiHandler := api.NewHandler(userService, checkinService)

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(middleware.Timeout(60 * time.Second))

	router.Use(coreMiddleware())

	validatorOptions := netHttpMiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		},
		ErrorHandler: api.ValidateErrorHandlerFunc,
	}

	// register handler on router
	api.HandlerWithOptions(apiHandler, api.ChiServerOptions{
		BaseRouter:  router,
		Middlewares: []api.MiddlewareFunc{netHttpMiddleware.OapiRequestValidatorWithOptions(swagger, &validatorOptions), api.AuthMiddleware()},
	})

	router.Get("/websocket", websocket.CreateHandler(ws))

	return router
}

func coreMiddleware() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"X-Filename"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}
