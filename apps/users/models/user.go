package models
import "time"

type Users struct {
    ID        int     `f:"number, primary, auto"`
    Email     string  `f:"text, unique"`
    Username  string  `f:"text, unique, notnull"`
    Active    bool    `f:"bool, default:true"` // if user not active, then user cannot login
    OTPHash   string  `f:"text"`
    OTPExpiry int64   `f:"number"`
    Online    bool    `f:"bool, default:false"` //user becomes online after verification
    Roles     []Roles  `f:"many2many:user_roles"` // many to many relationship with roles One user → many roles , One role → many users
}

type Roles struct {
	ID          int     `f:"number, primary, auto"`
	Name        string  `f:"text, unique, notnull"`
	Description string  `f:"text"`
	Permissions []Permissions `f:"many2many:role_permissions"`
	CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
	UpdatedAt   time.Time `f:"timestamp, default:current_timestamp"`
}


//Permissions read R, write W, delete D, update U,  all A
type Permissions struct {
	ID     int     `f:"number, primary, auto"`
	resource   string  `f:"text, unique, notnull"` // resource can be the name of the model or the name of the API endpoint or the name of the action
	Action string  `f:"text, unique, notnull"` // action can be R, W, D, U, A
    Description string  `f:"text"`
    CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time `f:"timestamp, default:current_timestamp"`
}
