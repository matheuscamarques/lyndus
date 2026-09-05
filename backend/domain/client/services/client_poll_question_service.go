package services

import (
	"bitbucket.org/lyndus/backend/domain/client/contracts"
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/domain/client/repository"
	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/global/basic"
)

type ClientPollQuestionService struct {
	repo contracts.ClientPollQuestionRepositoryInterface
}


var PollQuestion contracts.ClientPollQuestionServiceInterface

func NewClientPollQuestionService() *ClientPollQuestionService{
	repo := repository.NewClientPollQuestionRepository()
	return &ClientPollQuestionService{
			repo:repo,
	}
}

func (cpqs ClientPollQuestionService)Create(pollQuestion entity.ClientPollQuestion)(int , error){
	return cpqs.repo.Create(pollQuestion)
}

func (cpqs ClientPollQuestionService)Update(pollQuestion entity.ClientPollQuestion) (err error){
	return cpqs.repo.Update(pollQuestion)
}

func (cpqs ClientPollQuestionService)Replicate(pollQuestion entity.ClientPollQuestion)(int , error){
	if pollQuestion.QuestionText == ""{
		pollQuestion,err := cpqs.repo.GetById(pollQuestion.ID)
		if err != nil {
			return 0,err
		}
		return cpqs.repo.Create(pollQuestion)
	}
	return cpqs.repo.Create(pollQuestion)
}


func (cpqs ClientPollQuestionService) ValidatePoll(questionID, pollID int) (bool, error) {
	return cpqs.repo.ValidatePoll(questionID, pollID)
}

func (cpqs ClientPollQuestionService)GetAll(pollID int)(questions []basic.PollQuestion, err error){
	return cpqs.repo.GetAllByIdPoll(pollID)
}

func (cpsq ClientPollQuestionService)GetById(questionID int)(question aggregate.PollQuestionAggregate,err error){
	return cpsq.repo.GetAggreagateById(questionID)
}