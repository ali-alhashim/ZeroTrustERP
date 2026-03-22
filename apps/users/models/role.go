package models

import "time"

type Role struct {
	ID          int          
	Name        string       
	Description string       
	Permissions []Permission 
	CreatedAt   time.Time    
	UpdatedAt   time.Time    
}

//Permissions read, write, delete, update

type Permission struct {
	ID     int    
	Name   string 
	Action string 
}
