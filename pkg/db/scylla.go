package db

import (
	"fmt"

	"github.com/Manjit2003/samespace/pkg/utils"
	"github.com/gocql/gocql"
)

var log = utils.GetChildLogger("database")

type DBConfig struct {
	Hosts string
	Port  int
}

var ScyllaSession *gocql.Session

const (
	keyspace = "samespace_keyspace"
)

func InitDatabase(cfg DBConfig) {
	cluster := gocql.NewCluster(cfg.Hosts)
	cluster.Port = cfg.Port
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4

	var err error
	ScyllaSession, err = cluster.CreateSession()
	if err != nil {
		log.Error("could not connect to the databse: %v", err)
	}

	if err := createKeyspace(); err != nil {
		log.Error("could not create keyspace: %v", err)
	}

	if err := createSchema(); err != nil {
		log.Error("could not create schema: %v", err)
	}

	// now i am replacing the session with new session with the keyspace set
	cluster.Keyspace = keyspace
	ScyllaSession, err = cluster.CreateSession()
	if err != nil {
		log.Error("could not connect to the databse: %v", err)
	}
}

// a hacky migration setup for now, will upgrade later
func createSchema() error {
	queries := []string{
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.todos (
        	id UUID PRIMARY KEY,
        	user_id UUID,
        	title TEXT,
        	description TEXT,
        	status TEXT,
        	created TIMESTAMP,
        	updated TIMESTAMP
    	)`, keyspace),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_todo_user_id ON %s.todos (user_id)`, keyspace),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_todo_status ON %s.todos (status)`, keyspace),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_todo_created ON %s.todos (created)`, keyspace),
	}

	for _, query := range queries {
		if err := ScyllaSession.Query(query).Exec(); err != nil {
			return fmt.Errorf("failed to execute query %q: %w", query, err)
		}
	}

	return nil
}

func createKeyspace() error {
	query := fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {
        'class': 'SimpleStrategy',
        'replication_factor': 1 
    }`, keyspace)

	if err := ScyllaSession.Query(query).Exec(); err != nil {
		return fmt.Errorf("failed to create keyspace: %v", err)
	}
	return nil
}
