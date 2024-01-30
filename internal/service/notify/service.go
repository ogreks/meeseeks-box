package notify

import (
	"context"
	"github.com/ogreks/meeseeks-box/internal/domain/agreements"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	ser "github.com/ogreks/meeseeks-box/internal/service"
	"go.uber.org/zap"
)

type Service interface {
	i()
	// SaveAgreement edit/save agreement
	SaveAgreement(ctx context.Context, ag AgreementReq) (err error)
	// AgreementDetail get agreement details
	AgreementDetail(ctx context.Context, agno string) (*AgreementRsp, error)
	// FindAgreementVersionByNewTimeDetail Get the latest protocol version information
	FindAgreementVersionByNewTimeDetail(ctx context.Context, agv AgreementVersionReq) (*AgreementVersionRsp, error)
	// AgreementPaginate get agreement paginate list
	AgreementPaginate(ctx context.Context, req AgreementReq) (*ser.PaginateRsp, error)
	// AgreementDelete delete agreement
	AgreementDelete(ctx context.Context, agno string) error
}

type service struct {
	repo   orm.Repo
	logger *zap.Logger
	ag     agreements.AGreementDomain
}

func New(repo orm.Repo, logger *zap.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
		ag:     agreements.NewAGreements(repo, logger),
	}
}

func (s *service) i() {}
