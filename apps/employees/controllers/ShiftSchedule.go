package controllers

import (
	"fmt"
	"net/http"
	"zerotrusterp/apps/employees/models"
	"zerotrusterp/core"
	"strconv"
	
)



func ListShift(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	search     := query.Get("q")
	sortBy     := query.Get("sort")
	order      := query.Get("order")
	page       := query.Get("page")
	pageSize   := query.Get("pageSize")

	fmt.Printf("list logs: search=%s, sort=%s, order=%s, page=%s, pageSize=%s\n", search, sortBy, order, page, pageSize)
    
	
    totalRecords := core.GetCountRecords("shift_schedules")
    

	shiftSchedules:= GetShiftSchedulesFromDB(search, sortBy, order, page, pageSize)


	data := map[string]interface{}{
		"Query": search,
		"Sort":  sortBy,
		"Order": order,
		"Page":  page,
		"PageSize": pageSize,
		"TotalRecords":totalRecords,
		"ShiftSchedules": shiftSchedules,
	}

	core.RenderPage(w,r, "apps/employees/views/ShiftSchedule-list.html", data)
}



func GetShiftSchedulesFromDB(search, sort, order, page, pageSize string)[]models.ShiftSchedule{

	var shiftSchedules []models.ShiftSchedule

	query :="select id, name, from_time,to_time,from_date, to_date, monday,tuesday, wednesday,thursday,friday,saturday, sunday from shift_schedules  WHERE 1=1"
    args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += " AND (name ILIKE $" + strconv.Itoa(argIndex) +
                 " OR from_time ILIKE $"+strconv.Itoa(argIndex) +
			     " OR to_date ILIKE $" + strconv.Itoa(argIndex+1) + ")"

		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}

	allowedSort := map[string]string{
		"ID":         "id",
		"Name":       "name",
		"from_time":  "from_time",
		"to_time":     "to_time",
		"from_date":   "from_date",
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
		panic(err)
	}
	defer rows.Close()


	for rows.Next() {
    var shift models.ShiftSchedule
   

    // 2. Scan into these NullString variables
   err := rows.Scan(
    &shift.ID,
	&shift.Name,
	&shift.FromTime,
	&shift.ToTime,
	&shift.FromDate,
	&shift.ToDate, 
	&shift.Monday,
	&shift.Tuesday,
	&shift.Wednesday,
	&shift.Thursday,
	&shift.Friday,
	&shift.Saturday,
	&shift.Sunday,

    
)
    if err != nil {
        fmt.Printf("Scan Error: %v\n", err)
        continue
    }

  

    shiftSchedules = append(shiftSchedules, shift)
}

	return shiftSchedules

}



func CreateShift(w http.ResponseWriter, r *http.Request){


	if r.Method == http.MethodGet{

		data:=map[string]interface{}{
			"Title":"Create Shift",
		}


		core.RenderPage(w,r, "apps/employees/views/ShiftSchedule-create.html", data)

	}

	if r.Method == http.MethodPost{

	}

}