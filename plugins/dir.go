package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type FileNode struct {
	Name      string      `json:"name"`
	Path      string      `json:"-"`
	FileNodes []*FileNode `json:"-"`
	IsDir     bool        `json:"isDir"`
}

func listCurrentDirectory(dirname string) ([]FileNode, error) {
	fileNodes := []FileNode{}

	if err := isDirExist(dirname); err != nil {
		return fileNodes, err
	}

	files := listFiles(dirname)
	for _, filename := range files {
		fileNode := FileNode{
			Name:  filename,
			IsDir: false,
		}
		fpath := filepath.Join(dirname, filename)
		fio, _ := os.Lstat(fpath)
		if fio.IsDir() {
			fileNode.IsDir = true
		}
		fileNodes = append(fileNodes, fileNode)
	}

	return fileNodes, nil
}

func listFiles(dirname string) []string {
	f, _ := os.Open(dirname)
	names, _ := f.Readdirnames(-1)
	f.Close()

	sort.Strings(names)

	return names
}

func walk(path string, info os.FileInfo, node *FileNode) {
	files := listFiles(path)

	for _, filename := range files {
		fpath := filepath.Join(path, filename)
		fio, _ := os.Lstat(fpath)

		child := FileNode{
			Name:      filename,
			Path:      fpath,
			FileNodes: []*FileNode{},
			IsDir:     fio.IsDir(),
		}

		node.FileNodes = append(node.FileNodes, &child)

		// 奾B弾^~\輾A~M伾N~F漾Z~D弾S伾I~M弾V~G件弾X¯个漾[®弾U﻾L伾H~Y达[伾E¥该漾[®弾U达[纠¾L輾@~R弾R
		if fio.IsDir() {
			walk(fpath, fio, &child)
		}
	}

	return
}

func isDirExist(dir string) error {
	f, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("dir(%v) not exist", dir)
		} else {
			return err
		}
	}
	if !f.IsDir() {
		return fmt.Errorf("%s is exist,but not dir", dir)
	}

	return nil
}

func getDirTree(dir string) ([]*FileNode, error) {
	emptyFileNode := []*FileNode{}
	if err := isDirExist(dir); err != nil {
		return emptyFileNode, err
	}
	root := FileNode{"projects", dir, []*FileNode{}, true}
	fileInfo, _ := os.Lstat(dir)
	walk(dir, fileInfo, &root)

	return root.FileNodes, nil
}

func ensureDirExist(dir string) error {
	f, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(dir, os.FileMode(0755))
		} else {
			return err
		}
	}
	if !f.IsDir() {
		//已存在，但是不是文件夹
		return fmt.Errorf("path %s is exist,but not dir", dir)
	}

	return nil
}
