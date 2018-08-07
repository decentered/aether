// Toolbox
// This package provides a container for functions that is universally usable. This package does not import any app dependencies, thus it should be usable by any package.

package toolbox

import (
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
)

func Round(x, unit float64) float64 {
	r := math.Round(x/unit) * unit
	formatted, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", r), 64)
	return formatted
}

func DumpStack() string {
	_, file, line, _ := runtime.Caller(1)
	_, file2, line2, _ := runtime.Caller(2)
	_, file3, line3, _ := runtime.Caller(3)
	_, file4, line4, _ := runtime.Caller(4)
	_, file5, line5, _ := runtime.Caller(5)
	return fmt.Sprintf("\nSTACK TRACE\n%s:%d\n%s:%d\n%s:%d \n%s:%d \n%s:%d\n",
		file, line, file2, line2, file3, line3, file4, line4, file5, line5)
}

// This makes me sad
func Singular(entityType string) string {
	if entityType == "boards" {
		return "board"
	} else if entityType == "threads" {
		return "thread"
	} else if entityType == "posts" {
		return "post"
	} else if entityType == "votes" {
		return "vote"
	} else if entityType == "keys" {
		return "key"
	} else if entityType == "truststates" {
		return "truststate"
	} else if entityType == "addresses" {
		return "address"
	} else {
		return ""
	}
}

func Plural(entityType string) string {
	if entityType == "board" {
		return "boards"
	} else if entityType == "thread" {
		return "threads"
	} else if entityType == "post" {
		return "posts"
	} else if entityType == "vote" {
		return "votes"
	} else if entityType == "vote" {
		return "votes"
	} else if entityType == "truststate" {
		return "truststates"
	} else if entityType == "address" {
		return "addresses"
	} else {
		return ""
	}
}

func CreatePath(path string) {
	// fmt.Printf("CreatePath called for the path %#v\n", path)
	os.MkdirAll(path, 0755)
}

func Trace() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	result := fmt.Sprintf("%s,:%d %s", frame.File, frame.Line, frame.Function)
	return result
}

func DeleteFromDisk(path string) {
	// return
	// fmt.Printf("DeleteFromDisk called for the path %#v\n", path)
	err := os.RemoveAll(path)
	if err != nil {
		panic(err)
	}
}

func IndexOf(searchString string, stringSlice []string) int {
	for key, _ := range stringSlice {
		if stringSlice[key] == searchString {
			return key
		}
	}
	return -1
}

func IndexOfInt(searchInt int, intSlice []int) int {
	for key, val := range intSlice {
		if val == searchInt {
			return key
		}
	}
	return -1
}

// GetInsecureRand gets a random number within the given range.
// WARNING: GetRand is NOT cryptographically secure! Do not use it within, as an input of, as a way to process the output of, any cryptographic process.
func GetInsecureRand(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

func GetInsecureRands(max, count int) []int {
	if max < count {
		max = count
	}
	ints := []int{}
	for count > 0 {
		rnd := GetInsecureRand(max)
		if IndexOfInt(rnd, ints) == -1 {
			ints = append(ints, rnd)
			count--
		}
	}
	return ints
}

func CnvToCutoffDays(days int) int64 {
	return int64(time.Now().Add(-(time.Duration(days) * time.Hour * time.Duration(24))).Unix())
}

func CnvToCutoffMinutes(mins int) int64 {
	return int64(time.Now().Add(-(time.Duration(mins) * time.Minute)).Unix())
}

func CnvToFutureCutoffMinutes(mins int) int64 {
	return int64(time.Now().Add((time.Duration(mins) * time.Minute)).Unix())
}

func FileExists(filePath string) bool {
	fileInfo, _ := os.Stat(filePath)
	if fileInfo == nil {
		return false
	}
	return true
}

func SplitHostPort(addr string) (string, uint16) {
	host, portAsStr, _ := net.SplitHostPort(addr)
	portAsInt, _ := strconv.Atoi(portAsStr)
	return host, uint16(portAsInt)
}
