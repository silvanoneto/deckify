package user

type UserRepo interface {
	InsertOrUpdate(User)
	Remove(string) error
	GetByID(string) (User, error)
	GetAllActive(uint, uint) []User
}
