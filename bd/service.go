package bd

import (
	"gitlab.com/tunder-tunder/avito/models"
	// TODO: import models correctly; solve issue with docker config
	"database/sql"
)

func (db Database) GetAllServices() (*models.ServiceList, error) {
	list := &models.ServiceList{}
	rows, err := db.Conn.Query("SELECT * FROM service_table ORDER BY ID DESC")

	if err != nil {
		return list, err
	}

	for rows.Next() {
		var item models.Service

		err := rows.Scan(&item.ID, &item.ServiceName, &item.Price, &item.Availability)
		if err != nil {
			return list, err
		}

		list.Services = append(list.Services, item)
	}

	return list, nil
}

func (db Database) AddService(service *models.Service) error {
	var id int

	query := `INSERT INTO service_table (service_name, price, availability) VALUES ($1, $2, $3) RETURNING id`

	err := db.Conn.QueryRow(query, service.ServiceName, service.Price, service.Availability).Scan(&id)
	if err != nil {
		return err
	}
	service.ID = id

	return nil
}

func (db Database) GetServiceById(serviceId int) (models.Service, error) {
	item := models.Service{}

	query := `SELECT * FROM service_table WHERE id = $1;`
	row := db.Conn.QueryRow(query, serviceId)

	switch err := row.Scan(&item.ID, &item.ServiceName, &item.Price, &item.Availability); err {
	case sql.ErrNoRows:
		return item, ErrNoMatch
	default:
		return item, err
	}
}

func (db Database) DeleteService(serviceID int) error {
	query := `DELETE FROM service_table WHERE id = $1;`

	_, err := db.Conn.Exec(query, serviceID)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}
