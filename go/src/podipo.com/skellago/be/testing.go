package be

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
)

var (
	TEST_VERSION        = "0.T.0"
	TEST_PORT           = 44556677
	TEST_SESSION_COOKIE = "test_session"
	TEST_SESSION_SECRET = "NotVerySecret"
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
	req.Header.Add("Accept", AcceptHeaderPrefix+TEST_VERSION)
	return client.Do(req)
}

/*
	Data for a testable API httpd
*/
type TestAPI struct {
	API      *API
	Server   *negroni.Negroni
	Listener *StoppableListener
}

func (api TestAPI) URL() string {
	return "http://127.0.0.1:" + strconv.Itoa(TEST_PORT) + "/api/" + TEST_VERSION
}

func (api *TestAPI) Stop() {
	api.Listener.Stop()
	api.Listener.wg.Wait()
}

/*
	Creates a testing API server on port TEST_PORT
	Create and cleanup (synchronously) like so:
		testAPI, err := NewTestAPI()
		AssertNil(t, err)
		defer testAPI.Stop()
*/
func NewTestAPI() (*TestAPI, error) {
	// Set up the usual API + Negroni
	negServer := negroni.New() // add negroni.NewLogger() to see all requests
	store := sessions.NewCookieStore([]byte(TEST_SESSION_SECRET))
	negServer.Use(sessions.Sessions(TEST_SESSION_COOKIE, store))
	api := NewAPI("/api/"+TEST_VERSION, TEST_VERSION)
	negServer.UseHandler(api.Mux)

	// Set up a stoppable listener so we can clean up afterwards
	sl, err := NewStoppableListener("tcp", fmt.Sprintf(":%d", TEST_PORT))
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
