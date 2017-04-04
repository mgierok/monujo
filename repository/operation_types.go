package repository

import "github.com/mgierok/monujo/repository/entities"

func TransactionalOperationTypes() (entities.TransactionalOperationTypes, error) {
	types := entities.TransactionalOperationTypes{}
	err := Db().Select(&types, "SELECT type FROM transactional_operation_types ORDER BY type ASC")
	return types, err
}

func FinancialOperationTypes() (entities.FinancialOperationTypes, error) {
	types := entities.FinancialOperationTypes{}
	err := Db().Select(&types, "SELECT type FROM financial_operation_types ORDER BY type ASC")
	return types, err
}
