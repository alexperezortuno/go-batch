package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type BackupConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
	Format  string `yaml:"format"`
}

type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Host    string `yaml:"host"`
}

type PathsConfig struct {
	UserFile string `yaml:"user_file"`
}

type FileConfig struct {
	Paths                 PathsConfig       `yaml:"paths"`
	Type                  string            `yaml:"type"`
	ColumnMapping         map[string]string `yaml:"column_mapping"`
	Delimiter             string            `yaml:"delimiter"`
	Encoding              string            `yaml:"encoding"`
	Header                bool              `yaml:"header"`
	SkipRows              int               `yaml:"skip_rows"`
	SheetName             string            `yaml:"sheet_name"`
	SheetIndex            int               `yaml:"sheet_index"`
	RowLimit              int               `yaml:"row_limit"`
	ColumnLimit           int               `yaml:"column_limit"`
	BatchSize             int               `yaml:"batch_size"`
	BatchTimeout          time.Duration     `yaml:"batch_timeout"`
	BatchDelay            time.Duration     `yaml:"batch_delay"`
	BatchRetries          int               `yaml:"batch_retries"`
	BatchRetryDelay       time.Duration     `yaml:"batch_retry_delay"`
	BatchRetryMaxDelay    time.Duration     `yaml:"batch_retry_max_delay"`
	BatchRetryMaxAttempts int               `yaml:"batch_retry_max_attempts"`
	HasHeader             bool              `yaml:"has_header"`
}

type LoggerConfig struct {
	Level      string        `yaml:"level"`
	Caller     bool          `yaml:"caller"`
	TimeFormat string        `yaml:"time_format"`
	FilePath   string        `yaml:"file_path"`
	MaxSize    int           `yaml:"max_size"`
	MaxBackups int           `yaml:"max_backups"`
	MaxAge     int           `yaml:"max_age"`
	Compress   bool          `yaml:"compress"`
	JSONFormat bool          `yaml:"json_format"`
	Rotate     bool          `yaml:"rotate"`
	RotateTime time.Duration `yaml:"rotate_time"`
}

type Config struct {
	BatchSize      int           `yaml:"batch_size"`
	MaxConcurrency int           `yaml:"max_concurrency"`
	DB             DBConfig      `yaml:"db"`
	Backup         BackupConfig  `yaml:"backup"`
	Metrics        MetricsConfig `yaml:"metrics"`
	FileProcessing FileConfig    `yaml:"file_processing"`
	Logger         LoggerConfig  `yaml:"logger"`
}

type DBConfig struct {
	DSN             string        `yaml:"dsn"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	DisableIndex    bool          `yaml:"disable_index"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	Migrate         bool          `yaml:"migrate"`
}

func Load(configPath string) (*Config, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		fmt.Printf("Failed to get absolute path: %v", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
