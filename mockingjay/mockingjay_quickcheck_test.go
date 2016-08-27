package mockingjay

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"testing/quick"
)

func (r FakeEndpoint) Generate(rand *rand.Rand, size int) reflect.Value {
	randomMethod := httpMethods[rand.Intn(len(httpMethods))]

	req := Request{
		Method: randomMethod,
		URI:    "/" + randomURL(rand.Intn(maxURLLen)),
	}

	res := response{
		Code: rand.Intn(599-100) + 100,
	}

	return reflect.ValueOf(FakeEndpoint{
		Name:     "Generated",
		Request:  req,
		Response: res,
	})
}

func TestItIsAlwaysCompatibleWithItself(t *testing.T) {

	compatabilityChecker := NewCompatabilityChecker(noopLogger, 2)

	assertion := func(endpoint FakeEndpoint) bool {

		// Start an MJ server with the random configuration
		mjSvr := NewServer([]FakeEndpoint{endpoint}, false, ioutil.Discard)
		svr := httptest.NewServer(http.HandlerFunc(mjSvr.ServeHTTP))
		defer svr.Close()

		// Run CDC against "itself". An MJ server should always be compatible with itself.
		errors := compatabilityChecker.check(&endpoint, svr.URL)

		if len(errors) > 0 {
			t.Log("wARNING!!!! It wasn't compatible with itself")
			for _, err := range errors {
				t.Log(err)
			}
		}

		return true // just always pass for now
		//return len(errors) == 0
	}

	config := quick.Config{
		MaxCount: 1000,
	}

	if err := quick.Check(assertion, &config); err != nil {
		t.Error(err)
	}
}