package assignment

type Repository interface {
	List(kidID, parentID *uint) ([]Assignment, error)
	FindByID(id uint) (Assignment, error)
	Create(item Assignment) (Assignment, error)
	Update(item Assignment) (Assignment, error)
	CreateBatch(items []Assignment) error
	Delete(id uint) error
}
