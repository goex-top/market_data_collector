package storage

import (
	"fmt"
	"os"
	"testing"
)

func TestCompressAllCsv(t *testing.T) {
	csvs := GetAllFileName("../csv", "csv")
	for _, v := range csvs {
		fmt.Println(v)
		t.Log(os.Remove("../csv/" + v))
	}

}
