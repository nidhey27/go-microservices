package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	// mongoURL = "mongodb://admin:password@mongo:27017/logs"
	gRPCPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// Connect to DB
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Setup Config
	app := &Config{
		Models: data.New(client),
	}
	log.Printf("Starting Logger Service on PORT %s\n", webPort)

	err = rpc.Register(new(RPCServer))
	go app.rpcListen()

	go app.gRPCListen()

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC Sever on PORT:", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		log.Println(err)
		return err
	}
	defer listen.Close()
	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create connections options
	clientOption := options.Client().ApplyURI(mongoURL)
	clientOption.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return c, nil
}
