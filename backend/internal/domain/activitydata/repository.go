package activitydata

type Repository interface {
	List(activityID, assignmentID, kidID *uint, key string) ([]Item, error)
	FindByID(id uint) (Item, error)
	Create(item Item) (Item, error)
	Update(item Item) (Item, error)
	Delete(id uint) error
}
