package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

// Protected protect routes
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte("SECRET"),
		ErrorHandler: jwtError,
		ContextKey:   "jwt",
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari header Authorization
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: Missing token"})
		}

		// Validasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte("secret"), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: Invalid token"})
		}

		// Simpan klaim ke context
		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user", claims)

		return c.Next()
	}
}

//func validToken(t *jwt.Token, id string) bool {
//	n, err := strconv.Atoi(id)
//	if err != nil {
//		return false
//	}
//
//	claims := t.Claims.(jwt.MapClaims)
//	uid := int(claims["user_id"].(float64))
//
//	if uid != n {
//		return false
//	}
//
//	return true
//}
//
//func validUser(id string, p string) bool {
//	db := database.DB
//	var user model.User
//	db.First(&user, id)
//	if user.Username == "" {
//		return false
//	}
//	if !CheckPasswordHash(p, user.Password) {
//		return false
//	}
//	return true
//}
