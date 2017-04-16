package repository

import "github.com/mgierok/monujo/repository/entities"

func FinancialOperationTypes() (entities.FinancialOperationTypes, error) {
	types := entities.FinancialOperationTypes{}
	err := Db().Select(&types, "SELECT type FROM financial_operation_types ORDER BY type ASC")
	return types, err
}
