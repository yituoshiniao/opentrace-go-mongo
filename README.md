# opentracing-go-mongo   author: smgano@163.com


#git clone https://github.com/yituoshiniao/opentracing-go-mongo.git

## Install
go get -u github.com/yituoshiniao/opentracing-go-mongo


## Usage


Example:

```go
    #注意这里是演示使用，在正式项目中请使用生成trace链路追踪的ctx
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
```


## License

[MIT](LICENSE)