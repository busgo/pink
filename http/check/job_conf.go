package check

import (
	"errors"
	"github.com/busgo/pink/http/model"
	"github.com/busgo/pink/pkg/protocol"
	"github.com/robfig/cron"
	"strings"
)

func JobConfAddRequest(request *model.JobConfAddRequest) error {

	if strings.TrimSpace(request.Name) == "" {
		return errors.New("name is nil")
	}

	if strings.TrimSpace(request.Group) == "" {
		return errors.New("group is nil")
	}
	if strings.TrimSpace(request.Target) == "" {
		return errors.New("target is nil")
	}
	if !(protocol.JobState(request.State) == protocol.JobNormalState || protocol.JobState(request.State) == protocol.JobStopState) {
		return errors.New("state is error")
	}
	if strings.TrimSpace(request.Cron) == "" {
		return errors.New("cron is nil")
	}
	_, err := cron.Parse(request.Cron)
	return err
}

func JobConfUpdateRequest(request *model.JobConfUpdateRequest) error {

	if strings.TrimSpace(request.Id) == "" {
		return errors.New("id is nil")
	}
	if strings.TrimSpace(request.Name) == "" {
		return errors.New("name is nil")
	}

	if !(protocol.JobState(request.State) == protocol.JobNormalState || protocol.JobState(request.State) == protocol.JobStopState) {
		return errors.New("state is error")
	}
	if strings.TrimSpace(request.Group) == "" {
		return errors.New("group is nil")
	}
	if strings.TrimSpace(request.Target) == "" {
		return errors.New("target is nil")
	}

	if strings.TrimSpace(request.Cron) == "" {
		return errors.New("cron is nil")
	}
	_, err := cron.Parse(request.Cron)

	return err
}

func JobConfDeleteRequest(request *model.JobConfDeleteRequest) error {
	if strings.TrimSpace(request.Id) == "" {
		return errors.New("id is nil")
	}
	if strings.TrimSpace(request.Group) == "" {
		return errors.New("group is nil")
	}
	return nil
}
