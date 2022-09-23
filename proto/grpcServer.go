package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"math/rand"
	"net"
	"net/http"
	"rusprofileGrpcWrap/cmd/config"
	"rusprofileGrpcWrap/logging"
	"strings"
	"time"
)

type GrpcServer struct {
}

func StartGrpc(cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start grpc")

	s := grpc.NewServer()
	srv := &GrpcServer{}

	RegisterRusprofileGrpcWrapServer(s, srv)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port))
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("grpc is listening port %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)
	if err = s.Serve(listener); err != nil {
		logger.Fatal(err)
	}
}

const apiURL = "https://www.rusprofile.ru/ajax.php?query=%s&action=search&cacheKey=%.12f"

type apiResponse struct {
	UL []struct {
		INN  string `json:"inn"`
		OGRN string `json:"ogrn"` // В результате запроса отсутствует КПП, поэтому вставляю ОГРН
		Name string `json:"name"`
		CEO  string `json:"ceo_name"`
	} `json:"ul"`
}

func (receiver *GrpcServer) mustEmbedUnimplementedRusprofileGrpcWrapServer() {
	//TODO implement me
	panic("implement me")
}

func (receiver *GrpcServer) FirmInfoGet(ctx context.Context, req *InnRequest) (*InfoResponse, error) {
	logger := logging.GetLogger()

	if len(req.GetInn()) != 10 && len(req.GetInn()) != 12 && len(req.GetInn()) != 5 {
		return nil, errors.New("incorrect inn")
	}

	rand.Seed(time.Now().Unix())

	resp, err := http.Get(fmt.Sprintf(apiURL, req.GetInn(), rand.Float64()))
	if err != nil {
		logger.Fatal(err)
		return nil, err
	}
	defer resp.Body.Close()

	r := apiResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		logger.Fatal(err)
	}

	if len(r.UL) == 0 {
		return nil, errors.New("inn not found")
	}

	r.UL[0].INN = strings.Trim(r.UL[0].INN, "~!")

	return &InfoResponse{
		Inn:         r.UL[0].INN,
		Ogrn:        r.UL[0].OGRN,
		CompanyName: r.UL[0].Name,
		CeoName:     r.UL[0].CEO,
	}, nil
}
