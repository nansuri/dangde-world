package assignment

import (
	"dangde-world/backend/internal/domain/activity"
	"dangde-world/backend/internal/domain/user"
)

type Assignment struct {
	ID         uint              `json:"id"`
	ParentID   uint              `json:"parentId"`
	Parent     user.User         `json:"parent"`
	KidID      uint              `json:"kidId"`
	Kid        user.User         `json:"kid"`
	ActivityID uint              `json:"activityId"`
	Activity   activity.Activity `json:"activity"`
	Status     string            `json:"status"`
	Progress   int               `json:"progress"`
}
