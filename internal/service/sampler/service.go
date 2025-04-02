package sampler

import (
	"context"
	"go_project_template/internal/logger"
	"go_project_template/internal/repository"
)

type Service struct {
	ctx  context.Context
	log  logger.AppLogger
	repo *repository.Repo
}

func InitService(ctx context.Context, log logger.AppLogger, repo *repository.Repo) *Service {
	return &Service{
		ctx:  ctx,
		repo: repo,
		log:  log.With(logger.WithService("sampler")),
	}
}

func (s *Service) Stop() {
	s.log.Info("stopping service")
}
