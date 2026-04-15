package http

import (
	"net/http"
	"strconv"

	"dangde-world/backend/internal/application/auth"
	"dangde-world/backend/internal/application/catalog"
	"dangde-world/backend/internal/application/dashboard"
	"dangde-world/backend/internal/application/runtime"
	"dangde-world/backend/internal/domain/activity"
	"dangde-world/backend/internal/domain/activitydata"
	"dangde-world/backend/internal/infrastructure/persistence"
	"dangde-world/backend/internal/shared/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	auth      auth.Service
	catalog   catalog.Service
	dashboard dashboard.Service
	runtime   runtime.Service
}

type loginRequest struct {
	UserID uint `json:"userId"`
}

type activityRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Language    string `json:"language" binding:"required"`
	Difficulty  string `json:"difficulty" binding:"required"`
	AgeGroup    string `json:"ageGroup" binding:"required"`
	CategoryID  uint   `json:"categoryId" binding:"required"`
	Prompt      string `json:"prompt" binding:"required"`
	Icon        string `json:"icon"`
	HTMLCode    string `json:"htmlCode"`
	CSSCode     string `json:"cssCode"`
	JSCode      string `json:"jsCode"`
}

type assignmentRequest struct {
	ParentID   uint `json:"parentId" binding:"required"`
	KidID      uint `json:"kidId" binding:"required"`
	ActivityID uint `json:"activityId" binding:"required"`
}

type assignmentUpdateRequest struct {
	Status   string `json:"status"`
	Progress *int   `json:"progress"`
}

type activityDataRequest struct {
	ActivityID   uint   `json:"activityId" binding:"required"`
	AssignmentID uint   `json:"assignmentId" binding:"required"`
	KidID        uint   `json:"kidId" binding:"required"`
	Key          string `json:"key" binding:"required"`
	Value        string `json:"value" binding:"required"`
}

type activityDataUpdateRequest struct {
	Key   string `json:"key"`
	Value string `json:"value" binding:"required"`
}

