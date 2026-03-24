package product

import (
	"time"
	domainProduct "github.com/gbrayhan/microservices-go/src/domain/product"
)

type ProductMedia struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductID     string    `gorm:"column:product_id;type:uuid"`
	Product       *Product  `gorm:"foreignKey:ProductID"`
	MediaType     string    `gorm:"column:media_type"`
	MediaURL      string    `gorm:"column:media_url"`
	AltText       string    `gorm:"column:alt_text"`
	IsPrimary     bool      `gorm:"column:is_primary"`
	OrderPosition int       `gorm:"column:order_position"`
	FileSize      int64     `gorm:"column:file_size"`
	Duration      *int      `gorm:"column:duration"`
	CreatedAt     time.Time `gorm:"autoCreateTime:mili"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime:mili"`
}

func (ProductMedia) TableName() string { return "product_media" }

func (m *ProductMedia) toDomainMapper() *domainProduct.ProductMedia {
	return &domainProduct.ProductMedia{
		ID:            m.ID,
		ProductID:     m.ProductID,
		MediaType:     m.MediaType,
		MediaURL:      m.MediaURL,
		AltText:       m.AltText,
		IsPrimary:     m.IsPrimary,
		OrderPosition: m.OrderPosition,
		FileSize:      m.FileSize,
		Duration:      m.Duration,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func fromDomainProductMediaMapper(m *domainProduct.ProductMedia) *ProductMedia {
	return &ProductMedia{
		ID:            m.ID,
		ProductID:     m.ProductID,
		MediaType:     m.MediaType,
		MediaURL:      m.MediaURL,
		AltText:       m.AltText,
		IsPrimary:     m.IsPrimary,
		OrderPosition: m.OrderPosition,
		FileSize:      m.FileSize,
		Duration:      m.Duration,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
