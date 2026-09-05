package contracts

import "bitbucket.org/lyndus/backend/domain/bs/entity"

type BSBannerRepositoryInterface interface {
	CreateBanner(bsID int, name string) (id int, err error)
	GetBanners(bsID int) (banners []entity.Banner, err error)
	DeleteBanner(bsID, bannerID int) error
}
