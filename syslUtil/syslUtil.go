package syslutil

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type CompileOutput struct {
	FileName string
	Content  string
}

func writeFileToTempFile(input string, fileName string) (string, error) {
	uuid := uuid.Must(uuid.NewRandom())
	tmpDir := uuid.String()
	err := os.Mkdir(tmpDir, 0755)
	check(err)
	d1 := []byte(input)
	e := ioutil.WriteFile(tmpDir+"/"+fileName, d1, 0644)
	check(e)
	return tmpDir, nil
}

//Execute sysl command
func Execute(input string, args []string) ([]*CompileOutput, error) {
	app := "sysl"
	currentDir, _ := os.Getwd()
	inputFileName := getInputFileName(args)
	//create a temp folder for sysl and output
	tmpDir, err := writeFileToTempFile(input, inputFileName)
	check(err)
	//copy dependency for code gen
	if args[0] == "codegen" {
		index, _ := Find(args, "--transform")
		//copy the following transform file
		copy("transforms/"+args[index+1], tmpDir)
		//copy grammar file
		copy("grammars/go.gen.g", tmpDir)
	}

	workingDir := filepath.Join(currentDir, tmpDir)

	cmd := exec.Command(app, args[0:]...)
	cmd.Dir = workingDir
	_, e := cmd.Output()
	//clean up
	defer os.RemoveAll(tmpDir)
	check(e)
	resultExtension := findOutputExtension(args)
	if args[0] == "codegen" {
		resultExtension = ".go"
	}
	return readFilesFromFolder(tmpDir, resultExtension)
}

func getInputFileName(args []string) string {
	for i, item := range args {
		if strings.Contains(item, "--input=") || strings.Contains(item, "-i=") {
			fileName := strings.Split(item, "=")[1]
			//clean path
			fileName = clearPath(fileName)
			//reset input
			args[i] = "--input=" + fileName
			return fileName
		}
		if strings.Compare(item, "--input") == 0 || strings.Compare(item, "-i") == 0 {
			fileName := args[i+1]
			fileName = clearPath(fileName)
			//reset input
			args[i+1] = fileName
			return fileName
		}
	}
	//no input params, the last args will be the input name
	fileName := args[len(args)-1]
	fileName = clearPath(fileName)
	//reset input
	args[len(args)-1] = fileName
	return fileName
}

func clearPath(input string) string {
	ss := strings.Split(input, "/")
	return ss[len(ss)-1]
}

func findOutputExtension(args []string) string {
	for i, item := range args {
		if strings.Contains(item, "--output=") || strings.Contains(item, "-o=") {
			ss := strings.Split(item, ".")
			return "." + ss[len(ss)-1]
		}
		if strings.Compare(item, "--output") == 0 || strings.Compare(item, "-o") == 0 {
			ss := strings.Split(args[i+1], ".")
			return "." + ss[len(ss)-1]
		}
	}
	return ""
}

func copy(sourceFile, dstDir string) {
	input, err := ioutil.ReadFile(sourceFile)
	check(err)
	ss := strings.Split(sourceFile, "/")
	fileName := ss[len(ss)-1]
	destinationFile := dstDir + "/" + fileName
	err = ioutil.WriteFile(destinationFile, input, 0644)
	check(err)
}

func readFilesFromFolder(tmpDir string, resultExtension string) ([]*CompileOutput, error) {
	currentDir, _ := os.Getwd()
	workingDir := filepath.Join(currentDir, tmpDir)
	files, err := ioutil.ReadDir(workingDir)
	check(err)
	var result []*CompileOutput
	for _, file := range files {
		if !strings.Contains(file.Name(), resultExtension) {
			continue
		}
		binary, e := ioutil.ReadFile(workingDir + "/" + file.Name())
		check(e)
		content := string(binary)
		if strings.Contains(file.Name(), ".png") {
			content = base64.StdEncoding.EncodeToString(binary)
		}
		content = string(content)
		result = append(result, &CompileOutput{FileName: file.Name(), Content: content})
	}

	return result, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
