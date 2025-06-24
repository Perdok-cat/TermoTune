package shared

import (

)


type Task struct {
	Type int // download , search ... 
	Error string
}


type IntOrString struct {
	IntVal int 
	StrVal string 
	IsInt bool
}
