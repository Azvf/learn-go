package main

import (
	"bufio"
	"context"
	"fmt"
	pb "github.com/Azvf/learn-go/grpc-example/route"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"time"
)

func runFirst(client pb.RouteGuideClient) {
	feature, err := client.GetFeature(context.Background(), &pb.Point{Longitude: 121437403, Latitude: 310235000})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(feature)
}

func runSecond(client pb.RouteGuideClient) {
	serverStream, err := client.ListFeatures(context.Background(), &pb.Rectangle{
		Lo: &pb.Point{Longitude: 121358540, Latitude: 313374060},
		Hi: &pb.Point{Longitude: 121598790, Latitude: 311034130},
	})

	if err != nil {
		log.Fatalln(err)
	}

	for {
		feature, err := serverStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(feature)
	}

}

func runThird(client pb.RouteGuideClient) {
	points := []*pb.Point{
		{Latitude: 313374060, Longitude: 121358540},
		{Latitude: 311054130, Longitude: 121598790},
		{Latitude: 310235000, Longitude: 121437403},
	}

	clientStream, err := client.RecordRoute(context.Background())

	if err != nil {
		log.Fatalln(err)
	}

	for _, point := range points {
		if err := clientStream.Send(point); err != nil {
			log.Fatalln(err)
		}
		time.Sleep(time.Second)
	}

	summary, err := clientStream.CloseAndRecv()

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(summary)

}

func readIntFromCommandLine(reader *bufio.Reader, target *int32) {
	_, err := fmt.Fscanf(reader, "%d\n", target)
	if err != nil {
		log.Fatalln("Cannot scan", err)
	}
}

func runForth(client pb.RouteGuideClient) {
	stream, err := client.Recommend(context.Background())

	if err != nil {
		log.Fatalln(err)
	}

	// gorountine listen to the server
	go func() {
		feature, err2 := stream.Recv()

		if err2 != nil {
			log.Fatalln(err2)
		}

		fmt.Println("Recommend: ", feature)
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		// new a point for or the point in request would be nil
		request := pb.RecommendationRequest{Point: new(pb.Point)}
		var mode int32
		fmt.Print("Enter recommendation mode: (0 for farthest, 1 for nearest)")
		readIntFromCommandLine(reader, &mode)
		fmt.Print("Enter Latitude: ")
		readIntFromCommandLine(reader, &request.Point.Latitude)
		fmt.Print("Enter Longitude: ")
		readIntFromCommandLine(reader, &request.Point.Longitude)
		request.Mode = pb.RecommendationMode(mode)

		if err := stream.Send(&request); err != nil {
			log.Fatalln(err)
		}
		time.Sleep(time.Millisecond * 100)
	}

}

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		log.Fatalln("client cannot dial grpc server")
	}
	defer conn.Close()

	client := pb.NewRouteGuideClient(conn)

	runForth(client)
}
