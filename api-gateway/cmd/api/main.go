package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"E-commerce_micro/api-gateway/cmd/api/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Iniciando a API Gateway...")

	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	productsServiceURL := os.Getenv("PRODUCTS_SERVICE_URL")

	if jwtSecret == "" || authServiceURL == "" || productsServiceURL == "" {
		log.Fatal("As variaveis de ambiente devem ser definidas")
	}

	router := gin.Default()

	authMiddleware := auth.NewAuthMiddleware(jwtSecret)

	proxyToAuth := createReverseProxy(authServiceURL)
	proxyToProducts := createReverseProxy(productsServiceURL)

	apiV1 := router.Group("/api/v1")
	{
		authRoutes := apiV1.Group("/users")
		authRoutes.Any("/*proxyPath", proxyToAuth)

		productsRoutes := apiV1.Group("/products")
		productsRoutes.Use(authMiddleware)
		productsRoutes.Any("/*proxyPath", proxyToProducts)
	}

	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "UP"}) })

	port := ":8000"
	log.Printf("Api gateway esta ouvindo na porta %s", port)
	if err := router.Run(port); err != nil {
		log.Fatal("Falha ao iniciar o API Gateway: %v", err)
	}
}

func createReverseProxy(target string) gin.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatal("URL de destino do proxy invalida: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.Director = func(req *http.Request) {
		req.Host = targetURL.Host
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host

		ginContext := req.Context().Value(gin.ContextKey).(*gin.Context)
		req.URL.Path = ginContext.Param("proxyPath")
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
