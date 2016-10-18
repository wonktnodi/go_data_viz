package main

import (
	"encoding/json"
	//"fmt"
	"net/http"
	"strings"

	//"log"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
    "./storage"
    "./model"
    "fmt"
)

func renderLogin(c *gin.Context) {
	isAuthenticated, _ := c.Get("isAuthenticated")
	if is, ok := isAuthenticated.(bool); ok && is {
		var redirectURL string
		defaultReturnUrl, _ := c.Get("DefaultReturnUrl")
		redirectURL = defaultReturnUrl.(string)
		session := sessions.Default(c)
		returnURL := session.Get("returnURL")
		if returnURL != nil {
			redirectURL = returnURL.(string)
			session.Delete("returnURL")
			session.Save()
		}
		c.Redirect(http.StatusFound, redirectURL)
	} else {
		injectSocials(c)
		c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
	}
}

func login(c *gin.Context) {
	response := newResponse(c)

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail()
		return
	}

	// validate
	if len(body.Username) == 0 {
		response.ErrFor["username"] = "required"
	}
	if len(body.Password) == 0 {
		response.ErrFor["password"] = "required"
	}
	if response.HasErrors() {
		response.Fail()
		return
	}
	body.Username = strings.ToLower(body.Username)
	db := getMongoDBInstance()
	defer db.Session.Close()

	// abuseFilter
	collection := db.C(LOGINATTEMPTS)
	IpCountChan := make(chan int)
	IpUserCountChan := make(chan int)
	clientIP := c.ClientIP()
    //var cnt int64
    //cnt, err = storage.GetUserLoginCntByIP(clientIP)
    //fmt.Println("@@@1 cnt: ", cnt)
    //cnt, err = storage.GetUserLoginCntByIpAndName(clientIP, body.Username)
    //fmt.Println("@@@2 cnt: ", cnt)



	//go getCount(collection, IpCountChan, bson.M{
	//	"ip": clientIP,
	//})
	//go getCount(collection, IpUserCountChan, bson.M{
	//	"ip":   clientIP,
	//	"user": body.Username,
	//})
	//IpCount := <-IpCountChan
	//IpUserCount := <-IpUserCountChan

    go storage.GetUserLoginCntByIP(clientIP, IpCountChan)
    go storage.GetUserLoginCntByIpAndName(clientIP, body.Username, IpUserCountChan)
    IpCount := <-IpCountChan
    IpUserCount := <-IpUserCountChan

    fmt.Println("@@@@ :", IpCount, ", ", IpUserCount)
    if IpCount > config.LoginAttempts.ForIp || IpUserCount > config.LoginAttempts.ForIpAndUser {
		response.Errors = append(response.Errors, "You've reached the maximum number of login attempts. Please try again later.")
		response.Fail()
		return
	}

	// attemptLogin
	collection = db.C(USERS)
	user := User{}
	err = collection.Find(bson.M{"$or": []bson.M{
		bson.M{"username": body.Username},
		bson.M{"email": body.Username}, // instead username can be email
	}}).One(&user)

    var val *model.User
    val, err = storage.GetOne(body.Username)
    fmt.Printf("%+v\n", val)

	if err != nil {
		//if err == mgo.ErrNotFound {
			response.Errors = append(response.Errors, "check username or email")
			response.Fail()
			return
		//}
		//EXCEPTION(err)
	}
	//err = user.isPasswordOk(body.Password)
    err = val.IsPasswordOk(body.Password)
	if err != nil {
		attempt := LoginAttempt{}
		attempt.IP = clientIP
		attempt.User = body.Username
		collection = db.C(LOGINATTEMPTS)
		err = collection.Insert(attempt)
		if err != nil {
			EXCEPTION(err)
		}

        storage.InsertUserLoginAttempt(clientIP, body.Username)
		response.Errors = append(response.Errors, "check password")
		response.Fail()
		return
	}

	session := sessions.Default(c)
	session.Set("public", user.ID.Hex())
	if returnURL, ok := session.Get("returnURL").(string); ok {
		c.Redirect(http.StatusFound, returnURL)
	}
	session.Delete("returnURL")
	session.Save()

	response.Finish()
}
