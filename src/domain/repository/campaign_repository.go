package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/idprm/go-pass-tsel/src/domain/entity"
)

const (
	queryCountCampaign  = "SELECT COUNT(*) as count FROM campaigns WHERE service_id = $1 AND adnet = $2 AND total_mo >= limit_mo"
	queryUpdateCampaign = "UPDATE campaigns SET total_mo = total_mo + 1 WHERE service_id = $1 AND adnet = $2"
)

type CampaignRepository struct {
	db *sql.DB
}

type ICampaignRepository interface {
	Count(int, string) (int, error)
	Update(*entity.Campaign) error
}

func NewCampaignRepository(db *sql.DB) *CampaignRepository {
	return &CampaignRepository{
		db: db,
	}
}

func (r *CampaignRepository) Count(serviceId int, adnet string) (int, error) {
	var count int
	err := r.db.QueryRow(queryCountCampaign, serviceId, adnet).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (r *CampaignRepository) Update(c *entity.Campaign) error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.db.PrepareContext(ctx, queryUpdateCampaign)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, c.ServiceID, c.Adnet)
	if err != nil {
		log.Printf("Error %s when update row into campaigns table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d campaigns updated ", rows)

	return nil
}
