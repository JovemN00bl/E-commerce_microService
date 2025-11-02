package main

import (
	"log"
	"os"

	"E-commerce_micro/auth-service/internal/database"
	"E-commerce_micro/auth-service/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Iniciando o Serviço de Autenticação...")

	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Arquivo .env não encontrado, usando variáveis de ambiente do sistema.")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("FATAL: A variável de ambiente JWT_SECRET não está definida!")
	}

	dbPool, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("FATAL: Não foi possível conectar ao banco de dados: %v", err)
	}

	defer dbPool.Close()
	log.Println("Conexão com o banco de dados estabelecida.")

	userRepo := user.NewPostgresRepository(dbPool)
	log.Println("Repositório de usuário inicializado.")

	userService := user.NewService(userRepo, jwtSecret)
	log.Println("Serviço de usuário inicializado.")

	userHandler := user.NewHandler(userService)
	log.Println("Handler de usuário inicializado.")

	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		userHandler.RegisterRoutes(apiV1)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	port := ":8080"
	log.Printf("Servidor web ouvindo na porta %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("FATAL: Falha ao iniciar o servidor Gin: %v", err)
	}
}
