package base

import (
	"fmt"
	"github.com/ichaly/yugong/core/util"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Debug     bool        `mapstructure:"debug" jsonschema:"title=Debug"`
	App       *App        `mapstructure:"app" jsonschema:"title=App"`
	Workspace string      `mapstructure:"workspace" jsonschema:"title=Workspace"`
	Cache     *DataSource `mapstructure:"cache" jsonschema:"title=Cache"`
	Database  *DataSource `mapstructure:"database" jsonschema:"title=DataSource"`
	Proxy     *Proxy      `mapstructure:"proxy" jsonschema:"title=Proxy Config"`
}

type App struct {
	Name string `mapstructure:"name" jsonschema:"title=Application Name"`
	Port string `mapstructure:"port" jsonschema:"title=Application Port"`
	Host string `mapstructure:"host" jsonschema:"title=Application Host"`
}

type DataSource struct {
	Url      string       `json:"url"`
	Host     string       `json:"host"`
	Port     int          `json:"port"`
	Name     string       `json:"name"`
	Dialect  string       `json:"dialect"`
	Username string       `json:"username"`
	Password string       `json:"password"`
	Sources  []DataSource `json:"sources"`
	Replicas []DataSource `json:"replicas"`
}

type Proxy struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewConfig() (*Config, error) {
	return readInConfig(filepath.Join("../conf", "dev.yml"))
}

func readInConfig(configFile string) (*Config, error) {
	cp := filepath.Dir(configFile)
	vi := newViper(cp, filepath.Base(configFile))

	if err := vi.ReadInConfig(); err != nil {
		return nil, err
	}

	if pcf := vi.GetString("inherits"); pcf != "" {
		cf := vi.ConfigFileUsed()
		vi = newViper(cp, pcf)

		if err := vi.ReadInConfig(); err != nil {
			return nil, err
		}

		if v := vi.GetString("inherits"); v != "" {
			return nil, fmt.Errorf("inherited config (%s) cannot itself inherit (%s)", pcf, v)
		}

		vi.SetConfigFile(cf)

		if err := vi.MergeInConfig(); err != nil {
			return nil, err
		}
	}

	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "GJ_") || strings.HasPrefix(e, "SJ_") {
			kv := strings.SplitN(e, "=", 2)
			util.SetKeyValue(vi, kv[0], kv[1])
		}
	}

	c := &Config{}

	if err := vi.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed to decode config, %v", err)
	}

	return c, nil
}

func newViper(configPath, configFile string) *viper.Viper {
	vi := newViperWithDefaults()
	vi.SetConfigName(strings.TrimSuffix(configFile, filepath.Ext(configFile)))

	if configPath == "" {
		vi.AddConfigPath("./config")
	} else {
		vi.AddConfigPath(configPath)
	}

	return vi
}

func newViperWithDefaults() *viper.Viper {
	vi := viper.New()

	vi.SetDefault("debug", true)

	vi.SetDefault("app.port", "3000")

	vi.SetDefault("cache.dialect", "memory")

	vi.SetDefault("database.dialect", "postgres")
	vi.SetDefault("database.host", "localhost")
	vi.SetDefault("database.port", 5432)
	vi.SetDefault("database.username", "postgres")
	vi.SetDefault("database.password", "")
	vi.SetDefault("database.schema", "public")
	vi.SetDefault("database.pool_size", 10)

	vi.SetDefault("env", "development")

	_ = vi.BindEnv("env", "GO_ENV")
	_ = vi.BindEnv("host", "HOST")
	_ = vi.BindEnv("port", "PORT")

	return vi
}
