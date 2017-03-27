package main

import (
    /*
//    "errors"
	"flag"
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"github.com/dgiagio/getpass"
    
    
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
*/
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/go-ini/ini"
  "fmt"
  "bufio"
  "os/user"
  "log"
  "os"
	"os/exec"
  "path/filepath"
  "strings"
  "regexp"
  "sort"
  "flag"
)




func checkErr(err error) {
	if err != nil {
		//panic(err)
		log.Fatal(err)
	}
}

func myinput(content string) (result string) {
	inputReader := bufio.NewReader(os.Stdin)
	for result == "" {
		fmt.Println("\n>>>", content)
		input, err := inputReader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)
		if err != nil {
			fmt.Println("There ware errors reading, input again\n")
			continue
		}
		if input == "exit" {
			os.Exit(1)
		}
		result = input
		//fmt.Printf("Your input is %s", input)
	}
	return result
}

func myinput_compatibleEmpty(content string) (result string) {
	inputReader := bufio.NewReader(os.Stdin)
	//var err error = errors.New("this is a new error")

	for {
		fmt.Println("\n>>>", content)
		input, err := inputReader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)
		if err != nil {
			fmt.Println("There ware errors reading, input again\n")
			continue
		}
		if input == "exit" {
			os.Exit(1)
		}
		result = input
		break
		//fmt.Printf("Your input is %s", input)
	}
	return result
}

type User struct {
    Id_        bson.ObjectId  `bson:"_id"`
    Usr        string         `bson:"usr"`
    Name       string         `bson:"name"`
    Groups     string         `bson:"groups"`
    Email      string         `bson:"email`
    Pipelines  []string       `bson:"pipelines"`
}

type Bug struct {
	Id_        		   bson.ObjectId    `bson:"_id"`
	Submiter           string       	`bson:"Submiter"`
	Pipeline           string        	`bson:"Pipeline"`
	Type               string			`bson:"Type"`
	Path               string           `bson:"Path"`
	Detail             string           `bson:"Detail`
	IU                 string           `bson:"IU"`
	Debuger            string           `bson:"Debuger"`
	Estimated_time     string           `bson:"Estimated_time"`
	Confirmation_time  string           `bson:"Confirmation_time"`
}




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


func addUsr(usr_collection *mgo.Collection ,usr string,PBS_file string){
	PBS_cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, PBS_file)
	if err != nil {
		panic(err)
	}
	allgroups := PBS_cfg.Section("groups").KeysHash()
	var (
		name       string
		groups     string
		email      string
		pipelines   []string
	)
	name = myinput("请输入你的姓名:")
	groupss := sortSprintMap(allgroups)
	i := myinput(fmt.Sprintf("请输入你的组别:\n%s\n如果你的组不在列表里请直接输入组名",groupss))
	groups, ok := allgroups[i]
	if !ok {
		groups = i
		fmt.Println(fmt.Sprintf("你输入了一个新的组名:%s",i))
	}
    
	email_suffix := PBS_cfg.Section("email").Key("suffix").String()
	for email == "" {
		email = myinput("请输入你的邮箱:")
		match, _ := regexp.MatchString(email_suffix, email)
		if match == false {
			email = ""
		}
	}

	sections := PBS_cfg.SectionStrings()[3:]
	pipliness := make(map[string]string)
	var piplineStr string
	ii := -1
	for _, v := range sections{
 		piplineStr = fmt.Sprintf("%s\n%v\n", piplineStr,v)
 		allpip := PBS_cfg.Section(v).KeyStrings()
 		for _,v2 := range allpip{
			ii += 1
			pipliness[fmt.Sprintf("%d",ii)] = v2
			piplineStr = fmt.Sprintf("%s\t%d\t%v\n", piplineStr,ii, v2)
		}
	}
    
	is := myinput_compatibleEmpty(fmt.Sprintf("研发人员 请输入你负责的流程编号(以空格隔开):\n%s\n如果你负责的流程不在列表中请联系管理员更改更改PBS.ini\n分析人员请直接按回车键,",piplineStr))
    
	if is != ""{
		slice_is := strings.Split(is," ")
		for _, pip_key := range(slice_is){
			if pip,ok := pipliness[pip_key]; ok{
				pipelines = append(pipelines,pip)
			} else {
				panic("流程编号err")
			}
		}
	} else {
		fmt.Println("你没有选择任何流程，默认为分析人员")
	}
    
	err = usr_collection.Insert(&User{ Id_: bson.NewObjectId(),
                                        Usr: usr,
                                        Name:name,
                                        Groups:groups,
                                        Email:email,
                                        Pipelines:pipelines,})
	if err != nil {
		panic(err)
	}
	os.Exit(1)
}

