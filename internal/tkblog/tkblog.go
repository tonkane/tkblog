package tkblog

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tkane/tkblog/internal/pkg/log"
	"github.com/tkane/tkblog/pkg/version/verflag"
)

// 全局配置文件信息
var cfgFile string

func NewBlogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "miniblog",
		Short:        "A good Go practical project",
		Long:         "tkane's practical project!",
		SilenceUsage: true,
		// cmd.Execute() 时执行的run函数
		RunE: func(cmd *cobra.Command, args []string) error {
			// 打印版本信息
			verflag.PrintAndExitIfReq()
			// 初始化日志
			log.Init(logOptions())
			defer log.Sync()
			return run()
		},
		// 命令运行时，不需要指定命令行参数
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any args, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	// 初始化配置
	cobra.OnInitialize(initConfig)

	// 定义了一个命令行标志，允许用户通过 --config 或 -c 选项来指定一个配置文件的路径
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path to the blog configuration file. Empty string for no configuration file.")

	// 当用户在命令行中运行程序并提供 -t 或 --toggle 选项时，toggle 标志的值将被设置为 true。如果用户没有提供这个选项，它的值将保持默认的 false。
	// 例如，如果你的程序是一个服务器应用程序，--toggle 标志可能用来启用或禁用某个特定的功能，如调试模式。用户可以通过命令行参数 -t 来启用这个功能，而不需要修改程序的配置文件或环境变量。
	cmd.Flags().BoolP("toggle", "t", false, "help message for toggle")

	// 添加 --version 标志
	verflag.AddFlags(cmd.PersistentFlags())
	return cmd
}

func run() error {
	getConfigInfo()
	return nil
}
