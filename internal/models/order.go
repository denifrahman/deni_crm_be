package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CustomerID uint           `json:"customer_id"`
	Customer   Customer       `json:"customer" gorm:"foreignKey:CustomerID"`
	Total      float64        `json:"total"`
	DealID     uint           `json:"deal_id"`
	Deal       Deal           `json:"deal" gorm:"foreignKey:DealID"`
	Location   string         `json:"location" gorm:"not null"`
	Lat        string         `json:"lat" `
	Status     string         `json:"status"`
	Lng        string         `json:"lng" `
	OrderItems []OrderItem    `json:"order_items"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type OrderItem struct {
	ID          uint           `gorm:"primaryKey"`
	ProductID   uint           `json:"product_id"`
	OrderID     uint           `json:"order_id"`
	Order       Order          `json:"order" gorm:"foreignKey:OrderID"`
	Product     Product        `json:"product" gorm:"foreignKey:ProductID"`
	Qty         int            `json:"qty" gorm:"not null"`
	Price       float64        `json:"price" gorm:"not null"`
	ProductName string         `json:"product_name" gorm:"not null"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type OrderRequestCreate struct {
	CustomerID uint        `json:"customer_id"`
	Total      float64     `json:"total"`
	DealID     uint        `json:"deal_id" binding:"required"`
	Location   string      `json:"location"`
	Status     string      `json:"status"`
	Lat        string      `json:"lat"`
	Lng        string      `json:"lng"`
	OrderItems []OrderItem `json:"order_items"`
	UserId     int         `json:"user_id"`
}

func (r *OrderRequestCreate) ToModel() *Order {
	orderItems := make([]OrderItem, len(r.OrderItems))
	for i, item := range r.OrderItems {
		orderItems[i] = OrderItem{
			ProductID: item.ProductID,
			Qty:       item.Qty,
			Price:     item.Price,
			Status:    "new order",
		}
	}

	return &Order{
		DealID:     r.DealID,
		CustomerID: r.CustomerID,
		Total:      r.Total,
		Location:   r.Location,
		Status:     "new order",
		Lat:        r.Lat,
		Lng:        r.Lng,
		OrderItems: orderItems,
	}
}

type OrderRequestUpdate struct {
	Id         uint        `json:"id"`
	CustomerID uint        `json:"customer_id" binding:"required"`
	Total      float64     `json:"total"`
	Location   string      `json:"location" binding:"required"`
	Lat        string      `json:"lat"`
	Lng        string      `json:"lng"`
	OrderItems []OrderItem `json:"order_items" binding:"required"`
}

func (r *OrderRequestUpdate) ToModel() *Order {
	orderItems := make([]OrderItem, len(r.OrderItems))
	for i, item := range r.OrderItems {
		orderItems[i] = OrderItem{
			ID:        item.ID,
			ProductID: item.ProductID,
			Qty:       item.Qty,
			Price:     item.Price,
			OrderID:   r.Id,
		}
	}

	return &Order{
		ID:         r.Id,
		CustomerID: r.CustomerID,
		Total:      r.Total,
		Location:   r.Location,
		Lat:        r.Lat,
		Lng:        r.Lng,
		OrderItems: orderItems,
	}
}

type OrderResponse struct {
	Data  []Order `json:"data"`
	Count int64   `json:"count"`
}

func OrderToResponse(r OrderResponse) *OrderResponse {
	return &OrderResponse{
		Data:  r.Data,
		Count: r.Count,
	}
}
