package services

import (
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
)

type BsSocialService struct {
	//TODO Validar
	repo contracts.BsSocialRepositoryInterface
}

var BSSocialService contracts.BsSocialServiceInterface

func NewBSSocialService() *BsSocialService {
	return &BsSocialService{
		repo: repository.NewBSSocialRepository(),
	}
}

func (b BsSocialService) GetSocialByID(socialID int) (social entity.BsSocial, err error) {
	return b.repo.GetSocialByID(socialID)
}

func (b BsSocialService) GetBSSocialByID(bsID, socialID int) (social entity.BsSocial, err error) {
	return b.repo.GetBSSocialByID(bsID, socialID)
}

func (b BsSocialService) GetAll(bsID int) ([]entity.BsSocial, error) {
	return b.repo.GetAll(bsID)
}

func (b BsSocialService) GetSocials(bsID int) (socials []entity.Social, err error) {
	return b.repo.GetSocials(bsID)
}

func (b BsSocialService) Update(bsSocial entity.BsSocial) error {
	return b.repo.Update(bsSocial)
}

func (b BsSocialService) Create(bsSocial entity.BsSocial) (int, error) {
	return b.repo.Create(bsSocial)
}

func (b BsSocialService) Remove(bsID, social int) error {
	return b.repo.Remove(bsID, social)
}
