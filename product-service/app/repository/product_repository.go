package repository

import (
	"context"
	"database/sql"
	product "product-service/app"
	"product-service/app/models"
)

type productRepository struct {
	db *sql.DB
}

func NewproductRepository(db *sql.DB) product.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetProductList(ctx context.Context) ([]models.Product, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.name, p.description, p.price, 
			COALESCE(SUM(s.quantity - s.reserved), 0) as available_stock
		FROM products p
		LEFT JOIN stock s ON p.id = s.product_id
		LEFT JOIN warehouses w ON s.warehouse_id = w.id AND w.is_active = true
		GROUP BY p.id
		ORDER BY p.name
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
