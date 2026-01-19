package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kartikgujjar22/tokenvault/internal/database"
)

// StoreRequest is the JSON body we expect from the backend
type StoreRequest struct {
	Project string `json:"project" binding:"required"`
	Token   string `json:"token" binding:"required"`
	Tag     string `json:"tag"` // Optional: defaults to "default" if empty
}

// HandleStoreToken receives the token from the backend
// POST /store
func HandleStoreToken(c *gin.Context) {
	var req StoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}
	// Default tag if not provided
	if req.Tag == "" {
		req.Tag = "default"
	}
	// Call the NEW V2 database function
	err := database.StoreToken(req.Project, req.Tag, req.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt and store token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "securely stored", "project": req.Project, "tag": req.Tag})
}

// HandleFetchToken retrieves the token for Postman
// GET /fetch/:project?tag=admin
func HandleFetchToken(c *gin.Context) {
	project := c.Param("project")

	// Get "tag" from query param (e.g., ?tag=admin)
	tag := c.DefaultQuery("tag", "default")

	token, meta, err := database.FetchToken(project, tag)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found or decryption failed"})
		return
	}

	// Return the token + metadata
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"meta":  meta,
	})
}
