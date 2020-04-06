package usecases

import (
	"context"
	"database/sql"
	"github.com/muroon/datadog_sample/httpserver/jsonmodel"
	"encoding/json"
	"net/http"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/go-sql-driver/mysql"

	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type message struct {
	Id   uint32
	Text string
}

func initDB() {
	// Datadog
	sqltrace.Register("mysql",
		&mysql.MySQLDriver{},
		sqltrace.WithServiceName("db-service"),
	)
}

func openDB() error {
	var err error
	db, err = sqltrace.Open("mysql", "root:@/grpc_datadog")
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
	var err error

	span, ctx := tracer.StartSpanFromContext(ctx,
		"db",
		tracer.SpanType(ext.SpanTypeSQL),
		tracer.ServiceName("db-service"),
		tracer.ResourceName("get-messages"),
	)
	defer span.Finish(tracer.WithError(err))

	var rows *sql.Rows
	rows, err = db.QueryContext(ctx, "SELECT * FROM message")
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

	var err error

	span, ctx := tracer.StartSpanFromContext(ctx,
		"db",
		tracer.SpanType(ext.SpanTypeSQL),
		tracer.ServiceName("db-service"),
		tracer.ResourceName("post-message"),
	)
	defer span.Finish(tracer.WithError(err))

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
