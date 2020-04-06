package service

import (
	"context"

	"github.com/muroon/datadog_sample/grpcserver/model"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func GetMessages(ctx context.Context) ([]*model.Message, error) {
	list := make([]*model.Message, 0)
	var err error

	span, ctx := tracer.StartSpanFromContext(ctx,
		"db",
		tracer.SpanType(ext.SpanTypeSQL),
		tracer.ServiceName(dbServiceName),
		tracer.ResourceName("get-messages"),
	)
	defer span.Finish(tracer.WithError(err))

	rows, err := db.QueryContext(ctx, "SELECT * FROM message")
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
	var err error

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
