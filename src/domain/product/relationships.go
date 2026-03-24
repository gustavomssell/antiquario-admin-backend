package product

type ProductMaterial struct {
	ID         string // UUID
	ProductID  string // FK to Product
	MaterialID string // FK to Material
	Percentage float64
	Notes      string
}

type ProductTag struct {
	ID        string // UUID
	ProductID string // FK to Product
	TagID     string // FK to Tag
}
