package DB

import (
	"context"

	"github.com/DeltaCapstone/ChoiceMoversBackend/models"
)

////////////////////////////////////////////////
//Jobs

const (
	all       = "SELECT * FROM jobs WHERE start_time > CURRENT_DATE"
	pending   = "SELECT * FROM jobs WHERE start_time > CURRENT_DATE AND finalized = false"
	finalized = "SELECT * FROM jobs WHERE start_time > CURRENT_DATE AND finalized = true"
)

// TODO: Figure out error handling for address errors
func (pg *postgres) GetJobsByStatusAndRange(ctx context.Context, status string) ([]models.JobResponse, error) {
	var jobs []models.JobResponse
	var query string
	switch status {
	case "all":
		query = all
	case "pending":
		query = pending
	case "finalized":
		query = finalized
	}

	rows, err := pg.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		LoadAddrID   int
		UnloadAddrID int
		CustomerID   int
	)
	for rows.Next() {
		var j models.JobResponse
		if err := rows.Scan(
			&j.ID,
			&CustomerID,
			&LoadAddrID,
			&UnloadAddrID,
			&j.StartTime,
			&j.HoursLabor,
			&j.Finalized,
			&j.Rooms,
			&j.Pack,
			&j.Unpack,
			&j.Load,
			&j.Unload,
			&j.Clean,
			&j.Milage,
			&j.Cost,
		); err != nil {
			return nil, err
		}
		//need to figure out error handling here
		j.LoadAddr, _ = getAddr(ctx, LoadAddrID)
		j.UnloadAddr, _ = getAddr(ctx, UnloadAddrID)
		j.Customer, _ = pg.GetCustomerById(ctx, CustomerID)
		jobs = append(jobs, j)
	}
	return jobs, nil
}

const addrQuery = "SELECT * FROM addresses WHERE address_id = $1"

func getAddr(ctx context.Context, addrID int) (models.Address, error) {
	var a models.Address
	row := PgInstance.db.QueryRow(ctx, addrQuery, addrID)
	err := row.Scan(
		&a.AddressID,
		&a.Street,
		&a.City,
		&a.State,
		&a.Zip,
		&a.ResType,
		&a.Flights,
		&a.AptNum,
	)
	if err != nil {
		return a, err
	}
	return a, nil
}
