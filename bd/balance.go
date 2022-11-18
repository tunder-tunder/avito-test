package bd

import (
	"avito/models"
	"github.com/google/uuid"
	"time"
	// TODO: import models correctly; solve issue with docker config
	"database/sql"
)

const (
	stSuccess string = "SUCCESS"
	stFail    string = "FAILED"
)

func (db Database) GetAllBalances() (*models.BalanceList, error) {
	list := &models.BalanceList{}
	rows, err := db.Conn.Query("SELECT * FROM balance_table ORDER BY ID DESC")

	if err != nil {
		return list, err
	}

	for rows.Next() {
		var item models.Balance

		err := rows.Scan(&item.ID, &item.UserId, &item.Total, &item.Reserve, &item.OrderNumber,
			&item.ServiceId, &item.CreatedAt, &item.Status)
		if err != nil {
			return list, err
		}

		list.Balances = append(list.Balances, item)
	}

	return list, nil
}

func (db Database) InitBalance(userId int) error {
	var id int
	var createdAt time.Time = time.Now()
	// order number is basically uuid sorry
	orderNumber := uuid.New()

	query := `INSERT INTO balance_table (order_number, created_at, total, user_id) VALUES ($1, $2, $3, $4) RETURNING id`

	err := db.Conn.QueryRow(query, orderNumber, createdAt, 0, userId).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (db Database) AddBalance(userId int, topUp int) error {
	var id int
	var createdAt time.Time = time.Now()
	item := models.Balance{}
	// order number is basically uuid sorry
	orderNumber := uuid.New()

	query := `INSERT INTO balance_table (user_id, order_number, created_at, total, status) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := db.Conn.QueryRow(query, userId, orderNumber, createdAt, topUp, stSuccess).Scan(&id)
	if err != nil {
		return err
	}
	item.ID = id

	return nil
}

func (db Database) ReserveBalance(userId int, serviceId int, orderNumber string, price int, total int) (*models.Balance, error) {
	var id int
	var createdAt time.Time = time.Now()
	item := &models.Balance{}
	//orderNumber := uuid.New() goes into handler here to work with ger order by id and pay balance methods

	query := `INSERT INTO balance_table (user_id, order_number, created_at, total, status, reserve, service_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := db.Conn.QueryRow(query, userId, orderNumber, createdAt, total, stSuccess, price, serviceId).Scan(&id)
	if err != nil {
		return nil, err
	}
	item.ID = id

	return item, nil
}

// метод признания выручки
func (db Database) PayBalance(userId int, serviceId int, orderNumber string, total int) (*models.Balance, error) {
	var id int
	var createdAt time.Time = time.Now()
	item := &models.Balance{}

	query := `INSERT INTO balance_table (user_id, order_number, created_at, total, status, reserve, service_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := db.Conn.QueryRow(query, userId, orderNumber, createdAt, total, stSuccess, 0, serviceId).Scan(&id)
	if err != nil {
		return item, err
	}
	item.ID = id

	return item, nil
}

func (db Database) GetBalanceById(userId int) (models.Balance, error) {
	item := models.Balance{}
	//
	var total int
	query := `SELECT * FROM balance_table WHERE user_id = $1 ORDER BY ID DESC LIMIT 1;`
	row := db.Conn.QueryRow(query, userId)

	switch err := row.Scan(&item.ID, &item.UserId, &item.Total, &item.Reserve, &item.OrderNumber,
		&item.ServiceId, &item.CreatedAt, &item.Status); err {
	case sql.ErrNoRows:
		item.Total = total
		return item, ErrNoMatch
	default:
		return item, err
	}

}

func (db Database) GetOrderByNumber(orderNumber string) (models.Balance, error) {
	item := models.Balance{}
	//
	query := `SELECT * FROM balance_table WHERE order_number = $1 ORDER BY ID DESC LIMIT 1;`
	row := db.Conn.QueryRow(query, orderNumber)

	switch err := row.Scan(&item.ID, &item.UserId, &item.Total, &item.Reserve, &item.OrderNumber,
		&item.ServiceId, &item.CreatedAt, &item.Status); err {
	case sql.ErrNoRows:
		return item, ErrNoMatch
	default:
		return item, err
	}
}

func (db Database) DeleteBalance(balanceId int) error {
	query := `DELETE FROM balance_table WHERE id = $1;`

	_, err := db.Conn.Exec(query, balanceId)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}
