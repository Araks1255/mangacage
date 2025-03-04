package db

import (
	"context"
	//"log"

	//"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func MongoInit(mongoUrl string) (*mongo.Client, error) {
	// monitor := &event.CommandMonitor{
	// 	Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
	// 		log.Printf("[MONGO REQUEST] %s %s %s\n", evt.DatabaseName, evt.CommandName, evt.Command)
	// 	},
	// 	Succeeded: func(ctx context.Context, evt *event.CommandSucceededEvent) {
	// 		log.Printf("[MONGO RESPONSE] %s %s (duration: %s)\n", evt.CommandName, evt.Reply, evt.Duration)
	// 	},
	// 	Failed: func(ctx context.Context, evt *event.CommandFailedEvent) {
	// 		log.Printf("[MONGO ERROR] %s (duration: %s) %s\n", evt.CommandName, evt.Duration, evt.Failure)
	// 	},
	// }

	// clientOpts := options.Client().ApplyURI(mongoUrl).SetMonitor(monitor)

	client, err := mongo.Connect(context.TODO())
	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}
