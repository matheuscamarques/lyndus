package contracts

import "bitbucket.org/lyndus/backend/domain/client/entity"

type ClientCompanyRepositoryInterface interface {
	GetCompanyByID(targetID int) (entity.Company, error)
	UpdateCompany(company entity.Company) error
}