func updateUsr(usr_collection *mgo.Collection, usr User, PBS_conf string){
	PBS_cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, PBS_conf)
	if err != nil {
		panic(err)
	}

	var (
		name       string
		groups     string
		email      string
		Piplines   []string
		ok		   bool
	)
	name = myinput_compatibleEmpty(fmt.Sprintf("请输入你的姓名:%s", usr.Name))
	if name == "" {
		name = usr.Name
	}
	allgroups := PBS_cfg.Section("groups").KeysHash()
	groupss := sortSprintMap(allgroups)
	i := myinput_compatibleEmpty(fmt.Sprintf("请输入你的组别:\n%s\n如果你的组不在列表里请直接输入组名:%s",groupss,usr.Groups))
	if i == ""{
		groups = usr.Groups
	} else {
		groups, ok = allgroups[i]
		if !ok {
			groups = i
			fmt.Println(fmt.Sprintf("你输入了一个新的组名:%s",i))			
		}
	}

	email_suffix := PBS_cfg.Section("email").Key("suffix").String()
	for email == "" {
		email = myinput_compatibleEmpty(fmt.Sprintf("请输入你的邮箱:%s",usr.Email))
		if email == ""{
			email = usr.Email
		} else{
			match, _ := regexp.MatchString(email_suffix, email)
			if match == false {
				email = ""
			}
		}
	}

	sections := PBS_cfg.SectionStrings()[3:]
	pipliness := make(map[string]string)
	var piplineStr string
	ii := -1
	for _, v := range sections{
 		piplineStr = fmt.Sprintf("%s\n%v\n", piplineStr,v)
		allpip := PBS_cfg.Section(v).KeyStrings()
		for _,v2 := range allpip{
			ii += 1
			pipliness[fmt.Sprintf("%d",ii)] = v2
			piplineStr = fmt.Sprintf("%s\t%d\t%v\n", piplineStr,ii, v2)
		}
	}
    
	is := myinput_compatibleEmpty(fmt.Sprintf("研发人员请输入你负责的流程编号(以空格隔开):\n%s\n如果你负责的流程不在列表中请联系管理员更改更改PBS.ini\n分析人员请直接按回车键,",piplineStr))
    
	if is != ""{
		slice_is := strings.Split(is," ")
		for _, pip_key := range(slice_is){
			if pip,ok := pipliness[pip_key]; ok{
				Piplines = append(Piplines,pip)
			} else {
				panic("流程编号err")
			}
		}
	} else {
		fmt.Println("你没有选择任何流程，默认为分析人员")
	}
	colQuerier := bson.M{"usr": usr.Usr}
	change := bson.M{"$set": bson.M{"name": name, "groups": groups, "email":email, "pipelines":Piplines}}
	err = usr_collection.Update(colQuerier, change)
	if err != nil {
		panic(err)
	}
}

func addBug(bug_collection *mgo.Collection, usr string, PBS_file string){
	PBS_cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, PBS_file)
	if err != nil {
		panic(err)
	}
	allgroups := PBS_cfg.Section("groups").KeysHash()
	var (
		name       string
		groups     string
		email      string
		pipelines   []string
	)
	name = myinput("请输入你的姓名:")
	groupss := sortSprintMap(allgroups)
	i := myinput(fmt.Sprintf("请输入你的组别:\n%s\n如果你的组不在列表里请直接输入组名",groupss))
	groups, ok := allgroups[i]
	if !ok {
		groups = i
		fmt.Println(fmt.Sprintf("你输入了一个新的组名:%s",i))
	}
    
	email_suffix := PBS_cfg.Section("email").Key("suffix").String()
	for email == "" {
		email = myinput("请输入你的邮箱:")
		match, _ := regexp.MatchString(email_suffix, email)
		if match == false {
			email = ""
		}
	}

	sections := PBS_cfg.SectionStrings()[3:]
	pipliness := make(map[string]string)
	var piplineStr string
	ii := -1
	for _, v := range sections{
 		piplineStr = fmt.Sprintf("%s\n%v\n", piplineStr,v)
 		allpip := PBS_cfg.Section(v).KeyStrings()
 		for _,v2 := range allpip{
			ii += 1
			pipliness[fmt.Sprintf("%d",ii)] = v2
			piplineStr = fmt.Sprintf("%s\t%d\t%v\n", piplineStr,ii, v2)
		}
	}
    
	is := myinput_compatibleEmpty(fmt.Sprintf("研发人员 请输入你负责的流程编号(以空格隔开):\n%s\n如果你负责的流程不在列表中请联系管理员更改更改PBS.ini\n分析人员请直接按回车键,",piplineStr))
    
	if is != ""{
		slice_is := strings.Split(is," ")
		for _, pip_key := range(slice_is){
			if pip,ok := pipliness[pip_key]; ok{
				pipelines = append(pipelines,pip)
			} else {
				panic("流程编号err")
			}
		}
	} else {
		fmt.Println("你没有选择任何流程，默认为分析人员")
	}
    
	err = usr_collection.Insert(&User{ Id_: bson.NewObjectId(),
                                        Usr: usr,
                                        Name:name,
                                        Groups:groups,
                                        Email:email,
                                        Pipelines:pipelines,})
	if err != nil {
		panic(err)
	}
	os.Exit(1)
}



func usage() {
	fmt.Printf("\nProgram: PipelineBugSubmit (Tools for pipeline bug submit)\nVersion: 0.0.1-20170323\n注意输入不能有中文标点符号!\nXshell删除字符时请按住Ctrl键\n\nUsage:bugSubmitct <command> [options]\n\n")
	fmt.Printf("Command:\n")

	fmt.Printf("   -editusr               ReEdit your usr information\n")
}

func main() {
	editusr := flag.Bool("editusr", false, "edit your usr information")

 	flag.Parse()
	usr, _ := user.Current()
    
	file, _ := exec.LookPath(os.Args[0])
	filepaths, _ := filepath.Abs(file)
	bin := filepath.Dir(filepaths)    
	PBS_conf := bin + "/PBS.ini"  
	session, err := mgo.Dial("127.0.0.1:27017")
	defer session.Close()
	db := session.DB("PBS")
	usr_collection := db.C("users")
    
	people := User{}
	err = usr_collection.Find(bson.M{"usr": usr.Username}).One(&people)

	if err != nil{
		addUsr(usr_collection, usr.Username, PBS_conf)
	}
	
	usage()
	
	if *editusr == true{
		updateUsr(usr_collection, people, PBS_conf)
	}
}	
	
	
	
	
