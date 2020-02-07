package diskutil

import (
	"errors"
	"strconv"
	"strings"
)

func parseFiled(line, filed string, targettype int) (interface{}, error) {
	fileds := strings.SplitN(line, ":", 2)
	if len(fileds) != 2 {
		return nil, errors.New("format illegal: " + line)
	}

	data := strings.TrimSpace(fileds[1])
	if targettype == typeString {
		return data, nil
	} else if targettype == typeInt {
		value, err := strconv.ParseInt(data, 10, 0)
		if err != nil {
			return nil, err
		}
		return int(value), nil
	} else if targettype == typeUint64 {
		value, err := strconv.ParseUint(data, 10, 0)
		if err != nil {
			return nil, err
		}
		return value, nil
	}
	return nil, errors.New("type not supported")
}
