package group

import (
	"errors"
	"sync"

	"github.com/zmb3/spotify"
)

type groupRepoInMemoryImpl struct {
	groups map[spotify.ID]Group
	mu     sync.RWMutex
}

func NewGroupRepoInMemoryImpl() GroupRepo {
	return &groupRepoInMemoryImpl{
		groups: make(map[spotify.ID]Group),
	}
}

func (repo *groupRepoInMemoryImpl) InsertOrUpdate(group Group) {
	if group.ID == "" {
		return
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.groups[group.ID] = group
}

func (repo *groupRepoInMemoryImpl) Remove(ID spotify.ID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if group, ok := repo.groups[ID]; !ok {
		return errors.New("group not found")
	} else {
		group.Active = false
		repo.groups[ID] = group
		return nil
	}
}

func (repo *groupRepoInMemoryImpl) GetByID(ID spotify.ID) (Group, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	if group, ok := repo.groups[ID]; !ok {
		return group, errors.New("group not found")
	} else {
		return group, nil
	}
}

func (repo *groupRepoInMemoryImpl) GetAllActive(pageNumber, pageSize uint) []Group {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	lowerLimit := pageNumber * pageSize
	upperLimit := (pageNumber + 1) * pageSize
	groupRepoSize := uint(len(repo.groups))

	if lowerLimit >= groupRepoSize {
		return []Group{}
	}
	if upperLimit > groupRepoSize {
		upperLimit = groupRepoSize
	}

	groups := make([]Group, 0, groupRepoSize)
	for _, group := range repo.groups {
		if group.Active {
			groups = append(groups, group)
		}
	}

	return groups[lowerLimit:upperLimit]
}

func (repo *groupRepoInMemoryImpl) GetAllActiveByUserID(userID spotify.ID, pageNumber uint, pageSize uint) []Group {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	lowerLimit := pageNumber * pageSize
	upperLimit := (pageNumber + 1) * pageSize
	groupRepoSize := uint(len(repo.groups))

	if lowerLimit >= groupRepoSize {
		return []Group{}
	}
	if upperLimit > groupRepoSize {
		upperLimit = groupRepoSize
	}

	groups := make([]Group, 0, groupRepoSize)
	for _, group := range repo.groups {
		_, existsUser := group.Users[userID]
		if group.Active && existsUser {
			groups = append(groups, group)
		}
	}

	return groups[lowerLimit:upperLimit]
}
