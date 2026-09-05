package criteria

import "bitbucket.org/lyndus/backend/domain/bs/entity"

type CProductResponse struct{
	CResponse
	Items      []entity.Product `json:"items"`
}