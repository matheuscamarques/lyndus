package criteria

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
)

type ClientCategoryResponse struct {
	CResponse
	Items []entity.Category `json:"items"`
}
