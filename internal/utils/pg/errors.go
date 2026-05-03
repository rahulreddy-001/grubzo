package pg

type Error struct {
	Severity            string
	SeverityUnlocalized string
	Code                string
	Message             string
	Detail              string
	Hint                string
	Position            int32
	InternalPosition    int32
	InternalQuery       string
	Where               string
	SchemaName          string
	TableName           string
	ColumnName          string
	DataTypeName        string
	ConstraintName      string
	File                string
	Line                int32
	Routine             string
}

func (pe *Error) Error() string {
	return pe.Severity + ": " + pe.Message + " (SQLSTATE " + pe.Code + ")"
}
