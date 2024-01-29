package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"logger-service/logs/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Model data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, request *logs.LogRequest) (*logs.LogResponse, error) {
	input := request.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Model.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "Failed to insert log into Mongo",
		}
		return res, err
	}

	return &logs.LogResponse{Result: "Logged via gRPC"}, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPort))
	if err != nil {
		log.Fatalln(err)
	}

	server := grpc.NewServer()

	logs.RegisterLogServiceServer(server, &LogServer{Model: app.Models})

	log.Println("gRPC server started on PORT:", gRPCPort)

	if err := server.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
