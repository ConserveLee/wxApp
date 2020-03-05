package api

import (
	"time"

	qb "github.com/didi/gendry/builder"
	sc "github.com/didi/gendry/scanner"
	"github.com/gofiber/fiber"
)

const (
	tableAll   = "points"
	tableToday = "points_flow"
	noSession  = "登录失效"
)

type All struct {
	AppId   string `ddb:"app_id"`
	Points  int    `ddb:"points"`
	UseTime int    `ddb:"use_time"`
}

type Today struct{
	Id 		int		`ddb:"id"`
	AppId 	string 	`ddb:"app_id"`
	Points	int		`ddb:"points"`
	DayTime	string	`ddb:"day_time"`
}

var dayTime = time.Now().Format("2006-01-02")

func GetNewestPoints(c *fiber.Ctx)  {
	userId := GetJson(c, "userId")
	if userId == nil {
		c.JSON(fiber.Map{
			"code": FailCode,
			"msg":	noSession,
		})
		return
	}
	all, today := getNewestPointsAllFormDb(), getNewestPointsTodayFormDb()
	c.JSON(fiber.Map{
		"code": SuccessCode,
		"msg":	SuccessMsg,
		"data": map[string]interface{} {
			"all": 		all,
			"today": 	today,
		},
	})
	return
}

func getNewestPointsAllFormDb() int {
	//查询逻辑
	where := map[string]interface{} {
		"app_id": 		appId,
		"_limit": 		[]uint{1},
	}
	cond, vals, err := qb.BuildSelect(tableAll, where, []string{"points"})
	if err != nil {
		return 0
	}
	rows, err := db.Query(cond, vals...)
	defer rows.Close()
	if err != nil {
		return 0
	}

	//创建字段实例
	var all All
	sc.Scan(rows, &all)
	return all.Points
}

func getNewestPointsTodayFormDb() int {

	//查询逻辑
	where := map[string]interface{} {
		"app_id": 		appId,
		"day_time": 	dayTime,
		"_limit": 		[]uint{1},
	}
	cond, vals, err := qb.BuildSelect(tableToday, where, []string{"points"})
	if err != nil {
		return 0
	}
	rows, err := db.Query(cond, vals...)

	if err != nil {
		return 0
	}
	defer rows.Close()

	//创建字段实例
	var today Today
	sc.Scan(rows, &today)
	return today.Points
}

func AddPoints(c *fiber.Ctx) {
	points := GetJson(c, "points")
	if points == nil {
		c.JSON(fiber.Map{
			"code": FailCode,
			"msg":	MissingParams,
		})
		return
	}

	ret, msg := addPointsFromDb(int(points.(float64)))
	if msg == "" {
		c.JSON(fiber.Map{
			"code": SuccessCode,
			"msg":	SuccessMsg,
			"data": ret,
		})
		return
	} else {
		c.JSON(fiber.Map{
			"code": FailCode,
			"msg":	msg,
		})
	}
}

func addPointsFromDb(points int) (map[string]interface{}, string) {
	where := map[string]interface{} {
		"app_id": 		appId,
		"day_time": 	dayTime,
	}

	//先查有没有
	//没有就insert
	//有就update
	cond, vals, err := qb.BuildSelect(tableToday, where, []string{"id, points"})
	if err != nil {
		return nil, err.Error()
	}
	rows, err := db.Query(cond, vals...)

	if err != nil {
		return nil, err.Error()
	}
	defer rows.Close()

	//创建字段实例
	var today Today
	sc.Scan(rows, &today)
	ret := make(map[string]interface{})
	if today.Id > 0 {
		todayPoint   := today.Points+points
		ret["today"] = todayPoint
		//update
		update := map[string]interface{}{
			"points": todayPoint,
		}
		cond, vals, err := qb.BuildUpdate(tableToday, where, update)
		if err != nil {
			return nil, err.Error()
		}
		db.Exec(cond, vals...)//这里有错。。。
	}
	if today.Id == 0 {
		ret["today"] = points
		//insert
		var data []map[string]interface{}
		data = append(data, map[string]interface{}{
			"app_id": 		appId,
			"points": 		points,
			"day_time":		dayTime,
			"created_at": 	time.Now().Format("2006-01-02 15:04:05"),
		})
		cond, vals, err := qb.BuildInsert(tableToday, data)
		_, err = db.Exec(cond, vals...)
		if err != nil {
			return nil, err.Error()
		}
	}
	//再查总
	eldAll := getNewestPointsAllFormDb()
	//再加到总表里
	ret["all"] = points+eldAll
	updateAll := map[string]interface{}{
		"app_id": appId,
		"points": ret["all"],
	}
	cond, vals, err = qb.BuildUpdate(
		tableAll,
		map[string]interface{} {
		"app_id": 		appId,
		},
		updateAll)
	if err != nil {
		return nil, err.Error()
	}

	db.Exec(cond, vals...)
	return ret, ""
}