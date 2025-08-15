package models

import (
	"time"

	"gorm.io/gorm"
)

type Lead struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email"`
	Phone     string         `json:"phone" gorm:"not null"`
	Company   string         `json:"company"`
	Status    Status         `json:"status" gorm:"not null"`
	Address   string         `json:"address"`
	UserId    int            `json:"user_id" gorm:"not null"`
	Needs     string         `json:"needs" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Status string

const (
	StatusNew         Status = "new"
	StatusContacted   Status = "contacted"
	StatusQualified   Status = "qualified"
	StatusUnqualified Status = "unqualified"
	StatusInProgress  Status = "in_progress"
	StatusLost        Status = "lost"
)

func IsValidLeadStatus(status string) bool {
	switch Status(status) {
	case StatusNew, StatusContacted, StatusQualified,
		StatusUnqualified, StatusInProgress, StatusLost:
		return true
	}
	return false
}

type LeadRequestCreate struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email"`
	Phone   string `json:"phone" binding:"required"`
	Company string `json:"company"`
	Needs   string `json:"needs" binding:"required"`
	Status  Status `json:"status" binding:"required"`
	UserId  int    `json:"user_id"`
}

func (r *LeadRequestCreate) ToModel() *Lead {
	return &Lead{
		Name:    r.Name,
		Email:   r.Email,
		Phone:   r.Phone,
		Company: r.Company,
		Status:  r.Status,
		Needs:   r.Needs,
		UserId:  r.UserId,
	}
}

type LeadRequestUpdate struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Company   string `json:"company"`
	Needs     string `json:"needs"`
	Status    Status `json:"status"`
	UserId    int    `json:"user_id"`
	CreatedAt time.Time
}

func (r *LeadRequestUpdate) ToModel() *Lead {
	return &Lead{
		ID:        r.Id,
		Name:      r.Name,
		Email:     r.Email,
		Phone:     r.Phone,
		Company:   r.Company,
		Status:    r.Status,
		Needs:     r.Needs,
		CreatedAt: r.CreatedAt,
		UserId:    r.UserId,
	}
}

type LeadResponse struct {
	Data  []Lead `json:"data"`
	Count int64  `json:"count"`
}

func LeadToResponse(r LeadResponse) *LeadResponse {
	return &LeadResponse{
		Data:  r.Data,
		Count: r.Count,
	}
}

type LeadToDeal struct {
	Id     uint       `json:"id"`
	UserId int        `json:"user_id"`
	Items  []DealItem `json:"details" binding:"required"`
}
