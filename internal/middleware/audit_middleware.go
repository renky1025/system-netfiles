package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditMiddleware logs all API requests for audit purposes
func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read body for POST/PUT requests
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Process request
		c.Next()

		// Only log authenticated requests
		userID := c.GetUint("userID")
		if userID == 0 {
			return
		}

		// Log file operations asynchronously
		go logOperation(c.Copy(), userID, bodyBytes)
	}
}

func logOperation(c *gin.Context, userID uint, bodyBytes []byte) {
	path := c.Request.URL.Path
	method := c.Request.Method

	var opType string
	var fileID uint
	var details string

	// Parse operation type and extract file ID
	switch {
	case method == "POST" && containsPath(path, "/upload"):
		opType = "upload"
	case method == "POST" && containsPath(path, "/instant-upload"):
		opType = "instant_upload"
	case containsPath(path, "/download"):
		opType = "download"
		fileID = extractFileIDFromPath(path)
	case method == "DELETE" && containsPath(path, "/file/"):
		opType = "delete"
		fileID = extractFileIDFromPath(path)
	case method == "POST" && containsPath(path, "/move"):
		opType = "move"
		fileID = extractFileIDFromPath(path)
	case method == "POST" && containsPath(path, "/copy"):
		opType = "copy"
		fileID = extractFileIDFromPath(path)
	case method == "PUT" && containsPath(path, "/rename"):
		opType = "rename"
		fileID = extractFileIDFromPath(path)
	case method == "POST" && containsPath(path, "/batch/delete"):
		opType = "batch_delete"
		details = extractBatchIDs(bodyBytes)
	case method == "POST" && containsPath(path, "/batch/move"):
		opType = "batch_move"
		details = extractBatchIDs(bodyBytes)
	case method == "POST" && containsPath(path, "/batch/copy"):
		opType = "batch_copy"
		details = extractBatchIDs(bodyBytes)
	case method == "POST" && containsPath(path, "/share/create"):
		opType = "share_create"
	case method == "DELETE" && containsPath(path, "/share/"):
		opType = "share_delete"
	case method == "POST" && containsPath(path, "/restore"):
		opType = "restore"
		fileID = extractFileIDFromPath(path)
	default:
		return // Don't log other operations
	}

	if details == "" {
		details = path
	}

	// Create log entry
	log := &model.FileOpLog{
		UserID:    userID,
		FileID:    fileID,
		OpType:    opType,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Details:   details,
		CreatedAt: time.Now(),
	}

	db.DB.Create(log)
}

func containsPath(path, substr string) bool {
	return strings.Contains(path, substr)
}

func extractFileIDFromPath(path string) uint {
	// Match patterns like /file/123/download or /file/123
	re := regexp.MustCompile(`/file/(\d+)`)
	matches := re.FindStringSubmatch(path)
	if len(matches) > 1 {
		id, _ := strconv.ParseUint(matches[1], 10, 64)
		return uint(id)
	}

	// Match patterns like /recycle/123/restore
	re = regexp.MustCompile(`/recycle/(\d+)`)
	matches = re.FindStringSubmatch(path)
	if len(matches) > 1 {
		id, _ := strconv.ParseUint(matches[1], 10, 64)
		return uint(id)
	}

	return 0
}

func extractBatchIDs(bodyBytes []byte) string {
	if len(bodyBytes) == 0 {
		return ""
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return ""
	}

	if fileIDs, ok := data["file_ids"].([]interface{}); ok {
		ids := make([]string, len(fileIDs))
		for i, id := range fileIDs {
			ids[i] = strconv.FormatFloat(id.(float64), 'f', 0, 64)
		}
		return "file_ids: [" + strings.Join(ids, ",") + "]"
	}

	return ""
}
