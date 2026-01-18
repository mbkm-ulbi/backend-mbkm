package middleware

import (
	"strings"
	"time"

	"mbkm-go/config"
	"mbkm-go/database"
	"mbkm-go/internal/models"
	"mbkm-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a user
func GenerateToken(user *models.User) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.AppConfig.JWTExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

// JWTAuth returns a JWT authentication middleware
func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return utils.UnauthorizedError(c, "Missing authorization header")
		}

		// Check for Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return utils.UnauthorizedError(c, "Invalid authorization header format")
		}

		tokenString := parts[1]

		// Parse token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil {
			return utils.UnauthorizedError(c, "Invalid or expired token")
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok || !token.Valid {
			return utils.UnauthorizedError(c, "Invalid token claims")
		}

		// Fetch user from database
		var user models.User
		if err := database.DB.Preload("Roles").First(&user, claims.UserID).Error; err != nil {
			return utils.UnauthorizedError(c, "User not found")
		}

		// Store user in context
		c.Locals("user", &user)
		c.Locals("userId", claims.UserID)
		c.Locals("claims", claims)

		return c.Next()
	}
}

// OptionalJWTAuth returns an optional JWT authentication middleware
// It doesn't return error if token is not present, but sets user if valid token exists
func OptionalJWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Next()
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil {
			return c.Next()
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok || !token.Valid {
			return c.Next()
		}

		var user models.User
		if err := database.DB.Preload("Roles").First(&user, claims.UserID).Error; err != nil {
			return c.Next()
		}

		c.Locals("user", &user)
		c.Locals("userId", claims.UserID)
		c.Locals("claims", claims)

		return c.Next()
	}
}

// GetCurrentUser gets the current authenticated user from context
func GetCurrentUser(c *fiber.Ctx) *models.User {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return nil
	}
	return user
}

// GetCurrentUserID gets the current authenticated user ID from context
func GetCurrentUserID(c *fiber.Ctx) uint {
	userID, ok := c.Locals("userId").(uint)
	if !ok {
		return 0
	}
	return userID
}

// GetRoleIDs returns role IDs of current user
func GetRoleIDs(c *fiber.Ctx) []uint {
	user := GetCurrentUser(c)
	if user == nil {
		return []uint{}
	}

	roleIDs := make([]uint, len(user.Roles))
	for i, role := range user.Roles {
		roleIDs[i] = role.ID
	}
	return roleIDs
}

// HasRole checks if current user has specific role
func HasRole(c *fiber.Ctx, roleID uint) bool {
	roleIDs := GetRoleIDs(c)
	for _, id := range roleIDs {
		if id == roleID {
			return true
		}
	}
	return false
}

// IsAdmin checks if current user is admin (role ID 1)
func IsAdmin(c *fiber.Ctx) bool {
	return HasRole(c, 1)
}

// IsCDC checks if current user is CDC (role ID 3)
func IsCDC(c *fiber.Ctx) bool {
	return HasRole(c, 3)
}

// IsCompany checks if current user is Company/Mitra (role ID 4)
func IsCompany(c *fiber.Ctx) bool {
	return HasRole(c, 4)
}
