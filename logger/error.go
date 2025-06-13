package logger


type TermoTuneError struct {
	err string 
	Details [] error
}

func (e *TermoTuneError) TermoTuneError() string {
	errors := e.err 
	for _, detail := range e.Details {
		errors += "\n" + detail.Error()
	}
	return errors
} 

func (e *TermoTuneError) Error(error string , detail ... error) *TermoTuneError {
	return &TermoTuneError{
		err:   error,
		Details: detail,
	}
}



