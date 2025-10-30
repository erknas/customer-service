package controller

import (
	"context"

	"github.com/erknas/customer-service/internal/models"
	pb "github.com/erknas/customer-service/pkg/api/customer"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const layout = "2006-01-02"

type CustomerUseCase interface {
	RegisterCustomer(ctx context.Context, username, fullname, city, birthDate string) (models.Customer, error)
}

type controller struct {
	customerUseCase CustomerUseCase
	logger          *zap.Logger
	pb.UnimplementedCustomerServiceServer
}

func New(customerUseCase CustomerUseCase, logger *zap.Logger) *controller {
	return &controller{
		customerUseCase: customerUseCase,
		logger:          logger,
	}
}

func (c *controller) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error) {
	if err := req.ValidateAll(); err != nil {
		c.logger.Warn("validation falied", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := c.customerUseCase.RegisterCustomer(ctx, req.GetUserName(), req.GetFullName(), req.GetCity(), req.GetBirthDate())
	if err != nil {
		c.logger.Error("register", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create customer")
	}

	return &pb.CreateCustomerResponse{
		Id:        result.ID.String(),
		UserName:  result.Username,
		FullName:  result.Fullname,
		City:      result.City,
		BirthDate: result.BirthDate.Format(layout),
	}, nil
}
