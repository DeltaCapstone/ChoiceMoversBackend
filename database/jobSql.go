package DB

import (
	"context"
	"fmt"
)

////////////////////////////////////////////////
//Jobs

const (
	all       = "Select * from jobs where start_time > CURRENT_DATE"
	pending   = "Select * from jobs where start_time > CURRENT_DATE AND finalized = false"
	finalized = "Select * from jobs where start_time > CURRENT_DATE and finalized = true"
)

func (pg *postgres) GetJobsByStatus(ctx context.Context, status string) ([]Job, error) {
	var jobs []Job
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
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var j Job
		if err := rows.Scan(
			&j.ID,
			&j.CustomerId,
			&j.LoadAddr,
			&j.UnloadAddr,
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
			return nil, fmt.Errorf("error reading row: %v", err)
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}
