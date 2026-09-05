package contracts

import "bitbucket.org/lyndus/backend/domain/bs/contracts"

type BonderBSUserServiceInterface interface {
	contracts.BSUserServiceInterface
	GetMasterUserAuthID(id int) (int, error)
}
