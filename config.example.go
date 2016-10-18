package main

import "os"
import (
    "log"
    "io"
    "io/ioutil"
    "encoding/json"
    "fmt"
    "github.com/gorilla/sessions"
)

const defaultLocalMongoDBUrl = "mongodb://120.26.6.49:27017/data_viz"
const defaultPORT = "4000"
const ROOTGROUP = "root"

var config struct {
    Port                       string
    CompanyName                string
    ProjectName                string
    SystemEmail                string
    CryptoKey                  string
    RequireAccountVerification bool
    MongoDB                    string
    dbName                     string
    LoginAttempts              LoginAttempts
    SMTP                       SMTP
    Socials                    map[string]OAuth
}

type LoginAttempts struct {
    ForIp         int
    ForIpAndUser  int
    LogExpiration string
}

type SMTP struct {
    From        struct {
                    Name, Address string
                }
    Credentials struct {
                    User, Password, Host string
                    SSL                  bool
                }
}

type OAuth struct {
    Key, Secret string
}

// MySQLInfo is the details for the database connection
type MySQLInfo struct {
    Username  string
    Password  string
    Name      string
    Hostname  string
    Port      int
    Parameter string
}

// SQLiteInfo is the details for the database connection
type SQLiteInfo struct {
    Parameter string
}

type Databases struct {
    Type   string
    MySQL  MySQLInfo
    SQLite SQLiteInfo
}

// Session stores session level information
type Session struct {
    Options   sessions.Options `json:"Options"`   // Pulled from: http://www.gorillatoolkit.org/pkg/sessions#Options
    Name      string           `json:"Name"`      // Name for: http://www.gorillatoolkit.org/pkg/sessions#CookieStore.Get
    SecretKey string           `json:"SecretKey"` // Key for: http://www.gorillatoolkit.org/pkg/sessions#CookieStore.New
}

// Server stores the hostname and port number
type Server struct {
    CompanyName                string
    ProjectName                string
    SystemEmail                string
    RequireAccountVerification bool
    Hostname                   string `json:"Hostname"`  // Server name
    UseHTTP                    bool   `json:"UseHTTP"`   // Listen on HTTP
    UseHTTPS                   bool   `json:"UseHTTPS"`  // Listen on HTTPS
    HTTPPort                   int    `json:"HTTPPort"`  // HTTP port
    HTTPSPort                  int    `json:"HTTPSPort"` // HTTPS port
    CertFile                   string `json:"CertFile"`  // HTTPS certificate
    KeyFile                    string `json:"KeyFile"`   // HTTPS private key
}

type configuration struct {
    CompanyName   string
    ProjectName   string
    Database      Databases      `json:"Database"`
    Session       Session        `json:"Session"`
    Server        Server         `json:"Server"`
    LoginAttempts LoginAttempts  `json:"UserLogin"`
}

var new_config configuration

func loadConfig(configFile string) {
    var err error
    var input = io.ReadCloser(os.Stdin)
    if input, err = os.Open(configFile); err != nil {
        log.Fatalln(err)
    }
    // Read the config file
    jsonBytes, err := ioutil.ReadAll(input)
    input.Close()
    if err != nil {
        log.Fatalln(err)
    }

    // Parse the config
    err = json.Unmarshal(jsonBytes, &new_config)
    if err != nil {
        log.Fatalln("Could not parse %q: %v", configFile, err)
    }
    fmt.Printf("%+v\n", new_config)
}

func getEnvOrSetDef(envName, defValue string) (val string) {
    val, ok := os.LookupEnv(envName)
    if !ok {
        val = defValue
    }
    return
}

func InitConfig() {

    config.Port = getEnvOrSetDef("PORT", defaultPORT)

    config.CompanyName = "xxxx, Inc."
    config.ProjectName = "Data viz"
    config.SystemEmail = "norelay@gmail.com"
    config.CryptoKey = "k3yb0ardc4t"
    config.RequireAccountVerification = true

    config.MongoDB = getEnvOrSetDef(
        "MONGODB_URI",
        getEnvOrSetDef(
            "MONGOLAB_URI",
            getEnvOrSetDef(
                "MONGOHQ_URL",
                defaultLocalMongoDBUrl,
            )))

    if config.dbName == "" {
        config.dbName = getDBName(&config.MongoDB)
    }
    log.Println("db config: ", config.MongoDB, ", ", config.dbName)
    config.LoginAttempts.ForIp = 50
    config.LoginAttempts.ForIpAndUser = 7
    config.LoginAttempts.LogExpiration = "20m"

    config.SMTP.From.Name = getEnvOrSetDef("SMTP_FROM_NAME", config.ProjectName + " Website")
    config.SMTP.From.Address = getEnvOrSetDef("SMTP_FROM_ADDRESS", "your@email.addy")

    //config.SMTP.Credentials.User = getEnvOrSetDef("SMTP_USERNAME", "your@email.addy")
    //config.SMTP.Credentials.Password = getEnvOrSetDef("SMTP_PASSWORD", "bl4rg!")

    config.SMTP.Credentials.User = getEnvOrSetDef("SMTP_USERNAME", "welcome@sturfee.com")
    config.SMTP.Credentials.Password = getEnvOrSetDef("SMTP_PASSWORD", "sturfee_knoxville")
    config.SMTP.Credentials.Host = getEnvOrSetDef("SMTP_HOST", "smtp.gmail.com")
    config.SMTP.Credentials.SSL = true

    // I think it's ok. I use it only for "get". No modifying
    config.Socials = make(map[string]OAuth)

    ins := OAuth{}

    ins.Key = getEnvOrSetDef("TWITTER_OAUTH_KEY", "")
    ins.Secret = getEnvOrSetDef("TWITTER_OAUTH_SECRET", "")
    if len(ins.Key) != 0 {
        config.Socials["twitter"] = ins
    }

    ins.Key = getEnvOrSetDef("FACEBOOK_OAUTH_KEY", "")
    ins.Secret = getEnvOrSetDef("FACEBOOK_OAUTH_SECRET", "")
    if len(ins.Key) != 0 {
        config.Socials["facebook"] = ins
    }

    ins.Key = getEnvOrSetDef("GITHUB_OAUTH_KEY", "")
    ins.Secret = getEnvOrSetDef("GITHUB_OAUTH_SECRET", "")
    if len(ins.Key) != 0 {
        config.Socials["github"] = ins
    }

    ins.Key = getEnvOrSetDef("GOOGLE_OAUTH_KEY", "")
    ins.Secret = getEnvOrSetDef("GOOGLE_OAUTH_SECRET", "")
    if len(ins.Key) != 0 {
        config.Socials["google"] = ins
    }

    ins.Key = getEnvOrSetDef("TUMBLR_OAUTH_KEY", "")
    ins.Secret = getEnvOrSetDef("TUMBLR_OAUTH_SECRET", "")
    if len(ins.Key) != 0 {
        config.Socials["tumblr"] = ins
    }

}
