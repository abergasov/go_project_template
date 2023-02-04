package graph

import (
	"fmt"
	"go_project_template/internal/config"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type GraphDBConnect struct {
	db neo4j.Driver
}

func InitGraphDBConnect(cnf *config.GraphConf) (*GraphDBConnect, error) {
	dsn := fmt.Sprintf("neo4j://%s:%s", cnf.Address, cnf.Port)
	driver, err := neo4j.NewDriver(dsn, neo4j.BasicAuth(cnf.User, cnf.Pass, ""), func(connectConf *neo4j.Config) {
		connectConf.MaxConnectionPoolSize = cnf.MaxConnections
		connectConf.MaxConnectionLifetime = 1 * time.Minute
	})
	if err != nil {
		return nil, fmt.Errorf("error connect to db: %w", err)
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, fmt.Errorf("error ping to db: %w", err)
	}
	return &GraphDBConnect{db: driver}, nil
}

func (c *GraphDBConnect) Client() neo4j.Driver {
	return c.db
}
