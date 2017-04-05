package entities

type OperationType struct {
	Type string `db:"type"`
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
