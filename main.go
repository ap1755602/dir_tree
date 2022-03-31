package main

import (
	"io"
	"io/fs"
	"log"
	"os"
	"sort"
	"strconv"
)

//var ignoredFiles = []string{
//	".DS_Store",
//	".gitignore",
//	".idea",
//}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	var lvl []bool //map[int]bool{0: true}
	var s string
	err := dirTreeRec(out, path, printFiles, &lvl, &s, 0)
	if err != nil {
		return err
	}
	_, _ = io.WriteString(out, s)
	return nil

}

func dirTreeRec(out io.Writer, path string, printFiles bool, lvl *[]bool, s *string, num int) error {
	file, err := os.Open(path)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	files, err := file.ReadDir(-1)

	if err != nil {
		log.Fatal(err)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name()[0] < files[j].Name()[0] })
	if !printFiles {
		deleteFilesFromSl(&files)
	}
	l := len(files)
	if *lvl == nil {
		*lvl = make([]bool, 1)
		(*lvl)[0] = num+1 == l
	} else {
		*lvl = append(*lvl, num+1 == l)
	}
	if l == 0 {
		if *lvl != nil && len(*lvl) != 0 {
			*lvl = (*lvl)[0:(len(*lvl) - 1)]
		}
		return nil
	}
	for n, file := range files {
		if len(*lvl) > 0 {
			(*lvl)[len(*lvl)-1] = n+1 == l
		}
		if file.IsDir() {
			printDir(file.Name(), *lvl, s, nil)

			err := dirTreeRec(out, path+"/"+file.Name(), printFiles, lvl, s, n)
			if err != nil {
				return err
			}
		} else {
			info, _ := file.Info()
			printDir(file.Name(), *lvl, s, info)
			if n+1 == l && len(*lvl) != 0 {
				*lvl = (*lvl)[0:(len(*lvl) - 1)]
				return nil
			}
		}
	}
	if len(*lvl) != 0 {
		*lvl = (*lvl)[:(len(*lvl) - 1)]
	}
	return nil
}

func printDir(fileName string, lvl []bool, s *string, info fs.FileInfo) {
	name := fileName

	if info != nil {
		if info.Size() > 0 {
			name += " (" + strconv.Itoa(int(info.Size())) + "b)"
		} else {
			name += " (empty)"
		}
	}

	for i, b := range lvl {
		if i+1 == len(lvl) {
			if b == false {
				*s += `├───` + name + "\n"
				return
			} else {
				*s += `└───` + name + "\n"
				return
			}
		}
		if b == false {
			*s += `│`
		}
		*s += "\t"
	}
}

func deleteFilesFromSl(files *[]fs.DirEntry) {
	for i := 0; i < len(*files); i++ {
		if (*files)[i].IsDir() {
			continue
		}
		copy((*files)[i:], (*files)[i+1:])
		(*files)[len(*files)-1] = nil
		*files = (*files)[:len(*files)-1]
		i--
	}
}
