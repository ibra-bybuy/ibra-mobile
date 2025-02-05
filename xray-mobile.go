package XRay

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	_ "github.com/xtls/xray-core/main/distro/all"

	v2net "github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf/serial"
)

type Logger interface {
	LogInput(s string)
}

var coreInstance *core.Instance

func SetMemoryLimit() {
	debug.SetGCPercent(10)
	debug.SetMemoryLimit(30 * 1024 * 1024)
}

func Start(config []byte, logger Logger) error {
	conf, err := serial.DecodeJSONConfig(bytes.NewReader(config))
	if err != nil {
		logger.LogInput("Config load error: " + err.Error())
		return err
	}
	pbConfig, err := conf.Build()
	if err != nil {
		return err
	}
	instance, err := core.New(pbConfig)
	if err != nil {
		logger.LogInput("Create XRay error: " + err.Error())
		return err
	}
	err = instance.Start()
	if err != nil {
		logger.LogInput("Start XRay error: " + err.Error())
	}
	coreInstance = instance
	return nil
}

func Stop() {
	coreInstance.Close()
}

func GetVersion() string {
	return core.Version()
}

func MeasureDelay(url string) (int64, error) {
	delay, err := measureInstDelay(context.Background(), coreInstance, url)
	return delay, err
}

func MeasureOutboundDelay(ConfigureFileContent string, url string) (int64, error) {
	config, err := serial.LoadJSONConfig(strings.NewReader(ConfigureFileContent))
	if err != nil {
		return -1, err
	}

	// dont listen to anything for test purpose
	config.Inbound = nil
	// config.App: (fakedns), log, dispatcher, InboundConfig, OutboundConfig, (stats), router, dns, (policy)
	// keep only basic features
	config.App = config.App[:5]

	inst, err := core.New(config)
	if err != nil {
		return -1, err
	}

	inst.Start()
	delay, err := measureInstDelay(context.Background(), inst, url)
	inst.Close()
	return delay, err
}

func measureInstDelay(ctx context.Context, inst *core.Instance, url string) (int64, error) {
	if inst == nil {
		return -1, errors.New("core instance nil")
	}

	tr := &http.Transport{
		TLSHandshakeTimeout: 6 * time.Second,
		DisableKeepAlives:   true,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dest, err := v2net.ParseDestination(fmt.Sprintf("%s:%s", network, addr))
			if err != nil {
				return nil, err
			}
			return core.Dial(ctx, inst, dest)
		},
	}

	c := &http.Client{
		Transport: tr,
		Timeout:   12 * time.Second,
	}

	if len(url) <= 0 {
		url = "https://www.google.com/generate_204"
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	start := time.Now()
	resp, err := c.Do(req)
	if err != nil {
		return -1, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return -1, fmt.Errorf("status != 20x: %s", resp.Status)
	}
	resp.Body.Close()
	return time.Since(start).Milliseconds(), nil
}
