package vela_common_mb

import (
	_ "gitee.com/opengauss/openGauss-connector-go-pq"
	_ "github.com/go-playground/locales"
	_ "github.com/go-playground/universal-translator"
	_ "github.com/go-playground/validator/v10"
	_ "github.com/vela-ssoc/vela-common-mba"
	_ "github.com/xgfone/ship/v5"
	_ "go.uber.org/zap"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gen"
	_ "gorm.io/gorm"
	_ "gorm.io/plugin/dbresolver"
)
