package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"log"
	"net"

	pb "github.com/Azvf/learn-go/grpc-example/route"
)

type RouteGuideServer struct {
	features []*pb.Feature
	pb.UnimplementedRouteGuideServer
}

func (s *RouteGuideServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	for _, feature := range s.features {
		if proto.Equal(feature.Location, point) {
			return feature, nil
		}
	}
	fmt.Println("point not found")
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
	return &RouteGuideServer{
		features: []*pb.Feature{
			{Name: "上海交通大学闵行校区 上海市闵行区东川路800号", Location: &pb.Point{
				Longitude: 121437403,
				Latitude:  310235000,
			}},
			{Name: "复旦大学 上海市杨浦区五角场邯郸路220号", Location: &pb.Point{
				Longitude: 121503457,
				Latitude:  312978870,
			}},
			{Name: "华东理工大学 上海市徐汇区梅陇路130号", Location: &pb.Point{
				Longitude: 121424904,
				Latitude:  311416130,
			}},
		},
	}
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
