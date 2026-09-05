package basic

type PollChoice struct {
	ID                   int    `json:"id" db:"id"`
	ChoiceText           string `json:"choiceText" db:"choice_text"`
}
// CanCreate Verifica se deve ser criado um
func (cpc PollChoice)CanCreate()bool{
	return cpc.ID == 0
}

// CanReplicate Verifica se deve ser replicado
func (cpc PollChoice)CanReplicate() bool {
	return cpc.ID != 0
}

// CanUpdate Verifica se deve ser replicado
func (cpc PollChoice)CanUpdate() bool {
	return cpc.ID != 0
}