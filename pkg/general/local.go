package general

import (
	"errors"
	"fmt"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Local struct {
}

const (
	defaultFilePerm os.FileMode = 0666
)

var (
	forbiddenInFileNameSymbols, _ = regexp.Compile(`[:?!\\*/|<>]`)
	invalidationPeriod            = time.Hour * 4
)

func (l *Local) Read(path string) ([]byte, error) {
	if path == "" {
		log.Errorf("Path for config shouldn't be empty")
		return nil, errors.New("config path can't be empty")
	}
	if !l.FileExists(path) {
		log.Errorf("file not found: %v", path)
		return nil, fmt.Errorf("file not found: %v", path)
	}
	// read config from file
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("problem on file read: %v", err)
		return nil, err
	}
	return bytes, err
}

func (*Local) Write(fileName string, body []byte) {
	err := ioutil.WriteFile(fileName, body, defaultFilePerm)
	if err != nil {
		log.Errorf("error occurred on writing file: %v", err)
	}
}

// BuildFileName according to stored file name and version
func (*Local) BuildFileName(file *types.File, suffix string) string {
	fileNames := strings.SplitN(file.Filename, ".", 2)
	if file.Version != "" {
		if suffix != "" {
			suffix = suffix + "_" + file.Version
		} else {
			suffix = file.Version
		}
	}
	resultName := file.Filename
	if suffix != "" {
		if len(fileNames) > 1 {
			resultName = fileNames[0] + "_" + suffix + "." + fileNames[1]
		}
		resultName = file.Filename + "_" + suffix
	}
	resultName = forbiddenInFileNameSymbols.ReplaceAllString(resultName, "")
	return resultName
}

// FileExists checks for file existence
func (*Local) FileExists(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && !stat.IsDir()
}

func (c *Local) QordobaHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Errorf("error occurred on home dir retrieval: %v", err)
		return "", err
	}
	return home + string(os.PathSeparator) + ".qordoba", nil
}

func (l *Local) PutInHome(fileName string, body []byte) {
	qordobaHome, err := l.QordobaHome()
	if err != nil {
		return
	}
	err = os.MkdirAll(qordobaHome, os.ModePerm)
	if err != nil {
		log.Errorf("error occurred on creating qordoba's folder: %v", err)
	}
	path := qordobaHome + string(os.PathSeparator) + fileName
	err = ioutil.WriteFile(path, body, 0644)
	if err != nil {
		log.Errorf("error occurred on writing config: %v", err)
	}
}

func (l *Local) LoadCached(cachedFileName string) ([]byte, error) {
	qordobaHome, err := l.QordobaHome()
	if err != nil {
		return nil, err
	}
	workspaceFilePath := qordobaHome + string(os.PathSeparator) + cachedFileName
	file, err := os.Stat(workspaceFilePath)
	if err != nil {
		return nil, err
	}
	modifiedtime := file.ModTime()
	// don't use cached workspace if 1 day has came
	if modifiedtime.Add(invalidationPeriod).Before(time.Now()) {
		return nil, errors.New("outdated file")
	}

	if !l.FileExists(workspaceFilePath) {
		return nil, fmt.Errorf("cached workspace file was not found")
	}
	return l.Read(workspaceFilePath)
}

func (l *Local) FilesInFolder(filePath string) []string {
	fileMap := make(map[string]bool)
	matches, err := filepath.Glob(filePath)
	result := make([]string, 0)
	if err == nil {
		for _, match := range matches {
			fileAbsPath, err := filepath.Abs(match)
			if err == nil {
				addIfFiles(fileAbsPath, fileMap)
			}
		}
	} else {
		log.Infof("err occurred on files add: %v", err)
		return result
	}
	addFilesInMap(filePath, fileMap)
	for k, _ := range fileMap {
		result = append(result, k)
	}

	return result
}

func addFilesInMap(path string, fileMap map[string]bool) {
	pathStat, err := os.Stat(path)
	if err != nil {
		return
	}
	if pathStat.IsDir() {
		infos, err := ioutil.ReadDir(path)
		if err != nil {
			return
		}
		for _, info := range infos {
			fileName := filepath.Join(path, info.Name())
			fileMap[fileName] = true
		}
	} else {
		fileMap[path] = true
	}
}

func addIfFiles(path string, fileMap map[string]bool) {
	pathStat, err := os.Stat(path)
	if err != nil {
		return
	}
	if !pathStat.IsDir() {
		fileMap[path] = true
	}
}
