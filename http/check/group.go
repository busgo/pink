package check

import (
	"errors"
	"github.com/busgo/pink/http/model"
	"strings"
)

func AddGroupRequest(request *model.AddGroupRequest) error {

	if strings.TrimSpace(request.Name) == "" {
		return errors.New("name is nil")
	}
	if strings.TrimSpace(request.Remark) == "" {
		return errors.New("remark is nil")
	}

	return nil
}
