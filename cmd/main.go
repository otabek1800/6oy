package main

import (
	"google_docs_user/api"
	"google_docs_user/api/handler"
	"google_docs_user/config"
	logs "google_docs_user/pkg/logger"
	"google_docs_user/service"
	"google_docs_user/storage/postgres"
	"log"
	"net"

	pb "google_docs_user/genproto/user"

	"google.golang.org/grpc"
)

func main() {
	logger := logs.InitLogger(config.Load())
	listener, err := net.Listen("tcp", config.Load().USER_SERVICE)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	db, err := postgres.ConnectionDb()
	if err != nil {
		log.Fatal(err)
	}

	servicee := service.NewUserService(db, logger)

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, servicee)

	go func() {
		err := server.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}
	}()

	hand := handler.NewHandler(logger)
	router := api.NewRouter(hand)
	err = router.Run(config.Load().USER_ROUTER)
	if err != nil {
		log.Fatal(err)
	}
}
