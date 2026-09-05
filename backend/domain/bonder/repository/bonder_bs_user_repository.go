package repository

import "bitbucket.org/lyndus/backend/domain/bs/repository"

type BonderBSUserRepository struct{
	repository.BSUserRepository
}



func NewBSUserRepository() (repo BonderBSUserRepository){
	repo.UsePostgres()
	return repo
}




