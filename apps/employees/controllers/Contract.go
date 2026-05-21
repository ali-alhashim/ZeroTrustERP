package controllers

import(
	"fmt"
	"net/http"
	"zerotrusterp/apps/employees/models"
	"zerotrusterp/core"
	"strconv"
	"time"
)


func ListContract(w http.ResponseWriter, r *http.Request) {


	query := r.URL.Query()

	search     := query.Get("q")
	sortBy     := query.Get("sort")
	order      := query.Get("order")
	page       := query.Get("page")
	pageSize   := query.Get("pageSize")

	fmt.Printf("list logs: search=%s, sort=%s, order=%s, page=%s, pageSize=%s\n", search, sortBy, order, page, pageSize)
    
    totalRecords := core.GetCountRecords("contracts")

	contracts:= GetContractsFromDB(search, sortBy, order, page, pageSize)


	data := map[string]interface{}{
		"Title": "Contract",
		"Query": search,
		"Sort":  sortBy,
		"Order": order,
		"Page":  page,
		"PageSize": pageSize,
		"TotalRecords":totalRecords,
		"Contracts": contracts,
		
	}

	core.RenderPage(w,r, "apps/employees/views/Contract-list.html", data)
}



func GetContractsFromDB(search, sort, order, page, pageSize string) []models.Contract {

	query := "SELECT id, employee_id, name, start_date, end_date,active, created_at, shift_schedule_id, status FROM contracts WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += " AND (employee_id ILIKE $" + strconv.Itoa(argIndex) +
			     " OR name ILIKE $" + strconv.Itoa(argIndex+1) + ")"

		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}


	// 🔒 Safe sorting
	allowedSort := map[string]string{
		"ID":        "id",
		"employee_id":     "employee_id",
		"Name":    "name",
		"LocalName":  "local_name",
		"Manager":  "manager_id",
		"Active":    "active",
		
	}

	if col, ok := allowedSort[sort]; ok {
		query += " ORDER BY " + col
		if order == "desc" {
			query += " DESC"
		} else {
			query += " ASC"
		}

	}

	


		// 📄 Pagination (page + pageSize)
	p, _ := strconv.Atoi(page)
	ps, _ := strconv.Atoi(pageSize)

	// defaults
	if p <= 0 {
		p = 1
	}
	if ps <= 0 || ps > 100 {
		ps = 10
	}

	offset := (p - 1) * ps

	query += " LIMIT $" + strconv.Itoa(argIndex) +
		" OFFSET $" + strconv.Itoa(argIndex+1)

	args = append(args, ps, offset)


		// ✅ Execute
	rows, err := core.DB.Query(query, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()


	



	var contracts []models.Contract
	
   

	for rows.Next() {
		
		var EmployeeId *string
		var ShiftId *string
		var contract models.Contract
	
		err := rows.Scan(&contract.ID, &EmployeeId, &contract.Name, &contract.StartDate, &contract.EndDate, &contract.Active, &contract.CreatedAt, &ShiftId, &contract.Status)
		if err != nil {
			panic(err)
		}

		if EmployeeId != nil {
			fmt.Print(" \n Get employee by ID = "+ *EmployeeId +"\n")
			theEmployee := GetEmployeeById(*EmployeeId)
			fmt.Print(" \n Get Employee  Contract  = "+ theEmployee.Name +"\n")
			contract.Employee = &theEmployee
		} else {
			contract.Employee = nil
		}


		if ShiftId != nil {
			fmt.Print(" \n Get ShiftId by ID = "+ *ShiftId +"\n")
			theShift := GetShiftById(*ShiftId)
			fmt.Print(" \n Get theShift    = "+ theShift.Name +"\n")
			contract.ShiftSchedule = &theShift
		} else {
			contract.ShiftSchedule = nil
		}


		//Get the Salary ContractSalaryLine and set to this contract also the compoenents
		contract.SalaryLines = GetSalaryLinesByContractId(contract.ID)


        
		contracts = append(contracts, contract)

	}

	return contracts
}


func GetSalaryLinesByContractId(id int) models.ContractSalaryLine{
	var salaryLine models.ContractSalaryLine

	query := "select id,  base_salary, effective_date, net_salary from contract_salary_lines where contract_id = $1"
	err := core.DB.QueryRow(query, id).Scan(&salaryLine.ID, &salaryLine.BaseSalary, &salaryLine.EffectiveDate, &salaryLine.NetSalary)
	if err !=nil{
		fmt.Print(err)
	}


	return salaryLine
}



func isEmployeeHaveAnyValidContract(id string) bool {
    var count int

   
    query := "SELECT COUNT(*) FROM contracts WHERE employee_id = $1 AND end_date IS NULL"

    // Execute the query and Scan the count directly into our variable
    err := core.DB.QueryRow(query, id).Scan(&count)
    if err != nil {
        // Log the error or handle it based on your core package setup
        return false
    }

    // If count is greater than 0, they have an active contract
    return count > 0
}


func CreateContract(w http.ResponseWriter, r *http.Request){

	if r.Method == http.MethodGet{

		data := map[string]interface{}{
		"Title": "Create Contract",
		}

	core.RenderPage(w,r, "apps/employees/views/Contract-create.html", data)

	} //--End Get

	if r.Method == http.MethodPost{
		//insert contract to db

		contractName:= r.PostFormValue("name")
		employeeId:= r.PostFormValue("employee")


		// before you create contract make sure the employee dose not have any valid contract all his contract shuld first be end if not 
		// don't allow to create more that contract for the same employee HR must terminate exist contract first 
		if isEmployeeHaveAnyValidContract(employeeId) {
    // 1. Set the correct HTTP Status code (400 Bad Request)
    w.WriteHeader(http.StatusBadRequest)

    // 2. Prepare your user-friendly data payload
    data := map[string]interface{}{
        "Title":   "Contract Conflict",
        "Message": "This employee already has an active contract. An employee cannot have multiple valid contracts at the same time. Please terminate the current contract before creating a new one.",
    }

    // 3. Render your error page view
    core.RenderPage(w, r, "apps/employees/views/error.html", data)
    return
}

        start_date, errStr := time.Parse("2006-01-02", r.PostFormValue("start_date"))
		if errStr !=nil{
			fmt.Print(errStr)
		}

		var end_date *time.Time
		endDateStr := r.PostFormValue("end_date")
		if endDateStr != "" {
			parsedDate, errStr2 := time.Parse("2006-01-02", endDateStr)
			if errStr2 !=nil{
				fmt.Print(errStr2)
			} else {
				end_date = &parsedDate
			}
		}
		
		
		totalServiceYears:=r.PostFormValue("totalServiceYears")
		jobTitleId:=r.PostFormValue("jobTitleId")
		IBAN:=r.PostFormValue("IBAN")
		BankName:=r.PostFormValue("BankName")
        AbsenseBalance:=r.PostFormValue("AbsenseBalance")
		YearlyTotalAllocationDays:= r.PostFormValue("YearlyTotalAllocationDays")
        AccrualRatePerDay:= r.PostFormValue("AccrualRatePerDay")
		Status:=r.PostFormValue("Status")
		ShiftSchedule:=r.PostFormValue("ShiftSchedule")
		workLocation:=r.PostFormValue("workLocation")
		note:=r.PostFormValue("note")
		active:=r.PostFormValue("active") =="on"

		//----Salary
		BaseSalary :=r.PostFormValue("BaseSalary")
		NetSalary  :=r.PostFormValue("NetSalary")

		EffectiveDate, errStr3 :=time.Parse("2006-01-02", r.PostFormValue("EffectiveDate"))
		if errStr3 !=nil{
			fmt.Print(errStr3)
		}


		SalaryComponentType:=r.PostForm["SalaryComponentType"]
		Amount:=r.PostForm["Amount"]

		


		 // if this first contract for the employee then we create record for him in ServicePeriod
		 query :=`insert into contracts 
		         (employee_id,
				   name,
				   start_date, 
				   end_date, 
				   active, 
				   iban, 
				   bank_name, 
				   absense_balance, 
				   status, 
				   shift_schedule_id,
				   job_title_id, 
				   work_location,
				   note,
				   yearly_total_allocation_days,
				   accrual_rate_per_day,
				   total_service_years
					) 
					values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15, $16) RETURNING id`

				var ContractId int
				err := core.DB.QueryRow(query,employeeId,contractName,start_date,end_date,active,IBAN,BankName,AbsenseBalance,Status,ShiftSchedule,jobTitleId,workLocation,note,YearlyTotalAllocationDays,AccrualRatePerDay,totalServiceYears).Scan(&ContractId)
				
				if err != nil {
					fmt.Print(err)
				}

				fmt.Printf("\n Contract Created with ID: %d \n ", ContractId)


				// Insert Salary lines row in contract_salary_lines(contract_id, base_salary, effective_date)
                
				var salaryLinesId int
				query = "insert into contract_salary_lines (contract_id, base_salary, net_salary, effective_date) values($1,$2,$3,$4) RETURNING id"
				err = core.DB.QueryRow(query, ContractId, BaseSalary,NetSalary, EffectiveDate).Scan(&salaryLinesId)

				if err != nil {
					fmt.Print(err)
				}

				 
				// Insert salary_component_values ( salary_line_id, type_id, amount) this may has many so we use loop
				for i:= range SalaryComponentType{

					query = "insert into salary_component_values (salary_line_id, type_id, amount) values ($1, $2, $3)"
					_,err = core.DB.Exec(query, salaryLinesId, SalaryComponentType[i], Amount[i])
					if err !=nil{
						fmt.Print(err)
					}


				}


		 if isThisFirstContract(employeeId){
			//Insert record in service_periods (employee_id, hire_date, note)
			query = "insert into service_periods (employee_id, hire_date, note) values ($1, $2, $3)"
			_,err = core.DB.Exec(query, employeeId, start_date, "hire date base on first contract created [Start Date of first Contract]")
           if err !=nil{
						fmt.Print(err)
					}

		 }




	} //end post request

}


func isThisFirstContract(employeeId string) bool{
var count int
	query := "select count(*) from contracts where employee_id = $1"

	// Execute the query and scan the result into the count variable
	err := core.DB.QueryRow(query, employeeId).Scan(&count)
	if err != nil {
		// Handle the error according to your application's logging/error strategy
		fmt.Printf("Error checking employee contracts: %v", err)
		return false 
	}

	// If count is 0, it is their first contract
	return count == 0

}


func GetContractById(id string) models.Contract {
    var contract models.Contract
    var salaryLine models.ContractSalaryLine
	var Employee models.Employee
	var Shift models.ShiftSchedule
    
    query := `
        SELECT 
            c.id, c.employee_id, c.name, c.start_date, c.end_date, c.active, c.created_at, c.shift_schedule_id, c.status,
            sl.base_salary, sl.net_salary, sl.effective_date,
            scv.type_id, scv.amount
        FROM contracts c
        LEFT JOIN contract_salary_lines sl ON c.id = sl.contract_id
        LEFT JOIN salary_component_values scv ON sl.id = scv.salary_line_id
        WHERE c.id = $1
    `

    rows, err := core.DB.Query(query, id)
    if err != nil {
        fmt.Print(err)
    }
    defer rows.Close()

    isFirstRow := true

    for rows.Next() {
        var comp models.SalaryComponentValue
        
        // Using pointers or sql.Null types for fields that could be NULL if a LEFT JOIN finds no match
        var baseSalary, netSalary *string
        var effectiveDate *time.Time
        var typeID *int
        var amount *string
		var employeeId string
		var shiftId string

        err := rows.Scan(
            &contract.ID, &employeeId, &contract.Name, &contract.StartDate, &contract.EndDate, &contract.Active, &contract.CreatedAt, &shiftId, &contract.Status,
            &baseSalary, &netSalary, &effectiveDate,
            &typeID, &amount,
        )
        if err != nil {
            fmt.Print(err)
        }

		Employee = GetEmployeeById(employeeId)
		Shift    = GetShiftById(shiftId)
		contract.Employee = &Employee
		contract.ShiftSchedule = &Shift

        // On the very first row, grab the salary line details if they exist
        if isFirstRow && baseSalary != nil {
            salaryLine.BaseSalary = *baseSalary
            salaryLine.NetSalary = *netSalary
            salaryLine.EffectiveDate = *effectiveDate
            isFirstRow = false
        }

        // If a component exists in this row, append it to our temporary slice
        if typeID != nil && amount != nil {
            comp.ID = *typeID
            comp.Amount = *amount
            salaryLine.Components = append(salaryLine.Components, comp)
        }
    }

    // Finally, nest the populated structs

	


    contract.SalaryLines = salaryLine

    return contract
}




func ContractDetails(w http.ResponseWriter, r *http.Request){

	id:= r.PathValue("id")


	if r.Method == http.MethodGet{

		contract:= GetContractById(id)


		data := map[string]interface{}{
		"Title": "Contract",
	
		"contract": contract,
		
	}

	core.RenderPage(w,r, "apps/employees/views/Contract-details.html", data)


	}



}