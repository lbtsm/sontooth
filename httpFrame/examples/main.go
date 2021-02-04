package main

import (
	"log"
	"net/http"
	"path"

	"github.com/lbtsm/gee"
)

func main() {
	log.Println(path.Join("/", "/"))
	g := gee.New()
	g.Use(func(ctx *gee.Context) {
		log.Println("root middle begin")
		ctx.Next()
		log.Println("root middle end")
	})
	g.Get("/", func(c *gee.Context) {
		c.Status(http.StatusOK)
		_, _ = c.String("get ok")
	})
	g.Post("/hello", func(c *gee.Context) {
		c.Status(http.StatusOK)
		_, _ = c.String("post ok")
	})
	group := g.Group("v1")
	group.Get("/api/hello", func(ctx *gee.Context) {
		_, _ = ctx.String("hello")
	})
	group.Use(func(ctx *gee.Context) {
		log.Println("group middle begin")
		ctx.Next()
		log.Println("group middle end")
	})

	err := g.Run(":8081")
	if err != nil {
		log.Println("run err", err)
	}
}
