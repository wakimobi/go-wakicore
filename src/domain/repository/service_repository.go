package repository

import (
	"database/sql"

	"github.com/idprm/go-pass-tsel/src/domain/entity"
)

const (
	queryCountServiceByCategory = "SELECT COUNT(*) as count FROM services WHERE category = $1"
	queryCountServiceByCode     = "SELECT COUNT(*) as count FROM services WHERE code = $1"
	querySelectIdService        = "SELECT id, category, code, name, price, program_id, sid, renewal_day, trial_day, url_telco, url_portal, url_notif_sub, url_notif_unsub, url_notif_renewal, url_postback FROM services WHERE id = $1 LIMIT 1"
	querySelectCodeService      = "SELECT id, category, code, name, price, program_id, sid, renewal_day, trial_day, url_telco, url_portal, url_notif_sub, url_notif_unsub, url_notif_renewal, url_postback FROM services WHERE code = $1 LIMIT 1"
)

type ServiceRepository struct {
	db *sql.DB
}

type IServiceRepository interface {
	CountByCategory(string) (int, error)
	CountByCode(string) (int, error)
	GetById(int) (*entity.Service, error)
	GetByCode(string) (*entity.Service, error)
}

func NewServiceRepository(db *sql.DB) *ServiceRepository {
	return &ServiceRepository{
		db: db,
	}
}

func (r *ServiceRepository) CountByCategory(category string) (int, error) {
	var count int
	err := r.db.QueryRow(queryCountServiceByCategory, category).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (r *ServiceRepository) CountByCode(code string) (int, error) {
	var count int
	err := r.db.QueryRow(queryCountServiceByCode, code).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (r *ServiceRepository) GetById(id int) (*entity.Service, error) {
	var s entity.Service
	err := r.db.QueryRow(querySelectIdService, id).Scan(&s.ID, &s.Category, &s.Code, &s.Name, &s.Price, &s.ProgramId, &s.Sid, &s.RenewalDay, &s.TrialDay, &s.UrlTelco, &s.UrlPortal, &s.UrlNotifSub, &s.UrlNotifUnsub, &s.UrlNotifRenewal, &s.UrlPostback)
	if err != nil {
		return &s, err
	}
	return &s, nil
}

func (r *ServiceRepository) GetByCode(code string) (*entity.Service, error) {
	var s entity.Service
	err := r.db.QueryRow(querySelectCodeService, code).Scan(&s.ID, &s.Category, &s.Code, &s.Name, &s.Price, &s.ProgramId, &s.Sid, &s.RenewalDay, &s.TrialDay, &s.UrlTelco, &s.UrlPortal, &s.UrlNotifSub, &s.UrlNotifUnsub, &s.UrlNotifRenewal, &s.UrlPostback)
	if err != nil {
		return &s, err
	}
	return &s, nil
}
