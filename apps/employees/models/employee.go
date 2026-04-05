package models

import "time"



type Employee struct {
    ID         int        `f:"number, primary, auto"`
    BadgeID    string     `f:"text, unique, notnull"`
    Name       string     `f:"text, notnull"`
    Department *Department `f:"many2one:"`   
    LocalName  string     `f:"text"`
    JobTitle   *JobTitle   `f:"many2one:"`   

    CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time `f:"timestamp, default:current_timestamp"`
}

// OrgUnit represents a top-level organizational unit that can contain multiple departments
type OrgUnit struct {
    ID        int        `f:"number, primary, auto"`
    Name      string     `f:"text, unique, notnull"`
    Code      string     `f:"text, unique, notnull"`
    Departments []Department `v:"true"`
    CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    Manager  *Employee    `f:"one2one:employees"`
}


type Department struct {
    ID        int        `f:"number, primary, auto"`
    Name      string     `f:"text, unique, notnull"`
    Code      string     `f:"text, unique, notnull"`
    Employees []Employee `v:"true"`
    Manager   *Employee  `f:"one2one:employees"` 
    CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time `f:"timestamp, default:current_timestamp"`
}

// ExManagerDepartment is a join table to track historical manager assignments for departments
type ExManagerDepartment struct {
    ID         int        `f:"number, primary, auto"`
    Employee   *Employee  `f:"many2one:employees"`
    Department *Department `f:"many2one:departments"`
    CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    From   time.Time `f:"timestamp,"`
    To   time.Time `f:"timestamp,"`
}

type JobTitle struct {
    ID          int    `f:"number, primary, auto"`
    Name        string `f:"text, unique, notnull"`
    Code        string `f:"text, unique, notnull"`
    Description string `f:"text"`
    Employees []Employee `v:"true"`
    CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time `f:"timestamp, default:current_timestamp"`
}