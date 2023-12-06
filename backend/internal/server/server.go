package server

import (
	"context"
	"github.com/d-rk/checkin-system/internal/api"
	"github.com/d-rk/checkin-system/internal/checkin"
	"github.com/d-rk/checkin-system/internal/database"
	"github.com/d-rk/checkin-system/internal/user"
	"github.com/d-rk/checkin-system/internal/websocket"
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

func Run() {

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

	if err := checkinService.DeleteOldCheckIns(context.Background()); err != nil {
		panic(err)
	}

	router := setupRouter(userService, checkinService, ws)

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
	router.Use(netHttpMiddleware.OapiRequestValidatorWithOptions(swagger, &validatorOptions))

	// register handler on router
	api.HandlerWithOptions(apiHandler, api.ChiServerOptions{
		BaseRouter:  router,
		Middlewares: []api.MiddlewareFunc{api.AuthMiddleware()},
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
