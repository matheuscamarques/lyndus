package contracts

import "bitbucket.org/lyndus/backend/domain/appuser/entity"

type AppUserCreateServiceInterface interface {
	CreatePerson(auc entity.AppUserCreate) (id int, err error)
	UpdatePerson(auc entity.AppUserCreate) error
	GetPersonIDByCPF(a *entity.AppUserCreate) (n int, err error)
	GetAppUserIDByCPF(a *entity.AppUserCreate) (n int, err error)
}
