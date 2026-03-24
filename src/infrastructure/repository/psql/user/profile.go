package user

import (
	"time"

	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
)

type UserProfile struct {
	ID             string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID         int        `gorm:"column:user_id;uniqueIndex"`
	Phone          string     `gorm:"column:phone;type:jsonb"`
	Address        string     `gorm:"column:address;type:jsonb"`
	DocumentType   string     `gorm:"column:document_type"`
	DocumentNumber string     `gorm:"column:document_number"`
	BirthDate      *time.Time `gorm:"column:birth_date"`
	CreatedAt      time.Time  `gorm:"autoCreateTime:mili"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime:mili"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}

func (p *UserProfile) toDomainMapper() *domainUser.UserProfile {
	return &domainUser.UserProfile{
		ID:             p.ID,
		UserID:         p.UserID,
		Phone:          p.Phone,
		Address:        p.Address,
		DocumentType:   p.DocumentType,
		DocumentNumber: p.DocumentNumber,
		BirthDate:      p.BirthDate,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

func fromDomainProfileMapper(p *domainUser.UserProfile) *UserProfile {
	return &UserProfile{
		ID:             p.ID,
		UserID:         p.UserID,
		Phone:          p.Phone,
		Address:        p.Address,
		DocumentType:   p.DocumentType,
		DocumentNumber: p.DocumentNumber,
		BirthDate:      p.BirthDate,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}
