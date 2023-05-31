package utils

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TableExperiment struct {
	Exp_id   string                 `json:"exp_id"`
	Labels   []string               `json:"labels"`
	Co2      float64                `json:"co2"`
	Metadata map[string]interface{} `json:"metadata"`
}

type ExperimentJSON struct {
	Exp_id   string                 `json:"exp_id" validate:"required"`
	Labels   []string               `json:"labels"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (exp ExperimentJSON) ToTable() TableExperiment {
	table_experiment := TableExperiment{
		Exp_id: exp.Exp_id,
		Labels: exp.Labels,
		//...
		Metadata: make(map[string]interface{}),
	}
	if co2_str, ok := exp.Metadata["co2"]; ok {
		if co2, err := strconv.ParseFloat(co2_str.(string), 32); err == nil {
			table_experiment.Co2 = co2
			delete(exp.Metadata, "co2")
		}
	}
	// delete already used keys

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
	Exp_id             string  `json:"exp_id" validate:"required"`
	Config_name        string  `json:"config_name" validate:"required"`
	Created_at         string  `json:"created_at"`
	Extension          string  `json:"extension" validate:"required"`
	Lossless           bool    `json:"lossless" validate:"required"`
	Nan_value_encoding int     `json:"nan_value_encoding" validate:"required"`
	Threshold          float32 `json:"threshold" validate:"required"`
	Chunks             int     `json:"chunks" validate:"required,gte=0"`
	Rx                 float64 `json:"rx" validate:""`
	Ry                 float64 `json:"ry" validate:""`
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
	Name       string                 `json:"name" validate:"required"`
	Paths_ts   []string               `json:"paths_ts" validate:"required,filepath"`
	Paths_mean []string               `json:"paths_mean" validate:"required,filepath"`
	Levels     int                    `json:"levels" validate:"required,gte=0"`
	Timesteps  int                    `json:"timesteps" validate:"required,gte=0"`
	Xsize      int                    `json:"xsize" validate:"required,gte=0"`
	Xfirst     float32                `json:"xfirst" validate:"required"`
	Xinc       float32                `json:"xinc" validate:"required"`
	Ysize      int                    `json:"ysize" validate:"required,gte=0"`
	Yfirst     float32                `json:"yfirst" validate:"required"`
	Yinc       float32                `json:"yinc" validate:"required"`
	Metadata   map[string]interface{} `json:"metadata" validate:"required,json"`
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
