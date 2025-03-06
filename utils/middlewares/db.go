package middlewares

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/yuchanns/kong-exercise-microservices/utils/helpers"
	"gorm.io/gorm"
)

// UseTenantDB is a middleware that sets up a tenant-specific database connection
// based on the x-organization-id header.
// It then injects the database connection into the request context.
// The database connection is closed after the request is processed.
// To retrieve the database connection, use the helpers.GetTenantDB function.
func UseTenantDB() gin.HandlerFunc {
	isDebug := os.Getenv("GO_DEBUG") == "1"

	return func(ctx *gin.Context) {
		orgID := ctx.GetHeader("x-organization-id")
		if orgID == "" {
			orgID = ctx.Query("organization_id")
		}
		if orgID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing x-organization-id header"})
			return
		}

		db, err := gorm.Open(sqlite.Open("catalog_"+orgID+".db"), &gorm.Config{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if isDebug {
			db = db.Debug()
		}

		defer func() {
			sqlDB, err := db.DB()
			if err == nil && sqlDB != nil {
				_ = sqlDB.Close()
			}
		}()

		c := helpers.WithTenantDB(ctx.Request.Context(), db)
		ctx.Request = ctx.Request.WithContext(c)

		ctx.Next()
	}
}
