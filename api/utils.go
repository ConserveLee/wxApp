package api

import (
	"encoding/json"
	"github.com/gofiber/fiber"
)

func GetJson(c *fiber.Ctx, param string) interface{} {
	params := map[string]interface{}{}
	json.Unmarshal(c.Fasthttp.PostBody(), &params)
	return params[param]
}