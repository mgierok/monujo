package repository

import "github.com/mgierok/monujo/repository/entity"

func FinancialOperationTypes() (entity.FinancialOperationTypes, error) {
	types := entity.FinancialOperationTypes{}
	err := Db().Select(&types, "SELECT type FROM financial_operation_types ORDER BY type ASC")
	return types, err
}
