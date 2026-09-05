package contracts

import (
	"bitbucket.org/lyndus/backend/domain/client/entity"
)

type ClientBenefitServiceInterface interface {
	CreateBenefitUserCategory(buc *entity.BenefitUserCategory) (int, error)
	CreateBenefit(b *entity.Benefit) error
	CreateBenefitUser(bu entity.BenefitUser) (int, error)
	UpdateBenefitValue(b *entity.Benefit) error
	UpdateBenefitUserValue(bu *entity.BenefitUser) error
	UpdateBenefitUserAdditional(bu *entity.BenefitUser) error
	GetBenefitStatus(b *entity.Benefit) (status int, err error)
	GetBenefit(b entity.Benefit) (benefitResp entity.BenefitRespFull, err error)
	UpdateBenefitDesc(b entity.Benefit) error
	UpdateBenefitStatus(b entity.Benefit) error
	GetBenefitByStatus(b entity.Benefit) error
	GetBenefits(clientID int) (benefitsResp []entity.BenefitResp, err error)
	GetBenefitUser(b entity.Benefit) (benefitUser entity.BenefitUserResp, err error)
	GetBenefitsUsers(b entity.Benefit) (benefitsUser []entity.BenefitUserResp, err error)
	GetBenefitsUserCategory(bu entity.BenefitUserResp) (categories []string, err error)
	GetBenefitsUsersValues(bu entity.BenefitUser) (benefitsUser []entity.BenefitUser, err error)
	GetBenefitsV2(activePage, itemsPerPage, clientID, statusID int, search, orderBy string, sortDesc bool) ([]entity.BenefitResp, int, int, error)
}
