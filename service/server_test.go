package service

import (
	"strings"
	"testing"

	"github.com/cloudfoundry-community/go-cfenv"
)

func TestResolvesProperRabbitURL(t *testing.T) {
	fakeVCAP := []string{
		`VCAP_APPLICATION={}`,
		`VCAP_SERVICES={"cloudamqp": [{"credentials": {"uri": "amqp://foo.bar"}, "label": "p-rabbitmq", "name": "rabbit", "syslog_drain_url": "", "tags": []}]}`,
	}

	testEnv := cfenv.Env(fakeVCAP)
	cfenv, err := cfenv.New(testEnv)

	if cfenv == nil || err != nil {
		t.Errorf("cfenv is nil: %v", err.Error())
	}

	url := resolveAMQPURL(cfenv)
	if strings.Compare(url, "fake://foo") == 0 {
		t.Errorf("Got the fake URL when we should've gotten the proper URL.")
	}
}

func TestFallsBackToFakeURLWhenNoBoundService(t *testing.T) {
	fakeVCAP := []string{
		`VCAP_APPLICATION={}`,
		`VCAP_SERVICES={}`,
	}

	testEnv := cfenv.Env(fakeVCAP)
	cfenv, err := cfenv.New(testEnv)

	if cfenv == nil || err != nil {
		t.Errorf("cfenv is nil: %v", err.Error())
	}

	url := resolveAMQPURL(cfenv)
	if strings.Compare(url, "fake://foo") != 0 {
		t.Errorf("Should have gotten the fake url, but didn't, got %s instead.", url)
	}
}
