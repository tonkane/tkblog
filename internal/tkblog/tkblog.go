package tkblog

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tkane/tkblog/internal/pkg/known"
	"github.com/tkane/tkblog/internal/pkg/log"
	"github.com/tkane/tkblog/internal/tkblog/controller/v1/user"
	"github.com/tkane/tkblog/internal/tkblog/store"
	"github.com/tkane/tkblog/pkg/token"
	"github.com/tkane/tkblog/pkg/version/verflag"

	"github.com/gin-gonic/gin"

	mw "github.com/tkane/tkblog/internal/pkg/middleware"

	pb "github.com/tkane/tkblog/pkg/proto/tkblog/v1"
	"google.golang.org/grpc"
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

	if err := initStore(); err != nil {
		return err
	}

	// 初始化 token 包蜜月
	token.Init(viper.GetString("jwt-secret"), known.XusernameKey)

	gin.SetMode(viper.GetString("runmode"))

	g := gin.New()

	// 中间件设置
	mws := []gin.HandlerFunc{gin.Recovery(), mw.NoCache, mw.Cors, mw.Secure, mw.RequestID()}

	g.Use(mws...)

	if err := installRouters(g); err != nil {
		return err
	}

	// http 实例
	// httpsrv := &http.Server{Addr: viper.GetString("addr"), Handler: g}
	// // 日志打印
	// log.Infow("start to listening the incoming requests on http address", "addr", viper.GetString("addr"))
	
	// go func() {
	// 	if err := httpsrv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
	// 		log.Fatalw(err.Error())
	// 	}
	// }()

	// 启动 http 服务器
	httpsrv := startInsecureServer(g)
	// 启动 https 服务器
	httpssrv := startSecureServer(g)
	// 启动 grpc 服务器
	grpcsrv := startGRPCServer()
	
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

	if err := httpssrv.Shutdown(ctx); err != nil {
		log.Errorw("secure server forced to shutdown", "err", err)
		return err
	}

	// 优雅关闭 grpc
	grpcsrv.GracefulStop()

	log.Infow("server exiting")
	return nil
}


// 启动 HTTP 服务器
func startInsecureServer(g *gin.Engine) *http.Server {
	// http 实例
	httpsrv := &http.Server{Addr: viper.GetString("addr"), Handler: g}

	// 日志打印
	log.Infow("start to listening the incoming requests on http address", "addr", viper.GetString("addr"))
	
	go func() {
		if err := httpsrv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())
		}
	}()

	return httpsrv
}

// 启动 HTTPS 服务器
func startSecureServer(g *gin.Engine) *http.Server {
	// http 实例
	httpssrv := &http.Server{Addr: viper.GetString("tls.addr"), Handler: g}

	// 日志打印
	log.Infow("start to listening the incoming requests on https address", "addr", viper.GetString("tls.addr"))
	cert, key := viper.GetString("tls.cert"), viper.GetString("tls.key")
	if cert != "" && key != "" {
		go func() {
			if err := httpssrv.ListenAndServeTLS(cert, key); err != nil && errors.Is(err, http.ErrServerClosed) {
				log.Fatalw(err.Error())
			}
		}()
	}
	
	return httpssrv
}

// 启动 rpc 服务器
func startGRPCServer() *grpc.Server {
	lis, err := net.Listen("tcp", viper.GetString("grpc.addr"))
	if err != nil {
		log.Fatalw("failed to listen", "err", err)
	}

	// grpc 实例
	grpcsrv := grpc.NewServer()
	pb.RegisterTkBlogServer(grpcsrv, user.New(store.S, nil))

	log.Infow("start to listening the incoming requests on grpc address", "addr", viper.GetString("grpc.addr"))
	go func() {
		if err := grpcsrv.Serve(lis); err != nil {
			log.Fatalw(err.Error())
		}
	}()

	return grpcsrv
}