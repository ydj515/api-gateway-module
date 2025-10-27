package router

import (
	"api-gateway-module/app/client"
	"api-gateway-module/config"
	"api-gateway-module/types/http"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
)

type Router struct {
	port string
	cfg  config.App

	engine *fiber.App

	client client.HttpClient
}

func NewRouter(cfg config.App, clients map[string]client.HttpClient) Router {
	r := Router{
		cfg:    cfg,
		port:   fmt.Sprintf(":%s", cfg.App.Port),
		client: clients[cfg.App.Name],
	}

	r.engine = fiber.New()
	r.engine.Use(recover2.New())
	r.engine.Use(cors.New(cors.Config{
		AllowMethods: strings.Join([]string{"GET", "POST", "PUT", "DELETE"}, ","),
		//AllowOrigins:
		//AllowMethods:
		//MaxAge: 86400,
	}))

	for _, v := range cfg.Http.Router {
		r.registerRouter(v)
	}

	return r
}

func (r Router) registerRouter(v config.Router) {
	switch v.Method {
	case http.GET:
		handler := AddGet(v, r.client)
		r.engine.Get(v.Path, handler)
	case http.POST:
		handler := AddPost(v, r.client)
		r.engine.Post(v.Path, handler)
	case http.DELETE:
		handler := AddDelete(v, r.client)
		r.engine.Delete(v.Path, handler)
	case http.PUT:
		handler := AddPut(v, r.client)
		r.engine.Put(v.Path, handler)
	default:
		panic("Failed to find router method")
	}
}

func (r Router) Run() error {
	return r.engine.Listen(r.port)
}
