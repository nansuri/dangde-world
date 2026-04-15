package category

type Repository interface {
	List(categoryType string) ([]Category, error)
	CreateBatch(items []Category) error
}
