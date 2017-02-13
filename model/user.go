package model

type UserRepository interface {
	FindById(id string) (*Task, error)
	Save(t *Task) error
}

type User struct {
	ID        string
	Firstname string
	Lastname  string
	Accounts  []string
}
