package repository

import (
	"github.com/mgierok/monujo/db"
	"github.com/mgierok/monujo/repository/entity"
)

func Currencies() (entity.Currencies, error) {
	currencies := entity.Currencies{}
	err := db.Connection().Select(&currencies, "SELECT currency FROM currencies ORDER BY currency ASC")
	return currencies, err
}
