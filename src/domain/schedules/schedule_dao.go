package schedules

import (
	"github.com/wakimobi/go-wakicore/src/datasource/pgsql/db"
	"github.com/wakimobi/go-wakicore/src/utils/errors"
	"github.com/wakimobi/go-wakicore/src/utils/pgsql_utils"
)

const (
	queryGetSchedule    = "SELECT id, name, publish_at, unlocked_at, is_unlocked FROM schedules WHERE name = $1 LIMIT 1"
	queryUpdateSchedule = "UPDATE schedules SET is_unlocked = $1 WHERE name = $2"
)

func (s *Schedule) Get() *errors.RestErr {
	stmt, err := db.Client.Prepare(queryGetSchedule)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	result := stmt.QueryRow(s.Name)
	if getErr := result.Scan(&s.ID, &s.Name, &s.PublishAt, &s.UnlockedAt, &s.IsUnlocked); getErr != nil {
		return pgsql_utils.ParseError(getErr)
	}
	return nil
}

func (s *Schedule) Update() *errors.RestErr {
	stmt, err := db.Client.Prepare(queryUpdateSchedule)
	if err != nil {
		return pgsql_utils.ParseError(err)
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(s.IsUnlocked, s.Name)
	if updateErr != nil {
		return pgsql_utils.ParseError(updateErr)
	}

	return nil
}
