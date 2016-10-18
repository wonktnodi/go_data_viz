package config

import (
    "github.com/gorilla/sessions"
    "io"
    "os"
    "log"
    "io/ioutil"
    "fmt"
    "encoding/json"
)
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

type LoginAttempts struct {
    ForIp         int
    ForIpAndUser  int
    LogExpiration string
}

type Configuration struct {
    CompanyName   string
    ProjectName   string
    Database      Databases      `json:"Database"`
    Session       Session        `json:"Session"`
    Server        Server         `json:"Server"`
    LoginAttempts LoginAttempts  `json:"UserLogin"`
}

func LoadConfig(configFile string) (Configuration){
    var ret Configuration
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
    err = json.Unmarshal(jsonBytes, &ret)
    if err != nil {
        log.Fatalln("Could not parse %q: %v", configFile, err)
    }
    fmt.Printf("%+v\n", ret)
    return ret
}