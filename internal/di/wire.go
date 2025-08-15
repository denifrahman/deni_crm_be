//go:build wireinject
// +build wireinject

package di

import (
	"deni-be-crm/internal/handlers"
	"deni-be-crm/internal/repositories"
	"deni-be-crm/internal/services"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(
	repositories.NewLeadsRepository,
	services.NewLeadsService,
	handlers.NewLeadHandler,

	repositories.NewUserRepository,
	services.NewAuthService,
	handlers.NewAuthHandler,

	repositories.NewOrdersRepository,
	services.NewOrdersService,
	handlers.NewOrderHandler,

	repositories.NewProductsRepository,
	services.NewProductsService,
	handlers.NewProductHandler,

	repositories.NewSubscriptionsRepository,
	services.NewSubscriptionsService,
	handlers.NewSubscriptionHandler,

	repositories.NewCustomersRepository,
	services.NewCustomersService,
	handlers.NewCustomerHandler,

	repositories.NewDealsRepository,
	services.NewDealsService,
	handlers.NewDealHandler,
)

func InitHandlers(db *gorm.DB) *Handlers {
	wire.Build(ProviderSet, wire.Struct(new(Handlers), "*"))
	return nil
}
