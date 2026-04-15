package runtime

import "dangde-world/backend/internal/domain/activitydata"

type Service struct {
	items activitydata.Repository
}

func NewService(items activitydata.Repository) Service {
	return Service{items: items}
}

func (s Service) List(activityID, assignmentID, kidID *uint, key string) ([]activitydata.Item, error) {
	return s.items.List(activityID, assignmentID, kidID, key)
}

func (s Service) Create(item activitydata.Item) (activitydata.Item, error) {
	return s.items.Create(item)
}

func (s Service) Update(id uint, key string, value string) (activitydata.Item, error) {
	item, err := s.items.FindByID(id)
	if err != nil {
		return activitydata.Item{}, err
	}
	if key != "" {
		item.Key = key
	}
	item.Value = value
	return s.items.Update(item)
}

func (s Service) Delete(id uint) error {
	return s.items.Delete(id)
}
