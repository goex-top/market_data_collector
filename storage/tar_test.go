package storage

import (
	"fmt"
	"testing"
)

func TestGetAllFileName(t *testing.T) {
	csvs := GetAllFileName("../csv/", "csv")
	for _, v := range csvs {
		fmt.Println(v)
		//t.Log(os.Remove("../csv/" + v))
	}
}

func TestCompressAllCsv(t *testing.T) {
	CompressAllCsv("../csv", "../tar/aaa.tar.gz")
}
