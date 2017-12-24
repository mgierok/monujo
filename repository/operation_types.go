package repository

import (
	"github.com/mgierok/monujo/repository/entity"
)

func (r *Repository) FinancialOperationTypes() (entity.FinancialOperationTypes, error) {
	types := entity.FinancialOperationTypes{}
	err := r.db.Select(&types, "SELECT type FROM financial_operation_types ORDER BY type ASC")
	return types, err
}
