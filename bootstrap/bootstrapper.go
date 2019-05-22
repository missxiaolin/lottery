package bootstrap

import (
	"github.com/gorilla/securecookie"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/sessions"
	"lottery/conf"
	"lottery/cron"
	"time"
)

type Configurator func(*Bootstrapper)

// 使用Go内建的嵌入机制(匿名嵌入)，允许类型之前共享代码和数据
// （Bootstrapper继承和共享 iris.Application ）
// 参考文章： https://hackthology.com/golangzhong-de-mian-xiang-dui-xiang-ji-cheng.html
type Bootstrapper struct {
	*iris.Application
	AppName      string
	AppOwner     string
	AppSpawnDate time.Time

	Sessions *sessions.Sessions
}

// New returns a new Bootstrapper.
func New(appName, appOwner string, cfgs ...Configurator) *Bootstrapper {
	b := &Bootstrapper{
		AppName:      appName,
		AppOwner:     appOwner,
		AppSpawnDate: time.Now(),
		Application:  iris.New(),
	}

	for _, cfg := range cfgs {
		cfg(b)
	}

	return b
}

// SetupViews loads the templates.
func (b *Bootstrapper) SetupViews(viewsDir string) {
	htmlEngine := iris.HTML(viewsDir, ".html").Layout("shared/layout.html")
	// 每次重新加载模版（线上关闭它）
	htmlEngine.Reload(true)
	// 给模版内置各种定制的方法
	htmlEngine.AddFunc("FromUnixtimeShort", func(t int) string {
		dt := time.Unix(int64(t), int64(0))
		return dt.Format(conf.SysTimeformShort)
	})
	htmlEngine.AddFunc("FromUnixtime", func(t int) string {
		dt := time.Unix(int64(t), int64(0))
		return dt.Format(conf.SysTimeform)
	})
	b.RegisterView(htmlEngine)
}

// SetupSessions initializes the sessions, optionally.
func (b *Bootstrapper) SetupSessions(expires time.Duration, cookieHashKey, cookieBlockKey []byte) {
	b.Sessions = sessions.New(sessions.Config{
		Cookie:   "SECRET_SESS_COOKIE_" + b.AppName,
		Expires:  expires,
		Encoding: securecookie.New(cookieHashKey, cookieBlockKey),
	})
}

func (b *Bootstrapper) SetupErrorHandlers() {
	b.OnAnyErrorCode(func(ctx iris.Context) {
		err := iris.Map{
			"app":     b.AppName,
			"status":  ctx.GetStatusCode(),
			"message": ctx.Values().GetString("message"),
		}

		if jsonOutput := ctx.URLParamExists("json"); jsonOutput {
			ctx.JSON(err)
			return
		}

		ctx.ViewData("Err", err)
		ctx.ViewData("Title", "Error")
		ctx.View("shared/error.html")
	})
}

// Configure accepts configurations and runs them inside the Bootstraper's context.
func (b *Bootstrapper) Configure(cs ...Configurator) {
	for _, c := range cs {
		c(b)
	}
}

// 启动计划任务服务
func (b *Bootstrapper) setupCron() {
	// 服务类应用
	if conf.RunningCrontabService {
		cron.ConfigueAppOneCron()
	}
	cron.ConfigueAppAllCron()
}

const (
	// StaticAssets is the root directory for public assets like images, css, js.
	StaticAssets = "./public/"
	// Favicon is the relative 9to the "StaticAssets") favicon path for our app.
	Favicon = "favicon.ico"
)

// 跨域请求解决
func (b *Bootstrapper) setCorsOptions() {
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	}
	corsWrapper := cors.New(corsOptions).ServeHTTP
	b.WrapRouter(corsWrapper)
}

// Bootstrap prepares our application.
//
// Returns itself.
func (b *Bootstrapper) Bootstrap() *Bootstrapper {
	b.SetupViews("./views")
	b.SetupErrorHandlers()
	b.Favicon(StaticAssets + Favicon)
	b.StaticWeb(StaticAssets[1:len(StaticAssets)-1], StaticAssets)
	b.setupCron()

	b.Use(recover.New())
	b.Use(logger.New())
	b.setCorsOptions()
	return b
}

// Listen starts the http server with the specified "addr".
func (b *Bootstrapper) Listen(addr string, cfgs ...iris.Configurator) {
	b.Run(iris.Addr(addr), cfgs...)
}