package entity

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"github.com/shopspring/decimal"
)

type ClientCategory struct {
	entity.Category
}

func DefaultClientCategory(clientID int) (category ClientCategory) {
	category.ClientID = clientID
	category.Value = decimal.Zero
	category.Name = "Padrão"
	category.Default = true
	return category
}
