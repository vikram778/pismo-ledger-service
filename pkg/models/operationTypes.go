package models

type Operations struct {
	OperationTypeID string `json:"operation_type_id" db:"operation_type_id"` // OperationTypeID is the operation Id which are constant in system
	Description     string `json:"description" db:"description"`             // Description associated with each operation type
}
