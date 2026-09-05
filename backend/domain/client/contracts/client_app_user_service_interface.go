package contracts

import "bitbucket.org/lyndus/backend/domain/client/entity"

type ClientAppUserServiceInterface interface {
	CreateAppUserBasic(user *entity.AppUser) error
	GetAppUserIDByCPF(user *entity.AppUser) error
}
