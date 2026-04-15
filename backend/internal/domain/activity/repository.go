package activity

type Repository interface {
	List(categoryID *uint) ([]Activity, error)
	FindByID(id uint) (Activity, error)
	Create(item Activity) (Activity, error)
	Update(item Activity) (Activity, error)
	CreateBatch(items []Activity) error
}
