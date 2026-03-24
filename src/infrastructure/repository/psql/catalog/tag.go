package catalog

import (
	"time"
	domainCatalog "github.com/gbrayhan/microservices-go/src/domain/catalog"
)

type Tag struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"column:name;unique"`
	Color       string    `gorm:"column:color"`
	Description string    `gorm:"column:description"`
	Active      bool      `gorm:"column:active;default:true"`
	CreatedAt   time.Time `gorm:"autoCreateTime:mili"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:mili"`
}

func (Tag) TableName() string {
	return "tags"
}

func (t *Tag) toDomainMapper() *domainCatalog.Tag {
	return &domainCatalog.Tag{
		ID:          t.ID,
		Name:        t.Name,
		Color:       t.Color,
		Description: t.Description,
		Active:      t.Active,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func fromDomainTagMapper(t *domainCatalog.Tag) *Tag {
	return &Tag{
		ID:          t.ID,
		Name:        t.Name,
		Color:       t.Color,
		Description: t.Description,
		Active:      t.Active,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
