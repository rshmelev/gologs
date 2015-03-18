package libgologs

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
)

func SliceToString(msg []interface{}) string {
	var buffer bytes.Buffer
	AddSliceToBuffer(&buffer, msg)
	return buffer.String()
}

func AddSliceToBuffer(buf *bytes.Buffer, msg []interface{}) {
	for _, v := range msg {
		buf.WriteString(fmt.Sprintf("%+v", v))
	}
}

func RegexExtract(rx, str string) (string, error) {
	r, err := regexp.Compile(rx)
	if err != nil {
		return "", err
	}
	sm := r.FindStringSubmatch(str)
	if len(sm) > 0 {
		return sm[1], nil
	}
	return "", errors.New("group unavailable")
}
