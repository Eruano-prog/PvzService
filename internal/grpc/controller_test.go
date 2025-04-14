package grpc_test

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/grpc"
	pvz_v1 "AvitoPvz/internal/grpc/api"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockPvzService struct {
	mock.Mock
}

func (m *mockPvzService) GetAllPvz(ctx context.Context) ([]models.PVZ, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.PVZ), args.Error(1)
}

func TestServer_GetPVZList(t *testing.T) {
	service := new(mockPvzService)
	server := grpc.NewServer(service)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		pvzID := uuid.New()
		registrationDate := time.Now()
		pvzs := []models.PVZ{
			{
				ID:               pvzID,
				City:             models.CityMSK,
				RegistrationDate: registrationDate,
			},
		}

		service.On("GetAllPvz", ctx).Return(pvzs, nil).Once()

		resp, err := server.GetPVZList(ctx, &pvz_v1.GetPVZListRequest{})

		require.NoError(t, err, "expected no error")
		require.NotNil(t, resp, "response should not be nil")
		require.Len(t, resp.Pvzs, 1, "expected one PVZ")

		pvz := resp.Pvzs[0]
		require.Equal(t, pvzID.String(), pvz.Id, "PVZ ID should match")
		require.Equal(t, string(models.CityMSK), pvz.City, "PVZ city should match")
		require.True(t, registrationDate.Sub(pvz.RegistrationDate.AsTime()) < time.Second, "registration date should match")

		service.AssertExpectations(t)
	})

	t.Run("empty result", func(t *testing.T) {
		service.On("GetAllPvz", ctx).Return([]models.PVZ{}, nil).Once()

		resp, err := server.GetPVZList(ctx, &pvz_v1.GetPVZListRequest{})

		require.NoError(t, err, "expected no error")
		require.NotNil(t, resp, "response should not be nil")
		require.Empty(t, resp.Pvzs, "expected empty PVZ list")

		service.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		serviceError := errors.New("database unavailable")

		service.On("GetAllPvz", ctx).Return([]models.PVZ{}, serviceError).Once()

		resp, err := server.GetPVZList(ctx, &pvz_v1.GetPVZListRequest{})

		require.Nil(t, resp, "response should be nil")
		require.Error(t, err, "expected error")
		st, ok := status.FromError(err)
		require.True(t, ok, "expected gRPC status error")
		require.Equal(t, codes.Internal, st.Code(), "expected Internal error code")
		require.Contains(t, st.Message(), serviceError.Error(), "error message should contain service error")

		service.AssertExpectations(t)
	})
}
