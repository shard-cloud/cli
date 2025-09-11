package http

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type Subscription struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

type MeResponse struct {
	Subscription *Subscription `json:"subscription"`
	User         *User         `json:"user"`
}

type IdResponse struct {
	ID uuid.UUID `json:"id"`
}

type App struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Ram       int       `json:"ram"`
	Vcpu      int       `json:"vcpu"`
}

type AppResponse struct {
	App          *App   `json:"app"`
	RealTimeRam  int    `json:"real_time_ram"`
	RealTimeVcpu int    `json:"real_time_vcpu"`
	Status       string `json:"status"`
}

type Backup struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created"`
}

type BackupsResponse struct {
	Backups []*Backup `json:"backups"`
}

type BackupResponse struct {
	Backup *Backup `json:"backup"`
}
