package sentry

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

type RavenConfig struct {
	Dsn     string
	User    string
	Pass    string
	Project string
	url     *url.URL
	client  *http.Client
}

type sentryEvent struct {
	Message      string           `json:"message"`
	Project      string           `json:"project"`
	Date         time.Time        `json:"timestamp"`
	EventId      string           `json:"event_id"`
	ServerName   string           `json:"server_name"`
	Level        string           `json:"level"`
	Stacktrace   *trace           `json:"sentry.interfaces.Stacktrace"`
	Platform     string           `json:"platform"`
	ContextLines []string         `json:"context_lines"`
	Exception    *sentryException `json:"sentry.interfaces.Exception"`
}

type sentryException struct {
	Value  string `json:"value"`
	Type   string `json:"type"`
	Module string `json:"module"`
}

type trace struct {
	Frames []*frame `json:"frames"`
}

type frame struct {
	AbsPath     string `json:"abs_path"`
	Filename    string `json:"filename"`
	Module      string `json:"module"`
	Function    string `json:"function"`
	Lineno      string `json:"lineno"`
	ContextLine string `json:"context_line"`
}

var (
	CLIENT         = "raven-gosentry"
	isoFormat      = "2013-01-01T12:12:12"
	LVL_ERROR      = "error"
	SENTRY_VERSION = "2.0"
	VERSION        = "0.0.1"
	PLATFORM       = "go"
)

func Setup(dsn string) (*RavenConfig, error) {
	rc := &RavenConfig{Dsn: dsn}
	err := rc.parseDsn(dsn)
	if err != nil {
		return nil, err
	}
	return rc, nil
}

func (rc *RavenConfig) parseDsn(dsn string) (err error) {
	u, e := url.Parse(dsn)
	if e != nil {
		return e
	}
	info := u.User
	if info == nil {
		e = errors.New("No user info parsed from dsn")
		return e
	}
	rc.User = info.Username()
	rc.Pass, _ = info.Password()
	rc.Project = strings.Replace(u.Path, "/", "", 1)
	rc.url = u
	rc.client = &http.Client{}
	return
}

func (rc *RavenConfig) apiPath() string {
	return fmt.Sprintf("%s://%s/api/%s/store/", rc.url.Scheme, rc.url.Host, rc.Project)
}

func (rc *RavenConfig) Message(panicMsg interface{}, msg string, ifaces ...interface{}) {
	stackString := string(debug.Stack())
	stack := strings.Split(stackString, "\n")
	frames := make([]*frame, 0)
	for _, line := range stack {
		// split := strings.Split(line, ":")
		// lineno := split[len(split)-1]
		frames = append(frames, &frame{AbsPath: line})
	}
	except := &sentryException{Type: fmt.Sprint(panicMsg), Value: stackString}
	date := time.Now().UTC()
	date.Format(isoFormat)
	hostname, _ := os.Hostname()
	uid, _ := uuid.NewV4()
	body := &sentryEvent{
		Message:      msg,
		Project:      rc.Project,
		Date:         date,
		EventId:      fmt.Sprint(uid),
		ServerName:   hostname,
		Level:        LVL_ERROR,
		Stacktrace:   &trace{frames},
		Platform:     PLATFORM,
		ContextLines: stack,
		Exception:    except,
	}
	fmt.Printf("In Raven: %+v\n", body)
	js, _ := json.Marshal(body)
	fmt.Printf("JSON BODY: %+v\n", string(js))
	if err := rc.sendMessage(body); err != nil {
		fmt.Println("ERR: ", err)
	}
}

func (rc *RavenConfig) sendMessage(body *sentryEvent) error {
	encodedBody, err := encodeBody(body)
	if err != nil {
		return err
	}
	// fmt.Println("Encoded body: ", string(encodedBody.Bytes()))
	// fmt.Println("API PATH: ", rc.apiPath())
	req, _ := http.NewRequest("POST", rc.apiPath(), encodedBody)
	timestamp := time.Now().UTC()
	client := fmt.Sprintf("%v/%v", CLIENT, VERSION)
	sentryAuth := fmt.Sprintf("Sentry sentry_version=%v", SENTRY_VERSION)
	sentryAuth += fmt.Sprintf(", sentry_client=%v", client)
	sentryAuth += fmt.Sprintf(", sentry_timestamp=%v", timestamp.Format(isoFormat))
	sentryAuth += fmt.Sprintf(", sentry_key=%v", rc.User)
	sentryAuth += fmt.Sprintf(", sentry_secret=%v", rc.Pass)

	req.Header.Add("X-Sentry-Auth", sentryAuth)
	req.Header.Add("User-Agent", client)
	req.Header.Add("Content-Type", "application/octet-stream")
	// req.Header.Add("Connection", "close")
	// req.Header.Add("Accept-Encoding", "identity")
	fmt.Println("HEADERS: ", req.Header)
	resp, err := rc.client.Do(req)
	bodyBuf := bytes.NewBuffer(make([]byte, 64000))
	resp.Body.Read(bodyBuf.Bytes())
	fmt.Printf("%+v %+v\n", resp.Status, resp.StatusCode)
	fmt.Println("RESPONSE: ", string(bodyBuf.Bytes()))
	if err != nil {
		fmt.Println("Do() ERR: ", err)
		return err
	}
	resp.Body.Close()
	return nil
}

func encodeBody(body *sentryEvent) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	base64Enc := base64.NewEncoder(base64.StdEncoding, buf)
	zlibWriter := zlib.NewWriter(base64Enc)
	jsonEnc := json.NewEncoder(zlibWriter)
	if err := jsonEnc.Encode(body); err != nil {
		return nil, err
	}
	if err := zlibWriter.Close(); err != nil {
		return nil, err
	}
	if err := base64Enc.Close(); err != nil {
		return nil, err
	}
	return buf, nil
}
