package contracts

import "bitbucket.org/lyndus/backend/domain/appuser/entity"

type AppUserBSBannerServiceInterface interface {
	GetBSBanners(bsID int) (banners []entity.BSBanner, err error)
}
