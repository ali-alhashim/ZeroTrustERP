package controllers

import (
	"net/http"
	"strconv"
	"fmt"
	"zerotrusterp/apps/employees/models"
	"zerotrusterp/core"
)


func ListEOSB(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	search     := query.Get("q")
	sortBy     := query.Get("sort")
	order      := query.Get("order")
	page       := query.Get("page")
	pageSize   := query.Get("pageSize")

	fmt.Printf("list : search=%s, sort=%s, order=%s, page=%s, pageSize=%s\n", search, sortBy, order, page, pageSize)
    
    totalRecords := core.GetCountRecords("eosb_records")

	eosb:= GetEOSBFromDB(search, sortBy, order, page, pageSize)


	data := map[string]interface{}{
		"Title": "EOSB Records",
		"Query": search,
		"Sort":  sortBy,
		"Order": order,
		"Page":  page,
		"PageSize": pageSize,
		"TotalRecords":totalRecords,
		"eosb": eosb,
		
	}

	core.RenderPage(w,r, "apps/employees/views/eosb-list.html", data)
}


func GetEOSBFromDB(search, sort, order, page, pageSize string) []models.EOSBRecord {

	query :="select id, employee_id, service_period_id, base_eosb, total_lines, final_payable, paid_eosb, status,payment_date, file_path, created_at, updated_at, contract_id, reason_modifier_multiplier, adjusted_final_eosb_reward, unused_leave_cashout, legal_rule_applied  from eosb_records where 1=1"

	args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += " AND (code ILIKE $" + strconv.Itoa(argIndex) +
			     " OR name ILIKE $" + strconv.Itoa(argIndex+1) + ")"

		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}

	allowedSort := map[string]string{
		"ID":          "id",
		"Code":        "code",
		"Name":        "name",
		"LocalName":   "local_name",
		"Description":  "description",
		
		
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
	if ps <= 0 || ps > 5000 {
		ps = 10
	}

	offset := (p - 1) * ps

	query += " LIMIT $" + strconv.Itoa(argIndex) +
		" OFFSET $" + strconv.Itoa(argIndex+1)

	args = append(args, ps, offset)

		// ✅ Execute
	rows, err := core.DB.Query(query, args...)
	if err != nil {
		fmt.Print(err)
	}
	defer rows.Close()


	var eosb []models.EOSBRecord

	for rows.Next() {
		var l models.EOSBRecord
		var employeeID, servicePeriodID, contractID string
		err := rows.Scan(&l.ID, &employeeID, &servicePeriodID, &l.BaseEOSB, &l.TotalLines, &l.FinalPayable, &l.PaidEOSB, &l.Status, &l.PaymentDate, &l.FilePath, &l.CreatedAt, &l.UpdatedAt, &contractID, &l.ReasonModifierMultiplier, &l.AdjustedFinalEOSBReward, &l.UnusedLeaveCashout, &l.LegalRuleApplied)
		if err != nil {
			fmt.Print(err)
		}

		//Get employee & Service Period and contract
		Contract      := GetContractById(contractID)
		Employee      := GetEmployeeById(employeeID)
		ServicePeriod := GetServicePeriodById(servicePeriodID)

		l.Contract = &Contract
		l.Employee = &Employee
		l.ServicePeriod = &ServicePeriod


		eosb = append(eosb, l)
	}

	return eosb


}




