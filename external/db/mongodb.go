package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host, port, username, password, database, authDB string) (db *mongo.Database, err error) {

	var mongoDBURL string
	var isAuth bool
	if username == "" && password == "" {
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s", host, port)
	} else {
		isAuth = true
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", host, port, username, password)
	}

	clientOptions := options.Client().ApplyURI(mongoDBURL)
	if isAuth {
		if authDB == "" {
			authDB = database
		}
		clientOptions.SetAuth(options.Credential{
			AuthSource: authDB,
			Username:   username,
			Password:   password,
		})
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("fucked up to connect to Mongo because of error is appeared: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("fucked up to ping to Mongo because of error is appeared: %v", err)
	}

	db = client.Database(database)

	db.Client().Database(database).Client().Database(database).Client().Database(database).Client()

	return client.Database(database), nil
}
