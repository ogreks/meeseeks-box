package user

import (
	"context"
	"errors"
	"github.com/ogreks/meeseeks-box/internal/domain"
	"time"

	"github.com/ogreks/meeseeks-box/internal/model"
	"github.com/ogreks/meeseeks-box/internal/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorAccountOrPassword = errors.New("account or password error")
	ErrorAccountNotEnable  = errors.New("account is not enable")
	ErrorAccountNotFound   = errors.New("account not found")
	ErrorUserNotFound      = errors.New("user not found")
)

type UserAccount struct {
	Aid         string    `json:"account_id"`
	LastLoginAt time.Time `json:"last_login_at"`
}

// CreateUserByUserName create user by userName
func (s *service) CreateUserByUserName(ctx context.Context, uid, userName, password string) error {
	modelAccount := model.Account{
		Aid:      uid,
		UserName: userName,
		Password: password,
	}

	modelUser := model.User{
		UserName: userName,
		Password: password,
	}

	return s.domain.CreateUser(ctx, modelAccount, modelUser)
}

// CreateUserByEmail create user by email
func (s *service) CreateUserByEmail(ctx context.Context, uid, email, password string) error {
	userName := utils.CreateUserName()
	modelAccount := model.Account{
		Aid:      uid,
		UserName: userName,
		Email:    email,
		Password: password,
	}

	modelUser := model.User{
		UserName: userName,
		Password: password,
	}

	return s.domain.CreateUser(ctx, modelAccount, modelUser)
}

// GetUserByUserName get user by userName
func (s *service) GetUserByUserName(ctx context.Context, userName string, password string) (*UserAccount, error) {
	account, err := s.domain.FindAccount(ctx, model.Account{
		UserName: userName,
	})

	if err != nil {
		if errors.Is(err, domain.NotFound) {
			return nil, ErrorAccountOrPassword
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		return nil, ErrorAccountOrPassword
	}

	if *account.IsEnabled == 0 {
		return nil, ErrorAccountNotEnable
	}

	// update last login time
	go func(aid uint64, t time.Time) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = s.domain.UpdateLastLoginAtByAccountId(ctx, aid, t)
	}(account.ID, time.Now())

	return &UserAccount{
		Aid:         account.Aid,
		LastLoginAt: *account.LastLoginAt,
	}, nil
}
