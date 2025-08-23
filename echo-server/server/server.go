package server

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

type EchoServer struct {
	echo.UnimplementedEchoServer
}

func getMetadata() (header metadata.MD, tailer metadata.MD) {
	header = getMetadataByMap(map[string]string{"time": time.Now().Format("2006-01-02 15:04:05"), "server_header_data": "true"})
	tailer = getMetadataByKV("server_header_data", "true")
	return
}

func (EchoServer) UnaryEcho(ctx context.Context, in *echo.EchoRequest) (*echo.EchoResponse, error) {
	// 响应请求发送元数据
	header, tailer := getMetadata()
	defer grpc.SetTrailer(ctx, tailer)
	grpc.SendHeader(ctx, header)

	// 获取请求中的元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("Server recv UnaryEcho failed: no metadata")
	} else {
		fmt.Println("Server recv UnaryEcho metadata", md)
	}
	fmt.Printf("Server recv: %v\n", in.Message)
	return &echo.EchoResponse{Message: "Server send msg"}, nil
}
func (EchoServer) ServerStreamingEcho(in *echo.EchoRequest, stream echo.Echo_ServerStreamingEchoServer) error {
	// 响应请求发送元数据
	header, tailer := getMetadata()
	// tailer,服务器调用结束后填充的数据
	defer stream.SetTrailer(tailer)
	// header,服务端开始调用时填充的数据
	stream.SendHeader(header)

	// 获取请求中的元数据
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		log.Println("Server recv UnaryEcho failed: no metadata")
	} else {
		fmt.Println("Server recv UnaryEcho metadata", md)
	}

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
	// 响应请求发送元数据
	header, tailer := getMetadata()
	// tailer,服务器调用结束后填充的数据
	defer stream.SetTrailer(tailer)
	// header,服务端开始调用时填充的数据
	stream.SendHeader(header)

	// 获取请求中的元数据
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		log.Println("Server recv UnaryEcho failed: no metadata")
	} else {
		fmt.Println("Server recv UnaryEcho metadata", md)
	}

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
	// 响应请求发送元数据
	header, tailer := getMetadata()
	// tailer,服务器调用结束后填充的数据
	defer stream.SetTrailer(tailer)
	// header,服务端开始调用时填充的数据
	stream.SendHeader(header)

	// 获取请求中的元数据
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		log.Println("Server recv UnaryEcho failed: no metadata")
	} else {
		fmt.Println("Server recv UnaryEcho metadata", md)
	}

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
