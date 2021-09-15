package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/tidwall/buntdb"
)

type Record struct {
	CreatedAt int64
	Name      string
	Odd       bool
}

func main() {
	db, err := buntdb.Open(":memory:")
	if err != nil {
		log.Fatal(err)
	}

	// Criacao de indice
	// Para classificar na data da criacao
	db.CreateIndex("index", "entry:*", buntdb.IndexJSON("CreatedAt"))
	// Para filtrar com livre impar
	db.CreateIndex("odd", "entry:*", buntdb.IndexJSON("Odd"))

	// Registro de dados
	db.Update(func(tx *buntdb.Tx) error {
		for i := 0; i < 10; i++ {
			odd := i%2 != 0
			rec := Record{Name: time.Now().String(), CreatedAt: time.Now().UnixNano(), Odd: odd}
			buf, err := json.Marshal(rec)
			if err != nil {
				return err
			}
			tx.Set("entry:"+strconv.Itoa(i), string(buf), nil)
			time.Sleep(time.Millisecond * 300)
		}
		return nil
	})

	// Dados registrados
	db.View(func(tx *buntdb.Tx) error {
		fmt.Println("CreatedAt decrescente")
		tx.Descend("index", func(key, value string) bool {
			fmt.Println(key, value)
			return true
		})

		fmt.Println("CreatedAt crescente")
		tx.Ascend("index", func(key, value string) bool {
			fmt.Println(key, value)
			return true
		})

		fmt.Println("Ordem pares")
		tx.DescendEqual("odd", `{"Odd":false}`, func(key, value string) bool {
			fmt.Println(key, value)
			return true
		})

		fmt.Println("Ordem impares")
		tx.AscendEqual("odd", `{"Odd":true}`, func(key, value string) bool {
			fmt.Println(key, value)
			return true
		})

		return nil
	})

	defer db.Close()
}
