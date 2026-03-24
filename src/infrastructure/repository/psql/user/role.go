package user

import (
	"time"

	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
)

type Role struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"column:name;unique"`
	Description string    `gorm:"column:description"`
	Permissions string    `gorm:"column:permissions;type:jsonb"`
	CreatedAt   time.Time `gorm:"autoCreateTime:mili"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:mili"`
}

func (Role) TableName() string {
	return "roles"
}

func (r *Role) toDomainMapper() *domainUser.Role {
	return &domainUser.Role{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: r.Permissions,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func fromDomainRoleMapper(r *domainUser.Role) *Role {
	return &Role{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: r.Permissions,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
