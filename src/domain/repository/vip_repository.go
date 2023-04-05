package repository

import (
	"database/sql"
)

const (
	queryCountVIP = "SELECT COUNT(*) as count FROM vips WHERE msisdn = $1"
)

type VIPRepository struct {
	db *sql.DB
}

type IVIPRepository interface {
	Count(string) (int, error)
}

func NewVIPRepository(db *sql.DB) *VIPRepository {
	return &VIPRepository{
		db: db,
	}
}

func (r *VIPRepository) Count(msisdn string) (int, error) {
	var count int
	err := r.db.QueryRow(queryCountVIP, msisdn).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}
