package services

import (
	"bitbucket.org/lyndus/backend/domain/client/contracts"
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/domain/client/repository"
)

type ClientPollChoiceService struct {
	repo contracts.ClientPollChoiceRepositoryInterface
}

var PollChoice contracts.ClientPollChoiceServiceInterface

func NewClientPollChoiceService() *ClientPollChoiceService{
	repo := repository.NewClientPollChoiceRepository()
	return &ClientPollChoiceService{
		repo:repo,
	}
}

func (cpcs ClientPollChoiceService)Create(pollChoice entity.ClientPollChoice)(int , error){
	return cpcs.repo.Create(pollChoice)
}

func (cpcs ClientPollChoiceService)Update(pollChoice entity.ClientPollChoice) error{
	return cpcs.repo.Update(pollChoice)
}


func (cpcs ClientPollChoiceService)Replicate(pollChoice entity.ClientPollChoice)(int , error){
	if pollChoice.ChoiceText == ""{
		pollChoice,err := cpcs.repo.GetById(pollChoice.ID)
		if err != nil {
			return 0,err
		}
		return cpcs.repo.Create(pollChoice)
	}
	return cpcs.repo.Create(pollChoice)
}