package helper

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	beego "github.com/beego/beego/v2/server/web"
	beegoContext "github.com/beego/beego/v2/server/web/context"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

func NewMockDB(typesConn string) (*gorm.DB, sqlmock.Sqlmock, error) {

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	if typesConn == "sql" {
		gormDB, err = gorm.Open(sqlserver.New(sqlserver.Config{
			Conn: sqlDB,
		}), &gorm.Config{})
		if err != nil {
			return nil, nil, err
		}
	}

	return gormDB, mock, nil
}

func GetValueAndColumnStructToDriverValue(value interface{}) ([]driver.Value, []string) {
	var result []driver.Value

	object := reflect.ValueOf(value)
	var column []string
	for i := 0; i < object.NumField(); i++ {

		// set data to driver.Value
		result = append(result, object.Field(i).Interface())

		// set column
		splitStringByTag := strings.Split(object.Type().Field(i).Tag.Get("gorm"), ";")
		for i := 0; i < len(splitStringByTag); i++ {
			if strings.Contains(splitStringByTag[i], "column") {
				column = append(column, strings.TrimPrefix(splitStringByTag[i], "column:"))
				break
			}
		}
	}

	return result, column
}

func PrepareHandler(handler *beego.Controller, request *http.Request, response http.ResponseWriter) {
	handler.Ctx = &beegoContext.Context{
		Request: request,
		ResponseWriter: &beegoContext.Response{
			ResponseWriter: response,
		},
	}
	body, _ := ioutil.ReadAll(handler.Ctx.Request.Body)
	handler.Ctx.Input = &beegoContext.BeegoInput{
		Context:     handler.Ctx,
		RequestBody: body,
	}
	handler.Ctx.Output = &beegoContext.BeegoOutput{
		Context: handler.Ctx,
	}
	handler.Data = map[interface{}]interface{}{}
}
