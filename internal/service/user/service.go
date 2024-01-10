package user

import (
	"context"

	"github.com/ogreks/meeseeks-box/internal/domain/user"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"go.uber.org/zap"
)

type Service interface {
	i()

	// CreateUserByUserName create user by userName
	CreateUserByUserName(ctx context.Context, uid, userName, password string) error
	// CreateUserByEmail create user by email
	CreateUserByEmail(ctx context.Context, uid, email, password string) error
	// GetUserByUserName get user by userName
	GetUserByUserName(ctx context.Context, userName string, password string) (*UserAccount, error)

	// LoginUserByGITHub login by github account or register user
	LoginUserByGITHub(ctx context.Context, account AccountPlatform) (*UserAccount, error)
}

type service struct {
	repo   orm.Repo
	logger *zap.Logger
	domain user.UDomain
}

func New(repo orm.Repo, logger *zap.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
		domain: user.New(repo, logger),
	}
}

func (s *service) i() {}
