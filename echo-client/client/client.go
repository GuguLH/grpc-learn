package client

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

func CallUnary(client echo.EchoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	in := &echo.EchoRequest{
		Message: "Client send msg",
		Time:    timestamppb.New(time.Now()),
	}
	res, err := client.UnaryEcho(ctx, in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Client recv: %v\n", res.Message)
}

func CallServerStream(client echo.EchoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	in := &echo.EchoRequest{
		Message: "Client send msg",
		Time:    timestamppb.New(time.Now()),
	}
	stream, err := client.ServerStreamingEcho(ctx, in)
	if err != nil {
		log.Fatal(err)
	}
	fileName := "echo-client/files/" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".png"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		file.Write(res.Bytes[:res.Length])
		fmt.Printf("Client recv: %v\n", res.Message)
	}
	stream.CloseSend()
}

func CallClientStream(client echo.EchoClient) {
	filePath := "echo-client/files/client.png"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := client.ClientStreamingEcho(ctx)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		stream.Send(&echo.EchoRequest{
			Message: "Client send file",
			Bytes:   buf[:n],
			Time:    timestamppb.New(time.Now()),
			Length:  int32(n),
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Client recv: %v\n", res.Message)
}

func CallDoubleStream(client echo.EchoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := client.DoubleStreamingEcho(ctx)
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		filePath := "echo-client/files/client.png"
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
				log.Fatal(err)
			}
			stream.Send(&echo.EchoRequest{
				Message: "Client send file",
				Bytes:   buf[:n],
				Time:    timestamppb.New(time.Now()),
				Length:  int32(n),
			})
		}
		stream.CloseSend()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fileName := "echo-client/files/" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".png"
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println(err)
				break
			}
			file.Write(res.Bytes[:res.Length])
			fmt.Printf("Client recv: %v\n", res.Message)
		}
	}()

	wg.Wait()
}
