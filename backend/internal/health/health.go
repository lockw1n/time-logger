package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/health", func(c *gin.Context) {
		sqlDB, _ := db.DB()
		if err := sqlDB.Ping(); err != nil {
			c.String(http.StatusServiceUnavailable, "database not reachable")
			return
		}
		c.String(http.StatusOK, "ok")
	})
	r.HEAD("/health", func(c *gin.Context) {
		sqlDB, _ := db.DB()
		if err := sqlDB.Ping(); err != nil {
			c.String(http.StatusServiceUnavailable, "database not reachable")
			return
		}
		c.String(http.StatusOK, "ok")
	})
}
