package products

type Products struct{
	products []models.Product
}

func New() *Products{
	return &Products{products: make([]models.Product,0)}
}

func (p Products) GetProducts()[]models.Product{
	return p.products
}

func (p Products) AddProduct(newProduct models.Product){
	
}