package demo

import (
	"encoding/json"
	"fmt"
	"math/rand"

	bucket "github.com/PoulIgorson/sub_engine_fiber/database"
	db "github.com/PoulIgorson/sub_engine_fiber/database"
)

type Car struct {
	ID    uint   `json:"id"`
	Model string `json:"model"`
	Color string `json:"color"`
	City  string `json:"city"`
}

func Create(db_ *db.DB, carStr string) *Car {
	car := &Car{}
	json.Unmarshal([]byte(carStr), car)
	return car
}

func (car Car) Create(db_ *db.DB, carStr string) db.Model {
	return *Create(db_, carStr)
}

func (car *Car) Save(bct *bucket.Bucket) error {
	return bucket.SaveModel(bct, car)
}

func CreateModels(db_ *db.DB) {
	models := []string{"BMW", "Volvo", "Porch", "WW", "Tesla"}
	colors := []string{"red", "green", "blue", "white", "black"}
	cities := []string{"Moscow", "SP", "Vladimir", "Paris", "Rostov"}

	carBct, _ := db_.Bucket("car", Car{})
	for i := 1; i < 11; i++ {
		car := Car{
			Model: models[rand.Int()%len(models)],
			Color: colors[rand.Int()%len(models)],
			City:  cities[rand.Int()%len(models)],
		}
		if err := car.Save(carBct); err != nil {
			fmt.Println(i, err)
		}
	}
}

func Run() {
	db_, err := db.Open("sub_engine_fiber_db.db")
	if err != nil {
		panic(err)
	}
	defer db_.Close()
	fmt.Println("App start")

	CreateModels(db_)

	carBct, _ := db_.Bucket("car", Car{})
	carBct.Delete(5)
	cars := carBct.Objects.Filter(db.Params{"Model": "Tesla"}, db.Params{"Color": "black", "City": "Moscow"}).All()
	fmt.Printf("%10v | %10v | %10v | %10v\n", "ID", "model", "color", "city")
	fmt.Printf("---------- | ---------- | ---------- | ----------\n")
	for _, carM := range cars {
		car := carM.(Car)
		fmt.Printf("%10v | %10v | %10v | %10v\n", car.ID, car.Model, car.Color, car.City)
	}
}