package storage

import (
    "github.com/jinzhu/gorm"
    "../model"
    "errors"
    "time"
)

type UserStorage struct {
    db *gorm.DB
}

var user_storage *UserStorage

func NewUserStorage(db *gorm.DB) {
    user_storage = &UserStorage{db}
}

func (s UserStorage) GetAll() ([]model.User, error) {
    var users []model.User
    s.db.Find(&users)
    if s.db.Error != nil {
        return nil, s.db.Error
    }

    //userMap := make(map[uint64]*model.User)
    //for i := range users {
    //    u := &users[i]
    //    userMap[u.Id] = u
    //}

    return users, nil
}

func GetAllUsers() ([]model.User, error) {
    var users []model.User
    var err error
    if user_storage == nil {
        return nil, errors.New("nvalid user storage")
    }

    users, err = user_storage.GetAll()
    return users, err
}

func GetOne(userNameOrEmail string) (*model.User, error) {
    if user_storage == nil {
        return nil, errors.New("invalid user storage")
    }

    var user model.User
    user_storage.db.Model(&user).Where("user_name = ? or email = ?", userNameOrEmail, userNameOrEmail).Scan(&user)

    return &user, user_storage.db.Error
}

func GetUserLoginCntByIP(ip string, c chan int) {
    var val int
    if user_storage == nil {
        val = 0
    } else {
       user_storage.db.Model(&model.LoginAttempt{}).Where("ip = ?", ip).Count(&val)
    }

    c <- val
}

func GetUserLoginCntByIpAndName(ip, name string, c chan int){
    var val int
    if user_storage == nil {
        val = 0
    } else {
        user_storage.db.Model(&model.LoginAttempt{}).Where("ip = ? and user = ?", ip, name).Count(&val)
    }

    c <- val
}

func InsertUserLoginAttempt(ip, name string) (error) {
    var val model.LoginAttempt
    val.IP = ip
    val.User = name
    val.Time = time.Now()

    if user_storage == nil {
        return errors.New("invalid user storage")
    } else {
        user_storage.db.Create(&val)
    }
    return user_storage.db.Error
}

