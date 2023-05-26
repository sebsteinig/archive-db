package utils

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type ParamValue struct {
	value    interface{}
	operator func(string, interface{}, *Placeholder) string
}
type Params map[string]ParamValue

func (params Params) ParseParams(c *fiber.Ctx, whitelist ...string) error {
	default_operator := func(key string, value interface{}, pl *Placeholder) string {
		return fmt.Sprintf("%s = %s ", key, pl.Get(value))
	}
	log.Default().Println(c.Query("config_name"))
	if value := c.Query("config_name", "##error##"); value != "##error##" && in("config_name", whitelist) {
		params["config_name"] = ParamValue{
			value:    value,
			operator: default_operator,
		}
	}
	if value := c.Query("extension", "##error##"); value != "##error##" && in("extension", whitelist) {
		params["extension"] = ParamValue{
			value:    value,
			operator: default_operator,
		}
	}
	if value := c.Query("lossless", "error"); value != "error" && in("lossless", whitelist) {
		params["lossless"] = ParamValue{
			value:    c.QueryBool("lossless"),
			operator: default_operator,
		}
	}
	if value := c.Query("threshold", "error"); value != "error" && in("threshold", whitelist) {
		params["threshold"] = ParamValue{
			value:    c.QueryFloat("threshold"),
			operator: default_operator,
		}
	}
	if value := c.Query("rx", "error"); value != "error" && in("rx", whitelist) {
		params["rx"] = ParamValue{
			value:    c.QueryFloat("rx"),
			operator: default_operator,
		}
	}
	if value := c.Query("ry", "error"); value != "error" && in("ry", whitelist) {
		params["ry"] = ParamValue{
			value:    c.QueryFloat("ry"),
			operator: default_operator,
		}
	}
	if value := c.Query("chunks", "error"); value != "error" && in("chunks", whitelist) {
		params["chunks"] = ParamValue{
			value:    c.QueryInt("chunks"),
			operator: default_operator,
		}
	}
	if value := c.Query("id", "##error##"); value != "error" && in("id", whitelist) {
		params["id"] = ParamValue{
			value: value,
			operator: func(key string, value interface{}, pl *Placeholder) string {
				return fmt.Sprintf("%s LIKE '%s%%' ", key, pl.Get(value))
			},
		}
	}
	return nil
}

func (params Params) ParamToSql(pl *Placeholder) string {
	res := " "
	for key, value := range params {
		res += fmt.Sprintf("AND %s", value.operator(key, value.value, pl))
	}
	return res
}

func in(value string, arr []string) bool {
	for _, e := range arr {
		if value == e {
			return true
		}
	}
	return false
}
