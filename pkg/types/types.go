package types

import "github.com/envelope-zero/backend/pkg/models"

type ResourceCerate struct {
	PermanentError bool
	ErrorCount     int
	Resource       any
}

type Transaction struct {
	Create             models.TransactionCreate
	Budget             string
	SourceAccount      string
	DestinationAccount string
	Envelope           string
}
