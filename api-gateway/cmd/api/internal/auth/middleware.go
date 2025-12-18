package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func NewAuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Erro: cabeçalho 'Authorization' ausente")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Autorização necessária"})
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			log.Println("Erro: Formato do cabeçalho 'Authorization' inválido")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "Formato de autorização inválido"})
			return
		}

		tokenString := headerParts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			log.Printf("Erro na validação do token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}

		if !token.Valid {
			log.Println("Erro: Token inválido")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}

		c.Next()

	}

}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// 1. Debug: Vamos ver exatamente o que está chegando
		if authHeader == "" {
			fmt.Println("❌ Middleware: Header Authorization vazio")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			return
		}

		// 2. Limpeza Robusta
		// strings.Fields remove todos os espaços extras, quebras de linha e tabulações
		parts := strings.Fields(authHeader)

		// O formato tem que ser ["Bearer", "eyJ..."]
		if len(parts) < 2 || strings.ToLower(parts[0]) != "bearer" {
			fmt.Printf("❌ Middleware: Formato inválido. Recebido: %v\n", parts)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Formato inválido. Use: Bearer <token>"})
			return
		}

		tokenString := parts[1] // Pega só o código

		// 3. Remover aspas extras (caso tenham sobrado do Postman/Frontend)
		tokenString = strings.Trim(tokenString, "\"")

		// Debug: Mostra o token limpo que será validado
		// fmt.Printf("🔍 Tentando validar token limpo: %s...\n", tokenString[:10])

		// 4. Validação
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método inesperado: %v", token.Header["alg"])
			}
			// ⚠️ GARANTA QUE ESTA CHAVE É A MESMA DO AUTH-SERVICE
			return []byte("sua_chave_secreta_super_secreta"), nil
		})

		if err != nil || !token.Valid {
			fmt.Printf("❌ Erro JWT: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
			return
		}

		// 5. Sucesso! Extrair dados e passar para frente
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if sub, ok := claims["sub"].(string); ok {
				// Adiciona o ID do usuário no header para o Products/Orders service saberem quem é
				c.Request.Header.Set("X-User-Id", sub)
			}
		}

		c.Next()
	}
}
