package repository

import (
	"github.com/mgierok/monujo/repository/entity"
)

func (r *Repository) Currencies() (entity.Currencies, error) {
	currencies := entity.Currencies{}
	err := r.db.Select(&currencies, "SELECT currency FROM currencies ORDER BY currency ASC")
	return currencies, err
}
