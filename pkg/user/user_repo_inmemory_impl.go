package user

import (
	"errors"
	"sync"

	"github.com/zmb3/spotify"
)

type userRepoInMemoryImpl struct {
	users map[spotify.ID]User
	mu    sync.RWMutex
}

func NewUserRepoInMemoryImpl() UserRepo {
	return &userRepoInMemoryImpl{
		users: make(map[spotify.ID]User),
	}
}

func (repo *userRepoInMemoryImpl) InsertOrUpdate(user User) {
	if user.ID == "" {
		return
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.users[user.ID] = user
}

func (repo *userRepoInMemoryImpl) Remove(ID spotify.ID) error {
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

func (repo *userRepoInMemoryImpl) GetByID(ID spotify.ID) (User, error) {
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
