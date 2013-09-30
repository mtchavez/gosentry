package sentry

import (
	"testing"
)

func TestParseDsn_BadURI(t *testing.T) {
	user, pass, proj, err := parseDsn("@cas.cc")
	if err == nil {
		t.Error("Bad DSN didn't raise error")
	}
	if user != "" {
		t.Errorf("Parsed user incorrectly got %+v", user)
	}

	if pass != "" {
		t.Errorf("Parsed password incorrectly got %+v", pass)
	}

	if proj != "" {
		t.Errorf("Parsed project incorrectly got %+v", proj)
	}
}

func TestParseDsn_NoPass(t *testing.T) {
	user, pass, proj, err := parseDsn("http://asdf123@cas.cc")
	if err != nil {
		t.Error("Incorrectly raised url parse error")
	}
	if user == "" {
		t.Errorf("Parsed user incorrectly got %+v", user)
	}

	if pass != "" {
		t.Errorf("Parsed password incorrectly got %+v", pass)
	}

	if proj != "" {
		t.Errorf("Parsed project incorrectly got %+v", proj)
	}
}

func TestParseDsn_Success(t *testing.T) {
	user, pass, proj, err := parseDsn("http://asdf123:my-pw@cas.cc/1234")
	if err != nil {
		t.Error("Incorrectly raised url parse error")
	}
	if user != "asdf123" {
		t.Errorf("Parsed user incorrectly got %+v", user)
	}

	if pass != "my-pw" {
		t.Errorf("Parsed password incorrectly got %+v", pass)
	}

	if proj != "1234" {
		t.Errorf("Parsed project incorrectly got %+v", proj)
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
