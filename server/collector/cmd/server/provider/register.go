package provider

import "github.com/samber/do/v2"

// Register registers all collector service providers with the injector.
func Register(injector do.Injector) {
	do.Provide(injector, ProvideEntClient)
	do.Provide(injector, ProvideFeedRepository)
	do.Provide(injector, ProvideAddFeedUsecase)
	do.Provide(injector, ProvideServer)
	do.Provide(injector, ProvideGRPCServer)
}
