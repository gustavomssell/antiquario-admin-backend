package product

import (
	"time"
	domainProduct "github.com/gbrayhan/microservices-go/src/domain/product"
)

type ProductDocument struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductID    string    `gorm:"column:product_id;type:uuid"`
	Product      *Product  `gorm:"foreignKey:ProductID"`
	DocumentType string    `gorm:"column:document_type"`
	DocumentURL  string    `gorm:"column:document_url"`
	Title        string    `gorm:"column:title"`
	Description  string    `gorm:"column:description"`
	FileSize     int64     `gorm:"column:file_size"`
	CreatedAt    time.Time `gorm:"autoCreateTime:mili"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime:mili"`
}

func (ProductDocument) TableName() string { return "product_documents" }

func (d *ProductDocument) toDomainMapper() *domainProduct.ProductDocument {
	return &domainProduct.ProductDocument{
		ID:           d.ID,
		ProductID:    d.ProductID,
		DocumentType: d.DocumentType,
		DocumentURL:  d.DocumentURL,
		Title:        d.Title,
		Description:  d.Description,
		FileSize:     d.FileSize,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

func fromDomainProductDocumentMapper(d *domainProduct.ProductDocument) *ProductDocument {
	return &ProductDocument{
		ID:           d.ID,
		ProductID:    d.ProductID,
		DocumentType: d.DocumentType,
		DocumentURL:  d.DocumentURL,
		Title:        d.Title,
		Description:  d.Description,
		FileSize:     d.FileSize,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}
