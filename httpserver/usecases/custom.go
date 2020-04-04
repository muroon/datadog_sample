package usecases

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"datadog_sample/httpserver/jsonmodel"
)

func Custom(w http.ResponseWriter, r *http.Request) {
	var b []byte

	ctx := r.Context()

	// parent Span

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

	// child logic
	time.Sleep(time.Millisecond * 300)

	return ctx
}
