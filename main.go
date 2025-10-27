package main

import (
	"api-gateway-module/app"
	"api-gateway-module/app/dependency"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		dependency.Cfg,
		dependency.HttpClient,
		dependency.Producer,
		dependency.Router,
		fx.Provide(app.NewApp),
		fx.Invoke(func(app.App) {}),
	).Run()
}
