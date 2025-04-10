package postgres

//TODO: add toModel methods everywhere
import (
	"AvitoPvz/internal/domain/models"
	"github.com/google/uuid"
	"time"
)

type userDTO struct {
	ID       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
	Role     string    `db:"role"`
}

func (u userDTO) toModel() (*models.User, error) {
	role, err := models.GetRoleFromString(u.Role)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		Role:     role,
	}, nil
}

type pvzDTO struct {
	ID              uuid.UUID `db:"id"`
	City            string    `db:"city"`
	RegitrationTime time.Time `db:"regitration_time"`
}

type receptionDTO struct {
	ID       uuid.UUID `db:"id"`
	PVZID    uuid.UUID `db:"pvz_id"`
	DateTime time.Time `db:"datetime"`
	Status   string    `db:"status"`
}

type productDTO struct {
	ID          uuid.UUID `db:"id"`
	ReceptionID uuid.UUID `db:"reception_id"`
	Type        string    `db:"type"`
	DateTime    time.Time `db:"datetime"`
}

func (p productDTO) toModel() (*models.Product, error) {
	t, err := models.GetProductTypeFromString(p.Type)
	if err != nil {
		return nil, err
	}

	return &models.Product{
		ID:          p.ID,
		ReceptionID: p.ReceptionID,
		Type:        t,
		DateTime:    p.DateTime,
	}, nil
}
