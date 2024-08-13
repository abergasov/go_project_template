package sampler

import (
	"go_project_template/internal/logger"
	"go_project_template/internal/repository"
)

type Service struct {
	log  logger.AppLogger
	repo *repository.Repo
}

func InitService(log logger.AppLogger, repo *repository.Repo) *Service {
	return &Service{
		repo: repo,
		log:  log.With(logger.WithService("sampler")),
	}
}
