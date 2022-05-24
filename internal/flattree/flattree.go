package flattree

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"obsidian-lens/internal/opts"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Leaf struct {
	path       string
	parentPath string
	isDir      bool
	info       *fs.FileInfo
}
type FlatFileTree struct {
	root    string
	leaves  []Leaf
	options *opts.Opts
}

func NewFlatFileTree(root string, opts *opts.Opts) *FlatFileTree {
	return &FlatFileTree{
		root:    root,
		leaves:  []Leaf{},
		options: opts,
	}
}

func (ft *FlatFileTree) Populate() {
	ft.recurse(ft.root)
}

func (ft *FlatFileTree) add(leaf Leaf) {
	ft.leaves = append(ft.leaves, leaf)
}

func (ft *FlatFileTree) Print() {
	for _, i := range ft.leaves {
		fmt.Println(i.path)
	}
}

// Recursive method to
func (ft *FlatFileTree) recurse(parentPath string) {
	items := ft.readDirFoldersFirst(parentPath)
	for _, f := range items {

		childpath := filepath.Join(parentPath, f.Name())
		if f.IsDir() {
			leaf := Leaf{
				path:       childpath,
				parentPath: parentPath,
				isDir:      true,
				info:       &f,
			}
			ft.add(leaf)
			ft.recurse(childpath)
		} else {
			leaf := Leaf{
				path:       childpath,
				parentPath: parentPath,
				isDir:      false,
				info:       &f,
			}
			ft.add(leaf)
		}
	}
}

func (ft *FlatFileTree) readDirFoldersFirst(path string) []fs.FileInfo {
	items, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("unable to read directory %s: %v", path, err)
		os.Exit(1)
	}
	sortFilesAscInsensitive(items)

	var dirs []fs.FileInfo
	var files []fs.FileInfo
	for _, f := range items {
		if f.IsDir() {
			if ft.options.IsDirAllowed(f) {
				dirs = append(dirs, f)
			}
		} else {
			if ft.options.IsFileAllowed(f) {
				files = append(files, f)
			}
		}
	}
	return append(dirs, files...)
}

func sortFilesAscInsensitive(files []os.FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name())
	})
}
