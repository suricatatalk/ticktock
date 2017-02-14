package user

type UserRepository interface {
	FindById(id string) (*User, error)
	Save(t *User) error
}

type User struct {
	ID        string
	Firstname string
	Lastname  string
	Accounts  []string
}
