package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
)

// config 管理所有通用參數，包含各種來源：env, json file等

var _env = &Env{}

func EnvInit() error {
	if err := env.Parse(_env); err != nil {
		return err
	}

	if err := validator.New().Struct(_env); err != nil {
		return err
	}

	return nil
}

func GetEnv() *Env {
	return _env
}
