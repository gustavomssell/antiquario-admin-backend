package catalog

import (
	"time"
	domainCatalog "github.com/gbrayhan/microservices-go/src/domain/catalog"
)

type Material struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"column:name;unique"`
	Description string    `gorm:"column:description"`
	Category    string    `gorm:"column:category"`
	CreatedAt   time.Time `gorm:"autoCreateTime:mili"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:mili"`
}

func (Material) TableName() string {
	return "materials"
}

func (m *Material) toDomainMapper() *domainCatalog.Material {
	return &domainCatalog.Material{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Category:    m.Category,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func fromDomainMaterialMapper(m *domainCatalog.Material) *Material {
	return &Material{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Category:    m.Category,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
