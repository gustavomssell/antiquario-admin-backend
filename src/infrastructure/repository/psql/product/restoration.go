package product

import (
	"time"
	domainProduct "github.com/gbrayhan/microservices-go/src/domain/product"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
)

type RestorationOrder struct {
	ID             string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductID      string     `gorm:"column:product_id;type:uuid"`
	Product        *Product   `gorm:"foreignKey:ProductID"`
	Status         string     `gorm:"column:status"`
	Priority       string     `gorm:"column:priority"`
	EstimatedCost  float64    `gorm:"column:estimated_cost"`
	ActualCost     *float64   `gorm:"column:actual_cost"`
	AssignedTo     *int       `gorm:"column:assigned_to"`
	AssignedUser   *user.User `gorm:"foreignKey:AssignedTo"`
	StartDate      *time.Time `gorm:"column:start_date"`
	CompletionDate *time.Time `gorm:"column:completion_date"`
	CreatedBy      int        `gorm:"column:created_by"`
	CreatedAt      time.Time  `gorm:"autoCreateTime:mili"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime:mili"`
}

func (RestorationOrder) TableName() string { return "restoration_orders" }

func (r *RestorationOrder) toDomainMapper() *domainProduct.RestorationOrder {
	return &domainProduct.RestorationOrder{
		ID:             r.ID,
		ProductID:      r.ProductID,
		Status:         r.Status,
		Priority:       r.Priority,
		EstimatedCost:  r.EstimatedCost,
		ActualCost:     r.ActualCost,
		AssignedTo:     r.AssignedTo,
		StartDate:      r.StartDate,
		CompletionDate: r.CompletionDate,
		CreatedBy:      r.CreatedBy,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

func fromDomainRestorationOrderMapper(r *domainProduct.RestorationOrder) *RestorationOrder {
	return &RestorationOrder{
		ID:             r.ID,
		ProductID:      r.ProductID,
		Status:         r.Status,
		Priority:       r.Priority,
		EstimatedCost:  r.EstimatedCost,
		ActualCost:     r.ActualCost,
		AssignedTo:     r.AssignedTo,
		StartDate:      r.StartDate,
		CompletionDate: r.CompletionDate,
		CreatedBy:      r.CreatedBy,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

type RestorationLog struct {
	ID                 string            `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RestorationOrderID string            `gorm:"column:restoration_order_id;type:uuid"`
	RestorationOrder   *RestorationOrder `gorm:"foreignKey:RestorationOrderID"`
	ActionDescription  string            `gorm:"column:action_description"`
	MaterialsUsed      string            `gorm:"column:materials_used;type:jsonb"`
	TechnicalHours     float64           `gorm:"column:technical_hours"`
	LoggedBy           int               `gorm:"column:logged_by"`
	CreatedAt          time.Time         `gorm:"autoCreateTime:mili"`
	UpdatedAt          time.Time         `gorm:"autoUpdateTime:mili"`
}

func (RestorationLog) TableName() string { return "restoration_logs" }

func (l *RestorationLog) toDomainMapper() *domainProduct.RestorationLog {
	return &domainProduct.RestorationLog{
		ID:                 l.ID,
		RestorationOrderID: l.RestorationOrderID,
		ActionDescription:  l.ActionDescription,
		MaterialsUsed:      l.MaterialsUsed,
		TechnicalHours:     l.TechnicalHours,
		LoggedBy:           l.LoggedBy,
		CreatedAt:          l.CreatedAt,
		UpdatedAt:          l.UpdatedAt,
	}
}

func fromDomainRestorationLogMapper(l *domainProduct.RestorationLog) *RestorationLog {
	return &RestorationLog{
		ID:                 l.ID,
		RestorationOrderID: l.RestorationOrderID,
		ActionDescription:  l.ActionDescription,
		MaterialsUsed:      l.MaterialsUsed,
		TechnicalHours:     l.TechnicalHours,
		LoggedBy:           l.LoggedBy,
		CreatedAt:          l.CreatedAt,
		UpdatedAt:          l.UpdatedAt,
	}
}
