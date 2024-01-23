package user

import (
	"context"
	"errors"
	"time"

	udomain "github.com/ogreks/meeseeks-box/internal/domain/user"
	"github.com/ogreks/meeseeks-box/internal/model"
	"github.com/ogreks/meeseeks-box/internal/pkg/utils"
)

var (
	ErrorAccountParamsNotFound  = errors.New("missing authentication account parameters")
	ErrorAccountConnectNotFound = errors.New("the authenticated user does not exist")
)

type AccountPlatform struct {
	Aid string

	PlatformID uint32
	AccountID  string

	Token                string
	RefreshToken         string
	RefreshTokenExpireAt *time.Time

	UserName string
	NickName string
	Email    string
	PurPhone *string
	Phone    string

	MoreJson string
}

func (s *service) CreatePlatformAccount(ctx context.Context, account AccountPlatform) error {
	_, err := s.domain.FindAccountConnect(ctx, account.PlatformID, account.AccountID)
	if err == nil {
		return nil
	}

	if !errors.Is(err, udomain.AccountConnectNotFound) {
		return err
	}

	if account.Aid == "" {
		return ErrorAccountConnectNotFound
	}

	// create connect account user
	countryCode := "1"
	waitDeleteAt := time.Now().AddDate(0, 0, 30)
	t := time.Now()

	userName := utils.CreateUserName()
	return s.domain.CreatePlatformAccount(
		ctx,
		model.Account{
			Aid:          account.Aid,
			Type:         3,
			CountryCode:  &countryCode,
			Email:        account.Email,
			UserName:     userName,
			WaitDelete:   1,
			WaitDeleteAt: &waitDeleteAt,
			LastLoginAt:  &t,
		},
		model.AccountConnect{
			ConnectPlatformID:    1,
			ConnectAccountID:     account.AccountID,
			ConnectToken:         account.Token,
			ConnectRefreshToken:  account.RefreshToken,
			ConnectUserName:      account.UserName,
			ConnectNickName:      account.NickName,
			MoreJSON:             account.MoreJson,
			RefreshTokenExpireAt: account.RefreshTokenExpireAt,
		},
		model.User{
			UserName:       userName,
			NickName:       account.NickName,
			LastActivityAt: &t,
			WaitDelete:     1,
			WaitDeleteAt:   &waitDeleteAt,
		},
	)
}

// LoginUserByGITHub login by github account
func (s *service) LoginUserByGITHub(ctx context.Context, account AccountPlatform) (*UserAccount, error) {
	if account.PlatformID == 0 || account.AccountID == "" {
		return nil, ErrorAccountParamsNotFound
	}

	err := s.CreatePlatformAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	accountConnect, err := s.domain.FindAccountConnect(ctx, account.PlatformID, account.AccountID)
	if err != nil {
		if errors.Is(err, udomain.AccountConnectNotFound) {
			return nil, ErrorAccountConnectNotFound
		}

		return nil, err
	}

	ac, err := s.domain.FindAccount(ctx, model.Account{
		ID: accountConnect.AccountID,
	})
	if err != nil {
		if errors.Is(err, udomain.AccountNotFound) {
			return nil, ErrorAccountNotEnable
		}

		return nil, err
	}

	// update last login time
	go func(aid uint64, t time.Time) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_ = s.domain.UpdateLastLoginAtByAccountId(ctx, aid, t)
	}(ac.ID, time.Now())

	return &UserAccount{
		Aid:         ac.Aid,
		LastLoginAt: *ac.LastLoginAt,
	}, nil
}

type User struct {
	AccountID string `json:"account_id"`
	Email     string `json:"email"`

	CountryCode string `json:"country_code"`
	PurePhone   string `json:"pure_phone"`
	Phone       string `json:"phone"`

	UserName string `json:"user_name"`

	RegisterTime time.Time `json:"register_time"`

	NickName string `json:"nick_name"`
	Gender   uint32 `json:"gender"`
	Bio      string `json:"bio"`
}

// GetUserByAccountAid get user info by account aid
func (s *service) GetUserByAccountAid(ctx context.Context, aid string) (*User, error) {
	account, err := s.domain.FindAccount(ctx, model.Account{Aid: aid})
	if err != nil {
		if errors.Is(err, udomain.AccountNotFound) {
			return nil, ErrorAccountNotFound
		}
		return nil, err
	}

	user, err := s.domain.FindUser(ctx, model.User{
		AccountID: account.ID,
	})
	if err != nil {
		if errors.Is(err, udomain.UserNotFound) {
			return nil, ErrorUserNotFound
		}
		return nil, err
	}

	return &User{
		AccountID:    account.Aid,
		Email:        account.Email,
		CountryCode:  *account.CountryCode,
		PurePhone:    account.PurePhone,
		Phone:        account.Phone,
		UserName:     account.UserName,
		RegisterTime: *account.CreatedAt,

		NickName: user.NickName,
		Gender:   *user.Gender,
		Bio:      user.Bio,
	}, nil
}
