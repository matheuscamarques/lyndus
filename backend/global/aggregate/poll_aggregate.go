package aggregate

import "bitbucket.org/lyndus/backend/global/basic"

type PollAgregate struct {
	basic.Poll
	Questions   []PollQuestionAggregate `json:"questions"`
	Category []int `json:"category"`
}
