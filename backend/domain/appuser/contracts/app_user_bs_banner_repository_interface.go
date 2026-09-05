package contracts

import "bitbucket.org/lyndus/backend/domain/appuser/entity"

type AppUserBSBannerRepositoryInterface interface {
	GetBSBanners(bsID int) (banners []entity.BSBanner, err error)
}
