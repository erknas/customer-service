package controller

import (
	"context"

	"github.com/erknas/customer-service/internal/models"
	pb "github.com/erknas/customer-service/pkg/api/customer"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CustomerUseCase interface {
	AddCustomer(ctx context.Context, username, fullname, city, birthDate string) (*models.Customer, error)
	GetCustomer(ctx context.Context, id string) (*models.Customer, error)
	UpdateCustomer(ctx context.Context, id string, userName, fullName, city, birthDateStr *string) (*models.Customer, error)
}

type controller struct {
	useCase CustomerUseCase
	logger  *zap.Logger
	pb.UnimplementedCustomerServiceServer
}

func New(useCase CustomerUseCase, logger *zap.Logger) *controller {
	return &controller{
		useCase: useCase,
		logger:  logger,
	}
}

func (c *controller) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error) {
	if err := req.ValidateAll(); err != nil {
		c.logger.Warn("validation failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	customer, err := c.useCase.AddCustomer(ctx, req.GetUserName(), req.GetFullName(), req.GetCity(), req.GetBirthDate())
	if err != nil {
		c.logger.Error("register", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateCustomerResponse{
		Customer: domainToPb(customer),
	}, nil
}

func (c *controller) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.GetCustomerResponse, error) {
	if err := req.ValidateAll(); err != nil {
		c.logger.Warn("validation failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	customer, err := c.useCase.GetCustomer(ctx, req.GetId())
	if err != nil {
		c.logger.Error("failed to get customer", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetCustomerResponse{
		Customer: domainToPb(customer),
	}, nil
}

func (c *controller) UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (*pb.UpdateCustomerResponse, error) {
	if err := req.ValidateAll(); err != nil {
		c.logger.Warn("validation failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	customer, err := c.useCase.UpdateCustomer(ctx, req.GetId(), req.UserName, req.FullName, req.City, req.BirthDate)
	if err != nil {
		c.logger.Error("failed to update customer", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateCustomerResponse{
		Customer: domainToPb(customer),
	}, nil
}
