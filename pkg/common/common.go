package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/signal"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/docker/docker/api/types"
)

var cutNewlineReplacer = strings.NewReplacer("\r", "", "\n", "")

// StructToJSON convert struct to json.
func StructToJSON(i interface{}) string {
	j, err := json.Marshal(i)
	if err != nil {
		return ""
	}

	out := new(bytes.Buffer)
	json.Indent(out, j, "", "    ")
	return out.String()
}

// SortKeys sort keys.
func SortKeys(keys []string) []string {
	sort.Strings(keys)
	return keys
}

// GetOSenv get os environment.
func GetOSenv(env string) string {
	keyval := strings.SplitN(env, "=", 2)
	if keyval[1][:1] == "$" {
		keyval[1] = os.Getenv(keyval[1][1:])
		return strings.Join(keyval, "=")
	}

	return env
}

// ParseDateToString parse date to string.
func ParseDateToString(unixtime int64) string {
	t := time.Unix(unixtime, 0)
	return t.Format("2006/01/02 15:04:05")
}

// ParseSizeToString parse size to string.
func ParseSizeToString(size int64) string {
	mb := float64(size) / 1024 / 1024
	return fmt.Sprintf("%.1fMB", mb)
}

// ParsePortToString parse port to string.
func ParsePortToString(ports []types.Port) string {
	var port string
	for _, p := range ports {
		if p.PublicPort == 0 {
			port += fmt.Sprintf("%d/%s ", p.PrivatePort, p.Type)
		} else {
			port += fmt.Sprintf("%s:%d->%d/%s ", p.IP, p.PublicPort, p.PrivatePort, p.Type)
		}
	}
	return port
}

// ParseRepoTag parse image repo and tag.
func ParseRepoTag(repoTag string) (string, string) {
	tmp := strings.Split(repoTag, ":")
	tag := tmp[len(tmp)-1]
	repo := strings.Join(tmp[0:len(tmp)-1], ":")
	return repo, tag
}

// ParseLabels parse image labels.
func ParseLabels(labels map[string]string) string {
	if len(labels) < 1 {
		return ""
	}

	var result string
	for label, value := range labels {
		result += fmt.Sprintf("%s=%s ", label, value)
	}

	return result
}

// DateNow return date time.
func DateNow() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

// CutNewline cut new line.
func CutNewline(i string) string {
	return cutNewlineReplacer.Replace(i)
}

func getTermSize(fd uintptr) (int, int) {
	var sz struct {
		rows uint16
		cols uint16
	}

	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL,
		fd, uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&sz)))
	return int(sz.cols), int(sz.rows)
}

// IsTerminalWindowSizeThanZero check terminal window size
func IsTerminalWindowSizeThanZero() bool {
	out, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return false
	}

	defer out.Close()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGWINCH, syscall.SIGINT)

	for {
		// check terminal window size
		termw, termh := getTermSize(out.Fd())
		if termw > 0 && termh > 0 {
			return true
		}

		select {
		case signal := <-signalCh:
			switch signal {
			// when the terminal window size is changed
			case syscall.SIGWINCH:
				continue
			// use ctrl + c to cancel
			case syscall.SIGINT:
				return false
			}
		}
	}
}

func FormatDecimalBytes(b int) string {
	n := float64(b)
	units := []string{"B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	for _, unit := range units {
		if n > math.Pow(10, 3) {
			n /= math.Pow(10, 3)
		} else {
			val := fmt.Sprintf("%.2f%s", n, unit)
			if val == "0.00B" {
				return "0B"
			}
			return val
		}
	}
	return "a lot"
}

// from Dave C's answer at https://stackoverflow.com/questions/20767724/converting-unknown-interface-to-float64-in-golang
func GetFloat(unk interface{}) (float64, error) {
	floatType := reflect.TypeOf(float64(0))
	stringType := reflect.TypeOf("")

	switch i := unk.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case string:
		return strconv.ParseFloat(i, 64)
	default:
		v := reflect.ValueOf(unk)
		v = reflect.Indirect(v)
		if v.Type().ConvertibleTo(floatType) {
			fv := v.Convert(floatType)
			return fv.Float(), nil
		} else if v.Type().ConvertibleTo(stringType) {
			sv := v.Convert(stringType)
			s := sv.String()
			return strconv.ParseFloat(s, 64)
		} else {
			return math.NaN(), fmt.Errorf("Can't convert %v to float64", v.Type())
		}
	}
}
