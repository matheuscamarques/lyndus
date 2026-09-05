package contracts

import "bitbucket.org/lyndus/backend/domain/client/contracts"

type BonderClientUserServiceInterface interface {
	contracts.ClientUserServiceInterface
	GetMasterUserAuthID(id int) (int, error)

}
