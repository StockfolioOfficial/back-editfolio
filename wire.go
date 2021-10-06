//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/stockfolioofficial/back-editfolio/core/app"
	"github.com/stockfolioofficial/back-editfolio/core/di"
)

// getApp returns a real app.
func getApp() app.App {
	wire.Build(di.DI)
	return nil
}