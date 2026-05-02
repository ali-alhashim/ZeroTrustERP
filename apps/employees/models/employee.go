package models

import (
	"time"

	
)



type Employee struct {
    ID         int        `f:"number, primary, auto"`
    BadgeID    string     `f:"text, unique, notnull"`
    Name       string     `f:"text, notnull"`
    Department *Department `f:"many2one:"`   
    LocalName   string     `f:"text"`
    JobTitle    *JobTitle   `f:"many2one:"`   
    Grade       string      `f:"text"` // to distinguish the level of employees like C level or Manager ...Supervisor or just use numbers
    CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    BirthDate   time.Time `f:"timestamp,"`
    Active      bool      `f:"bool, default:true"`
    InsurancePolicies []InsurancePolicy `v:"true"`
    Certifications    []Certification `v:"true"`
    EmergencyContacts []EmergencyContact `v:"true"`
    FamilyMembers     []FamilyMember `v:"true"`
    EmployeeDocuments []EmployeeDocument `v:"true"`
    GovermentID string      `f:"text"` // Iqama/ID
    Image       string      `f:"text"` 
    Email       *string      `f:"text, unique"`
    Nationality *string      `f:"text"`
    Gender       *string      `f:"text"` // "Male", "Female"
    MaritalStatus *string      `f:"text"` // "Single", "Married", "Divorced", "Widowed"
    PhoneNumber  *string      `f:"text"`
    Address      *string      `f:"text"`
    Status       *string      `f:"text"` // "Active", "On Leave", "Terminated"
    Education    *string      `f:"text"` // e.g., "Bachelor's in Computer Science"
    Major        *string      `f:"text"` // e.g., "Computer Science", "Business Administration"
    Religion     *string      `f:"text"` // e.g., "Islam", "Christianity", "Hinduism"
}

// OrgUnit represents a top-level organizational unit that can contain multiple departments
type OrgUnit struct {
    ID          int        `f:"number, primary, auto"`
    Name        string     `f:"text, unique, notnull"`
    Code        string     `f:"text, unique, notnull"`
    Departments []Department `v:"true"`
    CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    Manager     *Employee    `f:"one2one:employees"`
}


type Department struct {
    ID          int        `f:"number, primary, auto"`
    Name        string     `f:"text, unique, notnull"`
    LocalName   string     `f:"text"`
    Code        string     `f:"text, unique, notnull"`
    Employees   []Employee `v:"true"`
    Manager     *Employee  `f:"one2one:employees"` 
    CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    Active      bool       `f:"bool, default:true"`
}

// ExManagerDepartment is a join table to track historical manager assignments for departments
type ExManagerDepartment struct {
    ID         int        `f:"number, primary, auto"`
    Employee   *Employee  `f:"many2one:employees"`
    Department *Department `f:"many2one:departments"`
    CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    FromDate       time.Time `f:"timestamp,"`
    ToDate         time.Time `f:"timestamp,"`
}

// is a join table to track historical jobtitle
type ExJobTitle struct {
    ID         int         `f:"number, primary, auto"`
    Employee   *Employee   `f:"many2one:employees"`
    JobTitle   *JobTitle   `f:"many2one:job_titles"`
    CreatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    FromDate       time.Time   `f:"timestamp,"`
    ToDate         time.Time   `f:"timestamp,"`
}

type JobTitle struct {
    ID          int    `f:"number, primary, auto"`
    Name        string `f:"text, unique, notnull"`
    LocalName   string     `f:"text"`
    Code        string `f:"text, unique, notnull"`
    Description string `f:"text"`
    Employees   []Employee `v:"true"`
    CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    Active      bool       `f:"bool, default:true"`
}


type ShiftSchedule struct {
    ID        int       `f:"number, primary, auto"`
    Name      string    `f:"text, unique, notnull"`
    CreatedAt time.Time `f:"timestamp, default:current_timestamp"`
    UpdatedAt time.Time `f:"timestamp, default:current_timestamp"`
    FromTime  time.Time `f:"timestamp"`
    ToTime    time.Time `f:"timestamp"`
    FromDate  time.Time `f:"timestamp"`
    ToDate    time.Time `f:"timestamp"`

    Monday      bool      `f:"bool, default:false"`
    Tuesday     bool      `f:"bool, default:false"`
    Wednesday   bool      `f:"bool, default:false"`
    Thursday    bool      `f:"bool, default:false"`
    Friday      bool      `f:"bool, default:false"`
    Saturday    bool      `f:"bool, default:false"`
    Sunday      bool      `f:"bool, default:false"`
    Employees   []Employee `v:"true"`
}


