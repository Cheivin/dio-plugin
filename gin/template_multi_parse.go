package web

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
)

func ParseFSDir(t *template.Template, fsys fs.FS, root string, patterns ...string) (*template.Template, error) {
	var filenames []string
	for _, pattern := range patterns {
		var list []string
		err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				if match, err := filepath.Match(pattern, d.Name()); err != nil {
					return err
				} else if match {
					list = append(list, path)
				}
			}
			return err
		})
		if err != nil {
			return nil, err
		}
		if len(list) == 0 {
			return nil, fmt.Errorf("template: pattern matches no files: %#q", pattern)
		}
		filenames = append(filenames, list...)
	}
	return parseFilesTrimRoot(t, root, readFileFSTrimRoot(fsys), filenames...)
}

func readFileFSTrimRoot(fsys fs.FS) func(string, string) (string, []byte, error) {
	return func(root string, file string) (name string, b []byte, err error) {
		name, err = filepath.Rel(root, file)
		fmt.Println(name)
		b, err = fs.ReadFile(fsys, file)
		return
	}
}

func ParseDir(t *template.Template, root string, filePattern ...string) (*template.Template, error) {
	var filenames []string
	if err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
		outer:
			for _, pattern := range filePattern {
				if match, err := filepath.Match(pattern, info.Name()); err != nil {
					return err
				} else if match {
					filenames = append(filenames, path)
					break outer
				}
			}

		}
		return err
	}); err != nil {
		return nil, err
	}

	if len(filenames) == 0 {
		return nil, fmt.Errorf("html/template: pattern matches no files: %#q", filePattern)
	}
	return parseFilesTrimRoot(t, root, readFileOSTrimRoot, filenames...)
}

func parseFilesTrimRoot(t *template.Template, root string, readFile func(string, string) (string, []byte, error), filenames ...string) (*template.Template, error) {
	if len(filenames) == 0 {
		// Not really a problem, but be consistent.
		return nil, fmt.Errorf("html/template: no files named in call to ParseFiles")
	}
	for _, filename := range filenames {
		name, b, err := readFile(root, filename)
		if err != nil {
			return nil, err
		}
		s := string(b)
		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func readFileOSTrimRoot(root string, file string) (name string, b []byte, err error) {
	name, err = filepath.Rel(root, file)
	b, err = os.ReadFile(file)
	return
}
