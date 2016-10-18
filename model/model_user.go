package model

import (
    "time"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    Id             uint64 `gorm:"primary_key;AUTO_INCREMENT"`
    UserName       string `gorm:"unique"`
    Password       string `gorm:"type:varchar(256)"`
    Email          string `gorm:"unique"`
    GroupID        int
    Active         bool
    CreateTime     time.Time
    PwdResetToken  string
    PwdResetExpire time.Time
}

type LoginAttempt struct {
    ID   uint64    `gorm:"primary_key;AUTO_INCREMENT"`
    IP   string
    User string
    Time time.Time
}

func (user *User) IsPasswordOk(password string) error {
    return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
