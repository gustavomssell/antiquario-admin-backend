package catalog

import (
	"time"
	domainCatalog "github.com/gbrayhan/microservices-go/src/domain/catalog"
)

type Style struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name          string    `gorm:"column:name;unique"`
	Description   string    `gorm:"column:description"`
	PeriodID      *string   `gorm:"column:period_id;type:uuid"`
	Period        *Period   `gorm:"foreignKey:PeriodID"`
	OriginCountry string    `gorm:"column:origin_country"`
	CreatedAt     time.Time `gorm:"autoCreateTime:mili"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime:mili"`
}

func (Style) TableName() string {
	return "styles"
}

func (s *Style) toDomainMapper() *domainCatalog.Style {
	return &domainCatalog.Style{
		ID:            s.ID,
		Name:          s.Name,
		Description:   s.Description,
		PeriodID:      s.PeriodID,
		OriginCountry: s.OriginCountry,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}

func fromDomainStyleMapper(s *domainCatalog.Style) *Style {
	return &Style{
		ID:            s.ID,
		Name:          s.Name,
		Description:   s.Description,
		PeriodID:      s.PeriodID,
		OriginCountry: s.OriginCountry,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}
