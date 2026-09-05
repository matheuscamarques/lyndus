package contracts

import "bitbucket.org/lyndus/backend/domain/client/entity"

type  ClientAppUserRepositoryInterface interface {
	CreateAppUserBasic(user *entity.AppUser) error
	GetAppUserIDByCPF(user *entity.AppUser)error
}
