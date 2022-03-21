package main

import (
	"fmt"
	"image/png"
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/fogleman/gg"
)

func indexHandler(ctx *gin.Context) {
	var ipas []Ipa
	ipas, _ = SQLiteGetIpas()
	for id, ipa := range ipas {
		ipas[id].URL = fmt.Sprintf("%s/ipa/%s/%s", cfg.Service.Url, ipa.SHA256, ipa.FileName)
	}

	e := casbin.NewEnforcer("./model.conf", "./policy.csv")
	//log.Println(ctx.Value("user"))

	// Admin rights
	if user := ctx.Value("user"); user != nil {

		if e.Enforce(user, "index", "write") {
			ctx.HTML(http.StatusOK, "index", gin.H{
				"title":       "IPA Manager",
				"version":     version,
				"ipas":        ipas,
				"admin":       1,
				"service_url": cfg.Service.Url,
			},
			)
			return
		}
	}

	// Guest rights
	ctx.HTML(http.StatusOK, "index", gin.H{
		"title":       "IPA Manager",
		"version":     version,
		"ipas":        ipas,
		"admin":       0,
		"service_url": cfg.Service.Url,
	},
	)
}

func removeHandler(ctx *gin.Context) {
	var ipa Ipa
	var id = ctx.PostForm("id")
	ipa, _ = SQLiteGetIpa(id)
	RemoveDir(fmt.Sprintf(".\\ipa\\%s", ipa.SHA256))
	SQLiteDelIpa(ipa)
	ctx.Redirect(http.StatusMovedPermanently, cfg.Service.Url)
	log.Println("Ipa delete has completed")
}

func versionHandler(ctx *gin.Context) {
	var ipa Ipa
	var sha256 = ctx.Param("sha256")

	ipa, _ = SQLiteFindIpa(sha256)
	ipa.URL = fmt.Sprintf("%s/ipa/%s/%s", cfg.Service.Url, ipa.SHA256, ipa.FileName)

	ctx.HTML(http.StatusOK, "version/index", gin.H{
		"title":       "IPA Manager",
		"ipa":         ipa,
		"service_url": cfg.Service.Url,
	},
	)
}

func qrHandler(ctx *gin.Context) {
	dataString := ctx.PostForm("url")
	CFBundleIdentifier := ctx.PostForm("CFBundleIdentifier")
	version := CFBundleIdentifier + " - " + ctx.PostForm("version")

	qrCode, _ := qr.Encode(dataString, qr.L, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 600, 600)

	im := qrCode

	dc := gg.NewContext(600, 626)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("arial.ttf", 16); err != nil {
		panic(err)
	}

	dc.DrawRoundedRectangle(0, 0, 600, 626, 0)
	dc.DrawImage(im, 0, 0)
	dc.DrawStringAnchored(version, 300, 615, 0.5, 0)
	dc.Clip()

	png.Encode(ctx.Writer, dc.Image())

	ctx.String(http.StatusOK, "Done")
}

func postIpaHandler(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}
	if !strings.HasSuffix(file.Filename, ".ipa") {
		ctx.JSON(http.StatusBadRequest, gin.H{"responce": "Invalid file extension"})
		return
	}

	//log.Println(file.Filename)

	err = ctx.SaveUploadedFile(file, "./temp/"+file.Filename)
	if err != nil {
		log.Fatal(err)
	}

	ipaProcessor("./temp", file.Filename)

	ctx.JSON(http.StatusOK, gin.H{"responce": "File processed"})
}
