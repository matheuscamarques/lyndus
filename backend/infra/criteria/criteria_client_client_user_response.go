package criteria

import "bitbucket.org/lyndus/backend/domain/client/entity"

type ClientClientUserResponse struct {
	CResponse
	Items []entity.ClientUser `json:"items"`
}
