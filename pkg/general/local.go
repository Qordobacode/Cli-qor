package general

import (
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/qordobacode/cli-v2/pkg/general/log"
	"github.com/qordobacode/cli-v2/pkg/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	defaultFilePerm os.FileMode = 0666
)

var (
	forbiddenInFileNameSymbols, _ = regexp.Compile(`[:?!\\*/|<>]`)
	invalidationPeriod            = time.Hour * 4
)

// Local implements pkg.Local
type Local struct {
	Config *types.Config
}

// Read function reads file locally with specified path
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

// Write function store body parameter as a file locally
func (l *Local) Write(fileName string, body []byte) {
	err := ioutil.WriteFile(fileName, body, defaultFilePerm)
	if err != nil {
		log.Errorf("error occurred on writing file: %v", err)
	}
}

// BuildDirectoryFilePath according to stored file name and version
func (l *Local) BuildDirectoryFilePath(j *types.File2Download, matchFilepathName []string, suffix string, isFilePathPattern bool) string {
	if j.File.Version != "" {
		if suffix != "" {
			suffix = j.File.Version + "_" + suffix
		} else {
			suffix = j.File.Version
		}
	}
	resultName := ""
	if isFilePathPattern {
		// handle via file path pattern
		resultName = l.buildFileName(j.File, suffix)
		fileDir := l.buildDirName(j, matchFilepathName)
		resultName = filepath.Join(fileDir, resultName)
	} else {
		// handle via `download.target
		resultName = l.buildTargetFileName(j, suffix)
		fileDir := ""
		resultName = filepath.Join(fileDir, resultName)
	}
	return resultName
}

func (l *Local) buildDirName(j *types.File2Download, matchFilepathNames []string) string {
	fileDir := j.File.Filepath
	wasUpdated := false
	for _, matchFilepathName := range matchFilepathNames {
		if strings.Contains(fileDir, "/"+matchFilepathName+"/") {
			fileDir = strings.ReplaceAll(fileDir, "/"+matchFilepathName+"/", "/"+j.ReplaceIn+"/")
			wasUpdated = true
			break
		} else if strings.Contains(fileDir, "/"+matchFilepathName+"-") {
			fileDir = strings.ReplaceAll(fileDir, "/"+matchFilepathName+"-", "/"+j.ReplaceIn+"-")
			wasUpdated = true
			break
		} else if strings.Contains(fileDir, "-"+matchFilepathName+"/") {
			fileDir = strings.ReplaceAll(fileDir, "-"+matchFilepathName+"/", "-"+j.ReplaceIn+"/")
			wasUpdated = true
			break
		}
	}
	if !wasUpdated {
		log.Infof("In fileDir %s source patterns '%s' was not found. Don't update it then.", fileDir, matchFilepathNames)
	}
	fileDir = strings.ReplaceAll(fileDir, "/", string(filepath.Separator))
	fileDir = filepath.Dir(fileDir)
	return fileDir
}

func (*Local) buildFileName(file *types.File, suffix string) string {
	fileNames := strings.Split(file.Filename, ".")
	if len(fileNames) > 1 && suffix != "" {
		fileNames[len(fileNames)-2] = fileNames[len(fileNames)-2] + "_" + suffix
	}
	fileName := strings.Join(fileNames, ".")
	fileName = forbiddenInFileNameSymbols.ReplaceAllString(fileName, "")
	return fileName
}

// FileExists checks for file existence
func (*Local) FileExists(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && !stat.IsDir()
}

// QordobaHome returns path to qordoba's home folder
func (l *Local) QordobaHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Errorf("error occurred on home dir retrieval: %v", err)
		return "", err
	}
	return home + string(os.PathSeparator) + ".qordoba", nil
}

// PutInHome put file in qordoba's home directory. Used for response caching
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

// LoadCached load stored in home directory file. If `invalidationPeriod` hasn't passed -> file is returned, else - return nil
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

// FilesInFolder returns all files with specified file path
func (l *Local) FilesInFolder(filePath string, isRecursive bool) []string {
	fileMap := make(map[string]bool)
	matches, err := filepath.Glob(filePath)
	result := make([]string, 0)
	if matches == nil || len(matches) == 0 {
		// filePath is not regexp
		res, err := filepath.Abs(filePath)
		if err == nil {
			log.Infof("No files were found for %s", res)
			return result
		}
		log.Infof("No files were found for %s", filePath)
		return result
	}
	if err == nil {
		for _, match := range matches {
			addFilesInMap(match, fileMap, isRecursive)
		}
	} else {
		log.Infof("err occurred on files add: %v", err)
		return result
	}
	addFilesInMap(filePath, fileMap, isRecursive)
	for k := range fileMap {
		result = append(result, k)
	}
	return result
}

// recursive function for finding all files in map
func addFilesInMap(path string, fileMap map[string]bool, isRecursive bool) {
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
			if isRecursive {
				addFilesInMap(fileName, fileMap, isRecursive)
			} else if !info.IsDir() {
				thisFilePath, err := filepath.Abs(fileName)
				if err == nil {
					fileMap[thisFilePath] = true
				}
			}
		}
	} else {
		path, err = filepath.Abs(path)
		if err == nil {
			fileMap[path] = true
		}
	}
}

// RenderTable2Stdin takes header and data together and print out in STDOUT as a table
func (*Local) RenderTable2Stdin(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(data)
	table.Render() // Send output
}

func (l *Local) buildTargetFileName(file2Download *types.File2Download, s string) string {
	//i18n/<language_code>/translations.json  ->  i18n/zh-cn/translations.json
	//folder1/values-<language_lang_code>/strings.xml  ->  folder1/values-en/strings.xml
	//config/locales/server.<language_code>.yml  ->  config/locales/server.fr-fr.yml
	//folder2/<language_name>/strings.xml  ->  folder2/Chinese/strings.xml
	//folder3/strings.<language_name_cap>  ->  folder3/strings.French
	//<language_name_allcap>.locale  ->  FRENCH.locale
	//./<language_code>/<filename>.<extension> -> downloads to current location of CLI  zh-cn/translations.json
	filename := file2Download.File.Filename
	fileNames := strings.Split(filename, ".")
	mimeType := ""
	if len(fileNames) > 1 {
		mimeType = fileNames[len(fileNames)-1]
		filename = strings.Join(fileNames[0:len(fileNames)-1], ".")
	}
	resultName := l.Config.Download.Target
	if resultName == "" {
		return file2Download.File.Filename
	}
	resultName = strings.ReplaceAll(resultName, "<extension>", mimeType)
	resultName = strings.ReplaceAll(resultName, "<filename>", filename)

	for k, v := range file2Download.ReplaceMap {
		resultName = strings.ReplaceAll(resultName, k, v)
	}
	return resultName
}
