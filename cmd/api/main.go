package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/adifsndri2/Aplikasi-Keuangan-Internal-Golang/tree/develop/cmd/internal/ap"
	"github.com/adifsndri2/Aplikasi-Keuangan-Internal-Golang/tree/develop/cmd/internal/auth"
	"github.com/adifsndri2/Aplikasi-Keuangan-Internal-Golang/tree/develop/cmd/internal/ar"
	"github.com/adifsndri2/Aplikasi-Keuangan-Internal-Golang/tree/develop/cmd/internal/model"
	"github.com/adifsndri2/Aplikasi-Keuangan-Internal-Golang/tree/develop/cmd/internal/reports"
)

func main() {
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&model.User{}, &model.Report{}, &model.ImportBatch{}, &model.APTransaction{}, &model.ARTransaction{})
	auth.SeedAdmin(db, os.Getenv("ADMIN_EMAIL"), os.Getenv("ADMIN_PASSWORD"))
	secret := os.Getenv("JWT_SECRET")
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.POST("/api/v1/auth/login", auth.Login(db, secret))

	api := r.Group("/api/v1", auth.Required(db, secret))
	users := api.Group("/users", auth.RequiredRole("admin"))
	users.GET("", auth.ListUsers(db))
	users.POST("", auth.CreateUser(db))
	users.PATCH("/:id", auth.PatchUser(db))

	ap := ap.New(db)
	api.POST("/ap/import/preview", ap.Preview)
	api.POST("/ap/import/commit", ap.Import)
	api.GET("/ap", ap.List)
	api.GET("/ap/:id", ap.Get)

	ar := ar.New(db)
	api.POST("/ar/import/preview", ar.Preview)
	api.POST("/ar/import/commit", ar.Commit)
	api.GET("/ar", ar.List)
	api.GET("/ar/:id", ar.Detail)

	api.POST("/reports/forecasting/generate", reports.Generate(db, "forecasting"))
	api.POST("/reports/profit-loss/generate", reports.Generate(db, "profit-loss"))
	api.POST("/reports/cash-flow/generate", reports.Generate(db, "cash-flow"))
	api.GET("/reports/:id", func(c *gin.Context) {
		if t, ok := reports.Load(db, c.Param("id")); ok {
			c.JSON(200, t)
			return
		}
		c.JSON(404, gin.H{"error": "Laporan Tidak Ditemukan"})
	})
	log.Fatal(r.Run(":8080"))
}