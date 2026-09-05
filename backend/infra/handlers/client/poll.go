package client

import (
	"bitbucket.org/lyndus/backend/domain/client"
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/domain/client/services"
	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"bitbucket.org/lyndus/backend/infra/types"
	"encoding/json"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

type PollCreateCase struct {
	StartTime   types.DateTime `json:"startTime"`
	EndTime     types.DateTime `json:"endTime"`
	AllEmployee bool           `json:"allEmployee"`
	Title       string         `json:"title"`
	Questions   []struct {
		entity.ClientPollQuestion
		Choices []entity.ClientPollChoice
	} `json:"questions"`
	Category []int `json:"category"`
}

type PollQuestionCase struct {
	entity.ClientPollQuestion
	Choices []entity.ClientPollChoice `json:"choices"`
}

func GetPollAll(w http.ResponseWriter, r *http.Request) {
	var err error
	id := r.Context().Value("id").(int)

	ClientID, err := services.AuthService.CheckPermission(id, client.POLL, client.VIEW)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}
	if ClientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	itemsPerPag, err := strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	polls, err := services.Poll.GetAll(ClientID, page, itemsPerPag)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(polls, http.StatusOK, w)
}

func GetPollById(w http.ResponseWriter, r *http.Request) {
	var err error
	id := r.Context().Value("id").(int)

	ClientID, err := services.AuthService.CheckPermission(id, client.POLL, client.VIEW)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	if ClientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	pollID, err := strconv.Atoi(chi.URLParam(r, "id"))
	aggregate, err := services.Poll.GetById(pollID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if aggregate.ID == 0 {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(aggregate, http.StatusOK, w)
}

func CreatePoll(w http.ResponseWriter, r *http.Request) {
	var err error
	poll := entity.ClientPoll{}

	id := r.Context().Value("id").(int)

	poll.ClientID, err = services.AuthService.CheckPermission(id, client.POLL, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	if poll.ClientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var pollCase PollCreateCase
	err = json.NewDecoder(r.Body).Decode(&pollCase)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	poll.Title = pollCase.Title
	poll.Status = client.POLLCREATED
	poll.CreatedBy = id
	poll.AllEmployee = pollCase.AllEmployee
	poll.StartTime = pollCase.StartTime
	poll.EndTime = pollCase.EndTime

	poll.ID, err = services.Poll.Create(poll)
	if err != nil || poll.ID == 0 {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	for i := range pollCase.Questions {
		pollQuestion := entity.ClientPollQuestion{
			ClientPollID:   poll.ID,
			MultipleChoice: pollCase.Questions[i].MultipleChoice,
			QuestionText:   pollCase.Questions[i].QuestionText,
		}
		var pollQuestionID int
		if !pollCase.Questions[i].CanCreate() && !pollCase.Questions[i].CanReplicate() {
			config.ResponsePerErr(w, err, config.INVALIDREQUEST)
			return
		}

		if pollCase.Questions[i].CanCreate() {
			pollQuestionID, err = services.PollQuestion.Create(pollQuestion)
		}
		if pollCase.Questions[i].CanReplicate() {
			pollQuestionID, err = services.PollQuestion.Replicate(pollQuestion)
		}

		if err != nil || pollQuestionID == 0 {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}

		for j := range pollCase.Questions[i].Choices {
			if !pollCase.Questions[i].Choices[j].CanCreate() && !pollCase.Questions[i].Choices[j].CanReplicate() {
				config.ResponsePerErr(w, err, config.INVALIDREQUEST)
				return
			}

			questionChoice := entity.ClientPollChoice{
				ClientPollQuestionID: pollQuestionID,
				ChoiceText:           pollCase.Questions[i].Choices[j].ChoiceText,
			}

			var pollChoiceID int
			if pollCase.Questions[i].Choices[j].CanCreate() {
				pollChoiceID, err = services.PollChoice.Create(questionChoice)
			}
			if pollCase.Questions[i].Choices[j].CanReplicate() {
				pollChoiceID, err = services.PollChoice.Replicate(questionChoice)
			}

			if err != nil || pollChoiceID == 0 {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}

		}

	}

	for i := range pollCase.Category {
		category, err := services.CategoryService.GetCategoryIDByID(poll.ClientID, pollCase.Category[i])
		if err != nil || category.ID == 0 {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}

		pollCategoryID, err := services.PollCategory.Create(
			entity.ClientPollCategory{
				ClientPollID:     poll.ID,
				ClientCategoryID: category.ID,
			},
		)
		if err != nil || pollCategoryID == 0 {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}

	}

	config.JSONResponse(response.ID{ID: poll.ID}, http.StatusOK, w)
}

func PutPoll(w http.ResponseWriter, r *http.Request) {
	var err error
	poll := entity.ClientPoll{}

	id := r.Context().Value("id").(int)

	poll.ClientID, err = services.AuthService.CheckPermission(id, client.POLL, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	if poll.ClientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var pollCase aggregate.PollAgregate
	err = json.NewDecoder(r.Body).Decode(&pollCase)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}


	validate , err := services.Poll.ValidateClient( pollCase.ID,poll.ClientID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if !validate {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	poll.Title = pollCase.Title
	poll.Status = pollCase.Status
	poll.AllEmployee = pollCase.AllEmployee
	poll.StartTime = pollCase.StartTime
	poll.EndTime = pollCase.EndTime
	poll.ID = pollCase.ID
	err = services.Poll.Update(poll)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	for i := range pollCase.Questions {
		pollQuestion := entity.ClientPollQuestion{
			ID:             pollCase.Questions[i].ID,
			ClientPollID:   poll.ID,
			MultipleChoice: pollCase.Questions[i].MultipleChoice,
			QuestionText:   pollCase.Questions[i].QuestionText,
		}
		var pollQuestionID int
		if !pollCase.Questions[i].CanCreate() && !pollCase.Questions[i].CanUpdate() {
			config.ResponsePerErr(w, err, config.INVALIDREQUEST)
			return
		}

		if pollCase.Questions[i].CanCreate() {
			pollQuestionID, err = services.PollQuestion.Create(pollQuestion)
			if pollQuestionID == 0 {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
		}
		if pollCase.Questions[i].CanUpdate() {
			pollQuestionID = pollQuestion.ID
			err = services.PollQuestion.Update(pollQuestion)
		}

		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}

		for j := range pollCase.Questions[i].Choices {
			if !pollCase.Questions[i].Choices[j].CanCreate() && !pollCase.Questions[i].Choices[j].CanUpdate() {
				config.ResponsePerErr(w, err, config.INVALIDREQUEST)
				return
			}

			questionChoice := entity.ClientPollChoice{
				ID:                   pollCase.Questions[i].Choices[j].ID,
				ClientPollQuestionID: pollQuestionID,
				ChoiceText:           pollCase.Questions[i].Choices[j].ChoiceText,
			}

			if pollCase.Questions[i].Choices[j].CanCreate() {
				pollChoiceID, err := services.PollChoice.Create(questionChoice)
				if pollChoiceID == 0 {
					config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
					return
				}
			}
			if pollCase.Questions[i].Choices[j].CanUpdate() {
				err = services.PollChoice.Update(questionChoice)
			}

			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}

		}

	}

	for i := range pollCase.Category {
		category, err := services.CategoryService.GetCategoryIDByID(poll.ClientID, pollCase.Category[i])
		if err != nil || category.ID == 0 {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
		pollCategory := entity.ClientPollCategory{
			ClientPollID:     poll.ID,
			ClientCategoryID: category.ID,
		}
		exist, err := services.PollCategory.Exist(pollCategory)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}

		if !exist {
			log.Printf("%+v %+v", pollCategory, !exist)
			pollCategoryID, err := services.PollCategory.Create(pollCategory)
			if err != nil || pollCategoryID == 0 {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
		}

	}

	config.JSONResponse(response.ID{ID: poll.ID}, http.StatusOK, w)
}

func PullAddQuestion(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.POLL, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}
	if clientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}
	var questionCase PollQuestionCase
	err = json.NewDecoder(r.Body).Decode(&questionCase)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	questionCase.ClientPollID, err = strconv.Atoi(chi.URLParam(r, "id"))
	validate , err := services.Poll.ValidateClient(questionCase.ClientPollID,questionCase.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if !validate {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	questionCase.ID, err = services.PollQuestion.Create(questionCase.ClientPollQuestion)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	for i := range questionCase.Choices {
		questionCase.Choices[i].ClientPollQuestionID = questionCase.ID
		choiceID, err := services.PollChoice.Create(questionCase.Choices[i])
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
		if choiceID == 0 {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
	}

	config.JSONResponse(response.ID{questionCase.ID}, http.StatusOK, w)
}

func GetAllPullQuestions(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	ClientID, err := services.AuthService.CheckPermission(id, client.POLL, client.VIEW)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	if ClientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	pollID, err := strconv.Atoi(chi.URLParam(r, "id"))
	validate , err := services.Poll.ValidateClient(pollID,ClientID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if !validate {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	questions, err := services.PollQuestion.GetAll(pollID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(questions, http.StatusOK, w)
}

func GetPullQuestion(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	ClientID, err := services.AuthService.CheckPermission(id, client.POLL, client.VIEW)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	if ClientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	pollID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	validate , err := services.Poll.ValidateClient(pollID,ClientID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if !validate {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	questionID, err := strconv.Atoi(chi.URLParam(r, "questionID"))

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	validate,err = services.PollQuestion.ValidatePoll(questionID,pollID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if !validate {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	question, err := services.PollQuestion.GetById(questionID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(question, http.StatusOK, w)
}

func PutPullQuestion(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	ClientID, err := services.AuthService.CheckPermission(id, client.POLL, client.VIEW)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}
	if ClientID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	pollID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	validate , err := services.Poll.ValidateClient(pollID,ClientID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if !validate {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var question PollQuestionCase
	err = json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}


	err = services.PollQuestion.Update(question.ClientPollQuestion)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	for i := range question.Choices {
		question.Choices[i].ClientPollQuestionID = question.ID
		if !question.Choices[i].CanCreate() && !question.Choices[i].CanUpdate() {
			config.ResponsePerErr(w, err, config.INVALIDREQUEST)
			return
		}

		if question.Choices[i].CanCreate() {
			choiceID, err := services.PollChoice.Create(question.Choices[i])
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}

			if choiceID == 0 {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}

		}

		if question.Choices[i].CanUpdate() {
			err = services.PollChoice.Update(question.Choices[i])
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
		}

	}
	config.JSONResponse(response.ID{ID: question.ID}, http.StatusOK, w)
}
