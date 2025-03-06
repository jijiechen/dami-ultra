package behavior_tests

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/yuchanns/kong-exercise-microservices/startup"
	"gorm.io/gorm"
)

type behaviorTest = func(assert *require.Assertions, b *reqBuilder)

func generateRandomInt(min, max int64) int64 {
	n := max - min + 1

	bigInt, _ := rand.Int(rand.Reader, big.NewInt(n))

	return min + bigInt.Int64()
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
}

func prepareSchema(assert *require.Assertions, orgID string) {
	db, err := gorm.Open(sqlite.Open("catalog_"+orgID+".db"), &gorm.Config{})
	assert.Nil(err)

	f, err := os.Open("./schema.sql")
	assert.Nil(err)
	defer f.Close()

	buf, err := io.ReadAll(f)
	assert.Nil(err)

	assert.Nil(db.Exec(string(buf)).Error)
}

func prepareRouter(assert *require.Assertions) *gin.Engine {
	r := gin.Default()

	err := startup.Register(r)
	assert.Nil(err)

	return r
}

func toBuffer(in any) *bytes.Buffer {
	if in == nil {
		return bytes.NewBuffer(nil)
	}
	b, _ := json.Marshal(in)
	return bytes.NewBuffer(b)
}

func parse[T any](assert *require.Assertions, b []byte) *T {
	var t T
	result := gjson.ParseBytes(b).Get("data")
	err := json.Unmarshal([]byte(result.String()), &t)
	assert.Nil(err)
	return &t
}

type reqBuilder struct {
	OrgID string

	Assert *require.Assertions

	Engine *gin.Engine
}

func (r *reqBuilder) Record(method, path string, body any) []byte {
	w := httptest.NewRecorder()
	req := r.newRequest(method, path, body)
	r.Engine.ServeHTTP(w, req)

	r.Assert.Equal(http.StatusOK, w.Code)

	return w.Body.Bytes()
}

func (r *reqBuilder) newRequest(method, path string, body any) *http.Request {
	path, err := url.JoinPath("/api/v1", path)
	r.Assert.Nil(err)

	if method == http.MethodGet && body != nil {
		jsonData, err := json.Marshal(body)
		r.Assert.Nil(err)

		var queryParams map[string]interface{}
		err = json.Unmarshal(jsonData, &queryParams)
		r.Assert.Nil(err)

		params := url.Values{}
		for key, value := range queryParams {
			params.Add(key, fmt.Sprintf("%v", value))
		}

		if len(params) > 0 {
			path = path + "?" + params.Encode()
		}
		body = nil
	}
	req, err := http.NewRequest(method, path, toBuffer(body))
	r.Assert.Nil(err)

	req.Header.Set("x-organization-id", r.OrgID)
	if method != http.MethodGet {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func TestBehavior(t *testing.T) {
	var tests = []behaviorTest{
		testServiceCreateAndGet,
		testServiceList,
		testServiceUpdate,
	}

	assert := require.New(t)
	orgID := "test"

	prepareSchema(assert, orgID)
	r := prepareRouter(assert)

	for i := range tests {
		test := tests[i]

		fullName := runtime.FuncForPC(reflect.ValueOf(test).Pointer()).Name()
		parts := strings.Split(fullName, ".")
		testName := strings.TrimPrefix(parts[len((parts))-1], "test")

		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			assert := require.New(t)
			b := &reqBuilder{OrgID: orgID, Assert: assert, Engine: r}

			test(assert, b)
		})
	}
}
