package catalog

import (
	"time"
	domainCatalog "github.com/gbrayhan/microservices-go/src/domain/catalog"
)

type Category struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"column:name;unique"`
	Description string    `gorm:"column:description"`
	ParentID    *string   `gorm:"column:parent_id;type:uuid"`
	Parent      *Category `gorm:"foreignKey:ParentID"`
	Level       int       `gorm:"column:level;default:1"`
	ImageURL    string    `gorm:"column:image_url"`
	Active      bool      `gorm:"column:active;default:true"`
	CreatedBy   int       `gorm:"column:created_by"`
	CreatedAt   time.Time `gorm:"autoCreateTime:mili"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:mili"`
}

func (Category) TableName() string {
	return "categories"
}

func (c *Category) toDomainMapper() *domainCatalog.Category {
	return &domainCatalog.Category{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		ParentID:    c.ParentID,
		Level:       c.Level,
		ImageURL:    c.ImageURL,
		Active:      c.Active,
		CreatedBy:   c.CreatedBy,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func fromDomainCategoryMapper(c *domainCatalog.Category) *Category {
	return &Category{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		ParentID:    c.ParentID,
		Level:       c.Level,
		ImageURL:    c.ImageURL,
		Active:      c.Active,
		CreatedBy:   c.CreatedBy,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
