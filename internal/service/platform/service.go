package platform

import (
	"context"
	pdomain "github.com/ogreks/meeseeks-box/internal/domain/platform"
	"github.com/ogreks/meeseeks-box/internal/model"
	"github.com/ogreks/meeseeks-box/internal/pkg/httpsign"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"go.uber.org/zap"
)

type Service interface {
	i()
	// DetailByAppID get session account by app id
	DetailByAppID(ctx context.Context, appid string) (model.SessionKey, error)
	// Secret give appid return secret
	Secret(ctx context.Context, appid httpsign.KeyID) (*httpsign.Secret, error)
	// SessionsKeysPaginate open api account paginate
	SessionsKeysPaginate(ctx context.Context, req SessionKeysPaginateReq) (*SessionKeyRsp, error)
	// SaveSessionKeys save open api account
	SaveSessionKeys(ctx context.Context, data *SessionKeysReq) (err error)
	// SetStatusSessionKeys set enable/disable by session no
	SetStatusSessionKeys(ctx context.Context, status uint32, sessionNos ...string) (int64, error)
}

type service struct {
	repo    orm.Repo
	logger  *zap.Logger
	pDomain pdomain.SDomain
}

func (s *service) i() {}

func New(repo orm.Repo, logger *zap.Logger) Service {
	return &service{
		repo:    repo,
		logger:  logger,
		pDomain: pdomain.NewSession(repo, logger),
	}
}
