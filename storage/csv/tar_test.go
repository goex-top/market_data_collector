package csv

import (
	"log"
	"testing"
)

func TestGetAllFileName(t *testing.T) {
	csvs := GetSrcFileName("../output/csv/", []string{"csv"})
	for _, v := range csvs {
		log.Println(v)
		//t.Log(os.Remove("../output/csv/" + v))
	}
}

func TestCompressAllCsv(t *testing.T) {
	CompressFile("../output/csv", []string{"2006-01-02"}, "../output/tar/aaa.tar.gz")
}
