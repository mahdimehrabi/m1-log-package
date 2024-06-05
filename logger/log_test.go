package logger

import (
	"errors"
	"fmt"
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
}

func BenchmarkLogger(b *testing.B) {
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
	for i := 0; i < 200000; i++ {
		lg.Error(errors.New(fmt.Sprintf("error %d", i)))
	}
	for {
		if len(lg.queue) > 0 {
			time.Sleep(1 * time.Millisecond)
		} else {
			break
		}
	}
	fmt.Printf("took %s\n", b.Elapsed().String())

	//waiting til log service process messages
	time.Sleep(30 * time.Second)
}
