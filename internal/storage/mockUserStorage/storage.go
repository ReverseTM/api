package mockUserStorage

import (
	"api/internal/domain/models"
	"api/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage struct {
	id    int64
	users map[string]models.User
}

func NewUserStorage() *UserStorage {
	var id int64 = 1
	users := make(map[string]models.User)

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("presale"), bcrypt.DefaultCost)

	user := models.User{
		ID:           id,
		Username:     "admin",
		PasswordHash: passwordHash,
	}
	users["admin"] = user

	id++

	return &UserStorage{
		id:    id,
		users: users,
	}
}

func (s *UserStorage) User(username string) (models.User, error) {
	if user, ok := s.users[username]; ok {
		return user, nil
	}

	return models.User{}, errors.ErrUserNotFound
}

func (s *UserStorage) SaveUser(username string, password []byte) (int64, error) {
	user := models.User{
		ID:           s.id,
		Username:     username,
		PasswordHash: password,
	}

	if _, ok := s.users[username]; ok {
		return 0, errors.ErrUserExists
	}

	s.users[username] = user

	var id = s.id
	s.id++

	return id, nil
}
