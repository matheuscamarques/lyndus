package services

import (
	"bitbucket.org/lyndus/backend/domain/bs"
	"bitbucket.org/lyndus/backend/domain/bs/contracts"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
	"bitbucket.org/lyndus/backend/infra/criteria"
	"github.com/shopspring/decimal"
	"strconv"
)

type BSProductService struct {
	repo contracts.BSProductRepositoryInterface
}

var ProductService contracts.BSProductServiceInterface

func NewBSProductService() *BSProductService {
	repo := repository.NewBSProductRepository()
	return &BSProductService{
		repo: &repo,
	}
}

func (bps BSProductService) GetCategories() (categories []entity.ProductCategory, err error) {
	return bps.repo.GetCategories()
}

func (bps BSProductService) Update(product entity.Product) (err error) {
	return bps.repo.Update(product)
}
func (bps BSProductService) GetByID(bsID, productID int) (product entity.Product, err error) {
	return bps.repo.GetByID(bsID, productID)
}

func (bps BSProductService) CountProducts(bsID int) (n int, err error) {
	return bps.repo.CountProducts(bsID)
}

func (bps BSProductService) GetPackages() (packages []entity.ProductPackage, err error) {
	return bps.repo.GetPackages()
}
func (bps BSProductService) GetUnits() (units []entity.ProductUnit, err error) {
	return bps.repo.GetUnits()
}
func (bps BSProductService) GetClasses() (classes []entity.ProductClass, err error) {
	return bps.repo.GetClasses()
}

func (bps BSProductService) Parameters() (resp map[string]interface{}, err error) {
	resp = make(map[string]interface{})

	categories, err := bps.GetCategories()
	if err != nil {
		return resp, err
	}
	packages, err := bps.GetPackages()
	if err != nil {
		return resp, err
	}
	units, err := bps.GetUnits()
	if err != nil {
		return resp, err
	}

	classes := bs.ProductClasses

	resp["categories"] = categories
	resp["packages"] = packages
	resp["units"] = units
	resp["classes"] = classes

	return resp, nil
}

func (bps BSProductService) Create(product entity.Product) (id int, err error) {
	return bps.repo.Create(product)
}

func (bps BSProductService) GetAllProducts(bsID, page, itemsPerPage int) (productR criteria.CProductResponse, err error) {
	c := criteria.Criteria{ActivePage: page, ItemsPerPage: itemsPerPage}
	c.CheckPages()

	productR.TotalPages, err = c.TotalPagesWithWhere(bps.repo, "WHERE bs_id = "+strconv.Itoa(bsID))
	productR.ActivePage = c.ActivePage
	productR.Items, err = bps.repo.GetAllProducts(bsID, c)
	return productR, err
}

func (bps BSProductService) GetProductsSimple(bsID int) (products []entity.ProductSimple, err error) {
	return bps.repo.GetProductsSimple(bsID)
}

func (bps BSProductService) StockSave(bsId, createdBy, id int, quantity int, costPrice, salesPrice, profitMargin decimal.Decimal) error {
	return bps.repo.StockSave(bsId, createdBy, id, quantity, costPrice, salesPrice, profitMargin)
}

func (bps BSProductService) GetProductSimple(bsID, productID int) (product entity.ProductSimple, err error) {
	return bps.repo.GetProductSimple(bsID, productID)
}
