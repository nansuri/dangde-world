package user

type Repository interface {
	FindByID(id uint) (User, error)
	List(role string, parentID *uint) ([]User, error)
	CreateBatch(items []User) error
	Count() (int64, error)
}
