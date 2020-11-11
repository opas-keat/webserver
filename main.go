package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func main() {
	app := fiber.New()

	app.Use("/beisApi", proxy.Balancer(proxy.Config{
		Servers: []string{
			"http://10.5.46.116:8002",
		},
		ModifyRequest: func(c *fiber.Ctx) error {
			c.Set("X-Real-IP", c.IP())
			return c.Redirect(c.OriginalURL())
		},
	}))

	app.Static("/", "./public", fiber.Static{
		Compress:  true, // default: false
		ByteRange: true, // default: false
	})

	app.Get("/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./public/index.html")
	})

	log.Fatal(app.Listen(":3000"))

}
