package user

import (
	"errors"
	"sync"
)

type userRepoInMemoryImpl struct {
	users map[string]User
	mu    sync.RWMutex
}

func NewUserRepoInMemoryImpl() UserRepo {
	return &userRepoInMemoryImpl{
		users: make(map[string]User),
	}
}

func (repo *userRepoInMemoryImpl) InsertOrUpdate(user User) {
	if user.UserInfo.ID == "" {
		return
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.users[user.UserInfo.ID] = user
}

func (repo *userRepoInMemoryImpl) Remove(ID string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if user, ok := repo.users[ID]; !ok {
		return errors.New("user not found")
	} else {
		user.Active = false
		repo.users[ID] = user
		return nil
	}
}

func (repo *userRepoInMemoryImpl) GetByID(ID string) (User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	if user, ok := repo.users[ID]; !ok {
		return user, errors.New("user not found")
	} else {
		return user, nil
	}
}

func (repo *userRepoInMemoryImpl) GetAllActive(pageNumber, pageSize uint) []User {

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	lowerLimit := pageNumber * pageSize
	upperLimit := (pageNumber + 1) * pageSize
	userRepoSize := uint(len(repo.users))

	if lowerLimit >= userRepoSize {
		return []User{}
	}
	if upperLimit > userRepoSize {
		upperLimit = userRepoSize
	}

	users := make([]User, 0, userRepoSize)
	for _, user := range repo.users {
		if user.Active {
			users = append(users, user)
		}
	}

	return users[lowerLimit:upperLimit]
}
