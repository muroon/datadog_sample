package usecases

import (
	"context"
	"database/sql"
	"datadog_sample/httpserver/jsonmodel"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type message struct {
	Id   uint32
	Text string
}

func openDB() error {
	var err error
	db, err = sql.Open("mysql", "root:@/grpc_datadog")
	return err
}

func closeDB() error {
	if db == nil {
		return nil
	}
	return db.Close()
}

func DBList(w http.ResponseWriter, r *http.Request) {
	var b []byte

	ctx := r.Context()

	mess, err := getMessages(ctx)
	if err != nil {
		renderErrorJSON(w, err)
		return
	}

	list := make([]*jsonmodel.Message, 0, len(mess))
	for _, mess := range mess {
		list = append(list, &jsonmodel.Message{ID: mess.Id, Text: mess.Text})
	}

	b, err = json.Marshal(list)
	if err != nil {
		renderErrorJSON(w, err)
		return
	}

	renderJSON(w, b)
}

func DBPost(w http.ResponseWriter, r *http.Request) {
	var b []byte

	ctx := r.Context()

	text := r.URL.Query().Get("text")

	_, err := postMessage(ctx, text)
	if err != nil {
		renderErrorJSON(w, err)
		return
	}

	res := &jsonmodel.PostMessageResult{Status: true}
	b, err = json.Marshal(res)
	if err != nil {
		renderErrorJSON(w, err)
		return
	}

	renderJSON(w, b)
}

func getMessages(ctx context.Context) ([]*message, error) {
	list := make([]*message, 0)
	rows, err := db.QueryContext(ctx, "SELECT * FROM message")
	defer rows.Close()
	if err != nil {
		return list, err
	}

	for rows.Next() {
		m := new(message)
		if err := rows.Scan(&m.Id, &m.Text); err != nil {
			return list, err
		}
		list = append(list, m)
	}

	return list, nil
}

func postMessage(ctx context.Context, text string) (*message, error) {
	m := new(message)

	stmt, err := db.Prepare("INSERT INTO message(text) VALUES(?)")
	defer stmt.Close()
	if err != nil {
		return m, err
	}

	res, err := stmt.ExecContext(ctx, text)
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
