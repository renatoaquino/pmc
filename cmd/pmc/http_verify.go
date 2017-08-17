package main

import (
	"context"
	"net/http"
	"time"
)

func httpVerify(ver verify) record {
	inf := record{Type: "HTTP", Host: ver.Config, Label: ver.Label}
	req, _ := http.NewRequest("HEAD", ver.Config, nil)
	var ctx = req.Context()
	if ver.Timeout.Duration > 0 {
		ctx, cancel := context.WithTimeout(req.Context(), ver.Timeout.Duration)
		req = req.WithContext(ctx)
		defer cancel()
	}

	inf.Starttime = time.Now()
	_, err := http.DefaultTransport.RoundTrip(req)
	inf.Endtime = time.Now()

	select {
	case <-ctx.Done():
	default:
		break
	}
	if e := ctx.Err(); e != context.Canceled && e != nil {
		inf.Status = statusTimeout
	} else {
		if err != nil {
			inf.Status = statusFail
		} else {
			inf.Status = statusOk
		}
	}

	return inf
}
