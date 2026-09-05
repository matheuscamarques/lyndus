package services

import (
	"bitbucket.org/lyndus/backend/domain/appuser/contracts"
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/repository"
)

type AppUserBSBannerService struct {
	repo contracts.AppUserBSBannerRepositoryInterface
}

var BSBannerService contracts.AppUserBSBannerServiceInterface

func NewAppUserBSBannerService() *AppUserBSBannerService {
	return &AppUserBSBannerService{
		repo: repository.NewAppUserBSBannerRepository(),
	}
}


func (aub AppUserBSBannerService) GetBSBanners(bsID int) (banners []entity.BSBanner, err error) {
	return aub.repo.GetBSBanners(bsID)
}
