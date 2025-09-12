package usecase

import (
	"errors"
	"time"

	product "github.com/evrintobing17/ecommerce-system/product-service/app"
	"github.com/evrintobing17/ecommerce-system/product-service/app/models"
)

var ErrInsufficientStock = errors.New("insufficient stock")

type productUsecase struct {
	productRepo product.ProductRepository
}

func NewProductUsecase(productRepo product.ProductRepository) product.ProductUsecase {
	return &productUsecase{productRepo: productRepo}
}

func (u *productUsecase) GetProducts(shopID int, page, limit int) ([]*models.Product, int64, error) {
	products, total, err := u.productRepo.FindAll(shopID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var result []*models.Product
	for _, product := range products {
		result = append(result, &models.Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			ShopID:      product.ShopID,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		})
	}

	return result, total, nil
}

func (u *productUsecase) GetProduct(id int) (*models.Product, error) {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &models.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		ShopID:      product.ShopID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

func (u *productUsecase) CreateProduct(name, description string, price float64, stock int32, shopID int) (*models.Product, error) {
	product := &models.Product{
		Name:        name,
		Description: description,
		Price:       price,
		ShopID:      shopID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := u.productRepo.Create(product)
	if err != nil {
		return nil, err
	}

	return &models.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		ShopID:      product.ShopID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

func (u *productUsecase) UpdateProduct(product *models.Product) error {
	existingProduct, err := u.productRepo.FindByID(product.ID)
	if err != nil {
		return err
	}

	existingProduct.Name = product.Name
	existingProduct.Description = product.Description
	existingProduct.Price = product.Price
	existingProduct.UpdatedAt = time.Now()

	return u.productRepo.Update(existingProduct)
}

func (u *productUsecase) DeleteProduct(id int) error {
	return u.productRepo.Delete(id)
}
