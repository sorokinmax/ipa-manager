package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
)

func deleteFile(path string) {
	log.Printf("Deleting %s ...\n", path)
	err := os.Remove(path)
	if isError(err) {
		return
	}
	log.Printf("File %s is deleted\n", path)
}

func filesEnum(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err)
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames
}

// ReadDir reads the directory named by dirname and returns
// a list of directory entries sorted by filename.
func ReadDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}

func CopyFile(src string, dst string, file string) bool {
	args := []string{src, dst, file}
	stdout, _ := exec.Command("robocopy", args...).Output()
	log.Println(Convert(41, stdout))
	return true
}

// ReadFileLines reads file by line
func ReadFileLines(path string) (lines []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)

	var line string
	for {
		line, err = reader.ReadString('\n')

		// Process the line here.
		re := regexp.MustCompile(`\r?\n`)
		lines = append(lines, re.ReplaceAllString(line, ""))

		if err != nil {
			break
		}
	}

	if err != io.EOF {
		log.Fatalf(" > Failed!: %v\n", err)
	}
	err = nil

	return lines, err
}

func RemoveDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	err = os.RemoveAll(dir)
	if err != nil {
		return err
	}

	return nil
}

func CopyDir(src string, dst string) bool {
	args := []string{src, dst, "/e"}
	exec.Command("robocopy", args...).Output()
	return true
}

func RenameFile(RenameFile string, newName string) error {
	err := os.Rename(RenameFile, newName)
	if err != nil {
		return err
	}
	return nil
}
