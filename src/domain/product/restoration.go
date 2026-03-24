package product

import "time"

type RestorationOrder struct {
	ID             string     // UUID
	ProductID      string     // FK to Product
	Status         string     // e.g., Aguardando, Em Progresso, Suspenso, Concluído, Cancelado
	Priority       string     // e.g., Baixa, Média, Alta, Crítica
	EstimatedCost  float64
	ActualCost     *float64
	AssignedTo     *int       // FK to User (Responsável pelo restauro)
	StartDate      *time.Time
	CompletionDate *time.Time
	CreatedBy      int        // FK to User
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type RestorationLog struct {
	ID                 string    // UUID
	RestorationOrderID string    // FK to RestorationOrder
	ActionDescription  string    // TEXT (Intervenções técnicas detalhadas)
	MaterialsUsed      string    // JSON structure for materials used
	TechnicalHours     float64
	LoggedBy           int       // FK to User
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type SearchResultRestorationOrder struct {
	Data       *[]RestorationOrder
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}
