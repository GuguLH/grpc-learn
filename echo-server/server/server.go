package server

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"grpc-learn/echo"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type EchoServer struct {
	echo.UnimplementedEchoServer
}

func (EchoServer) UnaryEcho(ctx context.Context, in *echo.EchoRequest) (*echo.EchoResponse, error) {
	fmt.Printf("Server recv: %v\n", in.Message)
	return &echo.EchoResponse{Message: "Server send msg"}, nil
}
func (EchoServer) ServerStreamingEcho(in *echo.EchoRequest, stream echo.Echo_ServerStreamingEchoServer) error {
	fmt.Printf("Server recv: %v\n", in.Message)
	filePath := "echo-server/files/server.png"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		stream.Send(&echo.EchoResponse{
			Message: "Server send file",
			Bytes:   buf[:n],
			Time:    timestamppb.New(time.Now()),
			Length:  int32(n),
		})
	}
	// 服务端流 return nil或者err 流结束
	return nil
}
func (EchoServer) ClientStreamingEcho(stream echo.Echo_ClientStreamingEchoServer) error {
	filePath := "echo-server/files/" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".png"
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return err
		}
		file.Write(req.Bytes[:req.Length])
		fmt.Printf("Server recv: %v\n", req.Message)
	}

	err = stream.SendAndClose(&echo.EchoResponse{
		Message: "Server send msg",
	})
	return err
}
func (EchoServer) DoubleStreamingEcho(stream echo.Echo_DoubleStreamingEchoServer) error {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		filePath := "echo-server/files/" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".png"
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println(err)
				return
			}
			file.Write(req.Bytes[:req.Length])
			fmt.Printf("Server recv: %v\n", req.Message)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		filePath := "echo-server/files/server.png"
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		buf := make([]byte, 1024)
		for {
			n, err := file.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				return
			}
			stream.Send(&echo.EchoResponse{
				Message: "Server send file",
				Bytes:   buf[:n],
				Time:    timestamppb.New(time.Now()),
				Length:  int32(n),
			})
		}
	}()

	wg.Wait()

	return nil
}
