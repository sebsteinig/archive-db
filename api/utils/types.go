package utils

import "time"

type Variable struct {
	Id                 int       `json:"id"`
	Name               string    `json:"name"`
	Exp_id             string    `json:"exp_id"`
	Paths              []string  `json:"paths"`
	Created_at         time.Time `json:"created_at"`
	Config_name        string    `json:"config_name"`
	Levels             int       `json:"levels"`
	Timesteps          int       `json:"timesteps"`
	Xsize              int       `json:"xsize"`
	Xfirst             float32   `json:"xfirst"`
	Xinc               float32   `json:"xinc"`
	Ysize              int       `json:"ysize"`
	Yfirst             float32   `json:"yfirst"`
	Yinc               float32   `json:"yinc"`
	Extension          string    `json:"extension"`
	Lossless           bool      `json:"lossless"`
	Nan_value_encoding int       `json:"nan_value_encoding"`
	Threshold          float32   `json:"threshold"`
	Chunks             int       `json:"chunks"`
	Metadata           string    `json:"metadata"`
}