// Contract represents the legal employment agreement for an employee.
type Contract struct {
    ID          int        `f:"number, primary, auto"`
    Employee    *Employee  `f:"many2one:employees, notnull"`
    Name        string     `f:"text"` // e.g., "Employment Agreement - Ali"
    StartDate   time.Time  `f:"timestamp, notnull"`
    EndDate     *time.Time `f:"timestamp"` // Nullable for open-ended contracts
    
    // Virtual field to see all salary updates linked to this contract
    SalaryLines []ContractSalaryLine `v:"true"`
    
    Active      bool       `f:"bool, default:true"`
    CreatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
}

// SalaryComponentType defines what the money is for (Housing, Transport, etc.)
type SalaryComponentType struct {
    ID          int    `f:"number, primary, auto"`
    Name        string `f:"text, unique, notnull"` // "Housing", "Transportation", "Telecom"
    Code        string `f:"text, unique, notnull"` // "HOU", "TRA", "TEL"
}

// ContractSalaryLine (Updated to support multiple allowances)
type ContractSalaryLine struct {
    ID            int       `f:"number, primary, auto"`
    Contract      *Contract `f:"many2one:contracts, notnull"`
    
    BaseSalary    float64   `f:"number, notnull"`
    
    // This connects to the individual allowances for this specific salary update
    Components    []SalaryComponentValue `v:"true"` 
    
    EffectiveDate time.Time `f:"timestamp, notnull"`
    CreatedAt     time.Time `f:"timestamp, default:current_timestamp"`
}

// SalaryComponentValue stores the actual amount for a specific employee
type SalaryComponentValue struct {
    ID          int                  `f:"number, primary, auto"`
    SalaryLine  *ContractSalaryLine  `f:"many2one:contract_salary_lines"`
    Type        *SalaryComponentType `f:"many2one:salary_component_types"`
    Amount      float64              `f:"number, notnull"`
}



// InsuranceGrade defines the levels like Class A, Class B, Class C
type InsuranceGrade struct {
    ID          int    `f:"number, primary, auto"`
    Grade       string `f:"text, unique, notnull"` // "A", "B", "C"
    Description string `f:"text"`                   // "Full coverage", "Standard", etc.
    CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
}

// InsurancePolicy links an employee to their specific insurance details
type InsurancePolicy struct {
    ID              int             `f:"number, primary, auto"`
    Employee        *Employee       `f:"many2one:employees, notnull"`
    Grade           *InsuranceGrade `f:"many2one:insurance_grades, notnull"`
    PolicyNumber    string          `f:"text, unique"`
    Provider        string          `f:"text"` // e.g., "Bupa", "Tawuniya"
    
    // Dates are critical for renewals
    StartDate       time.Time       `f:"timestamp"`
    ExpiryDate      time.Time       `f:"timestamp"`
    
    Active          bool            `f:"bool, default:true"`
    CreatedAt       time.Time       `f:"timestamp, default:current_timestamp"`
    UpdatedAt       time.Time       `f:"timestamp, default:current_timestamp"`
}


type Certification struct {
    ID          int        `f:"number, primary, auto"`
    Employee    *Employee  `f:"many2one:employees, notnull"`
    Name        string     `f:"text, notnull"` // e.g., "PMP", "AWS Certified Solutions Architect"
    Issuer      string     `f:"text"`         // e.g., "PMI", "Amazon"
    IssueDate   time.Time  `f:"timestamp"`
    ExpiryDate  time.Time  `f:"timestamp"`
    CreatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    FilePath    string     `f:"text"`  // Path to the stored certificate document
}


type EmergencyContact struct {
    ID          int        `f:"number, primary, auto"`
    Employee    *Employee  `f:"many2one:employees, notnull"`
    Name        string     `f:"text, notnull"`
    Relationship string     `f:"text"` // e.g., "Spouse", "Parent", "Sibling"
    Phone       string     `f:"text"`
    CreatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
}

type FamilyMember struct {
    ID          int        `f:"number, primary, auto"`
    Employee    *Employee  `f:"many2one:employees, notnull"`
    Name        string     `f:"text, notnull"`
    GovernmentId string     `f:"text"`  // Iqama/ID for the family member, if applicable
    ContactNumber string     `f:"text"`  // Phone number for the family member
    Relationship string     `f:"text"` // e.g., "Spouse", "Child", "Parent"
    BirthDate   time.Time  `f:"timestamp"`
    CreatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    FilePath    string     `f:"text"`  // Path to the stored document (e.g., birth certificate for a child)
}

type EmployeeDocument struct {
    ID          int        `f:"number, primary, auto"`
    Employee    *Employee  `f:"many2one:employees, notnull"`
    Name        string     `f:"text, notnull"` // e.g., "Passport", "Iqama"
    Type        string     `f:"text"`         // e.g., "Identification", "Work Permit"
    FilePath    string     `f:"text"`         // Path to the stored document
    ExpiryDate  time.Time  `f:"timestamp"`    // Important for documents like Iqama
    CreatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
}

