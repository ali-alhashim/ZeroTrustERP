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
    OnboardingProcesses []OnboardingProcess `v:"true"`
    OffboardingProcesses []OffboardingProcess `v:"true"`
    PerformanceReviews []PerformanceReview `v:"true"`
    TrainingSessions   []Training `v:"true"`  
    Tasks              []Task `v:"true"`
    CalendarEvents     []CalendarEvent `v:"true"`
    GovermentID string      `f:"text"` // Iqama/ID
    Image       *string      `f:"text"` 
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
    PersonalEmail *string     `f:"text"`
    PersonalMobile *string     `f:"text"`
    ServicePeriods []ServicePeriod `v:"true"`
    EOSBRecords   *[]EOSBRecord `v:"true"`
}

type ServicePeriod struct {
    ID               int        `f:"number, primary, auto"`
    Employee        *Employee  `f:"many2one:employees"`
    HireDate        time.Time `f:"timestamp"`
    TerminationDate *time.Time `f:"timestamp"`
    Reason          *string   `f:"text"`
    EOSBPaid        *bool   `f:"bool, default:false"`
    Note            *string     `f:"text"`
    EOSBRecord      *EOSBRecord `f:"one2one:eosb_records"`
}

type EOSBRecord struct {
    ID          int        `f:"number, primary, auto"`
    Employee    *Employee  `f:"many2one:employees"`
    ServicePeriod *ServicePeriod `f:"many2one:service_periods"`

    BaseEOSB      string            `f:"money"` // The initial calculated EOSB before adjustments
    ReasonModifierMultiplier string `f:"text"` // e.g., "100%", "50%", "0%" based on the reason for termination
    AdjustedFinalEOSBReward string  `f:"money"` // BaseEOSB adjusted by the ReasonModifierMultiplier
    UnusedLeaveCashout string       `f:"money"` // Cash value of any unused leave days that are added to the EOSB
    TotalLines    *string            `f:"money"` // Net sum of all additions and deductions
    FinalPayable  *string            `f:"money"` // Final amount to pay: BaseEOSB + TotalLines
    PaidEOSB      *string            `f:"money"` // The actual EOSB amount paid out
    LegalRuleApplied string            `f:"text"`  // Reference to the specific labor law article or company policy applied in this calculation
    Status        *string            `f:"text"`      // e.g., "draft", "approved", "paid"
    Lines         []*EOSBRecordLine `f:"one2many:eosb_record_lines"` // Detailed breakdown of additions and deductions

    PaymentDate *time.Time `f:"timestamp"` // When the EOSB was paid out

    FilePath    *string     `f:"text"`  // Path to the stored EOSB document (e.g., PDF of the EOSB calculation and payment details)
    CreatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    Contract    *Contract          `f:"many2one:contracts"` // Link to the employee's contract for reference
}

