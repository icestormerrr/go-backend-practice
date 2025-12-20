package repos

import (
	"errors"

	"github.com/icestormerrr/pz10-auth/internal/core"
	"golang.org/x/crypto/bcrypt"
)

type UserRecord struct {
	ID    int64
	Email string
	Role  string
	Hash  []byte
}

type UserInMemoryRepo struct{ users []UserRecord }

func NewUserInMemoryRepo() *UserInMemoryRepo {
	hash := func(s string) []byte { h, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost); return h }
	return &UserInMemoryRepo{
		users: []UserRecord{
			{ID: 1, Email: "admin@example.com", Role: "admin", Hash: hash("secret123")},
			{ID: 2, Email: "user@example.com", Role: "user", Hash: hash("secret123")},
		},
	}
}

var ErrNotFound = errors.New("user not found")
var ErrBadCreds = errors.New("bad credentials")

func (r *UserInMemoryRepo) GetById(id int64) (core.User, error) {
	for _, u := range r.users {
		if u.ID == id {
			return core.User{ID: u.ID, Email: u.Email, Role: u.Role}, nil
		}
	}

	return core.User{}, ErrNotFound
}

// В этом методе можем вернуть UserRecord, т.к. он приватный и нужен для внутреннего пользования репозитория
func (r *UserInMemoryRepo) getByEmail(email string) (UserRecord, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return UserRecord{}, ErrNotFound
}

func (r *UserInMemoryRepo) CheckPassword(email, pass string) (core.User, error) {
	u, err := r.getByEmail(email)
	if err != nil {
		return core.User{}, ErrNotFound
	}
	if bcrypt.CompareHashAndPassword(u.Hash, []byte(pass)) != nil {
		return core.User{}, ErrBadCreds
	}
	return core.User{ID: u.ID, Email: u.Email, Role: u.Role}, nil
}
