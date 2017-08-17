package main

import (
	"context"
	"net"
	"time"
)

func dnsVerify(ver verify) record {
	inf := record{Type: "DNS", Host: ver.Config, Label: ver.Label}
	var ctx context.Context
	var cancel func()
	if ver.Timeout.Duration > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), ver.Timeout.Duration)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	inf.Starttime = time.Now()
	_, err := net.DefaultResolver.LookupIPAddr(ctx, ver.Config)
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