type EOSBRecordLine struct {
	ID           int         `f:"number, primary, auto"`
	EOSBRecord   *EOSBRecord `f:"many2one:eosb_records"` // Links back to parent
	
	// Line Details
	Type         string      `f:"text"`  // "addition" or "deduction"
	Reason       string      `f:"text"`  // e.g., "Laptop non-return", "Company car damage", "Unpaid loan"
	Amount       string      `f:"money"` // The value of this specific adjustment
	
	// Optional: Asset tracking reference if applicable
	AssetSerialNumber      *string        `f:"text, nullable"` //laptop, phone, etc. serial number if this line is related to an asset or Car VID
	
	CreatedAt    time.Time   `f:"timestamp, default:current_timestamp"`
	UpdatedAt    time.Time   `f:"timestamp, default:current_timestamp"`
    FilePath     *string      `f:"text"`  // Path to the stored document related to this line (e.g., damage report, asset return receipt)
    Note         *string      `f:"text"`  // Additional notes about this specific line item
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

type ExDepartment struct {
    ID            int        `f:"number, primary, auto"`
    Department   *Department `f:"many2one:departments"`
    Employee     *Employee  `f:"many2one:employees"`
    CreatedAt    time.Time `f:"timestamp, default:current_timestamp"`
    UpdatedAt    time.Time `f:"timestamp, default:current_timestamp"`
    StartDate    time.Time  `f:"timestamp"` // When the employee started in this department
    EndDate      time.Time  `f:"timestamp"` // When the employee left this department 
}

// ExManagerDepartment is a join table to track historical manager assignments for departments
type ExManagerDepartment struct {
    ID         int        `f:"number, primary, auto"`
    Employee   *Employee  `f:"many2one:employees"`
    Department *Department `f:"many2one:departments"`
    CreatedAt   time.Time `f:"timestamp, default:current_timestamp"`
    StartDate       time.Time `f:"timestamp,"`
    EndDate         time.Time `f:"timestamp,"`
}

// is a join table to track historical jobtitle
type ExJobTitle struct {
    ID         int         `f:"number, primary, auto"`
    Employee   *Employee   `f:"many2one:employees"`
    JobTitle   *JobTitle   `f:"many2one:job_titles"`
    CreatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    StartDate       time.Time   `f:"timestamp,"`
    EndDate         time.Time   `f:"timestamp,"`
}

type JobTitle struct {
    ID          int    `f:"number, primary, auto"`
    Name        *string `f:"text"`
    LocalName   *string     `f:"text"`
    Code         string `f:"text, unique, notnull"`
    Description *string `f:"text"`
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
    NotificationDate *time.Time `f:"timestamp"` // Nullable for open-Notification Date the date to send notification for contract renewal or termination
    // Virtual field to see all salary updates linked to this contract
    SalaryLines []ContractSalaryLine `v:"true"`
    JobTitle    *JobTitle   `f:"many2one:"`   
    Active      bool       `f:"bool, default:true"`
    CreatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    IBAN        *string    `f:"text"` // Bank account number for salary payments
    BankName    *string    `f:"text"` // Bank name for salary payments
    AbsenseBalance *int    `f:"number"`
    YearlyTotalAllocationDays *int `f:"number"`
    AccrualRatePerDay *float64 `f:"number(12,2)"`
    Status      *string    `f:"text"`
    ShiftSchedule  *ShiftSchedule `f:"many2one:"`
    WorkLocation *string   `f:"text"`
    Note         *string   `f:"text"`

    CompanyAuthoritySignature *string `f:"text"`   //this path of image Signature
    EmployeeSignature         *string `f:"text"`   //this path of image Signature

    TotalServiceYears *float64 `f:"number(12,2)"`
    IBANChangeHistory *[]ExContractIBAN `v:"true"` // To track historical IBAN changes for this contract
    IBANFilePath *string `f:"text"` // Path to the stored document showing the current IBAN details
}

type ExContractIBAN struct {
    ID          int        `f:"number, primary, auto"`
    Contract    *Contract  `f:"many2one:contracts"`
    OldIBAN     string     `f:"text, notnull"`
    NewIBAN     string     `f:"text, notnull"`
    ChangedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    OldIBANFilePath *string    `f:"text"` // Path to the stored document showing the old IBAN details
    NewIBANFilePath *string    `f:"text"` // Path to the stored document showing the new IBAN details
    UpdatedBy   *Employee  `f:"many2one:employees"` // Who made the change
    VerificationStatus string `f:"text"` // "Pending", "Verified", "Rejected"
    ConfirmationBy *Employee `f:"many2one:employees"` // Who confirmed the change after verification
}

// Accrual Rate Per Day = Yearly Total Allocation Days / Days in year
// Cron Job
// for the year some time 366 & 365 we will stick to a fixed 365

//we have 2 documents end of service document this showing the total service years and the end of service benefits and the other one is the contract document showing the contract details and the signatures of both parties
// Clearnce document is for the employee when he is leaving the company and it will show the clearance status of the employee and the pending tasks if there are any and the final settlement amount if applicable
// also digital confirmation from department managers like IT, HR, Finance, etc. for the clearance process

type ClearanceTemplate struct {
    ID          int        `f:"number, primary, auto"`
    Name        string     `f:"text, unique, notnull"` // e.g., "Standard Clearance Template"
    Description string     `f:"text"`                 // Details about the template
    CreatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    Items       *[]ClearanceTemplateItem `v:"true"` // The specific clearance items that need to be completed for this template
}

type ClearanceTemplateItem struct {
    ID          int        `f:"number, primary, auto"`
    Template    *ClearanceTemplate `f:"many2one:clearance_templates, notnull"`
    Name        string     `f:"text, notnull"` // e.g., "IT Clearance", "HR Clearance", "Finance Clearance"
    Department  *Department `f:"many2one:departments"` // Optional link to department responsible for this clearance item
    CreatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    
}

type ClearanceDocument struct {
    ID                int                `f:"number, primary, auto"`
    Employee          *Employee          `f:"many2one:employees, notnull"`
    Template          *ClearanceTemplate `f:"many2one:clearance_templates, notnull"`
    Status            string             `f:"text, notnull, default:'Pending'"` // "Pending", "In Progress", "Completed", "Rejected"
    RequestedDate     time.Time          `f:"timestamp, default:current_timestamp"`
    CompletedDate     *time.Time         `f:"timestamp"` // Nullable until fully cleared
    CreatedAt         time.Time          `f:"timestamp, default:current_timestamp"`
    UpdatedAt         time.Time          `f:"timestamp, default:current_timestamp"`
    Contract          *Contract          `f:"many2one:contracts"` // Link to the employee's contract for reference
    ClearanceItems    *[]ClearanceDocumentItem `v:"true"` // The specific clearance items that need to be completed for this employee
}

type ClearanceDocumentItem struct {
    ID                  int                `f:"number, primary, auto"`
    Document            *ClearanceDocument `f:"many2one:clearance_documents, notnull"`
    Name                string             `f:"text, notnull"` // Copied from ClearanceTemplateItem.Name (e.g., "Return IT Assets")
    Department          *Department        `f:"many2one:departments"` 
    Status              string             `f:"text, notnull, default:'Pending'"` // "Pending", "In Progress", "Completed"
    ResponsibleEmployee *Employee          `f:"many2one:employees"` // Who needs to sign this specific off
    Note                *string            `f:"text"` // Notes added during the sign-off process
    FilePath            string             `f:"text"` // Path to uploaded proof/documents if required
    HandledAt           *time.Time         `f:"timestamp"` // When this specific item was signed off
    CreatedAt           time.Time          `f:"timestamp, default:current_timestamp"`
    UpdatedAt           time.Time          `f:"timestamp, default:current_timestamp"`
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
    
    BaseSalary    string   `f:"money, notnull"`
    GrossSalary   string   `f:"money"` // Base + Allowances
    NetSalary     string   `f:"money, notnull"`
    
    // This connects to the individual allowances for this specific salary update
    Components    []SalaryComponentValue `v:"true"` 
    
    EffectiveDate time.Time `f:"timestamp, notnull"`
    EndDate       *time.Time `f:"timestamp"` // When this salary line is no longer effective (e.g., after a raise)
    CreatedAt     time.Time `f:"timestamp, default:current_timestamp"`
}

// SalaryComponentValue stores the actual amount for a specific employee
type SalaryComponentValue struct {
    ID          int                  `f:"number, primary, auto"`
    SalaryLine  *ContractSalaryLine  `f:"many2one:contract_salary_lines"`
    Type        *SalaryComponentType `f:"many2one:salary_component_types"`
    Amount      string              `f:"money, notnull"`
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
    GovernmentId *string     `f:"text"`  // Iqama/ID for the family member, if applicable
    ContactNumber *string     `f:"text"`  // Phone number for the family member
    Relationship string     `f:"text"` // e.g., "Spouse", "Child", "Parent"
    Gender       string      `f:"text"`
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

type OnboardingProcess struct {
    ID          int        `f:"number, primary, auto"`
    Employee    *Employee  `f:"many2one:employees, notnull"`
    StepName    string     `f:"text, notnull"` // e.g., "Document Submission", "Orientation", "Training"
    Status      string     `f:"text"`         // e.g., "Pending", "Completed", "In Progress"
    DueDate     time.Time  `f:"timestamp"`    // Deadline for completing this step
    CreatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
}

type OffboardingProcess struct {
    ID          int        `f:"number, primary, auto"`
    Employee    *Employee  `f:"many2one:employees, notnull"`
    StepName    string     `f:"text, notnull"` // e.g., "Exit Interview", "Asset Return", "Final Settlement"
    Status      string     `f:"text"`         // e.g., "Pending", "Completed", "In Progress"
    DueDate     time.Time  `f:"timestamp"`    // Deadline for completing this step
    CreatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
}


type Training struct {
    ID          int        `f:"number, primary, auto"`
    Employee    *Employee  `f:"many2one:employees, notnull"`
    Name        string     `f:"text, notnull"` // e.g., "Safety Training", "Leadership Workshop"
    Description string     `f:"text"`         // Details about the training
    Date        time.Time  `f:"timestamp"`    // Date of the training session
    CreatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
}


type PerformanceReview struct {
    ID          int        `f:"number, primary, auto"`
    Employee    *Employee  `f:"many2one:employees, notnull"`
    Reviewer    *Employee  `f:"many2one:employees"` // The person conducting the review
    Date        time.Time  `f:"timestamp"`          // Date of the review
    Rating      int        `f:"number"`             // e.g., 1-5 rating
    Comments    string     `f:"text"`               // Feedback and notes from the review
    CreatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
}


type Task struct {
    ID          int        `f:"number, primary, auto"`
    Employee    *Employee  `f:"many2one:employees, notnull"`
    Title       string     `f:"text, notnull"` // e.g., "Complete Onboarding Paperwork"
    Description string     `f:"text"`         // Details about the task
    DueDate     time.Time  `f:"timestamp"`    // Deadline for the task
    Status      string     `f:"text"`         // e.g., "Pending", "Completed", "In Progress"
    CreatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
}



type CalendarEvent struct {
    ID          int        `f:"number, primary, auto"`
    Employee    *Employee  `f:"many2one:employees, notnull"`
    Title       string     `f:"text, notnull"` // e.g., "Team Meeting", "Performance Review"
    Description string     `f:"text"`         // Details about the event
    StartTime   time.Time  `f:"timestamp"`    // Event start time
    EndTime     time.Time  `f:"timestamp"`    // Event end time
    CreatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
    UpdatedAt   time.Time  `f:"timestamp, default:current_timestamp"`
}


 
