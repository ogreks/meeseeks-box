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
	ErrorAccountParamsNotFound  = errors.New("Missing authentication account parameters")
	ErrorAccountConnectNotFound = errors.New("The authenticated user does not exist")
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

	return s.domain.CreatePlatformAccount(
		ctx,
		model.Account{
			Aid:          account.Aid,
			Type:         3,
			CountryCode:  &countryCode,
			Email:        account.Email,
			UserName:     utils.CreateUserName(),
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
			UserName:       account.UserName,
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

	return &UserAccount{
		Aid: ac.Aid,
	}, nil
}
