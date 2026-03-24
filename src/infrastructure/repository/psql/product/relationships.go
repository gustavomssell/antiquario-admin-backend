package product

import (
	domainProduct "github.com/gbrayhan/microservices-go/src/domain/product"
)

type ProductMaterial struct {
	ID         string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductID  string  `gorm:"column:product_id;type:uuid"`
	MaterialID string  `gorm:"column:material_id;type:uuid"`
	Percentage float64 `gorm:"column:percentage"`
	Notes      string  `gorm:"column:notes"`
}

func (ProductMaterial) TableName() string { return "product_materials" }

func (pm *ProductMaterial) toDomainMapper() *domainProduct.ProductMaterial {
	return &domainProduct.ProductMaterial{
		ID:         pm.ID,
		ProductID:  pm.ProductID,
		MaterialID: pm.MaterialID,
		Percentage: pm.Percentage,
		Notes:      pm.Notes,
	}
}

func fromDomainProductMaterialMapper(pm *domainProduct.ProductMaterial) *ProductMaterial {
	return &ProductMaterial{
		ID:         pm.ID,
		ProductID:  pm.ProductID,
		MaterialID: pm.MaterialID,
		Percentage: pm.Percentage,
		Notes:      pm.Notes,
	}
}

type ProductTag struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductID string `gorm:"column:product_id;type:uuid"`
	TagID     string `gorm:"column:tag_id;type:uuid"`
}

func (ProductTag) TableName() string { return "product_tags" }

func (pt *ProductTag) toDomainMapper() *domainProduct.ProductTag {
	return &domainProduct.ProductTag{
		ID:        pt.ID,
		ProductID: pt.ProductID,
		TagID:     pt.TagID,
	}
}

func fromDomainProductTagMapper(pt *domainProduct.ProductTag) *ProductTag {
	return &ProductTag{
		ID:        pt.ID,
		ProductID: pt.ProductID,
		TagID:     pt.TagID,
	}
}
