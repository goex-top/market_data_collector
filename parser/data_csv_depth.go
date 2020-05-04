package parser

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/nntaoli-project/goex"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DepthDataFromCSVeData loads the market data from a csv file.
// It expands the underlying data struct.
type DepthDataFromCSVeData struct {
	FileDir string
}

func NewDepthDataParser(fileDir string) *DepthDataFromCSVeData {
	return &DepthDataFromCSVeData{FileDir: fileDir}
}

// Load single data events into a stream ordered by date (latest first).
func (d *DepthDataFromCSVeData) Load(filename string) ([]goex.Depth, error) {
	// check file location
	if len(d.FileDir) == 0 {
		return nil, errors.New("no directory for data provided: ")
	}

	if len(filename) == 0 {
		return nil, errors.New("no filename for data provided: ")
	}
	if !strings.Contains(filename, "depth") {
		return nil, errors.New("no a standard depth filename")
	}
	if d.FileDir[len(d.FileDir)-1] != '/' {
		d.FileDir += "/"
	}

	df := LoadDataFrameFromCSV(d.FileDir + filename)

	if df.Err != nil {
		return nil, df.Err
	}

	depths := make([]goex.Depth, 0)

	size := df.Nrow()
	lastTs := 0
	for i := 0; i < size; i++ {
		ts, _ := df.Elem(i, 0).Int()
		if ts == lastTs {
			continue
		}
		lastTs = ts

		str1 := df.Elem(i, 1).String()
		asks := make([][]float64, 0)
		err := json.Unmarshal([]byte(str1), &asks)
		if err != nil {
			panic(err)
		}
		str2 := df.Elem(i, 2).String()
		bids := make([][]float64, 0)
		err = json.Unmarshal([]byte(str2), &bids)
		if err != nil {
			panic(err)
		}

		asklist := make(goex.DepthRecords, 0)
		for _, v := range asks {
			asklist = append(asklist, goex.DepthRecord{
				Price:  v[0],
				Amount: v[1],
			})
		}
		bidlist := make(goex.DepthRecords, 0)
		for _, v := range bids {
			bidlist = append(bidlist, goex.DepthRecord{
				Price:  v[0],
				Amount: v[1],
			})
		}

		depths = append(depths, goex.Depth{
			UTime:   time.Unix(0, int64(ts)*int64(time.Millisecond)),
			AskList: asklist,
			BidList: bidlist,
		})
	}
	return depths, nil
}

func LoadDataFrameFromCSV(filename string) *dataframe.DataFrame {
	fmt.Println("load csv file: ", filename)
	f, err := os.Open(filename)
	defer f.Close()

	if err != nil {
		return nil
	}
	df := dataframe.ReadCSV(f)
	return &df
}

// readDepthFromCSVFile opens and reads a csv file line by line
// and returns a slice with a key/value map for each line.
func readDepthFromCSVFile(path string) (lines []map[string]string, err error) {
	log.Printf("Loading from %s.\n", path)
	// open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// create scanner on top of file
	reader := csv.NewReader(file)
	// set delimiter
	reader.Comma = ','
	// read first line for keys and fill in array
	keys, err := reader.Read()

	// read each line and create a map of values combined to the keys
	for line, err := reader.Read(); err == nil; line, err = reader.Read() {
		l := make(map[string]string)
		for i, v := range line {
			l[keys[i]] = v
		}
		// put found line as map into stream holder item
		lines = append(lines, l)
	}

	return lines, nil
}

// fetchFilesFromDir returns a map of all filenames in a directory,
func fetchFilesFromDir(dir string) (m map[string]string, err error) {
	// read filenames from directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return m, err
	}

	// initialise the map
	m = make(map[string]string)

	// read filenames from directory
	for _, file := range files {
		// file is directory
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		extension := filepath.Ext(filename)
		// file is not CSV
		if extension != ".csv" {
			continue
		}

		name := filename[0 : len(filename)-len(extension)]
		m[name] = filename
	}
	return m, nil
}
