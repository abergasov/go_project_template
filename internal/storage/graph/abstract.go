package graph

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type GraphDBConnector interface {
	Client() neo4j.Driver
}
