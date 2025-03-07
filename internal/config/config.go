package config

import "github.com/caarlos0/env"

type Config struct {
	HTTPAddress              string `json:"http_address" env:"GOPHERMART_PORTAL_HTTP_ADDRESS" envDefault:"0.0.0.0:8000"`
	BillingOrdersAddress     string `json:"billing_orders_address" env:"GOPHERMART_PORTALBILLING_ORDERS_ADDRESS" envDefault:"127.0.0.1:8010"`
	BillingAccountingAddress string `json:"billing_accounting_address" env:"GOPHERMART_PORTALBILLING_ACCOUNTING_ADDRESS" envDefault:"127.0.0.1:8020"`
	QueryOrdersAddress       string `json:"query_orders_address" env:"GOPHERMART_PORTALQUERY_ORDERS_ADDRESS" envDefault:"127.0.0.1:8030"`
	QueryAccountingAddress   string `json:"query_accounting_address" env:"GOPHERMART_PORTALQUERY_ACCOUNTING_ADDRESS" envDefault:"127.0.0.1:8040"`
	LogLevel                 int    `json:"log_level" env:"GOPHERMART_PORTAL_LOG_LEVEL" envDefault:"-1"`
	DatabaseDSN              string `json:"database_dsn" env:"GOPHERMART_PORTAL_DATABASE_DSN" envDefault:"postgres://postgres:secret@127.0.0.1:5432/gophermart_portal_development"`
	IamTokenTTL              int64  `json:"iam_token_ttl" env:"GOPHERMART_PORTAL_IAM_TOKEN_TTL" envDefault:"60"`
	IamSecretKey             string `json:"iam_secret_key" env:"GOPHERMART_PORTAL_IAM_SECRET_KEY" envDefault:"secret"`
}

func MustNewConfig() *Config {
	c := new(Config)
	env.Parse(c)

	return c
}
