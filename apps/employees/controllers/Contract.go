package controllers

import(
	"fmt"
	"net/http"
	"zerotrusterp/apps/employees/models"
	"zerotrusterp/core"
	"strconv"
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


        
		contracts = append(contracts, contract)

	}

	return contracts
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
		start_date:= r.PostFormValue("start_date")
		end_date:=r.PostFormValue("end_date")
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
		EffectiveDate:=r.PostFormValue("EffectiveDate")

		SalaryComponentType:=r.PostForm["SalaryComponentType"]
		Amount:=r.PostForm["Amount"]

		fmt.Printf("Create Contract Data :  \n%s \n %s \n %s \n %s \n %s \n %s\n %s\n %s\n %s \n%s \n%s \n%s\n %s \n%s \n%s\n %t \n%s\n %s\n %v\n %v\n", 
		contractName,
		employeeId,
		start_date,
		end_date,
		totalServiceYears,
		jobTitleId,
		IBAN,
		BankName,
		AbsenseBalance,
		YearlyTotalAllocationDays,
		AccrualRatePerDay,
		Status,
		ShiftSchedule,
		workLocation,
		note,
		active,
		BaseSalary,
		EffectiveDate,
		SalaryComponentType,
		Amount,
	     )


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
					values ($1,$2,$3)`

		 if isThisFirstContract(employeeId){
			//Insert record in service_periods
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