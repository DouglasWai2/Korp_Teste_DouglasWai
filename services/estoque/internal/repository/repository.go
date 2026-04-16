package repository

import(
	"database/sql"
)

type Repositories struct {
	Product interface {
		GetProducts()[]models.Product
		AddProduct()
		
	}
}

func CreateProduct(db *sql.DB, )