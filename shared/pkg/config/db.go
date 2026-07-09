package config

import (
	"fmt"
	"os"
)

// DBURI возвращает DSN текущего сервиса из переменной окружения DB_URI.
func DBURI() (string, error) {
	return requiredEnv("DB_URI")
}

// InventoryDBURI возвращает DSN inventory: INVENTORY_DB_URI или DB_URI.
func InventoryDBURI() (string, error) {
	if dsn := os.Getenv("INVENTORY_DB_URI"); dsn != "" {
		return dsn, nil
	}

	return requiredEnv("DB_URI")
}

// OrderDBURI возвращает DSN order: ORDER_DB_URI или DB_URI.
func OrderDBURI() (string, error) {
	if dsn := os.Getenv("ORDER_DB_URI"); dsn != "" {
		return dsn, nil
	}

	return requiredEnv("DB_URI")
}

func requiredEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("переменная окружения %s не задана", key)
	}

	return value, nil
}
