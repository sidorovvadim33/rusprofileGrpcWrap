FROM golang:1.19

WORKDIR ./cmd

RUN go version
ENV GOPATH=/

COPY ./ ./

# build go app
RUN go mod download
RUN go build -o rusprofile-grpc-wrap rusprofileGrpcWrap/cmd/


EXPOSE 8090

CMD ["./rusprofileGrpcWrap"]