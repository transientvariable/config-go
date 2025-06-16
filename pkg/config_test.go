package config

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	testDataDir         = "./testdata"
	testConfigFile      = testDataDir + "/application.yaml"
	testEnvVarTimeName  = "TEST_APP_TIME"
	testEnvVarTimeValue = ` 2025-05-13T01:38:00Z ` // format time.RFC3339
)

type ConfigTestSuite struct {
	suite.Suite
}

func NewConfigTestSuite() *ConfigTestSuite {
	return &ConfigTestSuite{}
}

func (t *ConfigTestSuite) SetupTest() {
	if err := os.Setenv(testEnvVarTimeName, testEnvVarTimeValue); err != nil {
		t.T().Fatal(err)
	}

	if err := Load(WithFilePath(testConfigFile)); err != nil {
		t.T().Fatal(err)
	}

	fmt.Printf("test configuration:\n%s\n", config.String())

	assert.Equal(t.T(), Root().String(), defaultRoot)
	assert.Equal(t.T(), 13, Size())
}

func TestConfig_TestSuite(t *testing.T) {
	suite.Run(t, NewConfigTestSuite())
}

func (t *ConfigTestSuite) TestConfig_Values() {
	r := Root()
	assert.Equal(t.T(), ValueMustResolve(r.Join("application.name").String()), "test-app")
	assert.Equal(t.T(), ValueMustResolve(r.Join("application.version").String()), "v1.0.0")

	v := r.Join("value") // <root>.value
	assert.True(t.T(), BoolMustResolve(v.Join("bool").String()))
	assert.Equal(t.T(), DurationMustResolve(v.Join("duration").String()).Milliseconds(), int64(30_000))
	assert.Equal(t.T(), FloatMustResolve(v.Join("float").String()), 1.168)
	assert.Equal(t.T(), IntMustResolve(v.Join("int").String()), 138)
	assert.Equal(t.T(), MultiaddrMustResolve(v.Join("multiaddr").String()).String(), "/dns4/example.com/tcp/9003")
	assert.Equal(t.T(), SizeBytesMustResolve(v.Join("sizeBytes").String()), int64(1048576))
	assert.Equal(t.T(), TimeMustResolve(v.Join("time").String()).Format(time.RFC3339), strings.TrimSpace(testEnvVarTimeValue))
	assert.Equal(t.T(), URLMustResolve(v.Join("url").String()).String(), "https://example.com:9003")
}
