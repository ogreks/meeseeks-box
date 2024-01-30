package user

import (
	"context"
	"errors"
	"github.com/ogreks/meeseeks-box/internal/domain"
	"time"

	"github.com/ogreks/meeseeks-box/internal/dao"
	"github.com/ogreks/meeseeks-box/internal/model"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"go.uber.org/zap"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

var _ UDomain = (*User)(nil)

type UDomain interface {
	CreateUser(ctx context.Context, account model.Account, user model.User) error
	CreatePlatformAccount(ctx context.Context, account model.Account, accountConnect model.AccountConnect, user model.User) error

	FindUser(ctx context.Context, user model.User) (*model.User, error)
	FindAccount(ctx context.Context, account model.Account) (*model.Account, error)
	FindAccountConnect(ctx context.Context, connectID uint32, connectAccountID string) (*model.AccountConnect, error)

	UpdateLastLoginAtByAccountId(ctx context.Context, Aid uint64, t time.Time) error
}

type User struct {
	dao *dao.Query
	log *zap.Logger
}

func New(repo orm.Repo, log *zap.Logger) *User {
	return &User{
		log: log,
		dao: dao.Use(repo.GetDB()),
	}
}

func (s *User) i() {}

// CreateUser create user by userName
func (s *User) CreateUser(ctx context.Context, account model.Account, user model.User) error {
	err := s.dao.Transaction(func(tx *dao.Query) error {
		err := tx.Account.WithContext(ctx).Create(&account)
		if err != nil {
			return err
		}

		if user.AccountID == 0 {
			user.AccountID = account.ID
		}

		return tx.User.WithContext(ctx).Create(&user)
	})

	return err
}

// CreatePlatformAccount create user by connect open account
func (s *User) CreatePlatformAccount(ctx context.Context, account model.Account, accountConnect model.AccountConnect, user model.User) error {
	return s.dao.Transaction(func(tx *dao.Query) error {
		err := tx.Account.WithContext(ctx).Create(&account)
		if err != nil {
			return err
		}

		if accountConnect.AccountID == 0 {
			accountConnect.AccountID = account.ID
		}

		err = tx.AccountConnect.WithContext(ctx).Create(&accountConnect)
		if err != nil {
			return err
		}

		if user.AccountID == 0 {
			user.AccountID = account.ID
		}

		return tx.User.WithContext(ctx).Create(&user)
	})
}

// UpdateLastLoginAtByAccountId update last login time by account id
func (s *User) UpdateLastLoginAtByAccountId(ctx context.Context, Aid uint64, t time.Time) error {
	err := s.dao.Transaction(func(tx *dao.Query) error {
		_, err := tx.Account.WithContext(ctx).
			Where(tx.Account.ID.Eq(Aid)).
			Update(tx.Account.LastLoginAt, t)

		if err != nil {
			return err
		}

		_, err = tx.User.WithContext(ctx).Where(tx.User.AccountID.Eq(Aid)).Update(tx.User.LastActivityAt, t)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// FindAccount get user by account
func (s *User) FindAccount(ctx context.Context, account model.Account) (*model.Account, error) {
	a, err := s.dao.Account.WithContext(ctx).Where(field.Attrs(account)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.NotFound
	}

	return a, err
}

// FindUser get user by user
func (s *User) FindUser(ctx context.Context, user model.User) (*model.User, error) {
	u, err := s.dao.User.WithContext(ctx).Where(field.Attrs(user)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.NotFound
	}

	return u, err
}

// FindAccountConnect get connect open account
func (s *User) FindAccountConnect(ctx context.Context, connectID uint32, connectAccountID string) (*model.AccountConnect, error) {
	accountConnect, err := s.dao.AccountConnect.
		WithContext(ctx).
		Where(s.dao.AccountConnect.ConnectPlatformID.Eq(connectID)).
		Where(s.dao.AccountConnect.ConnectAccountID.Eq(connectAccountID)).
		First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.NotFound
	}

	return accountConnect, err
}
