package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"howett.net/plist"
)

const version = "v.1.3.0"

var cfg Config

//Stand stands DB struct
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

func main() {

	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create("web.log")
	gin.DefaultWriter = io.MultiWriter(os.Stdout, f)
	defer f.Close()

	log.SetFlags(log.LstdFlags)
	lf, err := os.OpenFile("output.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer lf.Close()
	multi := io.MultiWriter(os.Stdout, lf)
	log.SetOutput(multi)

	readConfigFile(&cfg)

	SQLiteCreateDB(Ipa{})

	router := gin.Default()
	router.HTMLRender = ginview.Default()

	router.Use(static.Serve("/ipa", static.LocalFile("./ipa", false)))
	router.Use(static.Serve("/images", static.LocalFile("./images", false)))
	router.GET("/", indexHandler)
	router.POST("/action/qr", qrHandler)
	router.POST("/action/remove", removeHandler)

	go ipaScaner()

	log.Println("Web is available at " + cfg.Service.Url + ":" + strconv.Itoa(cfg.Service.Port))
	router.Run(":" + strconv.Itoa(cfg.Service.Port))
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
	var ipas []Ipa
	var ipa Ipa

	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		distrs := filesEnum(cfg.Paths.Distrs)
		for _, distr := range distrs {
			if strings.HasSuffix(distr, ".ipa") {
				ipaInfo, err := parseIpa(cfg.Paths.Distrs + "/" + distr)
				if err == nil {
					ipa.CFBundleIdentifier = fmt.Sprint(ipaInfo["CFBundleIdentifier"])
					ipa.CFBundleName = fmt.Sprint(ipaInfo["CFBundleName"])
					ipa.CFBundleDisplayName = fmt.Sprint(ipaInfo["CFBundleDisplayName"])
					ipa.CFBundleVersion = fmt.Sprint(ipaInfo["CFBundleVersion"])
					ipa.CFBundleShortVersionString = fmt.Sprint(ipaInfo["CFBundleShortVersionString"])
					ipa.DateTime = time.Now().Format("2006.01.02 15:04:05")
					ipa.URL = cfg.Service.Url + "/ipa/" + ipa.CFBundleName + "-" + ipa.CFBundleVersion + "/" + distr

					ipas, _ = SQLiteGetIpas()
					if containsIpas(ipas, ipa.CFBundleVersion) != true {
						CopyFile(cfg.Paths.Distrs, "./ipa/"+ipa.CFBundleName+"-"+ipa.CFBundleVersion, distr)
						CopyDir("./images", "./ipa/"+ipa.CFBundleName+"-"+ipa.CFBundleVersion)
						CreatePlist(ipa)
						SQLiteAddIpa(ipa)
						deleteFile(cfg.Paths.Distrs + "/" + distr)
						log.Printf("IPA %s is added\n", ipa.CFBundleVersion)
					} else {
						log.Printf("IPA %s is already exist\n", ipa.CFBundleVersion)
						deleteFile(cfg.Paths.Distrs + "/" + distr)
					}
				}
			}
		}
	}
}
