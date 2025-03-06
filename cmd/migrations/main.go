package main

import (
	"github.com/vysogota0399/gophermart_portal/internal/config"
	"github.com/vysogota0399/gophermart_portal/internal/storage"
)

func main() {
	storage.RunMigration(config.MustNewConfig())
}
