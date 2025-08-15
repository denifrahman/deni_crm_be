package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Name   string  `json:"name" gorm:"not null"`
	Hpp    float64 `json:"hpp" gorm:"not null"`
	Price  float64 `json:"price" gorm:"not null"`
	Margin int     `json:"margin" gorm:"not null"`
	Speed  string  `json:"speed" gorm:"not null"`

	Duration int `json:"duration"`

	Status    string         `json:"status" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type ProductRequestCreate struct {
	Name     string  `json:"name" binding:"required"`
	Hpp      float64 `json:"hpp" binding:"required"`
	Price    float64 `json:"price"`
	Speed    string  `json:"speed"`
	Duration int     `json:"duration" binding:"required"`
	Margin   int     `json:"margin" binding:"required"`
	Status   string  `json:"status" binding:"required"`
}

func (r *ProductRequestCreate) ToModel() *Product {
	return &Product{
		Name:     r.Name,
		Hpp:      r.Hpp,
		Speed:    r.Speed,
		Duration: r.Duration,
		Price:    r.Price,
		Margin:   r.Margin,
		Status:   r.Status,
	}
}

type ProductRequestUpdate struct {
	Id        uint    `json:"id"`
	Name      string  `json:"name" binding:"required"`
	Hpp       float64 `json:"hpp" binding:"required"`
	Price     float64 `json:"price"`
	Duration  int     `json:"duration"`
	Speed     string  `json:"speed"`
	Margin    int     `json:"margin" binding:"required"`
	Status    string  `json:"status" binding:"required"`
	CreatedAt time.Time
}

func (r *ProductRequestUpdate) ToModel() *Product {
	return &Product{
		ID:        r.Id,
		Name:      r.Name,
		Hpp:       r.Hpp,
		Speed:     r.Speed,
		Duration:  r.Duration,
		Price:     r.Price,
		Margin:    r.Margin,
		Status:    r.Status,
		CreatedAt: r.CreatedAt,
	}
}

type ProductResponse struct {
	Data  []Product `json:"data"`
	Count int64     `json:"count"`
}

func ProductToResponse(r ProductResponse) *ProductResponse {
	return &ProductResponse{
		Data:  r.Data,
		Count: r.Count,
	}
}
