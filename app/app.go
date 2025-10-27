package app

import (
	"api-gateway-module/app/router"
	"context"
	"log"

	"go.uber.org/fx"
)

type App struct {
	router map[string]router.Router
}

func NewApp(lc fx.Lifecycle, router map[string]router.Router) App {
	a := App{router: router}
	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			for _, r := range router {
				if err := r.Run(); err != nil {
					panic(err.Error())
				}
			}
			return nil
		},
		OnStop: func(c context.Context) error {
			log.Println("lifeCycle ended", c.Err())
			return nil
		},
	})

	return a
}
