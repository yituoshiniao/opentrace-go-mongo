# opentracing go-mongo



## Install


## Usage


Example:

```go
	ctx := context.TODO()
	uri := "mongodb://localhost:2701711"
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
```


## License

[MIT](LICENSE)