package main

import (
    "fmt"
    "github.com/go-ini/ini"
    "flag"
    "os"
	"os/exec"
	"path/filepath"
    "sort"
    "strings"
)

func sortSprintMap(mapin map[string]string) string{
    var sortMapOut string
    keys := make([]string, len(mapin))
    i := 0
    for k, _ := range mapin {
        keys[i] = k
        i++
    }
    sort.Strings(keys)
    for _, k := range keys {
        sortMapOut = fmt.Sprintf("%s\t%v\t%v\n", sortMapOut,k, mapin[k])
    }
    return sortMapOut
}


func main() {
    
	flag.Parse()
	file, _ := exec.LookPath(os.Args[0])
	filepaths, _ := filepath.Abs(file)
	bin := filepath.Dir(filepaths)

	Piplines_file := bin + "/PBS.ini"
    
//    cfg, err := ini.InsensitiveLoad(Piplines_file)
    cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, Piplines_file)
  
    if err != nil {
		panic(err)
	}    
    sections := cfg.SectionStrings()[2:]  
    piplines := make(map[string]string)
    var piplineStr string
    i := -1
    for _, v := range sections{
        piplineStr = fmt.Sprintf("%s\n%v\n", piplineStr,v)
        allpip := cfg.Section(v).KeyStrings()
        for _,v2 := range allpip{
            i += 1
            piplines[fmt.Sprintf("%d",i)] = v2
            piplineStr = fmt.Sprintf("%s\t%d\t%v\n", piplineStr,i, v2)
        }
    }
    
 //   groupss := sortSprintMap(allpip)
  

    fmt.Println(piplines)
    fmt.Println(piplineStr)
    
    var a string = "1,2 4 6"
    b := strings.Split(a, " ")
    fmt.Println(b)
  
    aa := make(map[string]string)
    aa["s"] = "wwew"
    aa["w"] = "ttt"
    fmt.Println(aa)
    v, err := aa["s"]
    fmt.Println(v,err)
    
    
    
    
  
}
