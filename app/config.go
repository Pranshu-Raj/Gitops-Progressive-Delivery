package main

import (
	"fmt"
	"os"
	"strconv"
)

type config struct {
	addr      string
	env       string
	failEvery int
	slowMs    int
}

func loadConfig() (config, error) {
	cfg := config{
		addr: getenv("ADDR", ":8080"),
		env:  getenv("APP_ENV", "local"),
	}

	var err error
	if cfg.failEvery, err = getenvInt("FAIL_EVERY"); err != nil {
		return config{}, err
	}
	if cfg.slowMs, err = getenvInt("SLOW_MS"); err != nil {
		return config{}, err
	}
	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getenvInt(key string) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return 0, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("%s: want a non-negative integer, got %q", key, v)
	}
	return n, nil
}
