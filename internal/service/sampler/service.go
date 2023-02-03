package sampler

import (
	"go_project_template/internal/logger"
	"go_project_template/internal/repository/sampler"

	"go.uber.org/zap"
)

type Service struct {
	log  logger.AppLogger
	repo *sampler.Repo
}

func InitService(log logger.AppLogger, repo *sampler.Repo) *Service {
	return &Service{
		repo: repo,
		log:  log.With(zap.String("service", "sampler")),
	}
}
