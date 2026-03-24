package product

import (
	"time"
	domainProduct "github.com/gbrayhan/microservices-go/src/domain/product"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/catalog"
)

type Product struct {
	ID                     string              `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Code                   string              `gorm:"column:code;unique"`
	QRCodeURL              string              `gorm:"column:qr_code_url"`
	ProductSetID           *string             `gorm:"column:product_set_id;type:uuid"`
	ProductSet             *ProductSet         `gorm:"foreignKey:ProductSetID"`
	SetPosition            *string             `gorm:"column:set_position"`
	Name                   string              `gorm:"column:name"`
	Description            string              `gorm:"column:description"`
	CategoryID             *string             `gorm:"column:category_id;type:uuid"`
	Category               *catalog.Category   `gorm:"foreignKey:CategoryID"`
	PeriodID               *string             `gorm:"column:period_id;type:uuid"`
	Period                 *catalog.Period     `gorm:"foreignKey:PeriodID"`
	StyleID                *string             `gorm:"column:style_id;type:uuid"`
	Style                  *catalog.Style      `gorm:"foreignKey:StyleID"`
	Dimensions             string              `gorm:"column:dimensions;type:jsonb"`
	Weight                 float64             `gorm:"column:weight"`
	ConditionRating        int                 `gorm:"column:condition_rating"`
	AcquisitionType        string              `gorm:"column:acquisition_type"`
	AcquisitionDate        *time.Time          `gorm:"column:acquisition_date"`
	AcquisitionPrice       float64             `gorm:"column:acquisition_price"`
	ConsignmentPercentage  float64             `gorm:"column:consignment_percentage"`
	ConsignmentDeadline    *time.Time          `gorm:"column:consignment_deadline"`
	SupplierID             *string             `gorm:"column:supplier_id;type:uuid"`
	EstimatedValue         float64             `gorm:"column:estimated_value"`
	SellingPrice           float64             `gorm:"column:selling_price"`
	CommissionRate         float64             `gorm:"column:commission_rate"`
	Status                 string              `gorm:"column:status"`
	ProvenanceStory        string              `gorm:"column:provenance_story"`
	HistoricalNotes        string              `gorm:"column:historical_notes"`
	IsSetItem              bool                `gorm:"column:is_set_item"`
	AvailableQuantity      int                 `gorm:"column:available_quantity"`
	CreatedBy              int                 `gorm:"column:created_by"`
	UpdatedBy              int                 `gorm:"column:updated_by"`
	CreatedAt              time.Time           `gorm:"autoCreateTime:mili"`
	UpdatedAt              time.Time           `gorm:"autoUpdateTime:mili"`
}

func (Product) TableName() string { return "products" }

func (p *Product) toDomainMapper() *domainProduct.Product {
	return &domainProduct.Product{
		ID:                    p.ID,
		Code:                  p.Code,
		QRCodeURL:             p.QRCodeURL,
		ProductSetID:          p.ProductSetID,
		SetPosition:           p.SetPosition,
		Name:                  p.Name,
		Description:           p.Description,
		CategoryID:            p.CategoryID,
		PeriodID:              p.PeriodID,
		StyleID:               p.StyleID,
		Dimensions:            p.Dimensions,
		Weight:                p.Weight,
		ConditionRating:       p.ConditionRating,
		AcquisitionType:       p.AcquisitionType,
		AcquisitionDate:       p.AcquisitionDate,
		AcquisitionPrice:      p.AcquisitionPrice,
		ConsignmentPercentage: p.ConsignmentPercentage,
		ConsignmentDeadline:   p.ConsignmentDeadline,
		SupplierID:            p.SupplierID,
		EstimatedValue:        p.EstimatedValue,
		SellingPrice:          p.SellingPrice,
		CommissionRate:        p.CommissionRate,
		Status:                p.Status,
		ProvenanceStory:       p.ProvenanceStory,
		HistoricalNotes:       p.HistoricalNotes,
		IsSetItem:             p.IsSetItem,
		AvailableQuantity:     p.AvailableQuantity,
		CreatedBy:             p.CreatedBy,
		UpdatedBy:             p.UpdatedBy,
		CreatedAt:             p.CreatedAt,
		UpdatedAt:             p.UpdatedAt,
	}
}

func fromDomainProductMapper(p *domainProduct.Product) *Product {
	return &Product{
		ID:                    p.ID,
		Code:                  p.Code,
		QRCodeURL:             p.QRCodeURL,
		ProductSetID:          p.ProductSetID,
		SetPosition:           p.SetPosition,
		Name:                  p.Name,
		Description:           p.Description,
		CategoryID:            p.CategoryID,
		PeriodID:              p.PeriodID,
		StyleID:               p.StyleID,
		Dimensions:            p.Dimensions,
		Weight:                p.Weight,
		ConditionRating:       p.ConditionRating,
		AcquisitionType:       p.AcquisitionType,
		AcquisitionDate:       p.AcquisitionDate,
		AcquisitionPrice:      p.AcquisitionPrice,
		ConsignmentPercentage: p.ConsignmentPercentage,
		ConsignmentDeadline:   p.ConsignmentDeadline,
		SupplierID:            p.SupplierID,
		EstimatedValue:        p.EstimatedValue,
		SellingPrice:          p.SellingPrice,
		CommissionRate:        p.CommissionRate,
		Status:                p.Status,
		ProvenanceStory:       p.ProvenanceStory,
		HistoricalNotes:       p.HistoricalNotes,
		IsSetItem:             p.IsSetItem,
		AvailableQuantity:     p.AvailableQuantity,
		CreatedBy:             p.CreatedBy,
		UpdatedBy:             p.UpdatedBy,
		CreatedAt:             p.CreatedAt,
		UpdatedAt:             p.UpdatedAt,
	}
}
