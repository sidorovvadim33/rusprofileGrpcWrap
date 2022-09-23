package main

import (
	"rusprofileGrpcWrap/cmd/config"
	"rusprofileGrpcWrap/logging"
	"rusprofileGrpcWrap/proto"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("start application")

	cfg := config.GetConfig()

	start(cfg)
}

func start(cfg *config.Config) {
	go rpc.StartGrpc(cfg)
	rpc.StartHttp(cfg)
}
