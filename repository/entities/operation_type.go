package entities

type OperationType struct {
	Type string `db:"type"`
}

type OperationTypes []OperationType

type FinancialOperationType struct {
	OperationType
}

type FinancialOperationTypes []FinancialOperationType
