package routers

import (
	"app/api"
	"github.com/gofiber/fiber"
	//"time"
)

func InitServer()  {
	app := fiber.New()

	app.Post("/api/login", api.Login)
	app.Post("/api/checkOpenId", api.CheckOpenId)
	app.Post("/api/getNewestPoints", api.GetNewestPoints)
	app.Post("/api/getLists", api.GetLists)
	app.Post("/api/addPoints", api.AddPoints)
	//app.Get("/api/test", func(c *fiber.Ctx) {
	//	c.Format(time.Now().Format("2006-01-02"))
	//})

	app.Listen(8080)
	//app.Listen(443, "./conf/ssl/app.crt", "./conf/ssl/app.key")//不方便管理,算了
}