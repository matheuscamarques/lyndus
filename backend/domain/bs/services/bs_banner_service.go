package services

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
)

type BSBannerService struct {
	repo contracts.BSBannerRepositoryInterface
}
var BannerService contracts.BSBannerServiceInterface


func NewBSBannerService() *BSBannerService {
	return &BSBannerService{
		repo: repository.NewBSBannerRepository(),
	}
}

func (bbs BSBannerService) CreateBanner(bsID int, name string) (id int, err error) {
	return bbs.repo.CreateBanner(bsID, name)
}

func (bbs BSBannerService) GetBanners(bsID int) (banners []entity.Banner, err error) {
	return bbs.repo.GetBanners(bsID)
}

func (bbs BSBannerService) DeleteBanner(bsID, bannerID int) error {
	return bbs.repo.DeleteBanner(bsID, bannerID)
}
