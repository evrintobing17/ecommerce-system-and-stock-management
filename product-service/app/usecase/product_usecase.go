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

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int32     `json:"stock"`
	ShopID      int       `json:"shop_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
			Stock:       product.Stock,
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
		Stock:       product.Stock,
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
		Stock:       stock,
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
		Stock:       product.Stock,
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

func (u *productUsecase) AddStock(id int, quantity int32) error {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return err
	}

	product.Stock += quantity
	product.UpdatedAt = time.Now()

	return u.productRepo.UpdateStock(id, product.Stock)
}

func (u *productUsecase) SubtractStock(id int, quantity int32) error {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return err
	}

	if product.Stock < quantity {
		return ErrInsufficientStock
	}

	product.Stock -= quantity
	product.UpdatedAt = time.Now()

	return u.productRepo.UpdateStock(id, product.Stock)
}

func (u *productUsecase) SetStock(id int, quantity int32) error {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return err
	}

	product.Stock = quantity
	product.UpdatedAt = time.Now()

	return u.productRepo.UpdateStock(id, product.Stock)
}

func (u *productUsecase) DeleteProduct(id int) error {
	return u.productRepo.Delete(id)
}
