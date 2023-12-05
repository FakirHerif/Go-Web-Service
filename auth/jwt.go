package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"example.com/webservice/models"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func Login(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GEÇERSİZ KİMLİK BİLGİLERİ"})
		return
	}

	users, err := models.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "KULLANICI VERİLERİ ALIRKEN HATA OLUŞTU"})
		return
	}

	var foundUser models.User
	for _, user := range users {
		if user.Username == creds.Username {
			foundUser = user
			break
		}
	}

	// Kullanıcı bulunamadıysa veya şifre eşleşmiyorsa hata ver
	if foundUser.Username == "" || foundUser.Password != creds.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID VEYA ŞİFRE YANLIŞ."})
		return
	}

	expirationTime := time.Now().Add(10 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		Role:     foundUser.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "TOKEN OLUŞTURULAMADI"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization BAŞLIĞI SAĞLANAMADI"})
			c.Abort()
			return
		}

		tokenString := authHeader[len("Bearer "):]

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "GEÇERSİZ TOKEN"})
			c.Abort()
			return
		}

		if (c.Request.Method == "DELETE" || c.Request.Method == "PUT") && claims.Role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Yetkisiz İşlem"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func SecuredEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GÜVENLİ UÇ NOKTAYA ERİŞTİNİZ !!!"})
}
