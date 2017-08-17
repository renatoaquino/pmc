package main

import (
	"log"
	"net/url"

	"github.com/influxdata/influxdb/client"
)

func influxdbRegister(reg register) func(rc record) error {
	host, err := url.Parse(reg.Config)
	if err != nil {
		log.Fatal(err)
	}

	username := host.User.Username()
	password, set := host.User.Password()

	conf := client.Config{
		URL:      *host,
		Username: username,
	}
	if set {
		conf.Password = password
	}

	con, err := client.NewClient(conf)
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = con.Ping()
	if err != nil {
		log.Fatal(err)
	}

	params := host.Query()
	if params["db"][0] == "" {
		log.Fatal("Influxdb register requires a ?db=name definition")
	}

	if params["series"][0] == "" {
		log.Fatal("Influxdb register requires a ?series=name definition")
	}

	return func(rc record) error {
		var status string
		if rc.Status == statusOk {
			status = "OK"
		} else if rc.Status == statusFail {
			status = "FAIL"
		} else if rc.Status == statusTimeout {
			status = "TIMEOUT"
		}
		var pts = make([]client.Point, 1)
		pts[0] = client.Point{
			Measurement: params["series"][0],
			Tags:        map[string]string{"type": rc.Type},
			Fields: map[string]interface{}{
				"label":      rc.Label,
				"host":       rc.Host,
				"latency_ns": rc.Endtime.Sub(rc.Starttime).Nanoseconds(),
				"status":     status,
			},
			Time:      rc.Starttime,
			Precision: "n",
		}
		bps := client.BatchPoints{
			Points:          pts,
			Database:        params["db"][0],
			RetentionPolicy: "default",
			Time:            rc.Starttime,
		}
		_, err = con.Write(bps)
		return err
	}
}
