package core

import (
	"clean_code_demo/repository"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/url"
	"os"
	"time"
)

type System struct {
	repo repository.Repo
}

type DocDbConf struct {
	ConnectionString  string
	Protocol          string
	Host              string
	DefaultDb         string
	Username          string
	Password          string
	ReplicaSet        string
	ReadPreference    string
	ConnectTimeoutMs  int
	SocketTimeoutMs   int
	ReconnectInterval int
	PoolSize          int
	BufferMaxEntries  int
	KeepAlive         bool
	BufferCommands    bool
	AutoReconnect     bool
}
type DocDB struct {
	client *mongo.Client
	config *DocDbConf
}

const (
	authenticationStringTemplate = "%s:%s@"
	connectionStringTemplate     = "%s://%s%s/%s?"
)

func ConnectWithDocDB(config *DocDbConf) (*DocDB, error) {

	docDb := &DocDB{
		config: config,
	}

	fmt.Printf("[DocumentDB] Try to connect document db\n")

	var clientOptions *options.ClientOptions
	if config.ConnectionString != "" {
		fmt.Printf("[DocumentDB] ConnectionStr: %s\n", config.ConnectionString)
		clientOptions = options.Client().ApplyURI(config.ConnectionString)
	} else {
		authenticationURI := ""
		if config.Username != "" {
			authenticationURI = fmt.Sprintf(
				authenticationStringTemplate,
				config.Username,
				config.Password,
			)
		}

		connectionURI := fmt.Sprintf(
			connectionStringTemplate,
			config.Protocol,
			authenticationURI,
			config.Host,
			config.DefaultDb,
		)

		connectUri, _ := url.Parse(connectionURI)
		connectQuery, _ := url.ParseQuery(connectUri.RawQuery)

		if config.ReplicaSet != "" {
			connectQuery.Add("replicaSet", config.ReplicaSet)
		}

		if config.ReadPreference != "" {
			connectQuery.Add("readpreference", config.ReadPreference)
		}

		connectUri.RawQuery = connectQuery.Encode()

		fmt.Printf("[DocumentDB] ConnectionURI: %s", connectionURI)
		clientOptions = options.Client().ApplyURI(connectUri.String())
	}

	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Failed to Create New Client, %s", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.ConnectTimeoutMs)*time.Millisecond)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		fmt.Printf("Failed to connect to cluster: %v", err)
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, readpref.SecondaryPreferred())
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Failed to ping cluster: %s", err))
	}

	fmt.Printf("[DocumentDB] Connected to DocumentDB!")

	docDb.client = client

	return docDb, nil
}

// Hardcode for demo
func connectMongo() (*mongo.Database, error) {
	fmt.Printf("try to connect to mongo db")
	dbConf := DocDbConf{
		ConnectionString:  "",
		Protocol:          "mongodb",
		Host:              "localhost",
		DefaultDb:         "demo",
		Username:          "developer",
		Password:          "developer",
		ReplicaSet:        "",
		ReadPreference:    "secondaryPreferred",
		ConnectTimeoutMs:  10000,
		SocketTimeoutMs:   0,
		ReconnectInterval: 20000,
		PoolSize:          15,
		BufferMaxEntries:  0,
		KeepAlive:         true,
		BufferCommands:    false,
		AutoReconnect:     true,
	}
	client, err := ConnectWithDocDB(&dbConf)
	if err != nil {
		fmt.Printf("[DocumentDB] Try To Connect Document Database Failed -> (%s)", err)
		os.Exit(102)
	}
	return client.client.Database(dbConf.DefaultDb), nil
}

func (s *System) GetRepository() *repository.Repo {
	return &s.repo
}

func (s *System) Initialize() {
	db, err := connectMongo()
	if err != nil {
		fmt.Printf("mongodb connection fail")
	}
	s.repo = repository.Repo{}
	s.repo.DB = db
}

func New() *System {
	return &System{}
}
