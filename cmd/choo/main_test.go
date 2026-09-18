package main

import (
	"testing"
	"time"
)

func TestParams(t *testing.T) {
	if got, _ := timeParam("9:05"); got != "0905" {
		t.Errorf("timeParam(9:05) = %q", got)
	}
	if _, err := timeParam("25:00"); err == nil {
		t.Error("timeParam(25:00) should fail")
	}
	yy := time.Now().Format("06")
	if got, _ := dateParam("3/12"); got != "0312"+yy {
		t.Errorf("dateParam(3/12) = %q", got)
	}
	if got, _ := dateParam(""); got != "" {
		t.Errorf("dateParam(\"\") = %q", got)
	}
	if _, err := dateParam("blah"); err == nil {
		t.Error("dateParam(blah) should fail")
	}
}
