package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os/exec"
	"time"
)

type UpsMon interface {
	Start()
	Shutdown(ctx context.Context) error
}

type upsMon struct {
	logger        *slog.Logger
	lastUpsStatus UpsStatus
	lastRunScript *time.Time
	cancel        chan struct{}
	client        http.Client
}

func NewUpsMon(logger *slog.Logger) UpsMon {
	upsmon := &upsMon{
		cancel: nil,
		client: http.Client{
			Timeout:   5 * time.Second,
			Transport: http.DefaultTransport,
		},
		lastRunScript: nil,
		lastUpsStatus: StatusUnknown,
		logger:        logger,
	}
	upsmon.client.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return upsmon
}

func (u *upsMon) Start() {
	u.logger.Info("Starting UpsMon service")
	u.logger.Debug("UpsMon configuration", "ups", UpsURL, "post", PostURL, "script", Script, "rate", Rate, "timeout", Timeout, "loglevel", LogLevel)
	var response []byte
	var upsResponse UpsResponse
	var resp *http.Response
	var err error
	var lastStatus UpsStatus
	for {
		u.logger.Debug("sleep", "duration", Rate)
		time.Sleep(Rate)
		u.logger.Debug("fetching ups status", "url", UpsURL, "timeout", u.client.Timeout)
		resp, err = u.client.Get(UpsURL.String())
		if err != nil {
			u.logger.Error("error fetching ups status", "err", err.Error())
			continue
		}
		if resp.StatusCode != 200 {
			u.logger.Error("error fetching ups status", "code", resp.StatusCode)
			continue
		}
		response, err = io.ReadAll(resp.Body)
		if err != nil {
			u.logger.Error("error reading ups status", "err", err.Error())
			continue
		}
		u.logger.Debug("ups response body", "body", string(response))
		if err = json.Unmarshal(response, &upsResponse); err != nil {
			u.logger.Error("error unmarshaling up status json", err, err.Error())
			continue
		}
		u.logger.Debug("ups response", "ups", fmt.Sprintf("%+v", upsResponse))
		lastStatus = u.lastUpsStatus
		u.lastUpsStatus = upsResponse.Status
		if lastStatus == upsResponse.Status {
			continue
		}
		u.logger.Warn("ups status change", "previous", lastStatus, "current", upsResponse.Status)

		if Script != "" {
			u.execOrCancelScript(&upsResponse.Status)
		}

		if PostURL != nil {
			if err = u.postStatus(response); err != nil {
				u.logger.Error("error posting ups status", "err", err.Error())
				continue
			}
		}

	}
}

func (u *upsMon) execOrCancelScript(status *UpsStatus) {
	if *status == StatusACFail && u.lastRunScript == nil {
		u.cancel = make(chan struct{}, 1)
		go u.runScript()
	} else if *status != StatusACFail && u.cancel != nil {
		u.cancel <- struct{}{}
		close(u.cancel)
		u.cancel = nil
		u.lastRunScript = nil
	}
}

func (u *upsMon) postStatus(body []byte) error {
	u.logger.Debug("posting status", "url", PostURL, "timeout", u.client.Timeout)
	resp, err := u.client.Post(PostURL.String(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("code %d", resp.StatusCode)
	}
	u.logger.Info("status posted successfully")
	return nil
}

func (u *upsMon) runScript() {
	now := time.Now()
	u.lastRunScript = &now
	u.logger.Warn("timeout to run AC Fail script started", "timeout", Timeout)
	select {
	case <-time.After(*Timeout):
		u.logger.Warn("running AC Fail script", "script", Script)
		cmd := exec.Command(Script)
		out, err := cmd.Output()
		if err != nil {
			u.logger.Error("AC Fail script failed", "output", string(out), "error", err.Error())
		} else {
			u.logger.Info("AC Fail script success", "output", string(out))
		}
	case <-u.cancel:
		u.logger.Warn("cancelling AC Fail script", "timeout duration", time.Since(now))
	}
}

func (u *upsMon) Shutdown(ctx context.Context) error {
	u.logger.Info("Shutting down UpsMon service")
	return nil
}
