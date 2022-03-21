package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// SQLiteCreateDB create/magrate DB
func SQLiteCreateDB(IpaStruct Ipa) error {
	db, err := gorm.Open("sqlite3", "./ipa.db")
	if err != nil {
		panic("Failed to open the SQLite database.")
	}
	defer db.Close()

	// Create the table from our struct.
	db.AutoMigrate(&IpaStruct)

	log.Println("Create/migrate registry-DB successfull")

	return nil
}

func SQLiteAddIpa(ipa Ipa) error {
	db, err := gorm.Open("sqlite3", "./ipa.db")
	if err != nil {
		panic("Failed to open the SQLite database.")
	}
	defer db.Close()

	db.Create(&ipa)

	return nil
}

func SQLiteDelIpa(ipa Ipa) error {
	db, err := gorm.Open("sqlite3", "./ipa.db")
	if err != nil {
		panic("Failed to open the SQLite database.")
	}
	defer db.Close()

	db.Where("id = ?", ipa.ID).Delete(&ipa)

	return nil
}

func SQLiteSaveIpa(ipa Ipa) error {
	db, err := gorm.Open("sqlite3", "./ipa.db")
	if err != nil {
		panic("Failed to open the SQLite database.")
	}
	defer db.Close()

	db.Where("id = ?", ipa.ID).Save(&ipa)

	return nil
}

func SQLiteGetIpa(id string) (Ipa, error) {
	db, err := gorm.Open("sqlite3", "./ipa.db")
	if err != nil {
		panic("Failed to open the SQLite database.")
	}
	defer db.Close()

	var ipa Ipa
	db.Where("id = ?", id).Find(&ipa)

	return ipa, nil
}

func SQLiteFindIpa(sha256 string) (Ipa, error) {
	db, err := gorm.Open("sqlite3", "./ipa.db")
	if err != nil {
		panic("Failed to open the SQLite database.")
	}
	defer db.Close()

	var ipa Ipa
	db.Where("sha256 = ?", sha256).Find(&ipa)

	return ipa, nil
}

func SQLiteGetIpas() ([]Ipa, error) {
	db, err := gorm.Open("sqlite3", "./ipa.db")
	if err != nil {
		panic("Failed to open the SQLite database.")
	}
	defer db.Close()

	var ipas []Ipa
	db.Order("date_time desc").Find(&ipas)

	return ipas, nil
}
