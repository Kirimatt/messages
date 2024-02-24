package db

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func init() {
	if os.Getenv("CASSANDRA_CONNECTION_URL") == "" {
		os.Setenv("CASSANDRA_CONNECTION_URL", "127.0.0.1")
	}
	if os.Getenv("CASSANDRA_CONNECTION_KEYSPACE") == "" {
		os.Setenv("CASSANDRA_CONNECTION_KEYSPACE", "zangerchat")
	}
	if os.Getenv("CASSANDRA_CONNECTION_PORT") == "" {
		os.Setenv("CASSANDRA_CONNECTION_PORT", "9042")
	}

	var err error
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_CONNECTION_URL"))
	port, err := strconv.Atoi(os.Getenv("CASSANDRA_CONNECTION_PORT"))
	if err != nil {
		panic(err)
	}
	cluster.Port = port
	cluster.Keyspace = os.Getenv("CASSANDRA_CONNECTION_KEYSPACE")
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("cassandra well initialized")
}
