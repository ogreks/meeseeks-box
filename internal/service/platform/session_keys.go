package platform

import (
	"context"
	"fmt"
	"github.com/ogreks/meeseeks-box/internal/dao"
	"github.com/ogreks/meeseeks-box/internal/model"
	"github.com/ogreks/meeseeks-box/internal/pkg/utils"
	common "github.com/ogreks/meeseeks-box/internal/service"
	"github.com/rs/xid"
	"gorm.io/gen"
	"time"
)

type SessionKeysReq struct {
	SessionNo string
	Name      string
	Type      uint32
	IsEnabled uint32
	Remark    string
}

type SessionKeysPaginateReq struct {
	common.PaginateReq

	ID        int64
	SessionNo string
	Name      string
	Type      uint32
	AppID     string
	IsEnabled uint32
	StartTime *time.Time
	EndTime   *time.Time
}

type SessionKeyRsp struct {
	common.PaginateRsp
}

type SessionKeyDetail struct {
	SessionNo string    `json:"session_no"`
	Name      string    `json:"name"`
	Type      uint32    `json:"type"`
	AppID     string    `json:"app_id"`
	IsEnabled uint32    `json:"is_enabled"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SaveSessionKeys save open api account
func (s *service) SaveSessionKeys(ctx context.Context, data *SessionKeysReq) (err error) {
	m := model.SessionKey{
		Name:      data.Name,
		Type:      data.Type,
		IsEnabled: data.IsEnabled,
		Remark:    data.Remark,
	}
	if data.SessionNo != "" {
		session, _ := s.pDomain.SessionKey(ctx, dao.SessionKey.SessionNo.Eq(data.SessionNo))
		if session != nil {
			m.AppID = session.AppID
			m.AppSecret = session.AppSecret
			m.SessionNo = session.SessionNo
		}
	}

	if m.AppID == "" || m.AppSecret == "" {
		m.AppID, err = utils.GetAppKey()
		if err != nil {
			return err
		}

		m.AppSecret, err = utils.GetSecret(m.AppID)
		if err != nil {
			return err
		}
	}

	if m.SessionNo == "" {
		m.SessionNo = xid.New().String()
	}

	return s.pDomain.SaveSessionKey(ctx, m)
}

// SetStatusSessionKeys set enable/disable by session no
func (s *service) SetStatusSessionKeys(ctx context.Context, status uint32, sessionNos ...string) (int64, error) {
	return s.pDomain.UpdateSessionKey(
		ctx,
		model.SessionKey{IsEnabled: status},
		dao.SessionKey.SessionNo.In(sessionNos...),
	)
}

// SessionsKeysPaginate open api account paginate
func (s *service) SessionsKeysPaginate(ctx context.Context, req SessionKeysPaginateReq) (*SessionKeyRsp, error) {
	var where []gen.Condition

	if req.SessionNo != "" {
		where = append(where, dao.SessionKey.SessionNo.Like(fmt.Sprintf("%%%s%%", req.SessionNo)))
	}

	if req.AppID != "" {
		where = append(where, dao.SessionKey.AppID.Eq(req.AppID))
	}

	if req.IsEnabled != 0 {
		where = append(where, dao.SessionKey.IsEnabled.Eq(req.IsEnabled))
	}

	if req.Name != "" {
		where = append(where, dao.SessionKey.Name.Like(fmt.Sprintf("%%%s%%", req.Name)))
	}

	if req.Type != 0 {
		where = append(where, dao.SessionKey.Type.Eq(req.Type))
	}

	if req.ID != 0 {
		where = append(where, dao.SessionKey.ID.Eq(req.ID))
	}

	if req.StartTime != nil {
		where = append(where, dao.SessionKey.CreatedAt.Gte(*req.StartTime))
	}

	if req.EndTime != nil {
		where = append(where, dao.SessionKey.CreatedAt.Lte(*req.EndTime))
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	list, total, err := s.pDomain.SessionKeyPaginate(ctx, req.Page, req.Limit, where...)
	if err != nil {
		return nil, err
	}

	var result = make([]SessionKeyDetail, len(list))

	for key, val := range list {
		result[key] = SessionKeyDetail{
			SessionNo: val.SessionNo,
			Name:      val.Name,
			Type:      val.Type,
			AppID:     val.AppID,
			IsEnabled: val.IsEnabled,
			Remark:    val.Remark,
			CreatedAt: *val.CreatedAt,
			UpdatedAt: *val.UpdatedAt,
		}
	}

	return &SessionKeyRsp{
		PaginateRsp: common.PaginateRsp{
			List:  result,
			Total: total,

			PaginateReq: common.PaginateReq{
				Limit: req.Limit,
				Page:  req.Page,
			},
		},
	}, nil
}
