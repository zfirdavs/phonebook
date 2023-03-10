package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	min = 0
	max = 26
)

type Entry struct {
	Name       string
	Surname    string
	Tel        string
	LastAccess string
}

func (e Entry) String() string {
	return fmt.Sprintf("Name: %s, Surname: %s, Phone: %s", e.Name, e.Surname, e.Tel)
}

var (
	data    = []Entry{}
	CSVFILE = "./csv.data"
	index   map[string]int
)

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

func readCSVFile(filepath string) error {
	_, err := os.Stat(filepath)
	if err != nil {
		return err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	// CSV file read all at once
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}

	for _, line := range lines {
		temp := Entry{
			Name:       line[0],
			Surname:    line[1],
			Tel:        line[2],
			LastAccess: line[3],
		}

		data = append(data, temp)
	}

	return nil
}

func saveCSVFile(filepath string) error {
	csvfile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer csvfile.Close()

	csvwriter := csv.NewWriter(csvfile)
	for _, row := range data {
		temp := []string{row.Name, row.Surname, row.Tel, row.LastAccess}
		_ = csvwriter.Write(temp)
	}
	csvwriter.Flush()
	return nil
}

func createIndex() error {
	index = make(map[string]int)
	for i, k := range data {
		key := k.Tel
		index[key] = i
	}
	return nil
}

func matchTel(s string) bool {
	t := []byte(s)
	return regexp.MustCompile(`\d+$`).Match(t)
}

func initS(N, S, T string) *Entry {
	if T == "" || S == "" {
		return nil
	}

	LastAccess := strconv.FormatInt(time.Now().Unix(), 10)
	return &Entry{
		Name:       N,
		Surname:    S,
		Tel:        T,
		LastAccess: LastAccess,
	}
}

func insert(e *Entry) error {
	// If tel is exist, return error
	_, ok := index[e.Tel]
	if !ok {
		return fmt.Errorf("%s already exists", e.Tel)
	}
	data = append(data, *e)

	// update the index
	_ = createIndex()

	if err := saveCSVFile(CSVFILE); err != nil {
		return err
	}
	return nil
}

func main() {
	args := os.Args
	if len(args) == 1 {
		exe := path.Base(args[0])
		fmt.Printf("Usage: %s insert|delete|search|list <arguments>\n", exe)
		return
	}

	_, err := os.Stat(CSVFILE)
	if err != nil {
		fmt.Println("Creating", CSVFILE)
		f, err := os.Create(CSVFILE)
		if err != nil {
			f.Close()
			fmt.Println(err)
			return
		}
		f.Close()
	}

	fileInfo, err := os.Stat(CSVFILE)
	mode := fileInfo.Mode()
	if !mode.IsRegular() {
		fmt.Println(CSVFILE, "not a regular file!")
		return
	}

	err = readCSVFile(CSVFILE)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = createIndex()
	if err != nil {
		fmt.Println("Cannot create index")
		return
	}

	// Differentiate between the commands
	switch args[1] {
	case "insert":
		if len(args) != 5 {
			fmt.Println("Usage: insert Name Surname Telephone")
			return
		}

		t := strings.ReplaceAll(args[4], "-", "")
		if !matchTel(t) {
			fmt.Println("Not a valid telephone number:", t)
			return
		}

		entry := initS(args[2], args[3], t)
		if entry != nil {
			if err := insert(entry); err != nil {
				fmt.Println(err)
				return
			}
		}

	case "search":
		if len(args) != 3 {
			fmt.Println("Usage: search Tel number")
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
