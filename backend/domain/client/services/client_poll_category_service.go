package services

import (
	"bitbucket.org/lyndus/backend/domain/client/contracts"
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/domain/client/repository"
)

type ClientPollCategoryService struct {
	repo contracts.ClientPollCategoryRepositoryInterface
}

var PollCategory contracts.ClientPollCategoryServiceInterface

func NewClientPollCategoryService() * ClientPollCategoryService{
	repo := repository.NewClientPollCategoryRepository()
	return &ClientPollCategoryService{
		repo : repo,
	}
}

func (cpcs ClientPollCategoryService)Create(pollCategory entity.ClientPollCategory) (id int, err error){
	return cpcs.repo.Create(pollCategory)
}


func (cpcs ClientPollCategoryService)Exist(pollCategory entity.ClientPollCategory) (exist bool, err error){
	return cpcs.repo.Exist(pollCategory)
}

