package api

import (
	"app/database"
	"fmt"
	"time"

	qb "github.com/didi/gendry/builder"
	sc "github.com/didi/gendry/scanner"
	"github.com/gofiber/fiber"
)

const (
	table            = "users"
	IllegalMsg       = "非法ID"
	IllegalId        = 0
	FailSelect		 = "查询失败"
)


type User struct{
	UserId 		int 	`ddb:"id"`
	AppId 		string 	`ddb:"app_id"`
	OpenId 		string 	`ddb:"user_open_id"`
	HasRule		bool 	`ddb:"has_rule"`
}

func getUserIdByOpenId(appId, openId string) int {
	db := database.Db

	//查询逻辑
	where := map[string]interface{} {
		"app_id": 		appId,
		"user_open_id": openId,
		"_limit": 		[]uint{1},
	}
	cond, vals, err := qb.BuildSelect(table, where, []string{"id"})
	if err != nil {
		return IllegalId
	}
	rows, err := db.Query(cond, vals...)
	if err != nil {
		return IllegalId
	}
	defer rows.Close()

	//创建user实例
	var user User
	_ = sc.Scan(rows, &user)
	if user.UserId == 0 {
		//新建user
		var data []map[string]interface{}
		data = append(data, map[string]interface{}{
			"app_id": 		appId,
			"user_open_id": openId,
			"has_rule": 	0,
			"created_at": 	time.Now().Format("2006-01-02 15:04:05"),
		})
		cond, vals, err := qb.BuildInsert(table, data)
		result, err := db.Exec(cond, vals...)
		if nil != err || result == nil {
			return IllegalId
		}
		id, _ := result.LastInsertId()
		user.UserId = int(id)
	}

	return user.UserId

}

func CheckOpenId(c *fiber.Ctx)  {
	openId := GetJson(c, "openId")
	if openId == nil {
		c.JSON(fiber.Map{
			"code": FailCode,
			"msg":	fmt.Sprintf("%s%s", "openId", MissingParams),
		})
		return
	}

	isIllegal, msg := checkOpenIdFormDb(appId, openId.(string))
	if isIllegal {
		c.JSON(fiber.Map{
			"code": SuccessCode,
			"msg":	SuccessMsg,
		})
		return
	} else {
		c.JSON(fiber.Map{
			"code": FailCode,
			"msg":	msg,
		})
		return
	}
}

func checkOpenIdFormDb(appId, openId string) (bool, string) {
	db := database.Db

	where := map[string]interface{} {
		"app_id": 		appId,
		"user_open_id": openId,
		"has_rule": 	1,
		"_limit": 		[]uint{1},
	}
	cond, vals, err := qb.BuildSelect(table, where, []string{"id"})
	if err != nil {
		return false, err.Error()
	}
	rows, err := db.Query(cond, vals...)

	if err != nil {
		return false, err.Error()
	}
	defer rows.Close()

	//创建user实例
	var user User
	_ = sc.Scan(rows, &user)
	if user.UserId > 0 {
		return true, ""
	}

	return false, FailSelect
}