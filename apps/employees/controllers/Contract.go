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
		"status":    "status",
		"start_date": "start_date",
		"end_date":   "end_date",
		
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


func GetSalaryLinesByContractId(id int) []models.ContractSalaryLine{
	var salaryLines []models.ContractSalaryLine

	query := "select id,  base_salary, effective_date, net_salary, gross_salary from contract_salary_lines where contract_id = $1"
	rows, err := core.DB.Query(query, id)
	if err != nil {
		fmt.Print(err)
		return salaryLines
	}
	defer rows.Close()

	for rows.Next() {
		var salaryLine models.ContractSalaryLine
		err := rows.Scan(&salaryLine.ID, &salaryLine.BaseSalary, &salaryLine.EffectiveDate, &salaryLine.NetSalary, &salaryLine.GrossSalary)
		if err != nil {
			fmt.Print(err)
			continue
		}
		salaryLines = append(salaryLines, salaryLine)
	}

	return salaryLines
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
		GrossSalary :=r.PostFormValue("GrossSalary")

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
				query = "insert into contract_salary_lines (contract_id, base_salary, net_salary, effective_date, gross_salary) values($1,$2,$3,$4,$5) RETURNING id"
				err = core.DB.QueryRow(query, ContractId, BaseSalary,NetSalary, EffectiveDate, GrossSalary).Scan(&salaryLinesId)

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


     http.Redirect(w, r, "/employees/contract", http.StatusSeeOther)

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
            c.id, c.employee_id, c.name, c.start_date, c.end_date, c.active, c.created_at, c.shift_schedule_id, c.status, c.iban, c.bank_name, c.absense_balance, c.yearly_total_allocation_days, c.accrual_rate_per_day, c.total_service_years,c.work_location,c.job_title_id,
            sl.base_salary, sl.net_salary, sl.effective_date, sl.gross_salary,
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
		var ComType models.SalaryComponentType
        
        // Using pointers or sql.Null types for fields that could be NULL if a LEFT JOIN finds no match
        var baseSalary, netSalary *string
        var effectiveDate *time.Time
        var typeID *string
        var amount *string
		var employeeId string
		var shiftId string
		var jobTitleId string

        err := rows.Scan(
            &contract.ID, &employeeId, &contract.Name, &contract.StartDate, &contract.EndDate, &contract.Active, &contract.CreatedAt, &shiftId, &contract.Status,&contract.IBAN,&contract.BankName,&contract.AbsenseBalance,&contract.YearlyTotalAllocationDays,&contract.AccrualRatePerDay,&contract.TotalServiceYears,&contract.WorkLocation,&jobTitleId,
            &baseSalary, &netSalary, &effectiveDate, &salaryLine.GrossSalary,
            &typeID, &amount,
        )
        if err != nil {
            fmt.Print(err)
        }

		Employee = GetEmployeeById(employeeId)
		Shift    = GetShiftById(shiftId)
		contract.Employee = &Employee
		contract.ShiftSchedule = &Shift

		JobTitle := GetJobTitleById(jobTitleId)
		contract.JobTitle = &JobTitle

        // On the very first row, grab the salary line details if they exist
        if isFirstRow && baseSalary != nil {
            salaryLine.BaseSalary = *baseSalary
            salaryLine.NetSalary = *netSalary
            salaryLine.EffectiveDate = *effectiveDate
            isFirstRow = false
        }

        // If a component exists in this row, append it to our temporary slice
        if typeID != nil && amount != nil {

			ComType = GetComponentTypeById(*typeID)
            comp.Type = &ComType
            comp.Amount = *amount
            salaryLine.Components = append(salaryLine.Components, comp)
        }
    }

    // Finally, nest the populated structs

	


    contract.SalaryLines = []models.ContractSalaryLine{salaryLine}

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


