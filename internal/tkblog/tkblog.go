package tkblog

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tkane/tkblog/internal/pkg/log"
	"github.com/tkane/tkblog/pkg/version/verflag"

	"github.com/gin-gonic/gin"

	mw "github.com/tkane/tkblog/internal/pkg/middleware"
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
	gin.SetMode(viper.GetString("runmode"))

	g := gin.New()

	// 中间件设置
	mws := []gin.HandlerFunc{gin.Recovery(), mw.NoCache, mw.Cors, mw.Secure, mw.RequestID()}

	g.Use(mws...)

	// 404 handler
	g.NoRoute(func (c *gin.Context)  {
		c.JSON(http.StatusOK, gin.H{"code": 10003, "message": "Page not found."})
	})
	// healthz handler
	g.GET("/healthz", func (c *gin.Context)  {
		// 打印 X-request-id
		log.C(c).Infow("healthz is called!")

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	// http 实例
	httpsrv := &http.Server{Addr: viper.GetString("addr"), Handler: g}
	// 日志打印
	log.Infow("start to listening the incoming requests on http address", "addr", viper.GetString("addr"))
	
	go func() {
		if err := httpsrv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())
		}
	}()
	
	quit := make(chan os.Signal, 1)
	// KILL = syscall.SIGTERM
	// KILL 2 = syscall.SIGINT = CTRL + C
	// KILL 9 = syscall.SIGKILL 无法捕获
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 这一行代码会阻塞，直到quit通道接收到一个信号
	log.Infow("shutting down server")
	// 设置一个超时上下文并尝试优雅地关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpsrv.Shutdown(ctx); err != nil {
		log.Errorw("insecure server forced to shutdown", "err", err)
		return err
	}

	log.Infow("server exiting")
	return nil
}
