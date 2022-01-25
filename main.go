package main

import (
	"context"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sorokinmax/websspi"
)

const version = "v.2.3.0"

var (
	cfg         Config
	auth        websspi.Authenticator
	UserInfoKey = "websspi-key-UserInfo"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create("web.log")
	gin.DefaultWriter = io.MultiWriter(os.Stdout, f)
	defer f.Close()

	log.SetFlags(log.LstdFlags)
	lf, err := os.OpenFile("output.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer lf.Close()
	multi := io.MultiWriter(os.Stdout, lf)
	log.SetOutput(multi)

	readConfigFile(&cfg)

	SQLiteCreateDB(Ipa{})

	router := gin.Default()
	router.HTMLRender = ginview.Default()
	config := websspi.NewConfig()
	auth, _ := websspi.New(config)

	//router.Use(MidAuth(auth))
	//router.Use(AddUserToCtx())
	router.StaticFile("favicon.ico", "./views/favicon.ico")
	router.Use(static.Serve("/ipa", static.LocalFile("./ipa", false)))
	router.Use(static.Serve("/images", static.LocalFile("./images", false)))
	router.GET("/", indexHandler)
	router.GET("/admin", MidAuth(auth), AddUserToCtx(), indexHandler)
	router.GET("/version/:version/", versionHandler)
	router.POST("/action/qr", qrHandler)
	router.POST("/action/remove", removeHandler)
	router.POST("/action/ipa", postIpaHandler)

	go ipaScaner()

	log.Println("Web is available at " + cfg.Service.Url + ":" + strconv.Itoa(cfg.Service.Port))
	router.Run(":" + strconv.Itoa(cfg.Service.Port))
}

func AddUserToCtx() gin.HandlerFunc {
	return func(c *gin.Context) {
		if ctxVars, ok := c.Request.Context().Value(UserInfoKey).(*websspi.UserInfo); ok {
			c.Set("user", ctxVars.Username)
		} else {
			//c.Set("user", "guest")
			c.Abort()
			//c.Next()
			return
		}
	}
}

func MidAuth(a *websspi.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {

		user, data, err := a.Authenticate(c.Request, c.Writer)
		if err != nil {
			a.Return401(c.Writer, data)
			return
		}

		// Add the UserInfo value to the reqest's context
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), UserInfoKey, user))
		// and to the request header with key Config.AuthUserKey
		if a.Config.AuthUserKey != "" {
			c.Request.Header.Set(a.Config.AuthUserKey, user.Username)
		}

		// The WWW-Authenticate header might need to be sent back even
		// on successful authentication (eg. in order to let the client complete
		// mutual authentication).
		if data != "" {
			a.AppendAuthenticateHeader(c.Writer, data)
		}

		c.Next()
	}
}
