package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

var ignoredFiles = []string{
	".DS_Store",
	".gitignore",
	".idea",
}

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
	err := dirTreeRec(out, path, printFiles, 0)
	return err
}

func dirTreeRec(out io.Writer, path string, printFiles bool, lvl int) error {
	//file, err := os.Open(path) // For read access.
	//if err != nil {
	//	log.Fatal(err)
	//}
	//data := make([]byte, 100)
	//count, err := file.Read(data)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("read %d bytes: %q\n", count, data[:count])

	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name()[0] < files[j].Name()[0] })
	l := len(files)

	var b bool
	for n, file := range files {
		if file.IsDir() {
			if n+1 == l {
				b = true
			}
			printer(file.Name(), lvl, b)
			//fmt.Println("dir: ", file.Name())
			err := dirTreeRec(out, path+"/"+file.Name(), printFiles, lvl+1)
			if err != nil {
				return err
			}
		}
		if printFiles {
			info, _ := file.Info()
			if info.Size() > 0 {
				fmt.Printf("File: %s (%db)\n", file.Name(), info.Size())
				continue
			}
			fmt.Println("File: ", file.Name(), "(empty)")
		}
	}

	return nil
}

func printer(name string, lvl int, b bool) {
	var branch string
	i := lvl
	if i > 0 {
		branch += "│\t"
		i--
	}
	for i > 0 {
		branch += "│\t"
		i--
	}
	if !b {
		branch += `├───` + name
	} else {
		branch += `└───` + name
	}
	fmt.Printf("%s\n", branch)
}
