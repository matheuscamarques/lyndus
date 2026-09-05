package criteria

import "bitbucket.org/lyndus/backend/global/basic"

type CPollQuestionResponse struct {
	CResponse
	Items []basic.PollQuestion `json:"items"`
}
