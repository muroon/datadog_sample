package usecases

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/muroon/datadog_sample/httpserver/jsonmodel"
)

func Custom(w http.ResponseWriter, r *http.Request) {
	var b []byte

	ctx := r.Context()

	// parent Span
	span, ctx := tracer.StartSpanFromContext(ctx,
		"custom",                             // operation Name
		tracer.ServiceName("custom-service"), // Service Name(必要であれば入れる)
		tracer.ResourceName("parent-logic"),  // Resource Name(必要であれば入れる)
	)
	defer span.Finish()

	// parent logic
	time.Sleep(time.Millisecond * 200)

	ctx = customChild(ctx)

	res := &jsonmodel.CustomResult{Status: true}
	b, err := json.Marshal(res)
	if err != nil {
		renderErrorJSON(w, err)
		return
	}

	renderJSON(w, b)
}

func customChild(ctx context.Context) context.Context {
	// child Span
	span, ctx := tracer.StartSpanFromContext(ctx,
		"custom",                             // operation Name
		tracer.ServiceName("custom-service"), // Service Name
		tracer.ResourceName("child-logic"),   // Resource Name
	)
	defer span.Finish()

	// child logic
	time.Sleep(time.Millisecond * 300)

	return ctx
}
