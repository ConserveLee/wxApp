package api

import (
	"app/database"

	qb "github.com/didi/gendry/builder"
	sc "github.com/didi/gendry/scanner"
	"github.com/gofiber/fiber"
)

const tableHouseWork = "houseworks"

var (
	db 	  = database.Db
	where = map[string]interface{} {
		"app_id": appId,
	}
)

type Homeworks struct{
	Id  	int		`ddb:"id"`
	AppId 	string 	`ddb:"app_id"`
	Points	int		`ddb:"points"`
	Title	string	`ddb:"title"`
	Image	string	`ddb:"img"`
}

func GetLists(c *fiber.Ctx)  {
	lists, msg := getListsFromDb()
	if msg == "" {
		c.JSON(fiber.Map{
			"code": SuccessCode,
			"msg":	SuccessMsg,
			"data": lists,
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

func getListsFromDb() (hsSet []Homeworks,msg string) {
	cond, vals, err := qb.BuildSelect(tableHouseWork, where, []string{"id", "app_id", "points", "title", "img"})
	if err != nil {
		return nil, err.Error()
	}
	rows, err := db.Query(cond, vals...)
	if err != nil {
		return nil, err.Error()
	}
	defer rows.Close()
	err = sc.Scan(rows, &hsSet)
	if err != nil {
		return nil, err.Error()
	}

	return hsSet, ""
}


//暂时不做
//func addItem()  {
//
//}
//
//func deleteItem()  {
//
//}