package product

import (
	"E-commerce_micro/product-service/internal/pb"
	"context"
	"fmt"
)

type GrcpServer struct {
	pb.UnimplementedProductServiceServer
	service Service
}

func NewGrcpServer(s Service) *GrcpServer {
	return &GrcpServer{service: s}
}

func (g *GrcpServer) CheckStock(ctx context.Context, req *pb.CheckStockRequest) (*pb.CheckStockResponse, error) {
	product, err := g.service.GetById(ctx, req.ProductId)
	if err != nil {
		return &pb.CheckStockResponse{Price: 0, InStock: false}, nil
	}

	inStock := product.StockQuantity > 0
	fmt.Printf("gRPC: Verificando estoque para %s. Tem? %v\n", product.Name, inStock)

	return &pb.CheckStockResponse{Price: product.Price, InStock: true}, nil
}
