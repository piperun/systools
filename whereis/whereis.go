package whereis

import (
	"os"
	"strings"
)


type Paths struct {
	Pathlist []string
	MANlist []string
}

type File struct {
	Name string
	Path string
}

func (p *Paths) SetPATHlist() {
	var Pathvar = strings.SplitAfter(os.Getenv("PATH"), ":")
	p.Pathlist = make([]string, len(Pathvar))
	for i := 0; i < len(Pathvar); i += 1 {
		p.Pathlist[i] = strings.Replace(Pathvar[i], ":", "", -1)
	}
}

func Add(dirlist []string) {
}


func findin() {

}


func Lookup(dirlist []string, pattern string) (File, error){
	var(
		fi os.FileInfo
		f File
		err error
	)
	for i := 0; i < len(dirlist); i += 1 {
		fi, err = os.Stat(dirlist[i] + "/" + pattern)
		if !os.IsNotExist(err) {
			f.Name = fi.Name()
			f.Path = dirlist[i] + "/"
			return f, nil
		}
	}
	return f, err
}

func main() {

}