func NewRouter(db *gorm.DB, cfg config.Config) *gin.Engine {
	userRepo := persistence.NewUserRepository(db)
	categoryRepo := persistence.NewCategoryRepository(db)
	activityRepo := persistence.NewActivityRepository(db)
	assignmentRepo := persistence.NewAssignmentRepository(db)
	activityDataRepo := persistence.NewActivityDataRepository(db)

	handler := Handler{
		auth:      auth.NewService(userRepo),
		catalog:   catalog.NewService(userRepo, categoryRepo, activityRepo, assignmentRepo),
		dashboard: dashboard.NewService(userRepo, assignmentRepo),
		runtime:   runtime.NewService(activityDataRepo),
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{cfg.FrontendURL, "http://127.0.0.1:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Content-Type"},
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	{
		api.POST("/auth/login", handler.login)
		api.GET("/users", handler.listUsers)
		api.GET("/categories", handler.listCategories)
		api.GET("/activities", handler.listActivities)
		api.POST("/activities", handler.createActivity)
		api.PUT("/activities/:id", handler.updateActivity)
		api.GET("/assignments", handler.listAssignments)
		api.POST("/assignments", handler.createAssignment)
		api.PATCH("/assignments/:id", handler.updateAssignment)
		api.GET("/activity-data", handler.listActivityData)
		api.POST("/activity-data", handler.createActivityData)
		api.PUT("/activity-data/:id", handler.updateActivityData)
		api.DELETE("/activity-data/:id", handler.deleteActivityData)
		api.GET("/stats/parent/:parentId", handler.parentStats)
	}

	return router
}

func (h Handler) login(c *gin.Context) {
	var payload loginRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.auth.Login(payload.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": item})
}

func (h Handler) listUsers(c *gin.Context) {
	var parentID *uint
	if raw := c.Query("parentId"); raw != "" {
		value, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parentId"})
			return
		}
		parsed := uint(value)
		parentID = &parsed
	}

	items, err := h.catalog.ListUsers(c.Query("role"), parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h Handler) listCategories(c *gin.Context) {
	items, err := h.catalog.ListCategories(c.Query("type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h Handler) listActivities(c *gin.Context) {
	var categoryID *uint
	if raw := c.Query("categoryId"); raw != "" {
		value, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid categoryId"})
			return
		}
		parsed := uint(value)
		categoryID = &parsed
	}

	items, err := h.catalog.ListActivities(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h Handler) createActivity(c *gin.Context) {
	var payload activityRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.catalog.CreateActivity(activity.Activity{
		Title:       payload.Title,
		Description: payload.Description,
		Language:    payload.Language,
		Difficulty:  payload.Difficulty,
		AgeGroup:    payload.AgeGroup,
		CategoryID:  payload.CategoryID,
		Prompt:      payload.Prompt,
		Icon:        payload.Icon,
		HTMLCode:    payload.HTMLCode,
		CSSCode:     payload.CSSCode,
		JSCode:      payload.JSCode,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"item": item})
}

func (h Handler) updateActivity(c *gin.Context) {
	var payload activityRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
		return
	}

	item, err := h.catalog.UpdateActivity(activity.Activity{
		ID:          uint(value),
		Title:       payload.Title,
		Description: payload.Description,
		Language:    payload.Language,
		Difficulty:  payload.Difficulty,
		AgeGroup:    payload.AgeGroup,
		CategoryID:  payload.CategoryID,
		Prompt:      payload.Prompt,
		Icon:        payload.Icon,
		HTMLCode:    payload.HTMLCode,
		CSSCode:     payload.CSSCode,
		JSCode:      payload.JSCode,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h Handler) listAssignments(c *gin.Context) {
	var kidID *uint
	var parentID *uint

	if raw := c.Query("kidId"); raw != "" {
		value, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kidId"})
			return
		}
		parsed := uint(value)
		kidID = &parsed
	}
	if raw := c.Query("parentId"); raw != "" {
		value, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parentId"})
			return
		}
		parsed := uint(value)
		parentID = &parsed
	}

	items, err := h.catalog.ListAssignments(kidID, parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h Handler) createAssignment(c *gin.Context) {
	var payload assignmentRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.catalog.CreateAssignment(payload.ParentID, payload.KidID, payload.ActivityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"item": item})
}

func (h Handler) updateAssignment(c *gin.Context) {
	var payload assignmentUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignment id"})
		return
	}

	item, err := h.catalog.UpdateAssignment(uint(value), payload.Status, payload.Progress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h Handler) parentStats(c *gin.Context) {
	value, err := strconv.ParseUint(c.Param("parentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent id"})
		return
	}

	item, err := h.dashboard.ParentStats(uint(value))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h Handler) listActivityData(c *gin.Context) {
	var activityID *uint
	var assignmentID *uint
	var kidID *uint

	if raw := c.Query("activityId"); raw != "" {
		value, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activityId"})
			return
		}
		parsed := uint(value)
		activityID = &parsed
	}
	if raw := c.Query("assignmentId"); raw != "" {
		value, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignmentId"})
			return
		}
		parsed := uint(value)
		assignmentID = &parsed
	}
	if raw := c.Query("kidId"); raw != "" {
		value, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kidId"})
			return
		}
		parsed := uint(value)
		kidID = &parsed
	}

	items, err := h.runtime.List(activityID, assignmentID, kidID, c.Query("key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h Handler) createActivityData(c *gin.Context) {
	var payload activityDataRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.runtime.Create(activitydata.Item{
		ActivityID:   payload.ActivityID,
		AssignmentID: payload.AssignmentID,
		KidID:        payload.KidID,
		Key:          payload.Key,
		Value:        payload.Value,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"item": item})
}

func (h Handler) updateActivityData(c *gin.Context) {
	var payload activityDataUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity data id"})
		return
	}

	item, err := h.runtime.Update(uint(value), payload.Key, payload.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h Handler) deleteActivityData(c *gin.Context) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity data id"})
		return
	}

	if err := h.runtime.Delete(uint(value)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
