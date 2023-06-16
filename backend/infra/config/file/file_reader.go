package file

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"

	"os"
	"path"
	"path/filepath"
	"strings"

	"code-shooting/infra/config/internal/log"
	"code-shooting/infra/config/internal/utils"
	"go.uber.org/config"
)

type FileHandler func(files []string) (map[string]interface{}, error)

func ParseYamlFile(files []string) (map[string]interface{}, error) {
	fs, err := parseFilePaths(files)
	if err != nil {
		return nil, errors.Wrapf(err, " parse file path")
	}
	configMap := make(map[string]interface{})
	for _, f := range fs {
		conf, err := parseYamlFile(f)
		if err != nil {
			return nil, errors.Wrapf(err, " parse yaml file")
		}
		for k, v := range conf {
			configMap[k] = v
		}
	}
	return configMap, nil
}

func parseFilePaths(files []string) ([]string, error) {
	var filePaths []string
	for _, file := range files {
		fs, err := getConfFiles(file)
		if err != nil {
			return nil, err
		}
		filePaths = append(filePaths, fs...)
	}
	if len(filePaths) == 0 {
		return nil, errors.New(" no supported config file")
	}
	return filePaths, nil
}

func parseYamlFile(filePath string) (map[string]interface{}, error) {
	content, _ := ioutil.ReadFile(filePath)
	yml, err := config.NewYAML(config.Source(strings.NewReader(string(content))), config.Expand(os.LookupEnv), config.Permissive())
	if err != nil {
		return nil, err
	}

	v := yml.Get("")
	conf := make(map[interface{}]interface{})
	err = v.Populate(&conf)
	if err != nil {
		return nil, fmt.Errorf("yaml unmarshal [%s] failed, %s", content, err)
	}
	return retrieveItems("", conf), nil
}

func buildConfigKey(prefix, field string) string {
	if prefix == "" {
		return field
	}
	return prefix + "#" + field
}

func retrieveItems(prefix string, subItems map[interface{}]interface{}) map[string]interface{} {

	result := map[string]interface{}{}
	for k, v := range subItems {
		k, ok := k.(string)
		if !ok {
			log.Warn("yaml tag is not string", k)
			continue
		}
		switch v.(type) {
		case map[interface{}]interface{}:
			subResult := retrieveItems(buildConfigKey(prefix, k), v.(map[interface{}]interface{}))
			for k, v := range subResult {
				result[k] = v
			}
		case []interface{}:
			arr := paresSlice(v.([]interface{}))
			result[buildConfigKey(prefix, k)] = arr
		default:
			result[buildConfigKey(prefix, k)] = v

		}

	}

	return result
}

func paresSlice(values []interface{}) interface{} {
	arr := make([]interface{}, 0)
	for _, v := range values {
		switch v.(type) {
		case []interface{}:
			arr = append(arr, paresSlice(v.([]interface{})))
		case map[interface{}]interface{}:
			arr = append(arr, retrieveItems("", v.(map[interface{}]interface{})))
		default:
			arr = append(arr, v)
		}

	}
	return arr
}

func getConfFiles(confPath string) ([]string, error) {

	confFileList := []string{}

	confInfo, err := getConfInfo(confPath)
	if err != nil {
		return []string{}, err
	}

	if confInfo.IsDir() {
		confFileList, err = getAllFilesPathBy(confPath, confFileList)
		if err != nil {
			return []string{}, err
		}
	} else if isYamlFile(confPath) {
		confFileList = []string{confPath}
	}
	return confFileList, nil
}

func getConfInfo(confPath string) (os.FileInfo, error) {

	if confPath == "" {
		return nil, errors.New("config path is empty")
	}

	fileInfo, err := os.Stat(confPath)
	if err != nil {
		return nil, err
	}
	return fileInfo, nil
}

func isYamlFile(filePath string) bool {
	ext := path.Ext(filePath)
	return ext == ".yaml" || ext == ".yml"
}

func getAllFilesPathBy(dirPath string, files []string) ([]string, error) {

	rd, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return files, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := filepath.Join(dirPath, fi.Name())
			files, err = getAllFilesPathBy(fullDir, files)
			if err != nil {
				return files, err
			}
		} else if fi.Mode().IsRegular() {
			if !isYamlFile(fi.Name()) {
				continue
			}
			fullName := filepath.Join(dirPath, fi.Name())
			files = append(files, fullName)
		}
	}
	return files, nil
}

func ConvertFile2ConfigMap(files []string) (map[string]interface{}, error) {
	configMap := make(map[string]interface{})
	for _, f := range files {
		content, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}
		key := utils.ConvertOutKeyToInner(f)
		configMap[key] = content
	}
	return configMap, nil
}
