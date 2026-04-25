package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"hospital-backend-api/handler"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthRequired เป็น middleware ตรวจสอบ JWT token
// ดึง token จาก Authorization: Bearer <token>
// ถ้า valid → แนบ staff_id, hospital_id ไว้ใน context
// ถ้าไม่ valid → return 401 Unauthorized
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. ดึง Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "กรุณาเข้าสู่ระบบก่อน (ไม่พบ token)"})
			c.Abort()
			return
		}

		// 2. ตรวจว่าเป็น Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "รูปแบบ token ไม่ถูกต้อง (ต้องเป็น Bearer <token>)"})
			c.Abort()
			return
		}
		tokenString := parts[1]

		// 2.5 ตรวจว่า token ถูก logout (blacklist) แล้วหรือยัง
		if handler.IsBlacklisted(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token นี้ถูก logout แล้ว กรุณาเข้าสู่ระบบใหม่"})
			c.Abort()
			return
		}

		// 3. Parse & validate token
		jwtSecret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// ตรวจสอบ signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token ไม่ถูกต้องหรือหมดอายุ"})
			c.Abort()
			return
		}

		// 4. ดึง claims แล้วแนบไว้ใน context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่สามารถอ่านข้อมูลจาก token ได้"})
			c.Abort()
			return
		}

		// แนบข้อมูลใน context เพื่อให้ handler ถัดไปเอาไปใช้ได้
		c.Set("staff_id", uint(claims["staff_id"].(float64)))
		c.Set("hospital_id", uint(claims["hospital_id"].(float64)))
		c.Set("username", claims["username"].(string))
		
		c.Next()
	}
}
