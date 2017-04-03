package repository

import "github.com/mgierok/monujo/repository/entities"

func Currencies() (entities.Currencies, error) {
	currencies := entities.Currencies{}
	err := Db().Select(&currencies, "SELECT currency FROM currencies ORDER BY currency ASC")
	return currencies, err
}
