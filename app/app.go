package app

import (
	"api-gateway-module/app/router"
	"context"
	"log"

	"go.uber.org/fx"
)

type App struct {
	router map[string]*router.Router
}

func NewApp(lc fx.Lifecycle, routers map[string]*router.Router) App {
	a := App{router: routers}
	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			for name, r := range routers {
				go func(name string, r *router.Router) {
					if err := r.Run(); err != nil {
						log.Printf("router %s stopped with error: %v", name, err)
					}
				}(name, r)
			}
			return nil
		},
		OnStop: func(c context.Context) error {
			for name, r := range routers {
				if err := r.Shutdown(c); err != nil {
					log.Printf("failed to shutdown router %s: %v", name, err)
				}
			}
			log.Println("lifeCycle ended", c.Err())
			return nil
		},
	})

	return a
}
