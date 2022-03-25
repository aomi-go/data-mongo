package datasource

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"os"
	"sync"
)

const (
	URI                     = "MONGODB_URI"
	HOST                    = "MONGODB_HOST"
	PORT                    = "MONGODB_PORT"
	AUTHENTICATION_DATABASE = "MONGODB_AUTHENTICATION_DATABASE"
	DATABASE                = "MONGODB_DATABASE"
	USERNAME                = "MONGODB_USERNAME"
	PASSWORD                = "MONGODB_PASSWORD"
)

type DataSource struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var (
	instance *DataSource
	once     sync.Once
)

func Instance() *DataSource {
	if instance != nil {
		return instance
	}
	once.Do(func() {
		var (
			clientOptions *options.ClientOptions
			database      string
		)
		uri := os.Getenv(URI)
		if uri == "" {
			loadBalanced := false
			clientOptions = &options.ClientOptions{
				Hosts: []string{os.Getenv(HOST) + ":" + os.Getenv(PORT)},
				Auth: &options.Credential{
					AuthSource:  os.Getenv(AUTHENTICATION_DATABASE),
					Username:    os.Getenv(USERNAME),
					Password:    os.Getenv(PASSWORD),
					PasswordSet: true,
				},
				LoadBalanced: &loadBalanced,
			}
			database = os.Getenv(DATABASE)
		} else {
			cs, err := connstring.ParseAndValidate(uri)
			if err != nil {
				panic(err.(any))
			}
			database = cs.Database
			clientOptions = options.Client().ApplyURI(uri)
		}
		client, err := mongo.Connect(context.TODO(), clientOptions)

		if err != nil {
			panic(err.(any))
		}

		instance = &DataSource{Client: client, Database: client.Database(database)}
	})
	return instance
}

func Client() *mongo.Client {
	return Instance().Client
}

func Database() *mongo.Database {
	return Instance().Database
}
