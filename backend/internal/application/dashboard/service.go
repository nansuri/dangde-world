package dashboard

import (
	"dangde-world/backend/internal/domain/assignment"
	"dangde-world/backend/internal/domain/user"
)

type ParentStats struct {
	KidsCount        int `json:"kidsCount"`
	AssignmentsCount int `json:"assignmentsCount"`
	CompletedCount   int `json:"completedCount"`
	InProgressCount  int `json:"inProgressCount"`
	AverageProgress  int `json:"averageProgress"`
}

type Service struct {
	users       user.Repository
	assignments assignment.Repository
}

func NewService(users user.Repository, assignments assignment.Repository) Service {
	return Service{users: users, assignments: assignments}
}

func (s Service) ParentStats(parentID uint) (ParentStats, error) {
	kids, err := s.users.List(user.RoleKid, &parentID)
	if err != nil {
		return ParentStats{}, err
	}

	items, err := s.assignments.List(nil, &parentID)
	if err != nil {
		return ParentStats{}, err
	}

	stats := ParentStats{
		KidsCount:        len(kids),
		AssignmentsCount: len(items),
	}

	totalProgress := 0
	for _, item := range items {
		totalProgress += item.Progress
		switch item.Status {
		case "completed":
			stats.CompletedCount++
		case "in_progress":
			stats.InProgressCount++
		}
	}

	if len(items) > 0 {
		stats.AverageProgress = totalProgress / len(items)
	}

	return stats, nil
}
