package main

import (
	"encoding/hex"
	"fmt"
	"github.com/shumon84/rainbow-table"
	"html/template"
	"log"
	"net/http"
	"os"
)

func HashHandler(w http.ResponseWriter, r *http.Request) {
	plain := r.FormValue("plain")
	hash := rainbow.Hash([]byte(plain))
	hashString := hex.EncodeToString(hash)
	templateFile := template.Must(template.ParseFiles("index.html"))
	templateFile.Execute(w, template.HTML(fmt.Sprintf("<p><b>%s</b>をハッシュ化すると</p><p><b>%s</b>になります</p>", plain, hashString)))
}

func RehashHandler(w http.ResponseWriter, r *http.Request) {
	hashString := r.FormValue("hash")
	hash, err := hex.DecodeString(hashString)
	templateFile := template.Must(template.ParseFiles("index.html"))
	if err != nil || len(hash) == 0 {
		log.Print(err)
		templateFile.Execute(w, template.HTML(fmt.Sprintf("<p><b>%s</b>は不正なハッシュ値です</p>", hashString)))
		return
	}
	plain := rainbow.ReHash(hash, rainbow.Hash, rainbow.Reduction)
	if string(plain) == rainbow.NotFound {
		templateFile.Execute(w, template.HTML(fmt.Sprintf("<p><b>%s</b>は複合できませんでした</p>", hashString)))
		return
	}
	templateFile.Execute(w, template.HTML(fmt.Sprintf("<p><b>%s</b>は</p><p><b>%s</b>をハッシュ化したものです</p>", hashString, string(plain))))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	templateFile := template.Must(template.ParseFiles("index.html"))
	templateFile.Execute(w, "")
}

func main() {
	file, err := os.Open("converted_table_4_20000_5000.txt")
	if err != nil {
		log.Fatal(err)
	}
	rainbow.ReadRainbowTable(file)
	file.Close()

	http.HandleFunc("/hash", HashHandler)
	http.HandleFunc("/rehash", RehashHandler)
	http.HandleFunc("/", IndexHandler)
	log.Fatal(http.ListenAndServe(":80", http.DefaultServeMux))
}
