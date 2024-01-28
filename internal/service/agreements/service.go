package agreements

import (
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"go.uber.org/zap"
)

type Service interface {
	i()
}

type service struct {
	repo   orm.Repo
	logger *zap.Logger
}

func New(repo orm.Repo, logger *zap.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) i() {}
