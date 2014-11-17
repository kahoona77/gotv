package service

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

//File represents af File in gotv
type File struct {
	Name   string `json:"name"`
	Folder string `json:"folder"`
	Dir    bool   `json:"dir"`
	Size   int64  `json:"size"`
}

// FilesService handls files
type FilesService struct {
	*Context
}

//GetFiles loads all files from the download and temp-dir
func (fs *FilesService) GetFiles() []File {
	var files []File
	settings := fs.GetSettings()
	//Downloads-Dir
	dir := filepath.FromSlash(settings.DownloadDir)
	currentFiles, _ := ioutil.ReadDir(dir)
	for _, f := range currentFiles {
		file := File{Name: f.Name(), Size: f.Size(), Folder: "Downloads", Dir: f.IsDir()}
		files = append(files, file)
	}

	//Temp-Dir
	dir = filepath.FromSlash(settings.TempDir)
	currentFiles, _ = ioutil.ReadDir(dir)
	for _, f := range currentFiles {
		file := File{Name: f.Name(), Size: f.Size(), Folder: "Temp", Dir: f.IsDir()}
		files = append(files, file)
	}

	return files
}

//DeleteFiles deletes the given files
func (fs *FilesService) DeleteFiles(files []File) error {
	settings := fs.GetSettings()
	for _, f := range files {
		var baseDir string
		if f.Folder == "Downloads" {
			baseDir = settings.DownloadDir
		} else {
			baseDir = settings.TempDir
		}

		absoluteFile := filepath.FromSlash(baseDir + "/" + f.Name)

		if f.Dir {
			err := os.RemoveAll(absoluteFile)
			if err != nil {
				return err
			}
		} else {
			err := os.Remove(absoluteFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//MoveFilesToMovies moves the given files to the MoviesFolder
func (fs *FilesService) MoveFilesToMovies(files []File) error {
	settings := fs.GetSettings()
	for _, f := range files {
		var baseDir string
		if f.Folder == "Downloads" {
			baseDir = settings.DownloadDir
		} else {
			baseDir = settings.TempDir
		}

		absoluteFile := filepath.FromSlash(baseDir + "/" + f.Name)
		destination := filepath.FromSlash(settings.MoviesFolder + "/" + f.Name)

		err := os.Rename(absoluteFile, destination)
		if err != nil {
			return err
		}
	}

	return nil
}
