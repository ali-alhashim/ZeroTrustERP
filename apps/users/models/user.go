package models

import (
	"time"
	"zerotrusterp/apps/employees/models"
)

type User struct {
    ID        int        `f:"number, primary, auto"`
    Email     string     `f:"text, unique"`
    Username  string     `f:"text, unique, notnull"`
    Active    bool       `f:"bool, default:true"` // if user not active, then user cannot login
    OTPHash   string     `f:"text"`
    SessionToken string  `f:"text"`
    SessionExpiry int64   `f:"timestamp"`
    OTPExpiry int64       `f:"timestamp"`
    LastLogin int64       `f:"timestamp"`
    IncorrectOtpAttempts int     `f:"number, default:0"` // number of incorrect OTP attempts, reset to 0 after successful login
    Online    bool     `f:"bool, default:false"` //user becomes online after verification -> WebSocket Approach (Real-time)

    Roles     *[]Role  `f:"many2many:"` // many to many relationship with roles One user → many roles , One role → many users

    
    RelatedEmployee *models.Employee `f:"one2one:employees"` //this skip by migration only object reference for query, not create foreign key in database

    CreatedAt time.Time `f:"timestamp, default:current_timestamp"`
    UpdatedAt time.Time `f:"timestamp, default:current_timestamp"`
}

// Role header can be admin, manager, employee, etc. Each role has many permissions, and each permission can be assigned to many roles. For example, admin role can have all permissions, manager role can have read and update permissions, employee role can have read permission only.
type Role struct {
	ID          int     `f:"number, primary, auto"`
	Name        string  `f:"text, unique, notnull"`
	Description string  `f:"text"`
	Permissions *[]Permission `f:"many2many:"`
	CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
	UpdatedAt   time.Time `f:"timestamp, default:current_timestamp"`
}


//Permissions read R, write W, delete D, update U,  all A
type Permission struct {
	ID          int     `f:"number, primary, auto"`
	Resource    string  `f:"text, unique, notnull"` // resource can be the name of the model or the name of the API endpoint or the name of the action
	Action      string  `f:"text, unique, notnull"` // action can be R, W, D, U, A
    Description string  `f:"text"`
    CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time `f:"timestamp, default:current_timestamp"`
}


// for audit-trail

type Log struct {
    ID        int       `f:"number, primary, auto"`
    UserID    int       `f:"number, notnull"` // foreign key to user
    Username  string    `f:"text, notnull"` // store username for easy query, 
    Email     string    `f:"text, notnull"` // store email for easy query,
    Resource  string    `f:"text, notnull"`
    Action    string    `f:"text, notnull"` // action can be login, logout, create, update, delete, etc.
    Timestamp time.Time `f:"timestamp, default:current_timestamp"`
}
