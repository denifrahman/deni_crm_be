package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	ID         uint `gorm:"primaryKey"`
	CustomerID uint
	ProductID  uint
	StartDate  time.Time
	EndDate    *time.Time
	IsActive   bool

	Customer Customer `gorm:"foreignKey:CustomerID"`
	Product  Product  `gorm:"foreignKey:ProductID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type SubscriptionRequestCreate struct {
	CustomerID uint       `json:"customer_id"`
	ProductID  uint       `json:"product_id"`
	StartDate  time.Time  `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
	IsActive   bool       `json:"is_active"`
}

func (t *SubscriptionRequestCreate) ToModel() *Subscription {
	return &Subscription{
		CustomerID: t.CustomerID,
		ProductID:  t.ProductID,
		StartDate:  t.StartDate,
		EndDate:    t.EndDate,
		IsActive:   t.IsActive,
	}
}

type SubscriptionRequestUpdate struct {
	Id         uint       `json:"id"`
	CustomerID uint       `json:"customer_id"`
	ProductID  uint       `json:"product_id"`
	StartDate  time.Time  `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
	IsActive   bool       `json:"is_active"`
}

func (t *SubscriptionRequestUpdate) ToModel() *Subscription {
	return &Subscription{
		CustomerID: t.CustomerID,
		ProductID:  t.ProductID,
		StartDate:  t.StartDate,
		EndDate:    t.EndDate,
		IsActive:   t.IsActive,
	}
}

type SubscriptionResponse struct {
	Data  []Subscription `json:"data"`
	Count int64          `json:"count"`
}

func SubscriptionToResponse(r SubscriptionResponse) *SubscriptionResponse {
	return &SubscriptionResponse{
		Data:  r.Data,
		Count: r.Count,
	}
}
