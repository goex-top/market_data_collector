package storage

import (
	"log"
	"testing"
)

func TestGetAllFileName(t *testing.T) {
	csvs := GetSrcFileName("../csv/", "csv")
	for _, v := range csvs {
		log.Println(v)
		//t.Log(os.Remove("../csv/" + v))
	}
}

func TestCompressAllCsv(t *testing.T) {
	CompressFile("../csv", []string{"2006-01-02"}, "../tar/aaa.tar.gz")
}
