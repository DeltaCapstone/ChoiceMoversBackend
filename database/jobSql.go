package DB

import (
	"context"
)

////////////////////////////////////////////////
//Jobs

// shell
func (pg *postgres) GetJobsByStatus(ctx context.Context, status string) ([]Job, error) {
	var jobs []Job
	return jobs, nil
}
