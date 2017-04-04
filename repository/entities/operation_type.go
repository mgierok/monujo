package entities

import "database/sql"

type OperationType struct {
	Type sql.NullString `db:"type"`
}

type OperationTypes []OperationType

type TransactionalOperationType struct {
	OperationType
}

type TransactionalOperationTypes []TransactionalOperationType

type FinancialOperationType struct {
	OperationType
}

type FinancialOperationTypes []FinancialOperationType
