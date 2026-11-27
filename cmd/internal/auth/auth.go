package auth

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/adifsndri2/Aplikasi-Keuangan-Internal-Golang/tree/develop/cmd/internal/model"
)

type Claims struct {
	UID uint `json:"uid"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func SeedAdmin(db *gorm.DB, email, pass string) {
	var n int64
	db.Model(&model.User{}).Where("role = ?", "admin").Count(&n)
	if n > 0 {
		return
	}
	h , _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	db.Create(&model.User{Name: "admin", Email: email, PasswordHash: string(h), Role: "admin", IsActive: true})
}

func Login(db *gorm.DB, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email"`
			Password string `json:"password"`
		}
		if c.ShouldBindJSON(&req) != nil {
			c.JSON(400, gin.H{"error": "Input Tidak Valid"})
			return
		}
		var u model.User
		if db.Where("email = ?", req.Email).First(&u).Error != nil ||
			bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
			c.JSON(401, gin.H{"error": "Email atau Password Salah"})
			return
		}
		if !u.IsActive {
			c.JSON(403, gin.H{"error": "Akun Tidak Aktif"})
			return
		}
		claims := Claims{UID: u.ID, Role: u.Role, RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour))}}
			tok, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
		c.JSON(200, gin.H{"token": tok, "role": u.Role})
	}
}

func Required(db *gorm.DB, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		claims := &Claims{}
		t, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		}, jwt.WithValidMethods([]string{"HS256"}))
		if err != nil || !t.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Token Tidak Valid atau Kadaluarsa"})
			return
		}
		var u model.User
		if db.First(&u, claims.UID).Error != nil || !u.IsActive {
			c.AbortWithStatusJSON(403, gin.H{"error": "Akun Tidak Aktif"})
			return
		}
		c.Set("uid", u.ID)
		c.Set("role", u.Role)
		c.Next()
	}
}

func RequiredRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.GetString("role")
		for _, x := range roles {
			if r == x {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(403, gin.H{"error": "Akses Ditolak"})
	}
}
