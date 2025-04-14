package grpc

import (
	pvz_v1 "AvitoPvz/internal/grpc/api"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pvz_v1.UnimplementedPVZServiceServer
	service PvzService
}

func NewServer(service PvzService) *Server {
	return &Server{service: service}
}

func (s Server) GetPVZList(ctx context.Context, in *pvz_v1.GetPVZListRequest) (*pvz_v1.GetPVZListResponse, error) {
	res, err := s.service.GetAllPvz(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var results []*pvz_v1.PVZ
	for _, v := range res {
		results = append(results, &pvz_v1.PVZ{
			Id:               v.ID.String(),
			City:             string(v.City),
			RegistrationDate: timestamppb.New(v.RegistrationDate),
		})
	}

	return &pvz_v1.GetPVZListResponse{
		Pvzs: results,
	}, nil
}
