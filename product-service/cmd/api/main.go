package main

import (
	"log"
	"net"

	"E-commerce_micro/product-service/internal/database"
	"E-commerce_micro/product-service/internal/pb"
	"E-commerce_micro/product-service/internal/product"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Iniciando o Serviço de Produtos...")

	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Arquivo .env não encontrado, usando vars de ambiente.")
	}

	dbPool, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("FATAL: Não foi possível conectar ao banco de dados: %v", err)
	}
	defer dbPool.Close()
	log.Println("Conexão com o banco de dados (products_db) estabelecida.")

	productRepo := product.NewPostgresRepository(dbPool)
	log.Println("Repositório de produto inicializado.")
	productService := product.NewService(productRepo)
	log.Println("Serviço de produto inicializado.")

	go func() {
		grpcPort := ":50051"
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalf("falha ao ouvir porta gRCP: %v", err)
		}
		grpcServer := grpc.NewServer()
		productGrcpServer := product.NewGrcpServer(productService)
		pb.RegisterProductServiceServer(grpcServer, productGrcpServer)

		log.Printf("Servidor gRPC rodando na porta %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("falha ao servir gRPC: %v", err)
		}

	}()

	productHandler := product.NewHandler(productService)
	log.Println("Handler de produto inicializado.")

	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		productHandler.RegisterRoutes(apiV1)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	port := ":8081"

	log.Printf("Servidor web ouvindo na porta %s", port)

	if err := router.Run(port); err != nil {
		log.Fatalf("FATAL: Falha ao iniciar o servidor Gin: %v", err)
	}
}
