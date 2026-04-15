package category

type Repository interface {
	List(categoryType string) ([]Category, error)
	ListAll() ([]Category, error)
	Create(item Category) (Category, error)
	Update(item Category) (Category, error)
	Delete(id uint) error
	CreateBatch(items []Category) error
}
