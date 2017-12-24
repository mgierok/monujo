package repository

import (
	"github.com/mgierok/monujo/repository/entity"
)

var sources = entity.Sources{
	{Name: "stooq"},
	{Name: "google"},
	{Name: "ingturbo"},
	{Name: "alphavantage"},
	{Name: "bankier"},
}

func (r *Repository) Sources() entity.Sources {
	return sources
}
