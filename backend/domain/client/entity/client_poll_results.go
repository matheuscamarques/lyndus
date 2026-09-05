package entity

type ClientPollResults struct {
	ID                   int
	ClientPollQuestionID int
	ClientPollChoiceID   int
	ClientPollEmployeeID int
	Responded            int
}
