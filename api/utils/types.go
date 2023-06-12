package utils

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TableExperiment struct {
	Exp_id        string                 `json:"exp_id" sql:"exp_id"`
	Labels        []string               `json:"labels"`
	Co2           float64                `json:"co2" sql:"co2"`
	Coast_Line_id int64                  `json:"coast_line_id" sql:"coast_line_id"`
	Gmst          float64                `json:"gmst" sql:"gmst"`
	Realistic     bool                   `json:"realistic" sql:"realistic"`
	Date_created  string                 `json:"date_created" sql:"date_wp_created"`
	Date_updated  string                 `json:"date_updated" sql:"date_wp_updated"`
	Metadata      map[string]interface{} `json:"metadata"`
}

type ExperimentJSON struct {
	Exp_id   string                 `json:"exp_id" validate:"required"`
	Labels   []string               `json:"labels"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (exp ExperimentJSON) ToTable() TableExperiment {
	table_experiment := TableExperiment{
		Exp_id:   exp.Exp_id,
		Labels:   exp.Labels,
		Metadata: make(map[string]interface{}),
	}
	if co2_str, ok := exp.Metadata["co2"]; ok {
		if co2, err := strconv.ParseFloat(co2_str.(string), 64); err == nil {
			table_experiment.Co2 = co2
			delete(exp.Metadata, "co2")
		}
	}
	if coast_line_id, ok := exp.Metadata["coast"]; ok {
		switch coast_line_id.(type) {
		case string:
			if coast_line_id, err := strconv.ParseInt(coast_line_id.(string), 10, 64); err == nil {
				table_experiment.Coast_Line_id = coast_line_id
				delete(exp.Metadata, "coast")
			}
		case int64:
			table_experiment.Coast_Line_id = coast_line_id.(int64)
			delete(exp.Metadata, "coast")
		}
	}
	if gmst, ok := exp.Metadata["gmst"]; ok {
		switch gmst.(type) {
		case string:
			if gmst, err := strconv.ParseFloat(gmst.(string), 64); err == nil {
				table_experiment.Gmst = gmst
				delete(exp.Metadata, "gmst")
			}
		case float64:
			table_experiment.Gmst = gmst.(float64)
			delete(exp.Metadata, "coast")
		}
	}
	if date_created, ok := exp.Metadata["date_original"]; ok {
		table_experiment.Date_created = date_created.(string)
		delete(exp.Metadata, "date_original")
	}
	if date_updated, ok := exp.Metadata["date_modified"]; ok {
		table_experiment.Date_updated = date_updated.(string)
		delete(exp.Metadata, "date_modified")
	}
	if realistic, ok := exp.Metadata["realistic"]; ok {
		switch realistic.(type) {
		case string:
			if realistic, err := strconv.ParseBool(realistic.(string)); err == nil {
				table_experiment.Realistic = realistic
				delete(exp.Metadata, "gmst")
			}
		case bool:
			table_experiment.Realistic = realistic.(bool)
			delete(exp.Metadata, "coast")
		}
	}

	for k, v := range exp.Metadata {
		table_experiment.Metadata[k] = v
	}
	return table_experiment
}

var validate = validator.New()

func validateStruct(obj interface{}) []fiber.Map {
	var errors []fiber.Map
	err := validate.Struct(obj)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			err_string := fmt.Sprint(err.StructNamespace(), "must be :", err.Tag())
			errors = append(errors, fiber.Map{
				err.StructNamespace(): err_string,
			})
		}
	}
	return errors
}

func (exp ExperimentJSON) Validate() (error, []fiber.Map) {
	errors := validateStruct(exp)
	if len(errors) > 0 {
		return fmt.Errorf("experiment validation error"), errors
	}
	return nil, errors
}

type NimbusExecution struct {
	Id                 int     `json:"id"`
	Exp_id             string  `json:"exp_id" validate:"required" sql:"exp_id"`
	Config_name        string  `json:"config_name" validate:"required" sql:"config_name"`
	Created_at         string  `json:"created_at" sql:"created_at"`
	Extension          string  `json:"extension" validate:"required" sql:"extension"`
	Lossless           bool    `json:"lossless" validate:"required" sql:"lossless"`
	Nan_value_encoding int     `json:"nan_value_encoding" validate:"required" sql:"nan_value_encoding"`
	Threshold          float32 `json:"threshold" validate:"required" sql:"threshold"`
	Chunks             int     `json:"chunks" validate:"required,gte=0" sql:"chunks"`
	Rx                 float64 `json:"rx" validate:"" sql:"rx"`
	Ry                 float64 `json:"ry" validate:"" sql:"ry"`
}

func (exp NimbusExecution) Validate() (error, []fiber.Map) {
	errors := validateStruct(exp)
	if len(errors) > 0 {
		return fmt.Errorf("NimbusExecution validation error"), errors
	}
	return nil, errors
}

type Variable struct {
	Id         int                    `json:"id"`
	Name       string                 `json:"name" validate:"required" sql:"name"`
	Paths_ts   []string               `json:"paths_ts" validate:"required,filepath" sql:"paths_ts"`
	Paths_mean []string               `json:"paths_mean" validate:"required,filepath" sql:"paths_mean"`
	Levels     int                    `json:"levels" validate:"required,gte=0" sql:"levels"`
	Timesteps  int                    `json:"timesteps" validate:"required,gte=0" sql:"timesteps"`
	Xsize      int                    `json:"xsize" validate:"required,gte=0" sql:"xsize"`
	Xfirst     float32                `json:"xfirst" validate:"required" sql:"xfirst"`
	Xinc       float32                `json:"xinc" validate:"required" sql:"xinc"`
	Ysize      int                    `json:"ysize" validate:"required,gte=0" sql:"ysize"`
	Yfirst     float32                `json:"yfirst" validate:"required" sql:"yfirst"`
	Yinc       float32                `json:"yinc" validate:"required" sql:"yinc"`
	Metadata   map[string]interface{} `json:"metadata" validate:"required,json" sql:"metadata"`
}

func (variable Variable) Validate() (error, []fiber.Map) {
	errors := validateStruct(variable)
	if len(errors) > 0 {
		return fmt.Errorf("variable validation error"), errors
	}
	return nil, errors
}

type RequestBody struct {
	Table_nimbus_execution NimbusExecution `json:"table_nimbus_execution"`
	Table_variable         []Variable      `json:"table_variable"`
	Table_experiment       TableExperiment `json:"-"`
	ExperimentJSON         ExperimentJSON  `json:"exp_metadata"`
}

type Request struct {
	Request RequestBody `json:"request"`
}