func GetContractByEmployeeId(employeeId string) models.Contract {
	var contract models.Contract	
	query := "SELECT id, employee_id, name, start_date, end_date,active, created_at, shift_schedule_id, status FROM contracts WHERE employee_id = $1 AND end_date IS NULL"

	row := core.DB.QueryRow(query, employeeId)

	var EmployeeId *string
	var ShiftId *string

	err := row.Scan(&contract.ID, &EmployeeId, &contract.Name, &contract.StartDate, &contract.EndDate, &contract.Active, &contract.CreatedAt, &ShiftId, &contract.Status)
	if err != nil {
		fmt.Print(err)
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

	return contract
}


func TerminateContract(w http.ResponseWriter, r *http.Request){

	id:= r.PathValue("id")
	contract:= GetContractById(id)

	if r.Method == http.MethodGet{

		

		lastSalaryLine := GetLastActiveSalaryDetailsByContractId(id)
		
		data := map[string]interface{}{
		"Title": "Terminate Contract",
	    "lastSalaryLine": lastSalaryLine,
		"contract": contract,
		
	}

	core.RenderPage(w,r, "apps/employees/views/Contract-end.html", data)

	}
		

	if r.Method == http.MethodPost{

		//first make sure the contract is active and not terminated before because we can not terminate contract more than one time if the contract already terminated we will show error message to user that this contract already terminated and you can not terminate it again
		if contract.Status != nil && *contract.Status == "terminated" {
			// 1. Set the correct HTTP Status code (400 Bad Request)
			w.WriteHeader(http.StatusBadRequest)

			// 2. Prepare your user-friendly data payload
			data := map[string]interface{}{
				"Title":   "Contract Already Terminated",
				"Message": "This contract has already been terminated. You cannot terminate a contract more than once.",
			}

			// 3. Render your error page view
			core.RenderPage(w, r, "apps/employees/views/error.html", data)
			return
		}

		//we will update the contract and we will generate clearance document and calculate EOSB if applicable
		//and we will open with response 2 pages one for clearance document and one for EOSB details if applicable
		// 1- update contract (contracts table)
		// 2- insert record in clerance_documents table with type end of service clearance and status draft
		// 3- insert record in eosb_records table if applicable (if the employee is eligible for EOSB based on his service period and the company policy) with status draft
		// 4- insert record in service_periods table with end_date = end date of the contract and note = "end date based on contract termination"
		// 5- redirect to clearance document details page and eosb details page if applicable
        end_date := r.PostFormValue("end_date")
        contract_duration_years := r.PostFormValue("contract_duration_years")
		notice_period := r.PostFormValue("notice_period")
		notification_date := r.PostFormValue("notification_date")
        reason := r.PostFormValue("reason")
		AbsenseBalance := r.PostFormValue("AbsenseBalance")
		ClearanceTemplate := r.PostFormValue("ClearanceTemplate") // ClearanceTemplate Id
		LegalRuleApplied := r.PostFormValue("LegalRuleApplied")
		TotalServiceDuration := r.PostFormValue("TotalServiceDuration")
		BaseAward := r.PostFormValue("BaseAward")
		ReasonModifierMultiplier := r.PostFormValue("ReasonModifierMultiplier")
		AdjustedFinalEOSBReward := r.PostFormValue("AdjustedFinalEOSBReward")
		UnusedLeaveCashout := r.PostFormValue("UnusedLeaveCashout")
		TotalNetSettlementAmount := r.PostFormValue("TotalNetSettlementAmount")

		fmt.Printf("\n End Date: %s, Contract Duration (Years): %s, Notice Period: %s, Notification Date: %s, Reason: %s, Absence Balance: %s, Clearance Template: %s, Legal Rule Applied: %s, Total Service Duration: %s, Base Award: %s, Reason Modifier Multiplier: %s, Adjusted Final EOSB Reward: %s, Unused Leave Cashout: %s, Total Net Settlement Amount: %s\n",
			end_date, contract_duration_years, notice_period, notification_date, reason, AbsenseBalance, ClearanceTemplate, LegalRuleApplied, TotalServiceDuration, BaseAward, ReasonModifierMultiplier, AdjustedFinalEOSBReward, UnusedLeaveCashout, TotalNetSettlementAmount)
		
			// 1- update contract (contracts table)
			query := "update contracts set end_date = $1, active = false, status = 'terminated', total_service_years = $2, notification_date = $3 where id = $4"
			_, err := core.DB.Exec(query, end_date, contract_duration_years, notification_date, id)
			if err != nil {
				fmt.Print(err)
			}

			 // 2- insert record in clerance_documents table with type end of service clearance and status draft
			 var clearanceDocId string
			 query = "insert into clearance_documents (employee_id, contract_id,  status, template_id, requested_date) values ($1, $2, 'draft', $3, $4) RETURNING id"
			 err = core.DB.QueryRow(query, contract.Employee.ID, id, ClearanceTemplate, notification_date).Scan(&clearanceDocId)
			 if err != nil {
				 fmt.Print(err)
			 }

			 //ok base on the tamplate we will generate clearance items for this clearance document
			 tempalateItems := GetClearanceTemplateItemsById(ClearanceTemplate)
			 for _, item := range tempalateItems {
				query:= "insert into clearance_document_items (document_id, name, department_id, status) values ($1, $2, $3, 'pending')"
				_, err = core.DB.Exec(query, clearanceDocId, item.Name, item.Department.ID)
			 }

			 // clearance document created with items with id clearanceDocId will show this for user to fill the missing details and submit for approval

			 // 3- insert record for eosb_records if applicable
			 // check TotalNetSettlementAmount if it is greater than 0 then we will create EOSB record otherwise we will not create EOSB record because the employee is not eligible for EOSB based on the company policy or his service period
			 totalNetSettlementAmountFloat, err := strconv.ParseFloat(TotalNetSettlementAmount, 64)
			 if err != nil {
				 fmt.Print(err)
			 }

			 var ServicePeriodId string
			 var EOSBRecordId string
			 query = "insert into service_periods (employee_id, hire_date, termination_date, note, reason) values ($1, $2, $3, $4, $5) RETURNING id"
			 err = core.DB.QueryRow(query, contract.Employee.ID, contract.StartDate, end_date, "Service period record for EOSB calculation based on contract termination", reason).Scan(&ServicePeriodId)
			 if err != nil {
				 fmt.Print(err)
			 }

             
			 if totalNetSettlementAmountFloat > 0 {
				// insert record in eosb_records table with status draft
				 query = "insert into eosb_records (employee_id, contract_id, service_period_id, base_eosb,reason_modifier_multiplier,adjusted_final_eosb_reward, unused_leave_cashout, final_payable,legal_rule_applied ) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"
				 err = core.DB.QueryRow(query, contract.Employee.ID, id, ServicePeriodId, BaseAward, ReasonModifierMultiplier, AdjustedFinalEOSBReward, UnusedLeaveCashout, TotalNetSettlementAmount, LegalRuleApplied).Scan(&EOSBRecordId)
				 if err != nil {
					 fmt.Print(err)
				 }
			 }

			 //update service_periods set eosb_record_id
			 query = "update service_periods set eosb_record_id = $1 where id = $2"
			 _, err = core.DB.Exec(query, EOSBRecordId, ServicePeriodId)
			 if err != nil {
				 fmt.Print(err)
			 }

			 // redirect to clearance document details page 
			 http.Redirect(w, r, "/clearance/details/"+clearanceDocId, http.StatusSeeOther)
			 




	} //end post request
} //end function


func GetLastActiveSalaryDetailsByContractId(contractId string) models.ContractSalaryLine {
	var salaryLine models.ContractSalaryLine
    query :=`select id, contract_id, base_salary, effective_date, end_date, net_salary, gross_salary 
	         from contract_salary_lines where contract_id = $1 and end_date is null order by effective_date desc limit 1`
	row := core.DB.QueryRow(query, contractId)

	var contId *string

	err := row.Scan(&salaryLine.ID, &contId, &salaryLine.BaseSalary, &salaryLine.EffectiveDate, &salaryLine.EndDate, &salaryLine.NetSalary, &salaryLine.GrossSalary)
	if err != nil {
		fmt.Print(err)
	}

	//Get the components for this salary line
	var salaryComponents []models.SalaryComponentValue

	query = "select id, salary_line_id, type_id, amount from salary_component_values where salary_line_id = $1"
	rows, err := core.DB.Query(query, salaryLine.ID)
	if err != nil {
		fmt.Print(err)
		return salaryLine
	}
	defer rows.Close()

	for rows.Next() {
		var comp models.SalaryComponentValue
		var typeId, salaryLineId *string

		err := rows.Scan(&comp.ID, &salaryLineId, &typeId, &comp.Amount)
		if err != nil {
			fmt.Print(err)
			continue
		}

		if typeId != nil {
			compType := GetComponentTypeById(*typeId)
			comp.Type = &compType
		}

		salaryComponents = append(salaryComponents, comp)
	}

	salaryLine.Components = salaryComponents

	return salaryLine
}