package basic

type PollQuestion struct {
	ID             int    `json:"id" db:"id"`
	MultipleChoice bool   `json:"multipleChoice" db:"multiple_choice"`
	QuestionText   string `json:"questionText" db:"question_text"`
}


// CanCreate Verifica se deve ser criado um
func (pq PollQuestion)CanCreate()bool{
	return pq.ID == 0
}

// CanReplicate Verifica se deve ser replicado
func (pq PollQuestion)CanReplicate() bool {
	return pq.ID != 0
}

// CanUpdate Verifica se deve ser replicado
func (pq PollQuestion)CanUpdate() bool {
	return pq.ID != 0
}