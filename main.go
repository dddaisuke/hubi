package main

import (
	"bytes"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io"
	"os"
	"os/user"
	"reflect"
	"strconv"
)

var timestampType = reflect.TypeOf(github.Timestamp{})

type tokenSource struct {
	token *oauth2.Token
}

// Token implements the oauth2.TokenSource interface
func (t *tokenSource) Token() (*oauth2.Token, error) {
	return t.token, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("%s [number of issue] [target repository]", os.Args[0])
		return
	}

	usr, _ := user.Current()
	conf, err := Parse(usr.HomeDir + "/.github/config.json")
	if err != nil {
		fmt.Println("~/.github/config.json が見つかりません。https://github.com/dddaisuke/hubi を参照して下さい。")
		os.Exit(1)
	}

	fmt.Printf("[move] manabo-inc/sandbox/issues/%s -> manabo-inc/%s/issues\n", os.Args[1], os.Args[2])
	var fromIssueNumber int
	fromIssueNumber, _ = strconv.Atoi(os.Args[1])
	var toRepositoryName = os.Args[2]

	ts := &tokenSource{
		&oauth2.Token{AccessToken: conf.AccessToken},
	}

	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	var from_issue *github.Issue

	from_issue, _, _ = client.Issues.Get("manabo-inc", "sandbox", fromIssueNumber)

	url := Stringify(from_issue.HTMLURL)
	body := Stringify(from_issue.Body)
	fmt.Printf("\n-------------------- [sadbox/issues/%s] --------------------\n", Stringify(from_issue.Number))
	fmt.Printf("  URL %s\n", url)
	fmt.Printf("  Status %s\n", Stringify(from_issue.State))
	fmt.Printf("  Title %s\n", Stringify(from_issue.Title))
	fmt.Println("-------------------------------------------------------------\n")

	issue := &github.IssueRequest{
		Title: from_issue.Title,
		Body:  github.String("Ref: " + url + "\n\n" + body),
	}

	client.Issues.Create("manabo-inc", toRepositoryName, issue)

	if len(os.Args) > 3 && os.Args[3] == "-c" {
		closeIssue := &github.IssueRequest{
			State: github.String("closed"),
		}

		client.Issues.Edit("manabo-inc", "sandbox", fromIssueNumber, closeIssue)
		fmt.Printf("[closed] manabo-inc/sandbox/issues/%d\n", fromIssueNumber)
	}
}

func Stringify(message interface{}) string {
	var buf bytes.Buffer
	v := reflect.ValueOf(message)
	stringifyValue(&buf, v)
	return buf.String()
}

func stringifyValue(w io.Writer, val reflect.Value) {
	if val.Kind() == reflect.Ptr && val.IsNil() {
		w.Write([]byte("<nil>"))
		return
	}

	v := reflect.Indirect(val)

	switch v.Kind() {
	case reflect.String:
		fmt.Fprintf(w, `%s`, v)
	case reflect.Slice:
		w.Write([]byte{'['})
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				w.Write([]byte{' '})
			}

			stringifyValue(w, v.Index(i))
		}

		w.Write([]byte{']'})
		return
	case reflect.Struct:
		if v.Type().Name() != "" {
			w.Write([]byte(v.Type().String()))
		}

		// special handling of Timestamp values
		if v.Type() == timestampType {
			fmt.Fprintf(w, "{%s}", v.Interface())
			return
		}

		w.Write([]byte{'{'})

		var sep bool
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				continue
			}
			if fv.Kind() == reflect.Slice && fv.IsNil() {
				continue
			}

			if sep {
				w.Write([]byte(", "))
			} else {
				sep = true
			}

			w.Write([]byte(v.Type().Field(i).Name))
			w.Write([]byte{':'})
			stringifyValue(w, fv)
		}

		w.Write([]byte{'}'})
	default:
		if v.CanInterface() {
			fmt.Fprint(w, v.Interface())
		}
	}
}
