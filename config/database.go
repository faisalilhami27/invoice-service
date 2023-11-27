package config

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"invoice-service/constant"

	"time"
)

type mongoClient struct {
	client *mongo.Client
}

var instance *mongoClient

func NewDatabaseConnection() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(Config.Database.Timeout)*time.Second,
	)
	defer cancel()

	appEnv := Config.AppEnv
	var config string
	if appEnv != constant.LOCAL {
		config = fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority",
			Config.Database.Username,
			Config.Database.Password,
			Config.Database.Host,
		)
	} else {
		config = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin&directConnection=true",
			Config.Database.Username,
			Config.Database.Password,
			Config.Database.Host,
			Config.Database.Port,
			Config.Database.Name,
		)
	}

	opts := options.Client().ApplyURI(config).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetReadPreference(readpref.Primary()).
		SetTimeout(time.Duration(Config.Database.Timeout) * time.Second)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Errorf("an error occurred when connect to mongoDB : %v", err)
		panic(err)
	}

	instance = &mongoClient{
		client: client,
	}

	return instance.client, nil
}
