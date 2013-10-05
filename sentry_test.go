package sentry

import (
	"testing"
)

func TestParseDsn_BadURI(t *testing.T) {
	rc := &RavenConfig{}
	err := rc.parseDsn("@cas.cc")
	if err == nil {
		t.Error("Bad DSN didn't raise error")
	}
	if rc.User != "" {
		t.Errorf("Parsed user incorrectly got %+v", rc.User)
	}

	if rc.Pass != "" {
		t.Errorf("Parsed password incorrectly got %+v", rc.Pass)
	}

	if rc.Project != "" {
		t.Errorf("Parsed project incorrectly got %+v", rc.Project)
	}
}

func TestParseDsn_NoPass(t *testing.T) {
	rc := &RavenConfig{}
	err := rc.parseDsn("http://asdf123@cas.cc")
	if err != nil {
		t.Error("Incorrectly raised url parse error")
	}
	if rc.User == "" {
		t.Errorf("Parsed user incorrectly got %+v", rc.User)
	}

	if rc.Pass != "" {
		t.Errorf("Parsed password incorrectly got %+v", rc.Pass)
	}

	if rc.Project != "" {
		t.Errorf("Parsed project incorrectly got %+v", rc.Project)
	}
}

func TestParseDsn_Success(t *testing.T) {
	rc := &RavenConfig{}	
	err := rc.parseDsn("http://asdf123:my-pw@cas.cc/1234")
	if err != nil {
		t.Error("Incorrectly raised url parse error")
	}
	if rc.User != "asdf123" {
		t.Errorf("Parsed user incorrectly got %+v", rc.User)
	}

	if rc.Pass != "my-pw" {
		t.Errorf("Parsed password incorrectly got %+v", rc.Pass)
	}

	if rc.Project != "1234" {
		t.Errorf("Parsed project incorrectly got %+v", rc.Project)
	}
}

func TestSetup_BadDsn(t *testing.T) {
	rc, err := Setup("dsn")
	if err == nil {
		t.Error("Didn't raise error for bad DSN")
	}
	if rc != nil {
		t.Errorf("Incorrectly returned RavenConfig %+v", rc)
	}
}
