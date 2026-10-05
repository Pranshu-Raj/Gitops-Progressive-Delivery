package main

import "testing"

func TestLoadConfigDefaults(t *testing.T) {
	for _, k := range []string{"ADDR", "APP_ENV", "FAIL_EVERY", "SLOW_MS"} {
		t.Setenv(k, "")
	}

	cfg, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.addr != ":8080" || cfg.env != "local" || cfg.failEvery != 0 || cfg.slowMs != 0 {
		t.Errorf("got %+v", cfg)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	t.Setenv("ADDR", ":9090")
	t.Setenv("APP_ENV", "prod")
	t.Setenv("FAIL_EVERY", "5")
	t.Setenv("SLOW_MS", "250")

	cfg, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.addr != ":9090" || cfg.env != "prod" || cfg.failEvery != 5 || cfg.slowMs != 250 {
		t.Errorf("got %+v", cfg)
	}
}

func TestLoadConfigRejectsBadValues(t *testing.T) {
	for _, tc := range []struct{ key, val string }{
		{"FAIL_EVERY", "-1"},
		{"FAIL_EVERY", "two"},
		{"SLOW_MS", "1.5"},
	} {
		t.Run(tc.key+"="+tc.val, func(t *testing.T) {
			t.Setenv(tc.key, tc.val)
			if _, err := loadConfig(); err == nil {
				t.Error("expected an error")
			}
		})
	}
}
