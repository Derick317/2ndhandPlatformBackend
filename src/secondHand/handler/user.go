package handler

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"secondHand/model"
	"secondHand/service"
	"secondHand/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var mySigningKey = []byte("secret")

func signinHandler(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": "Cannot decode user data from client: " + err.Error()})
		return
	}

	success, err := service.CheckUser(&user, user.Email, user.Password, nil)

	if err != nil && !errors.Is(err, util.ErrUserNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if errors.Is(err, util.ErrUserNotFound) || !success {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      strconv.FormatUint(user.ID, 10),
		"expTime": strconv.FormatInt(time.Now().Add(time.Hour*24).Unix(), 10),
	})

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "username": user.Username, "id": user.ID})
}

func signupHandler(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailRegexp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	// usernameRegexp := regexp.MustCompile(`^[a-zA-Z0-9]$`)
	if user.Password == "" || user.Username == "" || !emailRegexp.MatchString(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Invalid username or password or email"})
		return
	}

	success, err := service.AddUser(&user, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !success {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Email is already taken"})
		return
	}

	c.JSON(http.StatusOK, "success")
}

func querySellerListHandler(c *gin.Context) {
	sellerId, err := getUserIdFromGinContent(c)
	if err != nil {
		return
	}
	list, err := service.QuerySellerList(sellerId, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// authMiddleware returns a handler function which verifies whether the token is valid
// and sets user's id
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"status": "Invalid token"})
			return
		}

		token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "Unable to parse token: " + err.Error()})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !token.Valid || !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"status": "Invalid token"})
			return
		}

		if expTime, err := strconv.ParseInt(claims["expTime"].(string), 10, 64); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": fmt.Sprintf("Unable to Parse expire time: %v", err)})
			return
		} else if expTime < time.Now().Unix() {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"status": "Token expired"})
			return
		}
		if user_id, err := strconv.ParseUint(claims["id"].(string), 10, 64); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": fmt.Sprintf("Unable to Parse id: %v", err)})
			return
		} else {
			c.Set("user_id", user_id)
		}
		// Call the next handler
		c.Next()
	}
}

func getUserIdFromGinContent(c *gin.Context) (uint64, error) {
	userId, exists := c.Get("user_id")
	uintId, ok := userId.(uint64)
	if !exists || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": util.ErrUnexpected.Error()})
		return uintId, util.ErrUnexpected
	}
	return uintId, nil
}
