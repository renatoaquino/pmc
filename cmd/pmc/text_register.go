package main

import "fmt"

func textRegister(reg register) func(rc record) error {
	return func(rc record) error {
		fmt.Printf("%s\t%s\t%s\t%s", rc.Label, rc.Type, rc.Host, rc.Endtime.Sub(rc.Starttime))
		if rc.Status == statusOk {
			fmt.Printf("\tOK")
		} else if rc.Status == statusFail {
			fmt.Printf("\tFAIL")
		} else if rc.Status == statusTimeout {
			fmt.Printf("\tFAIL TIMEOUT")
		}
		fmt.Printf("\n")
		return nil
	}
}
