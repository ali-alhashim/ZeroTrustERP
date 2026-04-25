package controllers

import (
	"net/http"
	"zerotrusterp/core"
)

func ListJobs(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"Title": "Jobs",
		
	}

	core.RenderPage(w,r, "apps/employees/views/jobs-list.html", data)
}

func CreateJob(w http.ResponseWriter, r *http.Request){


	

	if r.Method == http.MethodGet {
		data := map[string]interface{}{
		"Title": "Jobs",
		
	}
		core.RenderPage(w,r, "apps/employees/views/jobs-create.html", data)
	}

	if r.Method == http.MethodPost{
		//create job
	}



}



