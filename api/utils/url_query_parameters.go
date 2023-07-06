package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const error_param = "###error###"

func setArray[T any](value string, field reflect.Value) error {
	a := []T{}
	err := json.Unmarshal([]byte(value), &a)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(a))
	return nil
}

type QueryParameters map[string]any

func BuildQueryParameters(c *fiber.Ctx, params_struct any) (QueryParameters, error) {
	elements := reflect.ValueOf(params_struct).Elem()

	query_parameters := make(QueryParameters)

	// keys := make(map[string]string)
	// for j := 0; j < elements.NumField(); j++ {
	// 	key, ok := elements.Type().Field(j).Tag.Lookup("param")
	// 	if ok {
	// 		keys[key] = ""
	// 	}
	// }
	// correct_query, wrong_key := true, ""
	// c.Context().QueryArgs().VisitAll(func(key, value []byte) {
	// 	_, ok := keys[string(key)]
	// 	if !ok {
	// 		correct_query, wrong_key = false, string(key)
	// 	}
	// })
	// if !correct_query {
	// 	return nil, fmt.Errorf("wrong key %s in query parameters", wrong_key)
	// }

	for i := 0; i < elements.NumField(); i++ {
		tag := elements.Type().Field(i).Tag
		key, ok := tag.Lookup("param")
		if !ok {
			continue
		}
		required := false
		if strings.Contains(key, "required") {
			key = strings.Replace(key, "required", "", -1)
			required = true
		}
		key = strings.Replace(key, ",", "", -1)
		value := c.Query(key, error_param)
		if value == error_param && required {
			return nil, fmt.Errorf("Parameter %s is required", key)
		}
		if value != error_param && elements.Field(i).CanSet() {
			interface_type := elements.Field(i).Interface()

			switch interface_type.(type) {
			case int:
				if number, err := strconv.ParseInt(value, 10, 64); err == nil {
					elements.Field(i).SetInt(number)
				} else {
					return nil, fmt.Errorf("Field(%s) must be of type int64", elements.Type().Field(i).Name)
				}
			case string:
				value = strings.TrimPrefix(value, "\"")
				value = strings.TrimSuffix(value, "\"")
				elements.Field(i).SetString(value)
			case float64:
				if number, err := strconv.ParseFloat(value, 64); err == nil {
					elements.Field(i).SetFloat(number)
				} else {
					return nil, fmt.Errorf("Field(%s) must be of type float64", elements.Type().Field(i).Name)
				}
			case bool:
				if v, err := strconv.ParseBool(value); err == nil {
					elements.Field(i).SetBool(v)
				} else {
					return nil, fmt.Errorf("Field(%s) must be of type bool", elements.Type().Field(i).Name)
				}
			case []string:
				err := setArray[string](value, elements.Field(i))
				if err != nil {
					return nil, err
				}
			case []int:
				err := setArray[int](value, elements.Field(i))
				if err != nil {
					return nil, err
				}
			case []float32:
				err := setArray[float32](value, elements.Field(i))
				if err != nil {
					return nil, err
				}
			case []float64:
				err := setArray[float64](value, elements.Field(i))
				if err != nil {
					return nil, err
				}
			case []bool:
				err := setArray[bool](value, elements.Field(i))
				if err != nil {
					return nil, err
				}
			case []byte:
				err := setArray[byte](value, elements.Field(i))
				if err != nil {
					return nil, err
				}
			case []any:
				err := setArray[any](value, elements.Field(i))
				if err != nil {
					return nil, err
				}

			}
			query_parameters[elements.Type().Field(i).Name] = elements.Field(i).Interface()
		}
	}
	return query_parameters, nil
}
