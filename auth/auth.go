package auth

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// func AccessToken(c *gin.Context) {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{ // Create a jwt token with the claims
// 		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(), // 5 minutes

// 	})

// 	ss, err := token.SignedString([]byte("==signature==")) // Symmetric key
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"token": ss})

// }

func AccessToken(signature string) gin.HandlerFunc {
	return  func(c *gin.Context) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{ // Create a jwt token with the claims
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(), // 5 minutes
		Audience: "GibGyb",

	})

	ss, err := token.SignedString([]byte(signature)) // Symmetric key
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": ss})

}

} 