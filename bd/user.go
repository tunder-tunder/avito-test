package bd

import (
	"avito/models"
	// TODO: import models correctly; solve issue with docker config
	"database/sql"
)

func (db Database) GetAllUsers() (*models.UserList, error) {
	list := &models.UserList{}
	rows, err := db.Conn.Query("SELECT * FROM user_table ORDER BY ID DESC")

	if err != nil {
		return list, err
	}

	for rows.Next() {
		var item models.User

		err := rows.Scan(&item.ID, &item.FirstName, &item.LastName, &item.BalanceId)
		if err != nil {
			return list, err
		}

		list.Users = append(list.Users, item)
	}

	return list, nil
}
func (db Database) AddUser(user *models.User) (int, error) {
	var id int

	//var createdAt time.Time = time.Now()

	query := `INSERT INTO user_table (first_name, last_name) VALUES ($1, $2) RETURNING id`

	err := db.Conn.QueryRow(query, user.FirstName, user.LastName).Scan(&id)
	if err != nil {
		return 0, err
	}
	user.ID = id

	return user.ID, nil
}
func (db Database) GetUserById(userId int) (models.User, error) {
	item := models.User{}

	query := `SELECT * FROM user_table WHERE id = $1;`
	row := db.Conn.QueryRow(query, userId)

	switch err := row.Scan(&item.ID, &item.FirstName, &item.LastName, &item.BalanceId); err {
	case sql.ErrNoRows:
		return item, ErrNoMatch
	default:
		return item, err
	}
}
func (db Database) DeleteUser(userId int) error {
	query := `DELETE FROM user_table WHERE id = $1;`

	_, err := db.Conn.Exec(query, userId)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}
func (db Database) UpdateUser(userId int, userData models.User) (models.User, error) {
	item := models.User{}
	query := `UPDATE user_table SET first_name=$1, last_name=$2 WHERE id=$3 RETURNING id, last_name, first_name;`
	err := db.Conn.QueryRow(query, userData.FirstName, userData.LastName, userId).Scan(&item.ID, &item.FirstName, &item.LastName)
	if err != nil {
		if err == sql.ErrNoRows {
			return item, ErrNoMatch
		}
		return item, err
	}
	return item, nil
}
