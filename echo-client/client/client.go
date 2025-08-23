package client

import (
	"context"
	"fmt"
	"grpc-learn/echo"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func getContext(ctx context.Context) context.Context {
	md := getMetadataByMap(map[string]string{"time": time.Now().Format("2006-01-02 15:04:05"), "header_data": "true"})
	// 将数据写入到ctx
	ctx = getOutgoingContext(ctx, md)
	// 将数据附加到ctx
	ctx = appendToOutgoingContext(ctx, "k1", "v1", "k2", "v2")
	//md1, _ := metadata.FromOutgoingContext(ctx)
	//fmt.Println(md1)
	return ctx
}

func CallUnary(client echo.EchoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ctx = getContext(ctx)

	in := &echo.EchoRequest{
		Message: "Client send msg",
		Time:    timestamppb.New(time.Now()),
	}

	var header, trailer metadata.MD

	res, err := client.UnaryEcho(ctx, in, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Client recv: %v\n", res.Message)

	fmt.Printf("Header: %v\n", header)
	fmt.Printf("Trailer: %v\n", trailer)
}

func CallServerStream(client echo.EchoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ctx = getContext(ctx)

	in := &echo.EchoRequest{
		Message: "Client send msg",
		Time:    timestamppb.New(time.Now()),
	}
	stream, err := client.ServerStreamingEcho(ctx, in)
	if err != nil {
		log.Fatal(err)
	}

	header, _ := stream.Header()
	fmt.Printf("Header: %v\n", header)

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
	trailer := stream.Trailer()
	fmt.Printf("Trailer: %v\n", trailer)
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

	ctx = getContext(ctx)

	stream, err := client.ClientStreamingEcho(ctx)
	if err != nil {
		log.Fatal(err)
	}

	header, _ := stream.Header()
	fmt.Printf("Header: %v\n", header)

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

	trailer := stream.Trailer()
	fmt.Printf("Trailer: %v\n", trailer)
}

func CallDoubleStream(client echo.EchoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ctx = getContext(ctx)

	stream, err := client.DoubleStreamingEcho(ctx)
	if err != nil {
		log.Fatal(err)
	}

	header, _ := stream.Header()
	fmt.Printf("Header: %v\n", header)

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
	trailer := stream.Trailer()
	fmt.Printf("Trailer: %v\n", trailer)
}
