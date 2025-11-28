package middleware

import (
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware checks if the user has admin privileges
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")
		if userID == 0 {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Check if user has admin role
		var user model.User
		if err := db.DB.Preload("Roles").First(&user, userID).Error; err != nil {
			c.JSON(404, gin.H{"error": "user not found"})
			c.Abort()
			return
		}

		isAdmin := false
		for _, role := range user.Roles {
			if role.Name == "admin" || role.Name == "super_admin" {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			c.JSON(403, gin.H{"error": "admin privileges required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
