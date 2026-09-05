package services

import (
	"bitbucket.org/lyndus/backend/domain/client/contracts"
	"bitbucket.org/lyndus/backend/domain/client/repository"
	"bitbucket.org/lyndus/backend/domain/client/services"
)

type ClientCategoryService struct{
	services.ClientCategoryService
}

var ClientCategory contracts.ClientCategoryServiceInterface

func NewClientCategoryService() (*ClientCategoryService) {
	service := &ClientCategoryService{}
	service.Repo = repository.NewClientCategoryRepository();
	return service;
}