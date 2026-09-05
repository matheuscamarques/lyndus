package entity

type ClientPollQuestion struct {
	ID int `json:"id" db:"id"`
	ClientPollID int `json:"clientPollID" db:"client_poll_id"`
	MultipleChoice bool `json:"multipleChoice" db:"multiple_choice"`
	QuestionText string `json:"questionText" db:"question_text"`
}
// CanCreate Verifica se deve ser criado um
func (cpq ClientPollQuestion)CanCreate()bool{
	return cpq.ID == 0
}

// CanReplicate Verifica se deve ser replicado
func (cpq ClientPollQuestion)CanReplicate() bool {
	return cpq.ID != 0
}

// CanUpdate Verifica se deve ser replicado
func (cpq ClientPollQuestion)CanUpdate() bool {
	return cpq.ID != 0
}