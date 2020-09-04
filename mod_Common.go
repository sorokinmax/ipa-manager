package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func isError(err error) bool {
	if err != nil {
		log.Println(err.Error())
	}
	return (err != nil)
}

func containsIpas(arr []Ipa, element string) bool {
	for _, a := range arr {
		if a.CFBundleVersion == element {
			return true
		}
	}
	return false
}

// MakeNetworkPath make network path from local path
func MakeNetworkPath(server string, directory string) string {
	path := strings.Replace("\\\\"+server+"\\"+directory, ":", "$", -1)
	return path
}

// Convert - convert executable output
func Convert(i int, s []byte) string {
	var reader *transform.Reader
	switch i {
	case 1:
		reader = transform.NewReader(bytes.NewReader(s), charmap.ISO8859_1.NewDecoder())
	case 2:
		reader = transform.NewReader(bytes.NewReader(s), charmap.ISO8859_2.NewDecoder())
	case 3:
		reader = transform.NewReader(bytes.NewReader(s), charmap.ISO8859_3.NewDecoder())
	case 4:
		reader = transform.NewReader(bytes.NewReader(s), charmap.ISO8859_4.NewDecoder())
	case 5:
		reader = transform.NewReader(bytes.NewReader(s), charmap.ISO8859_5.NewDecoder())
	case 6:
		reader = transform.NewReader(bytes.NewReader(s), charmap.ISO8859_6.NewDecoder())
	case 7:
		reader = transform.NewReader(bytes.NewReader(s), charmap.ISO8859_7.NewDecoder())
	case 8:
		reader = transform.NewReader(bytes.NewReader(s), charmap.ISO8859_8.NewDecoder())
	case 9:
		reader = transform.NewReader(bytes.NewReader(s), charmap.ISO8859_9.NewDecoder())
	case 10:
		reader = transform.NewReader(bytes.NewReader(s), charmap.ISO8859_10.NewDecoder())
	case 11:
		reader = transform.NewReader(bytes.NewReader(s), charmap.ISO8859_13.NewDecoder())
	case 12:
		reader = transform.NewReader(bytes.NewReader(s), charmap.ISO8859_14.NewDecoder())
	case 13:
		reader = transform.NewReader(bytes.NewReader(s), charmap.ISO8859_15.NewDecoder())
	case 14:
		reader = transform.NewReader(bytes.NewReader(s), charmap.ISO8859_16.NewDecoder())
	case 15:
		reader = transform.NewReader(bytes.NewReader(s), charmap.KOI8R.NewDecoder())
	case 16:
		reader = transform.NewReader(bytes.NewReader(s), charmap.KOI8U.NewDecoder())
	case 17:
		reader = transform.NewReader(bytes.NewReader(s), charmap.Macintosh.NewDecoder())
	case 18:
		reader = transform.NewReader(bytes.NewReader(s), charmap.MacintoshCyrillic.NewDecoder())
	case 19:
		reader = transform.NewReader(bytes.NewReader(s), charmap.Windows874.NewDecoder())
	case 20:
		reader = transform.NewReader(bytes.NewReader(s), charmap.Windows1250.NewDecoder())
	case 21:
		reader = transform.NewReader(bytes.NewReader(s), charmap.Windows1251.NewDecoder())
	case 22:
		reader = transform.NewReader(bytes.NewReader(s), charmap.Windows1252.NewDecoder())
	case 23:
		reader = transform.NewReader(bytes.NewReader(s), charmap.Windows1253.NewDecoder())
	case 24:
		reader = transform.NewReader(bytes.NewReader(s), charmap.Windows1254.NewDecoder())
	case 25:
		reader = transform.NewReader(bytes.NewReader(s), charmap.Windows1255.NewDecoder())
	case 26:
		reader = transform.NewReader(bytes.NewReader(s), charmap.Windows1256.NewDecoder())
	case 27:
		reader = transform.NewReader(bytes.NewReader(s), charmap.Windows1257.NewDecoder())
	case 28:
		reader = transform.NewReader(bytes.NewReader(s), charmap.Windows1258.NewDecoder())
	case 29:
		reader = transform.NewReader(bytes.NewReader(s), charmap.XUserDefined.NewDecoder())
	case 30:
		reader = transform.NewReader(bytes.NewReader(s), charmap.CodePage037.NewDecoder())
	case 31:
		reader = transform.NewReader(bytes.NewReader(s), charmap.CodePage437.NewDecoder())
	case 32:
		reader = transform.NewReader(bytes.NewReader(s), charmap.CodePage850.NewDecoder())
	case 33:
		reader = transform.NewReader(bytes.NewReader(s), charmap.CodePage852.NewDecoder())
	case 34:
		reader = transform.NewReader(bytes.NewReader(s), charmap.CodePage855.NewDecoder())
	case 35:
		reader = transform.NewReader(bytes.NewReader(s), charmap.CodePage858.NewDecoder())
	case 36:
		reader = transform.NewReader(bytes.NewReader(s), charmap.CodePage860.NewDecoder())
	case 37:
		reader = transform.NewReader(bytes.NewReader(s), charmap.CodePage862.NewDecoder())
	case 38:
		reader = transform.NewReader(bytes.NewReader(s), charmap.CodePage863.NewDecoder())
	case 39:
		reader = transform.NewReader(bytes.NewReader(s), charmap.CodePage865.NewDecoder())
	case 40:
		reader = transform.NewReader(bytes.NewReader(s), charmap.CodePage865.NewDecoder())
	case 41:
		reader = transform.NewReader(bytes.NewReader(s), charmap.CodePage866.NewDecoder())
	case 42:
		reader = transform.NewReader(bytes.NewReader(s), charmap.CodePage1047.NewDecoder())
	case 43:
		reader = transform.NewReader(bytes.NewReader(s), charmap.CodePage1140.NewDecoder())
	case 44:
		reader = transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	case 45:
		reader = transform.NewReader(bytes.NewReader(s), simplifiedchinese.HZGB2312.NewEncoder())
	case 46:
		reader = transform.NewReader(bytes.NewReader(s), simplifiedchinese.GB18030.NewEncoder())
	}

	d, _ := ioutil.ReadAll(reader)

	return string(d)
}
