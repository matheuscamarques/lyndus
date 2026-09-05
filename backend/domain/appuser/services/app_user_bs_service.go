package services

import (
	"bitbucket.org/lyndus/backend/global/aggregate"
	"bitbucket.org/lyndus/backend/infra/logger"
	"go.uber.org/zap"
	"strconv"
	"time"

	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/repository"
	"bitbucket.org/lyndus/backend/infra/criteria"
)

type AppUserBSService struct {
	repo contracts.AppUserBSRepositoryInterface
}

var BSService contracts.AppUserBSServiceInterface

func NewAppUserBSService() *AppUserBSService {
	return &AppUserBSService{
		repo: repository.NewAppUserBSRepository(),
	}
}
func (abs AppUserBSService) GetBSWeekDays(bsID int) (weekDays []entity.WeekDayBS, err error) {
	return abs.repo.GetBSWeekDays(bsID)
}
func (abs AppUserBSService) GetBSID(bsId int) (id int, err error) {
	return abs.repo.GetBSID(bsId)
}

func (abs AppUserBSService) GetBSName(bsId int) (name string, err error) {
	return abs.repo.GetBSName(bsId)
}

func (abs AppUserBSService) GetBS(bsId int) (bs entity.BS, err error) {
	bs, err = abs.repo.GetBS(bsId)
	if err != nil {
		return bs, err
	}
	bs.WeekDays, err = abs.repo.GetBsWeekDays(bs.ID)
	return bs, err
}

func (abs AppUserBSService) GetBSBMName(bsId, bmID int) (name string, err error) {
	return abs.repo.GetBSBMName(bsId, bmID)
}

func (abs AppUserBSService) GetBSDaysOff(bsID int, start, end time.Time) (days []entity.DayOff, err error) {
	return abs.repo.GetBSDaysOff(bsID, start, end)
}

func (abs AppUserBSService) GetBsIDServiceIDWeekDays(bsID, serviceId int) (weekDays []entity.WeekDayBS, err error) {
	return abs.repo.GetBsIDServiceIDWeekDays(bsID, serviceId)
}

func (abs AppUserBSService) GetLocation(bsID int) (location string, err error) {
	return abs.repo.GetLocation(bsID)
}

func (abs AppUserBSService) GetBSsMap(lat, lon float64) (bss []entity.BS, err error) {
	bss, err = abs.repo.GetBSsMap(lat, lon)
	if err != nil {
		return bss, err
	}
	for i := range bss {
		bss[i].WeekDays, err = abs.repo.GetBsWeekDays(bss[i].ID)
		if err != nil {
			//todo melhorar esse log de  erro abaixo
			logger.Error("error ", zap.String("error", err.Error()))
		}
	}
	return bss, err
}
func (abs AppUserBSService) GetBSsByCategoryMap(categoryID int, lat, lon float64) (bss []entity.BS, err error) {

	bss, err = abs.repo.GetBSsByCategoryMap(categoryID, lat, lon)
	if err != nil {
		return bss, err
	}
	for i := range bss {
		bss[i].WeekDays, err = abs.repo.GetBsWeekDays(bss[i].ID)
		if err != nil {
			//todo melhorar esse log de  erro abaixo
			logger.Error("error ", zap.String("error", err.Error()))
		}
	}
	return bss, err
}

func (abs AppUserBSService) GetBsWeekDays(bsID int) (weekDays []aggregate.WeekDay, err error) {
	return abs.repo.GetBsWeekDays(bsID)
}

func (abs AppUserBSService) GetBSs(page, itemsPerPage int, fantasyName string, lat, lon float64) (response criteria.BSResponse, err error) {
	c := criteria.Criteria{ActivePage: page, ItemsPerPage: itemsPerPage}

	response.ActivePage = c.ActivePage
	response.Items, response.TotalPages, err = abs.repo.GetBSs(fantasyName, lat, lon, c)
	if err != nil {
		return response, err
	}

	for i := range response.Items {
		response.Items[i].WeekDays, err = abs.repo.GetBsWeekDays(response.Items[i].ID)
		if err != nil {
			//todo melhorar esse log de  erro abaixo
			logger.Error("error ", zap.String("error", err.Error()))
		}
	}

	for i := range response.Items {
		response.Items[i].Categories, err = abs.repo.GetBSCategories(response.Items[i].ID)
		if err != nil {
			//todo melhorar esse log de  erro abaixo
			logger.Error("error ", zap.String("error", err.Error()))
		}
	}

	return response, err
}

func (abs AppUserBSService) GetBSsByCategory(categoryID, page, itemsPerPage int, fantasyName string, lat, lon float64) (response criteria.BSResponse, err error) {

	c := criteria.Criteria{ActivePage: page, ItemsPerPage: itemsPerPage}

	response.ActivePage = c.ActivePage
	response.Items, response.TotalPages, err = abs.repo.GetBSsByCategory(categoryID, fantasyName, lat, lon, c)
	if err != nil {
		return response, err
	}

	return response, err
}

func (abs AppUserBSService) GetBSsWithLocalization(lat, lon float64, page, itemsPerPage int) (response criteria.BSResponse, err error) {
	c := criteria.Criteria{ActivePage: page, ItemsPerPage: itemsPerPage}
	total, err := c.TotalPagesWithWhere(abs.repo, `INNER JOIN company ON bs.company_id = company.id`)
	if err != nil {
		return response, err
	}

	response.ActivePage = c.ActivePage
	response.TotalPages = total
	response.Items, err = abs.repo.GetBSsWithLocalization(c, lat, lon)

	if err != nil {
		return response, err
	}

	return response, err
}

func (abs AppUserBSService) GetBSsByCategoryLocalization(categoryID int, lat, lon float64, page, itemsPerPage int) (response criteria.BSResponse, err error) {
	c := criteria.Criteria{ActivePage: page, ItemsPerPage: itemsPerPage}
	total, err := c.TotalPagesWithWhere(abs.repo, `INNER JOIN company com ON bs.company_id = com.id
					INNER JOIN bs_bs_category bbc ON bs.id = bbc.bs_id
					WHERE bbc.bs_category_id=`+strconv.Itoa(categoryID))
	if err != nil {
		return response, err
	}

	response.ActivePage = c.ActivePage
	response.TotalPages = total
	response.Items, err = abs.repo.GetBSsByCategoryWithLocalization(categoryID, c, lat, lon)

	if err != nil {
		return response, err
	}

	return response, err
}
