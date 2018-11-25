package schemacheck

type FieldError struct {
	Field string
	Msg   string
}

type Error struct {
	Msg  string
	Errs []*FieldError
}

func (e *Error) Error() string {
	s := e.Msg + "; "
	for _, e := range e.Errs {
		s += e.Field + ": " + e.Msg + ";"
	}
	return s
}
