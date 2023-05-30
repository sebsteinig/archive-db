package utils

import (
	"fmt"
	"log"
	"strings"

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
	if value := c.Query("like", "##error##"); value != "error" && in("like", whitelist) {
		params["exp_id"] = ParamValue{
			value: value,
			operator: func(key string, value interface{}, pl *Placeholder) string {
				return fmt.Sprintf("%s LIKE %s ||'%%'", key, pl.Get(value))
			},
		}
	}
	if in("ids", whitelist) {
		if value := c.Query("ids", "###error###"); value != "###error###" {
			ids, ok := itemsToSlice(value)
			if ok {
				params["exp_id"] = ParamValue{
					value: ids,
					operator: func(key string, value interface{}, pl *Placeholder) string {
						exp_ids := value.([]string)
						var tab []string
						for i := 0; i < len(exp_ids); i++ {
							tab = append(tab, pl.Get(exp_ids[i]))
						}
						values := fmt.Sprintf("(%s)", strings.Join(tab, ","))
						return fmt.Sprintf("%s IN %s ", key, values)
					},
				}
			} else {
				return fmt.Errorf("ids not specified")
			}
		} else {
			return fmt.Errorf("ids not specified")
		}
	}

	if value := c.Query("variables", "###error###"); value != "###error###" && in("variables", whitelist) {
		variables, ok := itemsToSlice(value)
		if ok {
			params["variables"] = ParamValue{
				value: variables,
				operator: func(key string, value interface{}, pl *Placeholder) string {
					variables := value.([]string)
					var tab []string
					for i := 0; i < len(variables); i++ {
						tab = append(tab, pl.Get(variables[i]))
					}
					values := fmt.Sprintf("(%s)", strings.Join(tab, ","))
					return fmt.Sprintf("variable_name IN %s", values)
				},
			}
		} else {
			return fmt.Errorf("variables not correctly specifed")
		}
	}
	return nil
}

func (params Params) ParamToSql(pl *Placeholder) string {

	res := make([]string, 0, len(params))
	for key, value := range params {
		res = append(res, fmt.Sprintf("%s", value.operator(key, value.value, pl)))
	}
	return strings.Join(res, " AND ")
}

func in(value string, arr []string) bool {
	for _, e := range arr {
		if value == e {
			return true
		}
	}
	return false
}

func itemsToSlice(items string) ([]string, bool) {
	items, found_prefix := strings.CutPrefix(items, "[")
	items, found_suffix := strings.CutSuffix(items, "]")
	if !found_prefix || !found_suffix {
		log.Default().Println("items are not specified in the right format", items)
		return nil, false
	}
	items = strings.Trim(items, " ")
	res := strings.Split(items, ",")
	return res, true
}
