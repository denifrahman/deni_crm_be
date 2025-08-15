package di

import (
	"deni-be-crm/internal/handlers"
)

type Handlers struct {
	LeadHandler         *handlers.LeadHandler
	AuthHandler         *handlers.AuthHandler
	OrderHandler        *handlers.OrderHandler
	ProductHandler      *handlers.ProductHandler
	CustomerHandler     *handlers.CustomerHandler
	DealHandler         *handlers.DealHandler
	SubscriptionHandler *handlers.SubscriptionHandler
}
