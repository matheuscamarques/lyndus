package entity

type ClientPollChoice struct {
	ID                   int    `json:"id"`
	ClientPollQuestionID int    `json:"clientPollQuestionID"`
	ChoiceText           string `json:"choiceText"`
}

// CanCreate Verifica se deve ser criado um
func (cpc ClientPollChoice)CanCreate()bool{
	return cpc.ID == 0 && len(cpc.ChoiceText) >= 1
}

// CanReplicate Verifica se deve ser replicado
func (cpc ClientPollChoice)CanReplicate() bool {
	return cpc.ID != 0 && cpc.ClientPollQuestionID != 0
}

// CanUpdate Verifica se deve ser replicado
func (cpc ClientPollChoice)CanUpdate() bool {
	return cpc.ID != 0 && cpc.ClientPollQuestionID != 0
}