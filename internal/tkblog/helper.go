package tkblog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tkane/tkblog/internal/pkg/log"
	"github.com/tkane/tkblog/internal/tkblog/store"
	"github.com/tkane/tkblog/pkg/db"
)

const (
	defaultConfigDir = ".tkblog"

	defaultConfigName = "tkblog.yaml"
)

// 初始化配置，读取配置到viper中
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// 查找用户主目录
		home, err := os.UserHomeDir()
		// 如果主目录获取失败，打印错误，推出程序
		cobra.CheckErr(err)

		// 添加配置搜索路径
		viper.AddConfigPath(filepath.Join(home, defaultConfigDir))
		viper.AddConfigPath(".")

		// 格式设置
		viper.SetConfigType("yaml")
		// 默认名称
		viper.SetConfigName(defaultConfigName)
	}

	// 读取环境变量
	viper.AutomaticEnv()

	// 设置环境变量前缀
	viper.SetEnvPrefix("TKANE")

	// 替换掉特殊字符
	replacer := strings.NewReplacer(".", "_", "-", "_")
	viper.SetEnvKeyReplacer(replacer)

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// fmt.Fprintln(os.Stderr, err)
		log.Errorw("fail to read viper config", "err", err)
	}

	// 打印当前使用的配置文件
	// fmt.Fprintln(os.Stdout, "using config file:", viper.ConfigFileUsed())
	log.Infow("using config file:", "file", viper.ConfigFileUsed())
}

func getConfigInfo() {
	settings, _ := json.Marshal(viper.AllSettings())
	// fmt.Println(string(settings))
	log.Infow("using config:", "config", string(settings))
}

func logOptions() *log.Options {
	return &log.Options{
		DisableCaller: viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		Level: viper.GetString("log.level"),
		Format: viper.GetString("log.format"),
		OutputPaths: viper.GetStringSlice("log.output-paths"),
	}
}


// 读取db配置，创建gorm实例，初始化store层
func initStore() error {
	dbOptions := &db.MySQLOptions {
		Host: viper.GetString("db.host"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		Database: viper.GetString("db.database"),
		MaxIdleConnections: viper.GetInt("db.max-idle-connections"),
		MaxOpenConnections: viper.GetInt("db.max-open-connections"),
		MaxConnectionLifeTime: viper.GetDuration("db.max-connection-life-time"),
		LogLevel: viper.GetInt("db.log-level"),
	}

	ins, err := db.NewMySQL(dbOptions)

	if err != nil {
		return err
	}

	_ = store.NewStore(ins)

	log.Infow("connect to mysql!")
	return nil
}