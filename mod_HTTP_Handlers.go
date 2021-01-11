package main

import (
	"image/png"
	"log"
	"net/http"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/fogleman/gg"
)

func handlerError(c *gin.Context, err error) bool {
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "Error", "message": err.Error()})
		return true
	}
	return false
}

func handlerCustomError(c *gin.Context, err string) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "Error", "message": err})
}

func indexHandler(ctx *gin.Context) {
	var ipas []Ipa
	ipas, _ = SQLiteGetIpas()
	e := casbin.NewEnforcer("./model.conf", "./policy.csv")

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
	RemoveDir("./ipa/" + ipa.CFBundleName + "-" + ipa.CFBundleVersion)
	SQLiteDelIpa(ipa)
	ctx.Redirect(http.StatusMovedPermanently, cfg.Service.Url)
	log.Println("Ipa delete has completed")
}

func versionHandler(ctx *gin.Context) {
	var ipa Ipa
	var ver = ctx.Param("version")
	ipa, _ = SQLiteFindIpa(ver)

	ctx.HTML(http.StatusOK, "version/index", gin.H{
		"title":       "IPA Manager",
		"version":     version,
		"ipa":         ipa,
		"service_url": cfg.Service.Url,
	},
	)
}

func qrHandler(ctx *gin.Context) {
	dataString := ctx.PostForm("url")
	version := ctx.PostForm("version")

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
