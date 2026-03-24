package product

import (
	"time"
	domainProduct "github.com/gbrayhan/microservices-go/src/domain/product"
)

type ProductSet struct {
	ID                string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SetCode           string    `gorm:"column:set_code;unique"`
	Name              string    `gorm:"column:name"`
	Description       string    `gorm:"column:description"`
	TotalPieces       int       `gorm:"column:total_pieces"`
	CanSellSeparately bool      `gorm:"column:can_sell_separately"`
	CreatedBy         int       `gorm:"column:created_by"`
	CreatedAt         time.Time `gorm:"autoCreateTime:mili"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime:mili"`
}

func (ProductSet) TableName() string { return "product_sets" }

func (s *ProductSet) toDomainMapper() *domainProduct.ProductSet {
	return &domainProduct.ProductSet{
		ID:                s.ID,
		SetCode:           s.SetCode,
		Name:              s.Name,
		Description:       s.Description,
		TotalPieces:       s.TotalPieces,
		CanSellSeparately: s.CanSellSeparately,
		CreatedBy:         s.CreatedBy,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
}

func fromDomainProductSetMapper(s *domainProduct.ProductSet) *ProductSet {
	return &ProductSet{
		ID:                s.ID,
		SetCode:           s.SetCode,
		Name:              s.Name,
		Description:       s.Description,
		TotalPieces:       s.TotalPieces,
		CanSellSeparately: s.CanSellSeparately,
		CreatedBy:         s.CreatedBy,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
}
