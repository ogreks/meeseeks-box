package platform

import (
	"context"
	"errors"
	"github.com/ogreks/meeseeks-box/internal/dao"
	"github.com/ogreks/meeseeks-box/internal/domain"
	"github.com/ogreks/meeseeks-box/internal/model"
	"github.com/ogreks/meeseeks-box/internal/pkg/httpsign"
	"github.com/ogreks/meeseeks-box/internal/pkg/httpsign/secret"
)

var (
	ErrAppIdNotFound = errors.New("appid not found")
	ErrAppDisabled   = errors.New("appid is disable")
)

// DetailByAppID get session account by app id
func (s *service) DetailByAppID(ctx context.Context, appid string) (model.SessionKey, error) {
	data, err := s.pDomain.SessionKey(ctx, dao.SessionKey.AppID.Eq(appid))
	if errors.Is(err, domain.NotFound) {
		err = ErrAppIdNotFound
	}

	return *data, err
}

// Secret give appid return secret
func (s *service) Secret(ctx context.Context, appid httpsign.KeyID) (*httpsign.Secret, error) {
	data, err := s.DetailByAppID(ctx, string(appid))
	if err != nil {
		return nil, err
	}

	if data.IsEnabled == 0 {
		return nil, ErrAppDisabled
	}

	return &httpsign.Secret{
		Key:       data.AppSecret,
		Algorithm: &secret.HmacSha512{},
	}, nil
}
