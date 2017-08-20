package repository

import (
	"github.com/mgierok/monujo/repository/entity"
)

var sources = entity.Sources{
	{Name: "stooq"},
	{Name: "google"},
	{Name: "ingturbo"},
}

func Sources() entity.Sources {
	return sources
}
