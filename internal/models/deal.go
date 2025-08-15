package models

import (
	"time"

	"gorm.io/gorm"
)

type Deal struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `json:"name" gorm:"not null"`
	Email      string         `json:"email"`
	Phone      string         `json:"phone" gorm:"not null"`
	Company    string         `json:"company"`
	Address    string         `json:"address"`
	StatusDeal StatusDeal     `json:"status_deal" gorm:"not null"`
	UserId     int            `json:"user_id" gorm:"not null"`
	Needs      string         `json:"needs" gorm:"type:text"`
	Notes      string         `json:"notes" gorm:"type:text"`
	Items      []DealItem     `json:"items" gorm:"foreignKey:DealID"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type DealItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	DealID    uint    `json:"deal_id"`
	ProductID uint    `json:"product_id"`
	Approval  bool    `json:"approval"`
	Approved  bool    `json:"approved"`
	Notes     string  `json:"notes" gorm:"type:text"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Qty       int     `json:"qty"`
	Price     float64 `json:"price"`
}

type StatusDeal string

const (
	StatusDealQualified     StatusDeal = "qualified"
	StatusDealProposalSend  StatusDeal = "proposal_send"
	StatusDealInNegotiation StatusDeal = "negotiation"
	StatusDealWon           StatusDeal = "won"
	StatusDealLost          StatusDeal = "lost"
	StatusDealDone          StatusDeal = "done"
)

func IsValidDealStatusDeal(statusDeal string) bool {
	switch StatusDeal(statusDeal) {
	case StatusDealQualified,
		StatusDealProposalSend, StatusDealInNegotiation, StatusDealWon, StatusDealDone, StatusDealLost:
		return true
	}
	return false
}

type DealRequestCreate struct {
	Name       string     `json:"name" binding:"required"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone" binding:"required"`
	Company    string     `json:"company"`
	Needs      string     `json:"needs" binding:"required"`
	Items      []DealItem `json:"items"`
	StatusDeal StatusDeal `json:"status" binding:"required"`
	UserId     int        `json:"user_id"`
}

func (r *DealRequestCreate) ToModel() *Deal {
	return &Deal{
		Name:       r.Name,
		Email:      r.Email,
		Phone:      r.Phone,
		Company:    r.Company,
		Items:      r.Items,
		StatusDeal: r.StatusDeal,
		Needs:      r.Needs,
		UserId:     r.UserId,
	}
}

type DealRequestUpdate struct {
	Id         uint       `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	Company    string     `json:"company"`
	Notes      string     `json:"notes"`
	Items      []DealItem `json:"items"`
	Needs      string     `json:"needs" binding:"required"`
	StatusDeal StatusDeal `json:"status" binding:"required"`
	UserId     int        `json:"user_id"`
	CreatedAt  time.Time
	IsLeader   bool
}

type DealItemApproveRequestUpdate struct {
	DealItemId uint `json:"deal_item_id" binding:"required"`
	Approved   bool `json:"approved" binding:"required"`
	IsLeader   bool
}

func (r *DealItemApproveRequestUpdate) ToModel() *DealItem {
	return &DealItem{
		ID:       r.DealItemId,
		Approved: r.Approved,
	}
}

func (r *DealRequestUpdate) ToModel() *Deal {

	return &Deal{
		ID:         r.Id,
		Name:       r.Name,
		Email:      r.Email,
		Phone:      r.Phone,
		Items:      r.Items,
		Company:    r.Company,
		StatusDeal: r.StatusDeal,
		Notes:      r.Notes,
		Needs:      r.Needs,
		CreatedAt:  r.CreatedAt,
		UserId:     r.UserId,
	}
}

type DealResponse struct {
	Data  []Deal `json:"data"`
	Count int64  `json:"count"`
}

func DealToResponse(r DealResponse) *DealResponse {
	return &DealResponse{
		Data:  r.Data,
		Count: r.Count,
	}
}
