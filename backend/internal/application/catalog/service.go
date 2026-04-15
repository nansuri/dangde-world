package catalog

import (
	"dangde-world/backend/internal/domain/activity"
	"dangde-world/backend/internal/domain/assignment"
	"dangde-world/backend/internal/domain/category"
	"dangde-world/backend/internal/domain/user"
)

type Service struct {
	users       user.Repository
	categories  category.Repository
	activities  activity.Repository
	assignments assignment.Repository
}

func NewService(
	users user.Repository,
	categories category.Repository,
	activities activity.Repository,
	assignments assignment.Repository,
) Service {
	return Service{
		users:       users,
		categories:  categories,
		activities:  activities,
		assignments: assignments,
	}
}

func (s Service) ListUsers(role string, parentID *uint) ([]user.User, error) {
	return s.users.List(role, parentID)
}

func (s Service) ListCategories(categoryType string) ([]category.Category, error) {
	return s.categories.List(categoryType)
}

func (s Service) ListActivities(categoryID *uint) ([]activity.Activity, error) {
	return s.activities.List(categoryID)
}

func (s Service) CreateActivity(item activity.Activity) (activity.Activity, error) {
	return s.activities.Create(item)
}

func (s Service) UpdateActivity(item activity.Activity) (activity.Activity, error) {
	return s.activities.Update(item)
}

func (s Service) ListAssignments(kidID, parentID *uint) ([]assignment.Assignment, error) {
	return s.assignments.List(kidID, parentID)
}

func (s Service) CreateAssignment(parentID, kidID, activityID uint) (assignment.Assignment, error) {
	item := assignment.Assignment{
		ParentID:   parentID,
		KidID:      kidID,
		ActivityID: activityID,
		Status:     "assigned",
		Progress:   0,
	}
	return s.assignments.Create(item)
}

func (s Service) UpdateAssignment(id uint, status string, progress *int) (assignment.Assignment, error) {
	item, err := s.assignments.FindByID(id)
	if err != nil {
		return assignment.Assignment{}, err
	}
	if status != "" {
		item.Status = status
	}
	if progress != nil {
		item.Progress = *progress
	}
	return s.assignments.Update(item)
}
