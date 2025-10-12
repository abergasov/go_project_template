package routes

import (
	"go_project_template/internal/logger"
	"go_project_template/internal/service/sampler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	appAddr    string
	log        logger.AppLogger
	service    *sampler.Service
	httpEngine *fiber.App
}

// InitAppRouter initializes the HTTP Server.
func InitAppRouter(log logger.AppLogger, service *sampler.Service, address string, enableTelemetry bool) *Server {
	app := &Server{
		appAddr:    address,
		httpEngine: fiber.New(fiber.Config{}),
		service:    service,
		log:        log.With(logger.WithService("http")),
	}
	app.httpEngine.Use(recover.New())
	if enableTelemetry {
		reg := prometheus.NewRegistry()
		reg.MustRegister(
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
			collectors.NewBuildInfoCollector(),
		)
		app.httpEngine.Get("/metrics", adaptor.HTTPHandler(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))
	}
	app.initRoutes()
	return app
}

func (s *Server) initRoutes() {
	s.httpEngine.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
	})
}

// Run starts the HTTP Server.
func (s *Server) Run() error {
	s.log.Info("Starting HTTP server", logger.WithString("port", s.appAddr))
	return s.httpEngine.Listen(s.appAddr)
}

func (s *Server) Stop() error {
	return s.httpEngine.Shutdown()
}
