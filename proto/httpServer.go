package rpc

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"rusprofileGrpcWrap/cmd/config"
	"rusprofileGrpcWrap/logging"
)

var content embed.FS

func StartHttp(cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start http")

	grpcConn, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port), grpc.WithInsecure())
	if err != nil {
		logger.Fatal(err)
	}
	defer grpcConn.Close()

	grpcMux := runtime.NewServeMux()

	if err = RegisterRusprofileGrpcWrapHandler(context.Background(), grpcMux, grpcConn); err != nil {
		logger.Fatal(err)
	}
	mux := http.NewServeMux()

	mux.Handle("/inn/", grpcMux)
	mux.HandleFunc("/openapi.json", swaggerHandler)

	fs := http.FileServer(http.Dir("cmd/dist"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	logger.Infof("http server is listening port %s:%s", cfg.Listen.BindIp, cfg.Listen.ProxyPort)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Listen.ProxyPort), mux))
}

func swaggerHandler(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	b, err := ioutil.ReadFile("proto/server.swagger.json")
	if err != nil {
		logger.Info(err)

		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  err.Error(),
			"result": nil,
		})

		if err != nil {
			logger.Fatal(err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
