package repository

import (
	"github.com/mgierok/monujo/repository/entity"
)

func (r *Repository) StoreOperation(operation entity.Operation) (int64, error) {
	stmt, err := r.db.PrepareNamed(`
		INSERT INTO operations (portfolio_id, date, type, value, description, commision, tax)
		VALUES (:portfolio_id, :date, :type, :value, :description, :commision, :tax)
		RETURNING operation_id
	`)

	var operationId int64
	if nil == err {
		err = stmt.Get(&operationId, operation)
	}

	return operationId, err
}

func (r *Repository) PortfolioOperations(portfolio entity.Portfolio) (entity.Operations, error) {
	operations := entity.Operations{}
	err := r.db.Select(&operations,
		`SELECT
		operation_id,
		portfolio_id,
		date,
		type,
		value,
		COALESCE(description, '') AS description,
		commision,
		tax
	FROM operations
	WHERE portfolio_id = $1
	ORDER BY
		date ASC,
		operation_id ASC
	`,
		portfolio.PortfolioId)
	return operations, err
}

func (r *Repository) DeleteOperation(operation entity.Operation) error {
	_, err := r.db.Exec("DELETE FROM operations  WHERE portfolio_id = $1 AND operation_id = $2", operation.PortfolioId, operation.OperationId)
	return err
}
