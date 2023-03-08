package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
)

const (
	min = 0
	max = 26
)

type Entry struct {
	Name    string
	Surname string
	Tel     string
}

func (e Entry) String() string {
	return fmt.Sprintf("Name: %s, Surname: %s, Phone: %s", e.Name, e.Surname, e.Tel)
}

var data = []Entry{}

func search(key string) *Entry {
	for i, v := range data {
		if v.Tel == key {
			return &data[i]
		}
	}
	return nil
}

func list() {
	for _, v := range data {
		fmt.Println(v)
	}
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func getString(l int64) string {
	startChar := "A"
	temp := ""
	var i int64 = 1
	for {
		myRand := random(min, max)
		newChar := string(startChar[0] + byte(myRand))
		temp = temp + newChar
		if i == l {
			break
		}
		i++
	}
	return temp
}

func main() {
	args := os.Args
	if len(args) == 1 {
		exe := path.Base(args[0])
		fmt.Printf("Usage: %s search|list <arguments>\n", exe)
		return
	}

	data = append(data, Entry{"Mihalis", "Tsoukalos", "2109416471"})
	data = append(data, Entry{"Mary", "Doe", "2109416871"})
	data = append(data, Entry{"John", "Black", "2109416123"})

	// Differentia{te between the commands
	switch args[1] {
	case "search":
		if len(args) != 3 {
			fmt.Println("Usage: search Surname")
			return
		}

		result := search(args[2])
		if result == nil {
			fmt.Println("Entry not found:", args[2])
			return
		}
		fmt.Println(*result)
	case "list":
		list()
	default:
		fmt.Println("Not a valid option")
	}
}
