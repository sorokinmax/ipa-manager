package main

import (
	"fmt"
	"log"

	"github.com/beevik/etree"
)

func CreatePlist(ipa Ipa) bool {
	doc := GetXMLDoc("template.plist")

	el := doc.FindElement("//array[1]/dict[1]/string[2]")
	el.SetText(ipa.URL)

	el = doc.FindElement("//dict[2]/string[2]")
	el.SetText(cfg.Service.Url + "/ipa/" + ipa.CFBundleName + "-" + ipa.CFBundleVersion + "/display-image.png")

	el = doc.FindElement("//dict[3]/string[2]")
	el.SetText(cfg.Service.Url + "/ipa/" + ipa.CFBundleName + "-" + ipa.CFBundleVersion + "/full-size-image.png")

	el = doc.FindElement("//dict[1]/dict[1]/string[1]")
	el.SetText(ipa.CFBundleIdentifier)

	el = doc.FindElement("//dict[1]/dict[1]/string[2]")
	el.SetText(ipa.CFBundleVersion)

	el = doc.FindElement("//string[4]")
	el.SetText(ipa.CFBundleName)

	doc.WriteToFile(fmt.Sprintf(".\\ipa\\%s-%s.%s\\%s.plist", ipa.CFBundleName, ipa.CFBundleShortVersionString, ipa.CFBundleVersion, ipa.CFBundleName))
	return true
}

func GetXMLDoc(configPath string) *etree.Document {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(configPath); err != nil {
		log.Println("Error parse "+configPath+": ", err.Error())
		return nil
	}
	return doc
}
