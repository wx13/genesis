package installer

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/kardianos/osext"
)

func getFilesToArchive(allFiles []string, tmpdir string) []string {
	files := []string{}
	for _, file := range allFiles {
		if strings.HasPrefix(file, tmpdir) {
			p, err := filepath.Rel(tmpdir, file)
			if err == nil {
				files = append(files, p)
			}
		}
	}
	return files
}

func getExec() (string, []byte) {
	execname, _ := osext.Executable()
	execbody, err := ioutil.ReadFile(execname)
	if err != nil {
		fmt.Println("Error: cannot read executable (self):", execname, err)
	}
	return execname, execbody
}

func (inst *Installer) Build(dirs []string) {

	fmt.Println("Building the self-contained executable...")

	files := getFilesToArchive(inst.Files(), inst.Tmpdir)

	execname, execbody := getExec()

	// Create the zip archive.
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	w.SetOffset(int64(len(execbody)))
	addFilesToArchive(w, files, dirs)

	// Append zip to executable.
	execbody = append(execbody, buf.Bytes()...)

	// Write out executable.
	err := ioutil.WriteFile(execname+".x", execbody, 0755)
	if err != nil {
		fmt.Println("Error writing to zip file:", err)
		return
	}

	fmt.Println("Done building archive.")

}

func addFilesToArchive(w *zip.Writer, files, dirs []string) {

	fmt.Println("Adding files to archive:")
	for _, file := range files {
		fmt.Println("   ", file)

		f, err := w.Create(file)
		if err != nil {
			fmt.Println("Cannot add file to archive:", file, err)
			continue
		}
		body, err := readFile(file, dirs)
		if err != nil {
			fmt.Printf("Could not read file %s in directories %+v, because of %+v\n", file, dirs, err)
			continue
		}
		_, err = f.Write(body)
		if err != nil {
			fmt.Println("Cannot write file contents to archive:", file, err)
			continue
		}
	}

	err := w.Close()
	if err != nil {
		fmt.Println("Cannot close archive:", err)
	}

}

func readFile(file string, dirs []string) ([]byte, error) {
	var err error
	var body []byte
	if len(dirs) == 0 {
		dirs = []string{""}
	}
	for _, dir := range dirs {
		filename := path.Join(dir, file)
		body, err = ioutil.ReadFile(filename)
		if err == nil {
			return body, nil
		}
	}
	return body, err
}
