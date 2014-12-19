package be

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
)

var (
	TestVersion       = "0.T.0"
	TestPort          = 44556677
	TestSessionCookie = "test_session"
	TestSessionSecret = "NotVerySecret"
)

func AssertGetString(t *testing.T, url string) string {
	resp, err := connectToTestAPI("GET", url)
	if err != nil {
		t.Fatalf("AssertGet Failed: %s: %s", url, err.Error())
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("AssertGet Failed reading body: %s: %s", url, err.Error())
		return ""
	}
	if resp.StatusCode != 200 {
		t.Fatalf("AssertGet Received non-200 status: %d: %s", resp.StatusCode, url)
		return string(body)
	}
	return string(body)
}

func AssertStatus(t *testing.T, status int, method string, url string) {
	resp, err := connectToTestAPI(method, url)
	if err != nil {
		t.Fatalf("AssertStatus Failed: %s: %s", url, err.Error())
		return
	}
	if resp.StatusCode != status {
		t.Fatalf("AssertStatus for %d failed with status code %d: %s", status, resp.StatusCode, url)
		return
	}
}

func connectToTestAPI(method string, url string) (resp *http.Response, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", AcceptHeaderPrefix+TestVersion)
	return client.Do(req)
}

/*
TestAPI is data for a testable API httpd
*/
type TestAPI struct {
	API      *API
	Server   *negroni.Negroni
	Listener *StoppableListener
}

func (api TestAPI) URL() string {
	return "http://127.0.0.1:" + strconv.Itoa(TestPort) + "/api/" + TestVersion
}

func (api *TestAPI) Stop() {
	api.Listener.Stop()
	api.Listener.wg.Wait()
}

/*
NewTestAPI a testing API server on port TestPort
Create and cleanup (synchronously) like so:
	testAPI, err := NewTestAPI()
	AssertNil(t, err)
	defer testAPI.Stop()
*/
func NewTestAPI() (*TestAPI, error) {
	// Set up the usual API + Negroni
	negServer := negroni.New() // add negroni.NewLogger() to see all requests
	store := sessions.NewCookieStore([]byte(TestSessionSecret))
	negServer.Use(sessions.Sessions(TestSessionCookie, store))
	tempDir, err := ioutil.TempDir(os.TempDir(), "test-api-fs")
	if err != nil {
		return nil, err
	}
	fs, err := NewLocalFileStorage(tempDir)
	if err != nil {
		return nil, err
	}
	api := NewAPI("/api/"+TestVersion, TestVersion, fs)
	negServer.UseHandler(api.Mux)

	// Set up a stoppable listener so we can clean up afterwards
	sl, err := NewStoppableListener("tcp", fmt.Sprintf(":%d", TestPort))
	if err != nil {
		return nil, err
	}
	server := http.Server{
		Handler: negServer,
	}
	// Serve up the listener and set up the waitgroup so tests can wait until the server closes
	go func() {
		sl.wg.Add(1)
		defer sl.wg.Done()
		server.Serve(sl)
	}()

	return &TestAPI{
		API:      api,
		Server:   negServer,
		Listener: sl,
	}, nil
}

func tempFile(dir string, kilobytes int) (*os.File, error) {
	f, err := ioutil.TempFile(dir, "skella-test-file")
	if err != nil {
		return nil, err
	}
	if kilobytes > 0 {
		data := make([]byte, 1024)
		n := 0
		for i := 0; i < kilobytes; i++ {
			n, err = f.Write(data)
			if err != nil || n != len(data) {
				f.Close()
				return nil, err
			}
		}
		_, err = f.Seek(0, 0)
		if err != nil {
			f.Close()
			return nil, err
		}
	}
	return f, nil
}

func compareReaderData(file1 io.Reader, file2 io.Reader) bool {
	buf1 := make([]byte, 1024)
	n1 := 0
	buf2 := make([]byte, 1024)
	n2 := 0
	for {
		n1, _ = file1.Read(buf1)
		n2, _ = file2.Read(buf2)
		if n1 != n2 {
			logger.Print("Unbalanced read: ", n1, " ", n2)
			return false
		}
		if bytes.Compare(buf1[0:n1], buf2[0:n2]) != 0 {
			logger.Print("Different buffers: ", buf1[0:n1], " ", buf2[0:n2])
			return false
		}
		if n1 == 0 {
			return true
		}
	}
}
