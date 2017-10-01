package repository

import (
	"github.com/mgierok/monujo/db"
	"github.com/mgierok/monujo/repository/entity"
)

func StoreOperation(operation entity.Operation) (int64, error) {
	stmt, err := db.Connection().PrepareNamed(`
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

func PortfolioOperations(portfolio entity.Portfolio) (entity.Operations, error) {
	operations := entity.Operations{}
	err := db.Connection().Select(&operations,
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

func DeleteOperation(operation entity.Operation) error {
	_, err := db.Connection().Exec("DELETE FROM operations  WHERE portfolio_id = $1 AND operation_id = $2", operation.PortfolioId, operation.OperationId)
	return err
}
