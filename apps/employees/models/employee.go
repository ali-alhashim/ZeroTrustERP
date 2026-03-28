package models



type Employee struct {

	BadgeID	 string  `f:"text, unique, notnull"`
	Name	 string  `f:"text, notnull"`
	Department []Department `f:"many2one:employees_department"` // many employees can belong to one department, one department can have many employees
    LocalName string  `f:"text"` // local name of employee, can be in local language, optional
	JobTitle  JobTitle  `f:"many2one:employees_job_title"` // job title of employee, optional
	
}


type Department struct {

	Name	 string  `f:"text, unique, notnull"`
	Code	 string  `f:"text, unique, notnull"`
	Manager Employee `f:"one2one:employees"` // one department has one manager, one employee can be manager of one department
	
}


type JobTitle struct {

	Name	 string  `f:"text, unique, notnull"`
	Code	 string  `f:"text, unique, notnull"`
	Description string  `f:"text"` // description of job title, optional
	
}