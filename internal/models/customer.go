package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Name    string `json:"name" gorm:"not null"`
	Email   string `json:"email"`
	Phone   string `json:"phone" gorm:"not null"`
	Company string `json:"company"`

	Address string `json:"address" gorm:"not null"`

	Status string `json:"status" gorm:"not null"`
	UserId int    `json:"user_id" gorm:"not null"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type CustomerRequestCreate struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email"`
	Phone   string `json:"phone" binding:"required"`
	Address string `json:"address" binding:"required"`
	Company string `json:"company"`
	Status  string `json:"status" binding:"required"`
	UserId  int    `json:"user_id"`
}

func (r *CustomerRequestCreate) ToModel() *Customer {
	return &Customer{
		Name:    r.Name,
		Email:   r.Email,
		Phone:   r.Phone,
		Address: r.Address,
		Company: r.Company,
		Status:  r.Status,
		UserId:  r.UserId,
	}
}

type CustomerRequestUpdate struct {
	Id      uint   `json:"id"`
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email"`
	Phone   string `json:"phone" binding:"required"`
	Address string `json:"address" binding:"required"`
	Company string `json:"company"`
	Status  string `json:"status" binding:"required"`
	UserId  int    `json:"user_id"`
}

func (r *CustomerRequestUpdate) ToModel() *Customer {
	return &Customer{
		ID:      r.Id,
		Name:    r.Name,
		Email:   r.Email,
		Address: r.Address,
		Phone:   r.Phone,
		Company: r.Company,
		Status:  r.Status,
		UserId:  r.UserId,
	}
}

type CustomerResponse struct {
	Data  []Customer `json:"data"`
	Count int64      `json:"count"`
}

func CustomerToResponse(r CustomerResponse) *CustomerResponse {
	return &CustomerResponse{
		Data:  r.Data,
		Count: r.Count,
	}
}
