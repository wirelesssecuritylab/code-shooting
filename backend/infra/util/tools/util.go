package tools

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"time"

	"github.com/pkg/errors"

	"code-shooting/infra/errcode"
)

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func IsDockerEth(ip *net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		return ip4[0] == 172 && ip4[1] == 17 && ip4[2] == 0 && ip4[3] == 1
	}
	return false
}

func GetLocalIp() (net.IP, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
		return nil, err
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !IsDockerEth(&ipnet.IP) {
					if ipnet.IP.To4() != nil {
						return ipnet.IP, nil
					}
				}
			}
		}
	}
	return nil, nil
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func Any2Map(content interface{}) (res map[string]string, err error) {
	bytes, err := json.Marshal(content)
	if err != nil {
		return res, err
	}
	err1 := json.Unmarshal(bytes, &res)
	if err1 != nil {
		return res, err
	}
	return res, nil
}

func ListRemoveOne(items []string, item string) []string {
	res := make([]string, 0, len(items))
	for _, eachItem := range items {
		if eachItem != item {
			res = append(res, eachItem)
		}
	}
	return res
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func Sleep(t int64) {
	time.Sleep(time.Duration(t) * time.Second)
}

func IsIpListEqual(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	sort.Sort(IpSlice(a))
	sort.Sort(IpSlice(b))
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type IpSlice []string

func (p IpSlice) Len() int {
	return len(p)
}
func (p IpSlice) Less(i, j int) bool {
	return p[i] < p[j]
}
func (p IpSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func AppendInSliceWhenNotIn(all []string, subs []string) []string {
	var res = make([]string, len(all))
	copy(res, all)
	for _, sub := range subs {
		if IsContain(res, sub) {
			continue
		}
		res = append(res, sub)
	}
	return res
}

func GetCurGoFilePath() (string, error) {
	if _, filename, _, ok := runtime.Caller(1); ok {
		return path.Dir(filename), nil
	}
	return "", errors.New("GetCurGoFilePath failed")
}

func SaveFile(path string, data io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return errors.WithMessagef(errcode.ErrFileSystemError, "mkdir failed: %v", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return errors.WithMessagef(errcode.ErrFileSystemError, "create file failed: %v", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, data); err != nil {
		return errors.WithMessagef(errcode.ErrFileSystemError, "io copy failed: %v", err)
	}
	return nil
}

func RemoveFile(path string) error {
	if err := os.Remove(path); err != nil {
		return errors.WithMessagef(errcode.ErrFileSystemError, "remove file failed: %v", err)
	}
	return nil
}

/*
    移动文件处理策略：
    方案一：
	if err := os.Rename(srcPath, destPath); err != nil {  // 可能会有错误 invalid cross-device link
	 	return errors.WithMessagef(errcode.ErrFileSystemError, "move file failed: %v", err)
	}
	方案二：
	cmd := exec.Command("mv", srcPath, destPath)  // 命令操作系统不一定支持
	if _, err := cmd.Output(); err != nil {
		return errors.WithMessagef(errcode.ErrFileSystemError, "move file failed: %v", err)
	}
*/
func MoveFile(srcPath, destPath string) error {
	inputFile, err := os.Open(srcPath)
	if err != nil {
		return errors.WithMessagef(errcode.ErrFileSystemError, "open source file failed: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(destPath), 0750); err != nil {
		return errors.WithMessagef(errcode.ErrFileSystemError, "mkdir failed: %v", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return errors.WithMessagef(errcode.ErrFileSystemError, "open dest file failed: %v", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return errors.WithMessagef(errcode.ErrFileSystemError, "writing to output file failed: %v", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(srcPath)
	if err != nil {
		return errors.WithMessagef(errcode.ErrFileSystemError, "removing original file failed: %v", err)
	}
	return nil
}

func ListFileNames(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, errors.WithMessagef(errcode.ErrFileSystemError, "read dir failed: %v", err)
	}
	fileNames := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileNames = append(fileNames, file.Name())
	}
	return fileNames, nil
}

func Contains(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}
