package models

import (
	"time"
	"zerotrusterp/apps/employees/models"
)

type User struct {
    ID        int     `f:"number, primary, auto"`
    Email     string  `f:"text, unique"`
    Username  string  `f:"text, unique, notnull"`
    Active    bool    `f:"bool, default:true"` // if user not active, then user cannot login
    OTPHash   string  `f:"text"`
    SessionToken string  `f:"text"`
    OTPExpiry int64   `f:"timestamp"`
    LastLogin int64   `f:"timestamp"`
    IncorrectOtpAttempts int     `f:"number, default:0"` // number of incorrect OTP attempts, reset to 0 after successful login
    Online    bool    `f:"bool, default:false"` //user becomes online after verification

    Roles     *[]Role  `f:"many2many:user_roles"` // many to many relationship with roles One user → many roles , One role → many users

    
    RelatedEmployee *models.Employee `f:"one2one:employees"` //this skip by migration only object reference for query, not create foreign key in database

    CreatedAt time.Time `f:"timestamp, default:current_timestamp"`
    UpdatedAt time.Time `f:"timestamp, default:current_timestamp"`
}

type Role struct {
	ID          int     `f:"number, primary, auto"`
	Name        string  `f:"text, unique, notnull"`
	Description string  `f:"text"`
	Permissions *[]Permission `f:"many2many:role_permissions"`
	CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
	UpdatedAt   time.Time `f:"timestamp, default:current_timestamp"`
}


//Permissions read R, write W, delete D, update U,  all A
type Permission struct {
	ID          int     `f:"number, primary, auto"`
	resource    string  `f:"text, unique, notnull"` // resource can be the name of the model or the name of the API endpoint or the name of the action
	Action      string  `f:"text, unique, notnull"` // action can be R, W, D, U, A
    Description string  `f:"text"`
    CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time `f:"timestamp, default:current_timestamp"`
}
