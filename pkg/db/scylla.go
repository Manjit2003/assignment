package db

import (
	"fmt"

	"github.com/Manjit2003/samespace/pkg/config"
	"github.com/Manjit2003/samespace/pkg/utils"
	"github.com/gocql/gocql"
)

var log = utils.GetChildLogger("database")
var ScyllaSession *gocql.Session

func InitDatabase(cfg *config.APIConfig) {

	cluster := gocql.NewCluster(cfg.Database.Host)
	cluster.Port = cfg.Database.Port
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4

	var err error
	ScyllaSession, err = cluster.CreateSession()
	if err != nil {
		log.Error("could not connect to the databse: %v", err)
	}

	if err := createKeyspace(cfg); err != nil {
		log.Error("could not create keyspace: %v", err)
	}

	if err := createSchema(cfg); err != nil {
		log.Error("could not create schema: %v", err)
	}

	// now i am replacing the session with new session with the keyspace set
	cluster.Keyspace = cfg.Database.Keyspace
	ScyllaSession, err = cluster.CreateSession()
	if err != nil {
		log.Error("could not connect to the databse: %v", err)
	}

	log.Info("database initialized")
}

// a hacky migration setup for now, will upgrade later
func createSchema(cfg *config.APIConfig) error {
	queries := []string{
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.todos (
			user_id UUID,
			id UUID,
			title TEXT,
			description TEXT,
			status TEXT,
			created TIMESTAMP,
			updated TIMESTAMP,
			PRIMARY KEY (user_id, id)
		);`, cfg.Database.Keyspace),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_todo_status ON %s.todos (status)`, cfg.Database.Keyspace),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_todo_created ON %s.todos (created)`, cfg.Database.Keyspace),
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.users (
        	id UUID,
        	username TEXT,
        	hashed_password TEXT,
        	created TIMESTAMP,
        	updated TIMESTAMP,
			PRIMARY KEY (username, id)
    	)`, cfg.Database.Keyspace),
	}

	for _, query := range queries {
		if err := ScyllaSession.Query(query).Exec(); err != nil {
			return fmt.Errorf("failed to execute query %q: %w", query, err)
		}
	}

	log.Debug(fmt.Sprintf("migration complete, executed %d queries", len(queries)))

	return nil
}

func createKeyspace(cfg *config.APIConfig) error {
	query := fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {
        'class': 'SimpleStrategy',
        'replication_factor': 1 
    }`, cfg.Database.Keyspace)

	if err := ScyllaSession.Query(query).Exec(); err != nil {
		return fmt.Errorf("failed to create keyspace: %v", err)
	}

	log.Debug("keyspace created successfully")

	return nil
}
