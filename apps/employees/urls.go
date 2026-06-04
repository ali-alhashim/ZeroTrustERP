package employees

import (
	"net/http"

	"zerotrusterp/apps/employees/controllers"
	"zerotrusterp/core"
)


func EmployeeListRoutes(mux *http.ServeMux) {

	mux.Handle("GET /employees/list", core.AuthMiddleware(http.HandlerFunc(controllers.ListEmployees), "employees:R"))
	mux.Handle("GET /employees/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateEmployee), "employees:W"))
	mux.Handle("POST /employees/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateEmployee), "employees:W"))
	mux.Handle("GET /employees/generate-badge-id", core.AuthMiddleware(http.HandlerFunc(controllers.GenerateBadgeIdApi), "employees:W"))
	mux.Handle("GET /api/employees/list", core.AuthMiddleware(http.HandlerFunc(controllers.GetEmployeesListApi), "employees:R"))
	mux.Handle("GET /employees/details/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.GetEmployeeDetails), "employees:R"))
	
	mux.Handle("POST /employees/update/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.UpdateEmployee), "employees:W"))

	mux.Handle("GET /employees/Jobs", core.AuthMiddleware(http.HandlerFunc(controllers.ListJobs), "jobs:R"))
	mux.Handle("GET /employees/jobs/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateJob), "jobs:R"))
	mux.Handle("POST /employees/jobs/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateJob), "jobs:W"))
	mux.Handle("GET /api/job-titles/list", core.AuthMiddleware(http.HandlerFunc(controllers.GetJobTitleListApi), "jobs:R"))
	mux.Handle("GET /job/details/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.JobsDetails), "jobs:R"))
	mux.Handle("POST /job/update/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.UpdateJob), "job:W"))
	mux.Handle("GET /jobs/download/csv", core.AuthMiddleware(http.HandlerFunc(controllers.DownloadJobsCSV), "job:R"))
	mux.Handle("POST /jobs/upload/csv", core.AuthMiddleware(http.HandlerFunc(controllers.ImportCSVJobs), "job:W"))

	mux.Handle("GET /employees/departments", core.AuthMiddleware(http.HandlerFunc(controllers.ListDepartments), "departments:R"))
	mux.Handle("GET /employees/departments/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateDepartment), "departments:R"))
	mux.Handle("POST /employees/departments/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateDepartment), "departments:W"))
    mux.Handle("GET /departments/details/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.DepartmentsDetails), "departments:R"))
	mux.Handle("POST /department/update/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.UpdateDepartment), "departments:W"))

	mux.Handle("GET /api/departments/list", core.AuthMiddleware(http.HandlerFunc(controllers.GetDepartmentListApi), "departments:R"))
	mux.Handle("GET /department/download/csv", core.AuthMiddleware(http.HandlerFunc(controllers.DownloadDepartmentsCSV), "department:R"))
	mux.Handle("POST /department/upload/csv", core.AuthMiddleware(http.HandlerFunc(controllers.ImportCSVDepartments), "departments:W"))
	
	mux.Handle("GET /employees/media/employees/images/{imageName}", core.AuthMiddleware(http.HandlerFunc(controllers.EmployeeImageGET), "employees:R"))
	mux.Handle("GET /employees/media/employees/{documentType}/{badgeId}/{documentName}", core.AuthMiddleware(http.HandlerFunc(controllers.GetEmployeeAttachment), "employees:R"))


	mux.Handle("POST /employees/upload/image/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.UploadEmployeeImage), "employees:W"))
	
	mux.Handle("GET /employees/download/csv", core.AuthMiddleware(http.HandlerFunc(controllers.DownloadEmployeesCSV), "employees:R"))
	mux.Handle("POST /employees/upload/csv", core.AuthMiddleware(http.HandlerFunc(controllers.ImportCSVEmployees), "employees:W"))


	mux.Handle("GET /employees/ShiftSchedule", core.AuthMiddleware(http.HandlerFunc(controllers.ListShift), "ShiftSchedule:R"))
	mux.Handle("GET /employees/shift/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateShift), "ShiftSchedule:W"))
	mux.Handle("POST /employees/shift/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateShift), "ShiftSchedule:W"))

	mux.Handle("GET /ShiftSchedule/details/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.UpdateSchedules), "ShiftSchedule:R"))
	mux.Handle("POST /ShiftSchedule/details/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.UpdateSchedules), "ShiftSchedule:W"))
	mux.Handle("GET /api/ShiftSchedule/list", core.AuthMiddleware(http.HandlerFunc(controllers.GetShiftListApi), "ShiftSchedule:R"))
	


	mux.Handle("GET /employees/contract", core.AuthMiddleware(http.HandlerFunc(controllers.ListContract), "Contract:R"))
	mux.Handle("GET /employees/Contract/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateContract), "Contract:R"))
	mux.Handle("POST /employees/Contract/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateContract), "Contract:W"))
	mux.Handle("GET /employees/contract/end/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.TerminateContract), "Contract:R"))
	mux.Handle("POST /employees/contract/end/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.TerminateContract), "Contract:W"))


	mux.Handle("GET /salary/componentTypes", core.AuthMiddleware(http.HandlerFunc(controllers.ListSalaryComponentTypes), "Contract:R"))

	mux.Handle("GET /salary/componentType/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateSalaryComponentTypes), "Contract:W"))
	mux.Handle("POST /salary/componentType/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateSalaryComponentTypes), "Contract:W"))

	mux.Handle("GET /componentType/details/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.DetailsSalaryComponentTypes), "Contract:W"))
	mux.Handle("POST /componentType/details/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.DetailsSalaryComponentTypes), "Contract:W"))

    mux.Handle("GET /api/componentType/list", core.AuthMiddleware(http.HandlerFunc(controllers.GetComponentTypeListApi), "jobs:R"))


	mux.Handle("GET /employees/contract/details/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.ContractDetails), "Contract:R"))
	mux.Handle("POST /employees/contract/details/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.ContractDetails), "Contract:R"))



	mux.Handle("GET /employees/clearance", core.AuthMiddleware(http.HandlerFunc(controllers.ListClearance), "Clearance:R"))
	mux.Handle("GET /employees/clearance/details/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.GetClearanceDetails), "Clearance:R"))
	mux.Handle("GET /employees/Clearance-Templates", core.AuthMiddleware(http.HandlerFunc(controllers.ListClearanceTemplates), "Clearance:R"))
	mux.Handle("POST /employees/Clearance-Templates/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateClearanceTemplate), "Clearance:W"))
	mux.Handle("GET /employees/Clearance-Templates/create", core.AuthMiddleware(http.HandlerFunc(controllers.CreateClearanceTemplate), "Clearance:W"))

	mux.Handle("GET /Clearance-Template/details/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.DetailsClearanceTemplate), "Clearance:R"))
	mux.Handle("POST /Clearance-Template/details/{id}", core.AuthMiddleware(http.HandlerFunc(controllers.DetailsClearanceTemplate), "Clearance:R"))

	mux.Handle("GET /api/ClearanceTemplates/list", core.AuthMiddleware(http.HandlerFunc(controllers.GetClearanceTemplateListApi), "Clearance:R"))



     mux.Handle("GET /employees/eosb", core.AuthMiddleware(http.HandlerFunc(controllers.ListEOSB), "EOSB:R"))


	 
	 mux.Handle("GET /api/departments/{departmentId}/employees", core.AuthMiddleware(http.HandlerFunc(controllers.GetDepartmentEmployeesApiByDepartmentId), "employees:R"))


	

	
	
	

	
}