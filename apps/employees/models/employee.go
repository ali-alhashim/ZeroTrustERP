package models



type Employee struct {
    ID         int        `f:"number, primary, auto"`
    BadgeID    string     `f:"text, unique, notnull"`
    Name       string     `f:"text, notnull"`
    Department *Department `f:"many2one:"`   
    LocalName  string     `f:"text"`
    JobTitle   *JobTitle   `f:"many2one:"`   
}

type Department struct {
    ID        int        `f:"number, primary, auto"`
    Name      string     `f:"text, unique, notnull"`
    Code      string     `f:"text, unique, notnull"`
    Employees []Employee `v:"true"`
    Manager   *Employee  `f:"one2one:employees"` 
}

type JobTitle struct {
    ID          int    `f:"number, primary, auto"`
    Name        string `f:"text, unique, notnull"`
    Code        string `f:"text, unique, notnull"`
    Description string `f:"text"`
    Employees []Employee `v:"true"`
}