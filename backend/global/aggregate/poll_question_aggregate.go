package aggregate

import "bitbucket.org/lyndus/backend/global/basic"

type PollQuestionAggregate struct {
	basic.PollQuestion
	Choices []basic.PollChoice `json:"choices"`
}
