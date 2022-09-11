package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"math"
	"net"
	"time"

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

func inRange(point *pb.Point, rectangle *pb.Rectangle) bool {
	left := math.Min(float64(rectangle.Lo.Longitude), float64(rectangle.Hi.Longitude))
	right := math.Max(float64(rectangle.Lo.Longitude), float64(rectangle.Hi.Longitude))
	top := math.Max(float64(rectangle.Lo.Latitude), float64(rectangle.Hi.Latitude))
	bottom := math.Min(float64(rectangle.Lo.Latitude), float64(rectangle.Hi.Latitude))

	if float64(point.Longitude) > left &&
		float64(point.Longitude) < right &&
		float64(point.Latitude) > bottom &&
		float64(point.Latitude) < top {
		return true
	}

	return false
}

func (s *RouteGuideServer) ListFeatures(rectangle *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
	for _, feature := range s.features {
		if inRange(feature.Location, rectangle) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	return nil
}

func toRadians(num float64) float64 {
	return num * math.Pi / float64(180)
}

// calcDistance calculates the distance between two points using the "haversine" formula.
// The formula is based on http://mathforum.org/library/drmath/view/51879.html.
func calcDistance(p1 *pb.Point, p2 *pb.Point) int32 {
	const CordFactor float64 = 1e7
	const R = float64(6371000) // earth radius in metres
	lat1 := toRadians(float64(p1.Latitude) / CordFactor)
	lat2 := toRadians(float64(p2.Latitude) / CordFactor)
	lng1 := toRadians(float64(p1.Longitude) / CordFactor)
	lng2 := toRadians(float64(p2.Longitude) / CordFactor)
	dlat := lat2 - lat1
	dlng := lng2 - lng1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return int32(distance)
}

func (s *RouteGuideServer) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
	startTime := time.Now()
	var pointCount, distance int32
	var prevPoint *pb.Point
	for {
		point, err := stream.Recv()

		if err == io.EOF {
			// gen a summary and send back
			endTime := time.Now()
			return stream.SendAndClose(&pb.RouteSummary{
				PointCount:  pointCount,
				Distance:    distance,
				ElapsedTime: int32(endTime.Sub(startTime).Seconds())},
			)
		}

		if err != nil {
			log.Fatalln(err)
		}

		pointCount++
		if prevPoint != nil {
			distance += calcDistance(prevPoint, point)
		}
		prevPoint = point
	}
	return nil
}

func (s *RouteGuideServer) recommendOnce(request *pb.RecommendationRequest) (*pb.Feature, error) {
	var nearest, farthest *pb.Feature
	var nearestDistance, farthestDistance int32

	for _, feature := range s.features {
		distance := calcDistance(feature.Location, request.Point)
		if nearest == nil || distance < nearestDistance {
			nearestDistance = distance
			nearest = feature
		}
		if farthest == nil || distance > farthestDistance {
			farthestDistance = distance
			farthest = feature
		}
	}
	if request.Mode == pb.RecommendationMode_GetFarthest {
		return farthest, nil
	} else {
		return nearest, nil
	}
}

func (s *RouteGuideServer) Recommend(stream pb.RouteGuide_RecommendServer) error {
	for {
		request, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		recommended, err2 := s.recommendOnce(request)

		if err2 != nil {
			return err
		}

		return stream.Send(recommended)
	}

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
