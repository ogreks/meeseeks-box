package notify

import (
	"context"
	"errors"
	"fmt"
	"github.com/ogreks/meeseeks-box/internal/dao"
	"github.com/ogreks/meeseeks-box/internal/domain"
	"github.com/ogreks/meeseeks-box/internal/model"
	ser "github.com/ogreks/meeseeks-box/internal/service"
	"github.com/rs/xid"
	"gorm.io/gen"
	"time"
)

var (
	ErrDeleteAgreement = errors.New("agreement deletion failed")
)

type AgreementReq struct {
	ser.PaginateReq

	AgreementNo string     `json:"agreement_no"`
	Type        int32      `json:"type"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Status      uint32     `json:"status"`
	Version     string     `json:"version"`
	PublishAt   *time.Time `json:"publish_at"`

	StartPublishAt *time.Time `json:"start_time"`
	EndPublishAt   *time.Time `json:"end_time"`
}

type AgreementRsp struct {
	AgreementNo string     `json:"agreement_no"`
	Type        int32      `json:"type"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Status      uint32     `json:"status"`
	Version     string     `json:"version"`
	PublishAt   *time.Time `json:"publish_at"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type AgreementPagateRsp struct {
	ser.PaginateRsp
}

type AgreementVersionReq struct {
	AgreementNo string
	Version     string
}

type AgreementVersionRsp struct {
	AgreementNo string     `json:"agreement_no"`
	Type        int32      `json:"type"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Version     string     `json:"version"`
	CreatedAt   *time.Time `json:"created_at"`
}

// SaveAgreement edit/save agreement
func (s *service) SaveAgreement(ctx context.Context, ag AgreementReq) (err error) {
	m := model.Agreement{
		Type:      ag.Type,
		Title:     ag.Title,
		Content:   ag.Content,
		Status:    ag.Status,
		Version:   ag.Version,
		PublishAt: ag.PublishAt,
	}

	if ag.AgreementNo != "" {
		agr, err := s.ag.Agreement(ctx, dao.Agreement.AgreementNo.Eq(ag.AgreementNo))
		if err != nil {
			if errors.Is(err, domain.NotFound) {
				return ser.ErrDataNotFound
			}

			return err
		}

		m.ID = agr.ID
		m.AgreementNo = agr.AgreementNo
	}

	if m.AgreementNo == "" {
		m.AgreementNo = xid.New().String()
	}

	return s.ag.SaveAGreement(ctx, m)
}

// FindAgreementVersionByNewTimeDetail Get the latest protocol version information
func (s *service) FindAgreementVersionByNewTimeDetail(ctx context.Context, agv AgreementVersionReq) (*AgreementVersionRsp, error) {
	var where []gen.Condition

	if agv.AgreementNo != "" {
		where = append(where, dao.AgreementVersion.AgreementNo.Eq(agv.AgreementNo))
	}

	if agv.Version != "" {
		where = append(where, dao.AgreementVersion.Version.Eq(agv.Version))
	}

	if len(where) == 0 {
		return nil, ser.ErrMissRequiredParams
	}

	version, err := s.ag.AgreementVersionNew(ctx, where...)
	if err != nil {
		if errors.Is(err, domain.NotFound) {
			return nil, ser.ErrDataNotFound
		}

		return nil, err
	}

	return &AgreementVersionRsp{
		AgreementNo: version.AgreementNo,
		Type:        version.Type,
		Title:       version.Title,
		Content:     version.Content,
		Version:     version.Version,
		CreatedAt:   version.CreatedAt,
	}, nil
}

// AgreementDetail get agreement details
func (s *service) AgreementDetail(ctx context.Context, agno string) (*AgreementRsp, error) {
	agreement, err := s.ag.Agreement(ctx, dao.Agreement.AgreementNo.Eq(agno))
	if err != nil {
		if errors.Is(err, domain.NotFound) {
			return nil, ser.ErrDataNotFound
		}

		return nil, err
	}

	return &AgreementRsp{
		AgreementNo: agreement.AgreementNo,
		Type:        agreement.Type,
		Title:       agreement.Title,
		Content:     agreement.Content,
		Version:     agreement.Version,
		CreatedAt:   agreement.CreatedAt,
		UpdatedAt:   agreement.UpdatedAt,
	}, nil
}

// AgreementPaginate get agreement paginate list
func (s *service) AgreementPaginate(ctx context.Context, req AgreementReq) (*ser.PaginateRsp, error) {
	var where []gen.Condition

	if req.Type != 0 {
		where = append(where, dao.Agreement.Type.Eq(req.Type))
	}

	if req.Title != "" {
		where = append(where, dao.Agreement.Title.Like(fmt.Sprintf("%%%s%%", req.Title)))
	}

	if req.Version != "" {
		where = append(where, dao.Agreement.Version.Eq(req.Version))
	}

	if req.AgreementNo != "" {
		where = append(where, dao.Agreement.AgreementNo.Eq(req.AgreementNo))
	}

	//if req.PublishAt != nil {
	//	where = append(where,
	//		dao.Agreement.PublishAt.Gte(),
	//		dao.Agreement.PublishAt.Lte(),
	//	)
	//}

	if req.StartPublishAt != nil {
		where = append(where, dao.Agreement.PublishAt.Gte(*req.StartPublishAt))
	}

	if req.EndPublishAt != nil {
		where = append(where, dao.Agreement.PublishAt.Lte(*req.EndPublishAt))
	}

	data, total, err := s.ag.Paginate(ctx, req.Page, req.Limit, where...)
	if err != nil {
		return nil, err
	}

	var list = make([]*AgreementRsp, len(data))
	for key, agreement := range data {
		list[key] = &AgreementRsp{
			AgreementNo: agreement.AgreementNo,
			Type:        agreement.Type,
			Title:       agreement.Title,
			Version:     agreement.Version,
			CreatedAt:   agreement.CreatedAt,
			UpdatedAt:   agreement.UpdatedAt,
		}
	}

	return ser.NewPaginateRsp(list, total, req.PaginateReq), nil
}

// AgreementDelete delete agreement
func (s *service) AgreementDelete(ctx context.Context, agno string) error {
	n, err := s.ag.AgreementDelete(ctx, dao.Agreement.AgreementNo.Eq(agno))
	if err != nil {
		return err
	}

	if n <= 0 {
		return ErrDeleteAgreement
	}

	return nil
}
