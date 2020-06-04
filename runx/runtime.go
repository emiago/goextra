package helper

import (
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func GetFunctionName(fn interface{}) string {
	fnname := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	arr := strings.SplitN(fnname, ".", 2)
	if len(arr) > 1 {
		return arr[1]
	}
	return arr[0]
}

func TimeUnix() int64 {
	return time.Now().Unix()
	// strconv.FormatInt(time.Now().Unix(), 10)
}

func TimeUnixString() string {
	return strconv.FormatInt(TimeUnix(), 10)
}

/*
	Macro for checking is string numeric value
*/

func DumpStack() []byte {
	const size = 8192
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, false)]
	return buf
}

func DumpFullStack() []byte {
	const size = 32768
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, true)]
	return buf
}

func TimeTrack(start time.Time, cb func(time.Duration)) {
	elapsed := time.Since(start)
	cb(elapsed)
}

func GetGOMAXPROCS() int {
	return runtime.GOMAXPROCS(0)
}
