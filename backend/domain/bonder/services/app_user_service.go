package services

import (
	"bitbucket.org/lyndus/backend/domain/bonder/contracts"
	"bitbucket.org/lyndus/backend/domain/bonder/entity"
	"bitbucket.org/lyndus/backend/domain/bonder/repository"
	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"bitbucket.org/lyndus/backend/infra/types"
	"fmt"
)

type AppUserService struct {
	repo contracts.AppUserRepositoryInterface
}

var AppUser contracts.AppUserServiceInterface

func NewAppUserService() *AppUserService {
	repo := repository.NewAppUserRepository()
	return &AppUserService{&repo}
}

func (aps AppUserService) GetAppUserByCPF(cpf types.CPF, page int, pag int) (response criteria.AppUserResponse, err error) {
	criteria := criteria.Criteria{
		ActivePage:   page,
		ItemsPerPage: pag,
	}

	response.TotalPages, err = criteria.TotalPagesWithWhere(aps.repo, `WHERE cpf = '`+string(cpf)+`' `)

	if err != nil {
		return response, err
	}

	fmt.Println(criteria)
	response.Items, err = aps.repo.GetAppUserByCPF(cpf, criteria)
	if err != nil {
		return response, err
	}
	if response.Items == nil {
		response.Items = make([]aggregate.AppUserAggreate, 0)
	}
	return response, err
}

func (aps AppUserService) GetAppUserByName(name string, page int, pag int) (response criteria.AppUserResponse, err error) {
	criteria := criteria.Criteria{
		ActivePage:   page,
		ItemsPerPage: pag,
	}
	response.TotalPages, err = criteria.TotalPagesWithWhere(aps.repo, `INNER JOIN bs_person ON app_user.person_id = bs_person.id WHERE bs_person.name LIKE '%`+name+`%'`)

	if err != nil {
		return response, err
	}

	response.Items, err = aps.repo.GetAppUserByName(name, criteria)

	if err != nil {
		return response, err
	}
	if response.Items == nil {
		response.Items = make([]aggregate.AppUserAggreate, 0)
	}
	return response, err
}

func (aps AppUserService) GetAppUsers(activePage, itemsPerPage int, search, orderBy string, sortDesc, active bool) (appUsers []entity.AppUserPerson, totalItems, totalPages int, err error) {
	if activePage == 0 {
		activePage = 1
	}

	if itemsPerPage == 0 {
		itemsPerPage = 10
	} else if itemsPerPage > 100 {
		itemsPerPage = 100
	}

	appUsers, totalItems, totalPages, err = aps.repo.GetAppUsers(activePage, itemsPerPage, search, orderBy, sortDesc, active)
	if appUsers == nil {
		appUsers = make([]entity.AppUserPerson, 0)
	}

	return appUsers, totalItems, totalPages, err
}

func (aps AppUserService) GetAppUserByID(appUserID int) (appUsers entity.AppUser, err error) {

	return aps.repo.GetAppUserByID(appUserID)
}

func (aps AppUserService) UpdateLyndusBox(bonus entity.Bonus) error {
	return aps.repo.UpdateLyndusBox(bonus)
}

//
//func (aps AppUserService) UpdateBonus(bonus entity.Bonus) error {
//	if bonus.ID == 0 {
//		return errors.New("dados inconsistentes")
//	}
//
//	dbBonus, err := aps.repo.GetBalance(bonus.ID)
//
//	if err == sql.ErrNoRows {
//		dbBonus.ID, err = aps.repo.InsertBalance(bonus)
//		if err != nil {
//			return err
//		}
//	} else if err != nil {
//		return err
//	} else {
//
//		dbBonus.Balance = dbBonus.Balance + bonus.Balance
//
//		err = aps.repo.UpdateBalance(bonus.ID, dbBonus.Balance)
//		if err != nil {
//			return err
//		}
//	}
//
//	_, err = aps.repo.InsertBalanceHistory(bonus.ID, dbBonus.Balance, bonus.Balance)
//	if err != nil {
//		return err
//	}
//
//	//todo verificar datetime.
//	statements := entity.AppUserStatement{
//		AppUserID:        bonus.ID,
//		AppUserBalanceID: dbBonus.ID,
//		StatementID:      appuser.StatementBonus,
//		Value:            bonus.Balance,
//		Desc:             appuser.Statements[appuser.StatementBonus],
//	}
//	_, err = aps.repo.InsertStatement(statements)
//
//	return err
//}

func (aps AppUserService) GetAllAppUser(activePage, itemsPerPage int, search, orderBy string, sortDesc, active bool) (appUsers []entity.AppUserPerson, totalItems, totalPages int, err error) {
	if activePage == 0 {
		activePage = 1
	}

	if itemsPerPage == 0 {
		itemsPerPage = 10
	} else if itemsPerPage > 100 {
		itemsPerPage = 100
	}

	appUsers, totalItems, totalPages, err = aps.repo.GetAppUsers(
		activePage,
		itemsPerPage,
		search,
		orderBy,
		sortDesc,
		active)
	if appUsers == nil {
		appUsers = make([]entity.AppUserPerson, 0)
	}

	return appUsers, totalItems, totalPages, err
}
