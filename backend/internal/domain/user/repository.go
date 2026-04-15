package user

type Repository interface {
	FindByID(id uint) (User, error)
	List(role string, parentID *uint) ([]User, error)
	Create(item User) (User, error)
	Update(item User) (User, error)
	Delete(id uint) error
	CreateBatch(items []User) error
	Count() (int64, error)
}
