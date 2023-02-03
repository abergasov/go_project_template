package sampler

import "go_project_template/internal/storage/database"

type Repo struct {
	db database.DBConnector
}

func InitRepo(db database.DBConnector) *Repo {
	return &Repo{db: db}
}
