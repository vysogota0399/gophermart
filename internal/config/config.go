package config

import "github.com/caarlos0/env"

type Config struct {
	HTTPAddress              string `json:"http_address" env:"HTTP_ADDRESS" envDefault:"0.0.0.0:8000"`
	BillingOrdersAddress     string `json:"billing_orders_address" env:"BILLING_ORDERS_ADDRESS" envDefault:"127.0.0.1:8010"`
	BillingAccountingAddress string `json:"billing_accounting_address" env:"BILLING_ACCOUNTING_ADDRESS" envDefault:"127.0.0.1:8020"`
	QueryOrdersAddress       string `json:"query_orders_address" env:"QUERY_ORDERS_ADDRESS" envDefault:"127.0.0.1:8030"`
	QueryAccountingAddress   string `json:"query_accounting_address" env:"QUERY_ACCOUNTING_ADDRESS" envDefault:"127.0.0.1:8040"`
	LogLevel                 int    `json:"log_level" env:"LOG_LEVEL" envDefault:"-1"`
	DatabaseDSN              string `json:"database_dsn" env:"DATABASE_DSN" envDefault:"postgres://postgres:secret@127.0.0.1:5432/gophermart_portal_development"`
	IamTokenTTL              int64  `json:"iam_token_ttl" env:"IAM_TOKEN_TTL" envDefault:"60"`
	IamSecretKey             string `json:"iam_secret_key" env:"IAM_SECRET_KEY" envDefault:"secret"`
}

func MustNewConfig() *Config {
	c := &Config{}
	env.Parse(c)

	return c
}
