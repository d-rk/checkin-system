package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/d-rk/checkin-system/pkg/version"
	"github.com/d-rk/checkin-system/pkg/wifi"

	"github.com/d-rk/checkin-system/pkg/api"
	"github.com/d-rk/checkin-system/pkg/checkin"
	"github.com/d-rk/checkin-system/pkg/clock"
	"github.com/d-rk/checkin-system/pkg/database"
	"github.com/d-rk/checkin-system/pkg/user"
	"github.com/d-rk/checkin-system/pkg/websocket"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	netHttpMiddleware "github.com/oapi-codegen/nethttp-middleware"
)

const defaultMaxAge = 300
const defaultTimeout = 60 * time.Second
const defaultReadHeaderTimeout = 10 * time.Second

func NewDB(runMigration bool) *sqlx.DB {

	_ = godotenv.Load(".env")
	db := database.Connect()

	if runMigration {
		database.RunMigration(db)
	}

	return db
}

func NewRouter(ctx context.Context, db *sqlx.DB) chi.Router {

	ws := &websocket.Server{}

	userRepo := user.NewRepo(db)
	checkinRepo := checkin.NewRepo(db)

	userService := user.NewService(userRepo, ws)
	checkinService := checkin.NewService(checkinRepo, userService, ws)
	clockService := clock.NewService()
	wifiService := wifi.NewService()

	if err := checkinService.DeleteOldCheckIns(ctx); err != nil {
		slog.WarnContext(ctx, "failed to delete old checkins", "error", err)
	}

	return setupRouter(userService, checkinService, clockService, wifiService, ws)
}

func Run() {

	slog.Info("starting up", "version", version.Version, "build_time", version.BuildTime,
		"git_commit", version.GitCommit)

	db := NewDB(true)
	defer db.Close()

	router := NewRouter(context.Background(), db)

	srv := &http.Server{
		Handler:           router,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		Addr:              net.JoinHostPort("0.0.0.0", "8080"),
	}

	// And we serve HTTP until the world ends.
	err := srv.ListenAndServe()
	slog.Info("server stopped", "err", err)
}

func setupRouter(
	userService user.Service,
	checkinService checkin.Service,
	clockService clock.Service,
	wifiService wifi.Service,
	ws *websocket.Server,
) chi.Router {

	router := chi.NewRouter()

	swagger, err := api.GetSwagger()
	if err != nil {
		slog.Error("error loading swagger spec", "err", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	apiHandler := api.NewHandler(userService, checkinService, clockService, wifiService)

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(middleware.Timeout(defaultTimeout))

	router.Use(coreMiddleware())

	validatorOptions := netHttpMiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		},
		ErrorHandler: api.ValidateErrorHandlerFunc,
	}

	// register handler on router
	api.HandlerWithOptions(apiHandler, api.ChiServerOptions{
		BaseRouter: router,
		Middlewares: []api.MiddlewareFunc{
			netHttpMiddleware.OapiRequestValidatorWithOptions(swagger, &validatorOptions),
			api.AuthMiddleware(),
		},
	})

	router.Get("/websocket", websocket.CreateHandler(ws))

	_ = chi.Walk(router, func(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		slog.Info("registered route", "method", method, "route", route)
		return nil
	})

	return router
}

func coreMiddleware() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"X-Filename"},
		AllowCredentials: false,
		MaxAge:           defaultMaxAge,
	})
}
