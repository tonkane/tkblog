package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/tkane/tkblog/internal/pkg/log"
	pb "github.com/tkane/tkblog/pkg/proto/tkblog/v1"
)

var (
	addr = flag.String("addr", "localhost:9090", "the address to connect to.")
	limit = flag.Int64("limit", 10, "limit to list users")
)

func main() {
	flag.Parse()

	// 建立连接
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalw("did not connect", "err", err)
	}

	defer conn.Close()
	c := pb.NewTkBlogClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 请求list接口
	r, err := c.ListUser(ctx, &pb.ListUserRequest{Offset: 0, Limit: *limit})
	if err != nil {
		log.Fatalw("could not greet: %v", err)
	}

	// 打印结果
	fmt.Println("TotalCount:", r.TotalCount)
	for _, u := range r.Users {
		d, _ := json.Marshal(u)
		fmt.Println("result:", string(d))
	}
}