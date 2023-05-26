package utils

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Params map[string]interface{}

func (params Params) ParseParams(c *fiber.Ctx) error {
	log.Default().Println(c.Query("config_name"))
	if value := c.Query("config_name", "##error##"); value != "##error##" {
		params["config_name"] = value
	}
	if value := c.Query("extension", "##error##"); value != "##error##" {
		params["extension"] = value
	}
	if value := c.Query("lossless", "error"); value != "error" {
		params["lossless"] = c.QueryBool("lossless")
	}
	if value := c.Query("threshold", "error"); value != "error" {
		params["threshold"] = c.QueryFloat("threshold")
	}
	if value := c.Query("rx", "error"); value != "error" {
		params["rx"] = c.QueryFloat("rx")
	}
	if value := c.Query("ry", "error"); value != "error" {
		params["ry"] = c.QueryFloat("ry")
	}
	if value := c.Query("chunks", "error"); value != "error" {
		params["chunks"] = c.QueryInt("chunks")
	}

	return nil
}

func (params Params) ParamToSql(pl *Placeholder) string {
	res := " "
	for key, value := range params {
		res += fmt.Sprintf("AND %s = %s ", key, pl.Get(value))
	}
	return res
}
