package main

import (
	"context"
	"fmt"
	pb "github.com/Azvf/learn-go/grpc-example/route"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func runFirst(client pb.RouteGuideClient) {
	feature, err := client.GetFeature(context.Background(), &pb.Point{Longitude: 121437403, Latitude: 310235000})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(feature)

}

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		log.Fatalln("client cannot dial grpc server")
	}
	defer conn.Close()

	client := pb.NewRouteGuideClient(conn)

	runFirst(client)

}
