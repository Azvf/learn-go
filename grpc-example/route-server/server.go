package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "github.com/Azvf/learn-go/grpc-example/route"
)

type RouteGuideServer struct {
	pb.UnimplementedRouteGuideServer
}

func (s *RouteGuideServer) GetFeature(context.Context, *pb.Point) (*pb.Feature, error) {
	return nil, nil
}

func (s *RouteGuideServer) ListFeatures(*pb.Rectangle, pb.RouteGuide_ListFeaturesServer) error {
	return nil
}

func (s *RouteGuideServer) RecordRoute(pb.RouteGuide_RecordRouteServer) error {
	return nil
}

func (s *RouteGuideServer) Recommend(pb.RouteGuide_RecommendServer) error {
	return nil
}

func newServer() *RouteGuideServer {
	return &RouteGuideServer{}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		log.Fatalln("cannot create a listener")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRouteGuideServer(grpcServer, newServer())
	log.Fatalln(grpcServer.Serve(lis))

}
