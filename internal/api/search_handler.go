package api

import (
	"net/http"
	"netfilessys/internal/middleware"
	"netfilessys/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchService *service.SearchService
}

func NewSearchHandler() *SearchHandler {
	return &SearchHandler{
		searchService: service.NewSearchService(),
	}
}

// Search handles global search requests
func (h *SearchHandler) Search(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req SearchRequest
	if !middleware.BindAndValidate(c, &req) {
		return
	}

	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}
	if req.Type == "" {
		req.Type = "all"
	}

	results, err := h.searchService.Search(userID, req.Query, req.Type, req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Search failed: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": results,
	})
}

// SearchSuggestions returns search suggestions
func (h *SearchHandler) SearchSuggestions(c *gin.Context) {
	userID := c.GetUint("user_id")
	query := c.Query("q")

	if query == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": []string{},
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	suggestions, err := h.searchService.GetSearchSuggestions(userID, query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to get suggestions: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": suggestions,
	})
}
