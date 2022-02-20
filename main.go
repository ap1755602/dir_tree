package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"sort"
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
	var lvl map[int]bool //map[int]bool{0: true}
	return dirTreeRec(out, path, printFiles, &lvl)

}

func dirTreeRec(out io.Writer, path string, printFiles bool, lvl *map[int]bool) error {
	files, err := os.ReadDir(path)

	if err != nil {
		log.Fatal(err)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name()[0] < files[j].Name()[0] })
	if !printFiles {
		deleteFilesFromSl(&files)
	}
	l := len(files)

	//var b bool
	if l == 0 {
		if *lvl != nil {
			delete(*lvl, len(*lvl)-1)
		}
		return nil
	}
	for n, file := range files {
		if file.IsDir() {
			if *lvl == nil {
				*lvl = map[int]bool{0: n+1 == l}
			} else {
				(*lvl)[len(*lvl)] = n+1 == l
			}
			printDir(file.Name(), lvl)
			if n+1 == l {
				delete(*lvl, len(*lvl)-1)
				return nil
			}
			err := dirTreeRec(out, path+"/"+file.Name(), printFiles, lvl)
			if err != nil {
				return err
			}
		} else {
			if n+1 == l {
				delete(*lvl, len(*lvl)-1)
			}
			//	//if printFiles
			//	//info, _ := file.Info()
			//	if n+1 == l {
			//		return nil
			//		//delete(*lvl, len(*lvl))
			//	}
			//if info.Size() > 0 {
			//fmt.Printf("File: %s (%db)\n", file.Name(), info.Size())
			//continue
			//}
			//fmt.Println("File: ", file.Name(), "(empty)")
		}

	}

	return nil
}

func printDir(name string, lvl *map[int]bool) {
	var branch string
	//if len(*lvl) == 1 && (*lvl)[0] == false {
	//	branch += `└───` + name
	//	return
	//}
	for i, b := range *lvl {
		//fmt.Println("Yo")

		if i+1 == len(*lvl) {
			if b == false {
				branch += `├───` + name
				fmt.Printf("%s\n", branch)
				return
			} else {
				branch += `└───` + name
				fmt.Printf("%s\n", branch)
				return
			}
		}
		if b == false {
			branch += `│`
		}
		branch += "\t"
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
