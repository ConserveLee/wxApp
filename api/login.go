package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber"
)

const (
	SuccessCode   = 0
	FailCode      = -1
	SuccessMsg    = "ok"
	MissingParams = "参数缺失"
)

var (
	appId         = os.Getenv("APP_ID")
	secret        = os.Getenv("APP_SECRET")
)

func Login(c *fiber.Ctx) {
	code := GetJson(c, "code")
	if code == nil {
		c.JSON(fiber.Map{
			"code": FailCode,
			"msg":	fmt.Sprintf("%s%s", "code", MissingParams),
		})
		return
	}

	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appId, secret, code)
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return
	}
	result := map[string]interface{}{}
	json.NewDecoder(resp.Body).Decode(&result)
	if errMsg, ok := result["errmsg"]; ok {
		c.JSON(fiber.Map{
			"code": FailCode,
			"msg":	errMsg,
		})
		return
	} else {
		openId := result["openid"].(string)
		userId := getUserIdByOpenId(appId, openId)
		c.JSON(fiber.Map{
			"code": SuccessCode,
			"msg":	SuccessMsg,
			"data": map[string]interface{} {
				"userId": userId,
				"openId": openId,
			},
		})
		return
	}
}