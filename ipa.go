package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"howett.net/plist"
)

type Ipa struct {
	gorm.Model
	URL                        string
	DateTime                   string
	CFBundleIdentifier         string
	CFBundleName               string
	CFBundleDisplayName        string
	CFBundleVersion            string
	CFBundleShortVersionString string
}

//ParseIpa : It parses the given ipa and returns a map from the contents of Info.plist in it
func parseIpa(name string) (map[string]interface{}, error) {
	r, err := zip.OpenReader(name)
	if err != nil {
		log.Println("Error opening ipa/zip ", err.Error())
		return nil, err
	}
	defer r.Close()

	for _, file := range r.File {
		if strings.HasSuffix(file.Name, ".app/Info.plist") {
			rc, err := file.Open()
			if err != nil {
				log.Println("Error opening Info.plist in zip", err.Error())
				return nil, err
			}
			buf := make([]byte, file.FileInfo().Size())
			_, err = io.ReadFull(rc, buf)
			if err != nil {
				log.Println("Error reading Info.plist", err.Error())
				return nil, err
			}
			var info_map map[string]interface{}
			_, err = plist.Unmarshal(buf, &info_map)
			if err != nil {
				log.Println("Error reading Info.plist", err.Error())
				return nil, err
			}
			return info_map, nil
		}
	}
	return nil, errors.New("Info.plist not found")
}

func ipaScaner() {

	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		distrs := filesEnum(cfg.Paths.Distrs)
		for _, distr := range distrs {
			if strings.HasSuffix(distr, ".ipa") {
				ipaProcessor(cfg.Paths.Distrs, distr)
			}
		}
	}
}

func ipaProcessor(dirPath string, fileName string) {
	var ipas []Ipa
	var ipa Ipa

	filePath := dirPath + "/" + fileName

	ipaInfo, err := parseIpa(filePath)
	if err == nil {
		ipa.CFBundleIdentifier = fmt.Sprint(ipaInfo["CFBundleIdentifier"])
		ipa.CFBundleName = fmt.Sprint(ipaInfo["CFBundleName"])
		ipa.CFBundleDisplayName = fmt.Sprint(ipaInfo["CFBundleDisplayName"])
		ipa.CFBundleVersion = fmt.Sprint(ipaInfo["CFBundleVersion"])
		ipa.CFBundleShortVersionString = fmt.Sprint(ipaInfo["CFBundleShortVersionString"])
		ipa.DateTime = time.Now().Format("2006.01.02 15:04:05")
		ipa.URL = fmt.Sprintf("%s/ipa/%s-%s.%s/%s", cfg.Service.Url, ipa.CFBundleName, ipa.CFBundleShortVersionString, ipa.CFBundleVersion, fileName)

		ipas, _ = SQLiteGetIpas()
		if !containsIpas(ipas, ipa) {
			CopyFile(dirPath, fmt.Sprintf(".\\ipa\\%s-%s.%s", ipa.CFBundleName, ipa.CFBundleShortVersionString, ipa.CFBundleVersion), fileName)
			CopyDir(".\\images", fmt.Sprintf(".\\ipa\\%s-%s.%s", ipa.CFBundleName, ipa.CFBundleShortVersionString, ipa.CFBundleVersion))
			CreatePlist(ipa)
			SQLiteAddIpa(ipa)
			deleteFile(filePath)
			log.Printf("IPA %s is added\n", fmt.Sprintf("%s-%s.%s", ipa.CFBundleName, ipa.CFBundleShortVersionString, ipa.CFBundleVersion))
		} else {
			log.Printf("IPA %s is already exist\n", fmt.Sprintf("%s-%s.%s", ipa.CFBundleName, ipa.CFBundleShortVersionString, ipa.CFBundleVersion))
			deleteFile(filePath)
		}
	}

}
