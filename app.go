package main

import (
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/contrib/gzip"
    "github.com/gin-gonic/contrib/sessions"
    "github.com/gin-gonic/gin"
    "os"
    "strconv"
    new_cfg "./config"
    "./storage"
)

const VERSION = "0.1"
const ISOSTRING = "2006-01-02T15:04:05.99Z"

var store sessions.CookieStore
var cfg new_cfg.Configuration

var Router *gin.Engine

var year int

func init() {
    cfg = new_cfg.LoadConfig("./config" + string(os.PathSeparator) + "config.json")
    loadConfig("./config" + string(os.PathSeparator) + "config.json")
    InitConfig()
    //store = sessions.NewCookieStore([]byte("MFDQmJQ4TF"))
    store = sessions.NewCookieStore([]byte(new_config.Session.SecretKey))
    store.Options(sessions.Options{
        Path:     cfg.Session.Options.Path,
        MaxAge:   60 * 60 * 6,
        HttpOnly: cfg.Session.Options.HttpOnly,
    })
}

func main() {
    log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

    gin.SetMode(gin.DebugMode)

    err := storage.InitDB(cfg)
    if err != nil {
        log.Println(err)
        return
    }
    defer storage.CloseDB()

    Router = gin.New()
    Router.Use(gzip.Gzip(gzip.DefaultCompression))

    Router.StaticFile("/favicon.ico", "public/favicon.ico")
    Router.Static("/public", "public")
    Router.Static("/vendor", "vendor")

    // templates
    Router.HTMLRender = initTemplates()

    Router.Use(gin.Logger())
    Router.Use(checkRecover)

    Router.Use(sessions.Sessions("session", store))

    //refresh year every minute
    go func() {
        for {
            year, _, _ = time.Now().Date()
            time.Sleep(time.Minute)
        }
    }()

    Router.Use(func(c *gin.Context) {
        session := sessions.Default(c)
        oauthMessage, exist := session.Get("oauthMessage").(string)
        session.Delete("oauthMessage")
        session.Save()
        //session.Dump()

        c.Set("oauthMessage", oauthMessage)
        c.Set("oauthMessageExist", exist)
        c.Set("ProjectName", new_config.ProjectName)
        c.Set("CopyrightYear", year)
        c.Set("CopyrightName", new_config.CompanyName)
        c.Set("CacheBreaker", "br34k-01")
        c.Next()
    })
    Router.Use(IsAuthenticated)
    bindRoutes(Router) // --> cmd/go-getting-started/routers.go

    Router.Run(":" + strconv.Itoa(new_config.Server.HTTPPort))
}

// force redirect to https from http
// necessary only if you use https directly
// put your domain name instead of CONF.ORIGIN
func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
    //http.Redirect(w, req, "https://" + CONF.ORIGIN + req.RequestURI, http.StatusMovedPermanently)
}
