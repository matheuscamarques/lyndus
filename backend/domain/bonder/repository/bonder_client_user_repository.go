package repository

import (
	"bitbucket.org/lyndus/backend/domain/client/repository"
)

type BonderClientUserRepository struct {
	repository.ClientUserRepository
}


func NewBonderClientUserRepository() (repo BonderClientUserRepository){
	repo.UsePostgres()
	return repo
}

