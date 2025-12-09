// package main

// import (
// 	"fmt"
// 	"strings"
// )

// type Student struct {
// 	Name string
// 	Class string
// 	Member []string
// 	Scores map[string]int
// }

// func (p Student) info()  {
// 	fmt.Println("Hi, this is ", p.Name + " class ", p.Class + "\n")
// }

// func (p Student) ListClassStrings() {
// 	var builder strings.Builder
// 	var result string = ""

// 	for _, value := range p.Member {
// 		builder.WriteString(fmt.Sprintf("%v,", value))
// 	}

// 	result = builder.String()

// 	fmt.Printf("Result : %s\n", result)
// }

// func (p *Student) ListClassArray() []string {
//     var result []string

// 	result = append(result, p.Member...)
//     return result
// }

// func (p *Student) addStudent(name string)  {
// 	p.Member = append(p.Member, name)
// }

// func (p *Student) delStudent(name string) []string {
// 	var result []string

// 	for _, value := range p.Member {
// 		if value != name {
// 			result = append(result, value)
// 		}
// 	}

// 	p.Member = result
// 	return  result
// }

// func main(){
// 	point := map[string]int{
//         "math":    7,
//         "english": 8,
//     }

// 	//p := Student{Name: "David Scotch", Class: "12A1", Member: []string, Scores: point}
// 	p        := Student{}
// 	p.Name    = "David Scotch"
// 	p.Class   = "12A1"
// 	p.Member  = []string{"111"}
// 	p.Scores  = point
	
// 	p.info()
// 	p.addStudent("222")
// 	p.addStudent("4444")
// 	p.delStudent("111")
	
// 	fmt.Println(p.ListClassArray())
// 	fmt.Println(p.Member)
// }