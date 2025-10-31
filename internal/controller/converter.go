package controller

import (
	"github.com/erknas/customer-service/internal/models"
	pb "github.com/erknas/customer-service/pkg/api/customer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const layout = "2006-01-02"

func domainToPb(customer *models.Customer) *pb.Customer {
	return &pb.Customer{
		Id:        customer.ID.String(),
		UserName:  customer.Username,
		FullName:  customer.Fullname,
		City:      customer.City,
		BirthDate: customer.BirthDate.Format(layout),
		IsActive:  customer.IsActive,
		CreatedAt: timestamppb.New(customer.CreatedAt),
		UpdatedAt: timestamppb.New(customer.UpdatedAt),
	}
}
