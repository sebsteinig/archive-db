package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ParamValue struct {
	Value    interface{}
	operator func(string, interface{}, *Placeholder) string
}
type Params map[string]ParamValue

func (params Params) ParseParams(c *fiber.Ctx, whitelist ...string) error {
	default_operator := func(key string, value interface{}, pl *Placeholder) string {
		return fmt.Sprintf("%s = %s ", key, pl.Get(value))
	}
	if value := c.Query("config_name", "##error##"); value != "##error##" && in("config_name", whitelist) {
		params["config_name"] = ParamValue{
			Value:    value,
			operator: default_operator,
		}
	}
	if value := c.Query("extension", "##error##"); value != "##error##" && in("extension", whitelist) {
		params["extension"] = ParamValue{
			Value:    value,
			operator: default_operator,
		}
	}
	if value := c.Query("lossless", "error"); value != "error" && in("lossless", whitelist) {
		params["lossless"] = ParamValue{
			Value:    c.QueryBool("lossless"),
			operator: default_operator,
		}
	}
	if value := c.Query("threshold", "error"); value != "error" && in("threshold", whitelist) {
		params["threshold"] = ParamValue{
			Value:    c.QueryFloat("threshold"),
			operator: default_operator,
		}
	}
	if value := c.Query("rx", "error"); value != "error" && in("rx", whitelist) {
		params["rx"] = ParamValue{
			Value:    c.QueryFloat("rx"),
			operator: default_operator,
		}
	}
	if value := c.Query("ry", "error"); value != "error" && in("ry", whitelist) {
		params["ry"] = ParamValue{
			Value:    c.QueryFloat("ry"),
			operator: default_operator,
		}
	}
	if value := c.Query("chunks", "error"); value != "error" && in("chunks", whitelist) {
		params["chunks"] = ParamValue{
			Value:    c.QueryInt("chunks"),
			operator: default_operator,
		}
	}
	if value := c.Query("like", "##error##"); value != "##error##" && in("like", whitelist) {
		params["exp_id"] = ParamValue{
			Value: value,
			operator: func(key string, value interface{}, pl *Placeholder) string {
				return fmt.Sprintf("%s LIKE %s ||'%%'", key, pl.Get(value))
			},
		}
	}
	if value := c.Query("for", "##error##"); value != "##error##" && in("for", whitelist) {
		params["query"] = ParamValue{
			Value: strings.Fields(value),
			operator: func(key string, value interface{}, pl *Placeholder) string {
				return ""
			},
		}
	}
	if value := c.Query("with", "##error##"); value != "##error##" && in("with", whitelist) {
		labels, ok := idsToSlice(value)
		if ok {
			params["labels"] = ParamValue{
				Value: labels,
				operator: func(key string, value interface{}, pl *Placeholder) string {
					return ""
				},
			}
		}
	}
	if in("ids", whitelist) {
		if value := c.Query("ids", "###error###"); value != "###error###" {
			ids, ok := idsToSlice(value)
			if ok {
				params["exp_id"] = ParamValue{
					Value: ids,
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

	return nil
}

func (params Params) ParamToSql(pl *Placeholder) string {

	res := make([]string, 0, len(params))
	for key, value := range params {
		res = append(res, fmt.Sprintf("%s", value.operator(key, value.Value, pl)))
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

func idsToSlice(ids string) ([]string, bool) {
	ids, found_prefix := strings.CutPrefix(ids, "[")
	ids, found_suffix := strings.CutSuffix(ids, "]")
	if !found_prefix || !found_suffix {
		log.Default().Println("ids are not specified in the right format", ids)
		return nil, false
	}
	res := strings.Split(ids, ",")
	return res, true
}
