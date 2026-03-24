package catalog

import (
	"time"
	domainCatalog "github.com/gbrayhan/microservices-go/src/domain/catalog"
)

type Period struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"column:name;unique"`
	Description string    `gorm:"column:description"`
	StartYear   *int      `gorm:"column:start_year"`
	EndYear     *int      `gorm:"column:end_year"`
	CreatedAt   time.Time `gorm:"autoCreateTime:mili"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:mili"`
}

func (Period) TableName() string {
	return "periods"
}

func (p *Period) toDomainMapper() *domainCatalog.Period {
	return &domainCatalog.Period{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		StartYear:   p.StartYear,
		EndYear:     p.EndYear,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func fromDomainPeriodMapper(p *domainCatalog.Period) *Period {
	return &Period{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		StartYear:   p.StartYear,
		EndYear:     p.EndYear,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
