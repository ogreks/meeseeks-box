package agreements

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
	"time"
)

var _ AGreementDomain = (*AGreements)(nil)

type AGreementDomain interface {
	i()

	// SaveAGreement create agreement or publish agreement
	SaveAGreement(ctx context.Context, agreements model.Agreement) error
	// CreateAGreementVersion create agreement version
	CreateAGreementVersion(ctx context.Context, agv model.AgreementVersion) error
	// AgreementPublish publish agreement content
	AgreementPublish(ctx context.Context, agrenos ...string) error
	// Paginate get agreement page list
	Paginate(ctx context.Context, page, limit int, where ...gen.Condition) ([]*model.Agreement, int64, error)
	// PaginateVersions get agreement version page list
	PaginateVersions(ctx context.Context, page, limit int, where ...gen.Condition) ([]*model.AgreementVersion, int64, error)
	// Agreement get agreement detail
	Agreement(ctx context.Context, where ...gen.Condition) (*model.Agreement, error)
	// AgreementVersion get agreement version detail
	AgreementVersion(ctx context.Context, where ...gen.Condition) (*model.AgreementVersion, error)
	// AgreementVersionNew get agreement new version detail
	AgreementVersionNew(ctx context.Context, where ...gen.Condition) (*model.AgreementVersion, error)
	// AgreementDelete batch delete agreement
	AgreementDelete(ctx context.Context, where ...gen.Condition) (int64, error)
}

type AGreements struct {
	dao *dao.Query
	log *zap.Logger
}

func NewAGreements(repo orm.Repo, log *zap.Logger) *AGreements {
	return &AGreements{
		log: log,
		dao: dao.Use(repo.GetDB()),
	}
}

func (ag *AGreements) i() {}

func (ag *AGreements) CreateAGreementVersion(ctx context.Context, agv model.AgreementVersion) error {
	return ag.dao.AgreementVersion.WithContext(ctx).Create(&agv)
}

// SaveAGreement create agreement or publish agreement
func (ag *AGreements) SaveAGreement(ctx context.Context, agreements model.Agreement) error {
	return ag.dao.Transaction(func(tx *dao.Query) error {
		if agreements.Status == 1 {
			if agreements.PublishAt == nil {
				t := time.Now()
				agreements.PublishAt = &t
			}

			err := tx.WithContext(ctx).AgreementVersion.Create(&model.AgreementVersion{
				AgreementNo: agreements.AgreementNo,
				Type:        agreements.Type,
				Title:       agreements.Title,
				Content:     agreements.Content,
				Version:     agreements.Version,
			})
			if err != nil {
				return err
			}
		}
		return tx.WithContext(ctx).Agreement.Save(&agreements)
	})
}

// AgreementPublish publish agreement content
func (ag *AGreements) AgreementPublish(ctx context.Context, agrenos ...string) error {
	return ag.dao.Transaction(func(tx *dao.Query) error {
		data, err := tx.Agreement.WithContext(ctx).Where(dao.Agreement.AgreementNo.In(agrenos...)).Find()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.NotFound
			}
			return err
		}

		var argv = make([]*model.AgreementVersion, len(data))
		for k, agreement := range data {
			t := time.Now()
			agreement.PublishAt = &t
			agreement.Status = 1

			argv[k] = &model.AgreementVersion{
				AgreementNo: agreement.AgreementNo,
				Type:        agreement.Type,
				Title:       agreement.Title,
				Content:     agreement.Content,
				Version:     agreement.Version,
			}

			err := tx.Agreement.WithContext(ctx).Save(agreement)
			if err != nil {
				return err
			}
		}

		return tx.AgreementVersion.WithContext(ctx).Create(argv...)
	})
}

// Paginate get agreement page list
func (ag *AGreements) Paginate(ctx context.Context, page, limit int, where ...gen.Condition) ([]*model.Agreement, int64, error) {
	return ag.dao.Agreement.
		WithContext(ctx).
		Where(where...).
		FindByPage(
			(page-1)*limit,
			limit,
		)
}

// PaginateVersions get agreement version page list
func (ag *AGreements) PaginateVersions(ctx context.Context, page, limit int, where ...gen.Condition) ([]*model.AgreementVersion, int64, error) {
	return ag.dao.AgreementVersion.
		WithContext(ctx).
		Where(where...).
		FindByPage(
			(page-1)*limit,
			limit,
		)
}

// Agreement get agreement detail
func (ag *AGreements) Agreement(ctx context.Context, where ...gen.Condition) (*model.Agreement, error) {
	data, err := ag.dao.WithContext(ctx).Agreement.Where(where...).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = domain.NotFound
	}

	return data, err
}

// AgreementVersion get agreement version detail
func (ag *AGreements) AgreementVersion(ctx context.Context, where ...gen.Condition) (*model.AgreementVersion, error) {
	v, e := ag.dao.WithContext(ctx).AgreementVersion.Where(where...).First()
	if errors.Is(e, gorm.ErrRecordNotFound) {
		e = domain.NotFound
	}

	return v, e
}

// AgreementVersionNew get agreement new version detail
func (ag *AGreements) AgreementVersionNew(ctx context.Context, where ...gen.Condition) (*model.AgreementVersion, error) {
	v, e := ag.dao.WithContext(ctx).AgreementVersion.
		Order(dao.Agreement.CreatedAt.Desc(), dao.Agreement.Version.Desc()).
		Where(where...).
		First()
	if errors.Is(e, gorm.ErrRecordNotFound) {
		e = domain.NotFound
	}

	return v, e
}

// AgreementDelete batch delete agreement
func (ag *AGreements) AgreementDelete(ctx context.Context, where ...gen.Condition) (int64, error) {
	result, err := ag.dao.Agreement.WithContext(ctx).Where(where...).Delete()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected, result.Error
}
