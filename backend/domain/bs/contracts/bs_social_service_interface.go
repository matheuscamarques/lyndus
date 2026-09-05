package contracts

import "bitbucket.org/lyndus/backend/domain/bs/entity"

type BsSocialServiceInterface interface {
	GetSocialByID(socialID int) (social entity.BsSocial, err error)
	GetBSSocialByID(bsID, socialID int) (social entity.BsSocial, err error)
	GetSocials(bsID int) (socials []entity.Social, err error)
	GetAll(bsID int) ([]entity.BsSocial,error)
	Update(bsSocial entity.BsSocial) error
	Create(bsSocial entity.BsSocial) (int,error)
	Remove(bsID, social int) error
}
