package platform

import (
	"context"
	"errors"
	"github.com/ogreks/meeseeks-box/internal/dao"
	"github.com/ogreks/meeseeks-box/internal/domain"
	"github.com/ogreks/meeseeks-box/internal/model"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"go.uber.org/zap"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type SDomain interface {
	i()
	// SaveSessionKey create/update other open request account
	SaveSessionKey(ctx context.Context, data model.SessionKey) error
	// SessionKey get other open account
	SessionKey(ctx context.Context, where ...gen.Condition) (*model.SessionKey, error)
	// SessionKeyPaginate paginate list
	SessionKeyPaginate(ctx context.Context, page, limit int, where ...gen.Condition) ([]*model.SessionKey, int64, error)
	// UpdateSessionKey update open api request accounts
	UpdateSessionKey(ctx context.Context, data model.SessionKey, where ...gen.Condition) (int64, error)
}

type Session struct {
	dao *dao.Query
	log *zap.Logger
}

func NewSession(repo orm.Repo, log *zap.Logger) *Session {
	return &Session{
		dao: dao.Use(repo.GetDB()),
		log: log,
	}
}

func (s *Session) i() {}

// SaveSessionKey create/update other open request account
func (s *Session) SaveSessionKey(ctx context.Context, data model.SessionKey) error {
	return s.dao.SessionKey.WithContext(ctx).Save(&data)
}

// UpdateSessionKey update open api request accounts
func (s *Session) UpdateSessionKey(ctx context.Context, data model.SessionKey, where ...gen.Condition) (int64, error) {
	result, err := s.dao.SessionKey.WithContext(ctx).Where(where...).Updates(&data)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected, result.Error
}

// SessionKey get other open account
func (s *Session) SessionKey(ctx context.Context, where ...gen.Condition) (*model.SessionKey, error) {
	data, err := s.dao.WithContext(ctx).SessionKey.Where(where...).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = domain.NotFound
	}

	return data, err
}

// SessionKeyPaginate paginate list
func (s *Session) SessionKeyPaginate(ctx context.Context, page, limit int, where ...gen.Condition) ([]*model.SessionKey, int64, error) {
	return s.dao.WithContext(ctx).
		SessionKey.
		Where(where...).
		FindByPage((page-1)*limit, limit)
}
