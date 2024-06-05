package logger

import (
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	cc, err := grpc.NewClient("localhost:8000", opts...)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cc.Close()
	lg := NewLogger(cc)
	if err := lg.setup(); err != nil {
		log.Fatal(err)
	}

	lg.Error(errors.New("error"))
	time.Sleep(2 * time.Second)
}
