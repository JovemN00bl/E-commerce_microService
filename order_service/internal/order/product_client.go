package order

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"E-commerce_micro/order_service/internal/pb"
)

type ProductClient interface {
	CheckStock(ctx context.Context, productId string) (price float64, InStock bool, err error)
	Close()
}

type grpcProductClient struct {
	conn   *grpc.ClientConn
	client pb.ProductServiceClient
}

func NewProductClient(target string) (ProductClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Falha ao conectar via gRPC: %w", err)
	}

	client := pb.NewProductServiceClient(conn)
	log.Printf("Conectado ao serviço de produtos via gRPC!")

	return &grpcProductClient{conn: conn, client: client}, nil
}

func (c *grpcProductClient) CheckStock(ctx context.Context, productId string) (price float64, InStock bool, err error) {
	req := &pb.CheckStockRequest{ProductId: productId}

	resp, err := c.client.CheckStock(ctx, req)
	if err != nil {
		return 0, false, fmt.Errorf("erro na chamda gRPC: %w", err)
	}

	return resp.Price, resp.InStock, nil
}

func (c *grpcProductClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
