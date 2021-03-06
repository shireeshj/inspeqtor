package daemon

import (
	"io/ioutil"
	"testing"

	"github.com/mperham/inspeqtor/metrics"
	"github.com/stretchr/testify/assert"
)

func TestBadNginxConfig(t *testing.T) {
	t.Parallel()
	src, err := metrics.Sources["nginx"](map[string]string{"port": "885u"})
	assert.Nil(t, src)
	assert.NotNil(t, err)

	src, err = metrics.Sources["nginx"](map[string]string{"port": "8080"})
	assert.Nil(t, err)
	assert.NotNil(t, src)
}

func TestNginxCollection(t *testing.T) {
	t.Parallel()
	rs := testNginxSource(nil)
	rs.client = testNginxClient("fixtures/nginx.status.txt")
	assert.NotNil(t, rs)
	hash, err := rs.runCli()
	assert.Nil(t, err)
	assert.NotNil(t, hash)

	assert.Equal(t, metrics.Map{"Active_connections": 2, "requests": 3, "Waiting": 1}, hash)

	rs.Watch("bad_metric")
	hash, err = rs.runCli()
	assert.Nil(t, err)
	assert.NotNil(t, hash)

	assert.Equal(t, metrics.Map{"Active_connections": 2, "requests": 3, "Waiting": 1}, hash)
}

func TestRealNginxConnection(t *testing.T) {
	t.Parallel()
	rs := testNginxSource(nil)
	rs.Port = "8080"
	assert.NotNil(t, rs)
	hash, err := rs.Capture()
	assert.Nil(t, err)
	assert.NotNil(t, hash)

	// brew tap homebrew/homebrew-nginx
	// brew install nginx-full --with-status
	assert.True(t, hash["requests"] > 0, "This test will fail if you don't have nginx installed")
}

func testNginxSource(mets []string) *nginxSource {
	src, err := metrics.Sources["nginx"](map[string]string{})
	if err != nil {
		panic(err)
	}
	if mets == nil {
		mets = []string{"Active_connections", "requests", "Waiting"}
	}
	for _, x := range mets {
		src.Watch(x)
	}
	return src.(*nginxSource)
}

func testNginxClient(path string) func(string, string, string) ([]byte, error) {
	return func(host string, port string, ep string) ([]byte, error) {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}
