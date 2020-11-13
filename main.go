package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/helmet/v2"
)

func main() {
	app := fiber.New()

	app.Use(helmet.New())
	// app.Use(csrf.New())

	// Provide a custom compression level
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	// Or extend your config for customization
	app.Use(cache.New(cache.Config{
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))

	app.Use("/pxapi_ect", proxy.Balancer(proxy.Config{
		Servers: []string{
			"http://172.17.8.88:8080",
			"http://172.17.8.88:8081",
			"http://172.17.8.88:8082",
			"http://172.17.8.88:8083",
		},
		ModifyRequest: func(c *fiber.Ctx) error {
			c.Set("X-Real-IP", c.IP())
			return c.Redirect(c.OriginalURL())
		},
	}))

	app.Use("/report", proxy.Balancer(proxy.Config{
		Servers: []string{
			"http://172.17.8.87:8080/jasperserver/",
		},
		ModifyRequest: func(c *fiber.Ctx) error {
			c.Set("X-Real-IP", c.IP())
			return c.Redirect(c.OriginalURL())
		},
	}))

	app.Static("/document", "D:\\dbPraxticol\\Data\\Document", fiber.Static{
		Compress:  false, // default: false
		ByteRange: false, // default: false
	})

	app.Static("/", "./public", fiber.Static{
		Compress:  true, // default: false
		ByteRange: true, // default: false
	})

	app.Get("/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./public/index.html")
	})

	log.Fatal(app.Listen(":80"))

}
