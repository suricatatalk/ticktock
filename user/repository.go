package user

import (
	"github.com/sohlich/ticktock/config"
)

type UserRepository interface {
	FindById(id string) (*User, error)
	Save(t *User) error
}

var Repository UserRepository

type MgoUserRepository struct {
	config.MgoRepository
}

func InitializeRepository(db *config.Database) error {
	var err error
	Repository = &MgoUserRepository{
		config.MgoRepository{
			db.Database().C("users"),
		},
	}
	return err
}

func (r *MgoUserRepository) FindById(id string) (*User, error) {
	u := &User{}
	err := r.FindId(id).One(u)
	return u, err
}

func (r *MgoUserRepository) Save(u *User) error {
	_, err := r.UpsertId(u.ID, u)
	return err
}
