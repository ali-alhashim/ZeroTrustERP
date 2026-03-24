package models

type User struct {
    ID        int
    Email     string
    Username  string
    Active    bool    // if user not active, then user cannot login
    OTPHash   string
    OTPExpiry int64
    Online    bool //user becomes online after verification
}