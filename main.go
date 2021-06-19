package main

import (
	"context"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"opentracing-go-mongo/trace"
)

func main() {
	ctx := context.TODO()
	uri := "mongodb://localhost:27017"
	listener := &trace.Trace{}
	trace := &event.CommandMonitor{
		Started:   listener.HandleStartedEvent,
		Succeeded: listener.HandleSucceededEvent,
		Failed:    listener.HandleFailedEvent,
	}
	opts := options.Client().ApplyURI(uri).SetMonitor(trace)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)
	
}
