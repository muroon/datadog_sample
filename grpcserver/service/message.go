package service

import (
	"datadog_sample/grpcserver/model"
)

func GetMessages() ([]*model.Message, error) {
	list := make([]*model.Message, 0)
	rows, err := db.Query("SELECT * FROM message")
	defer rows.Close()
	if err != nil {
		return list, err
	}

	for rows.Next() {
		m := new(model.Message)
		if err := rows.Scan(&m.Id, &m.Text); err != nil {
			return list, err
		}
		list = append(list, m)
	}

	return list, nil
}

func PostMessage(text string) (*model.Message, error) {
	m := new(model.Message)

	stmt, err := db.Prepare("INSERT INTO message(text) VALUES(?)")
	defer stmt.Close()
	if err != nil {
		return m, err
	}

	res, err := stmt.Exec(text)
	if err != nil {
		return m, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return m, err
	}

	m.Id = uint32(id)
	m.Text = text

	return m, nil
}
