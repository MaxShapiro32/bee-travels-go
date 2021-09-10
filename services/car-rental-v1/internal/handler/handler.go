package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	"car-rental-v1/internal/data"
	"car-rental-v1/wrappers/server"
)

var Cars []data.Car

func GetData() {
	data, err := ioutil.ReadFile("./data/cars.json")
	if err != nil {
		fmt.Print(err)
	}

	err = json.Unmarshal([]byte(data), &Cars)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func GetCars() server.RequestHandler {
	return func(ctx server.RequestContext) {
		country := ctx.Params().Get("country")
		city := ctx.Params().Get("city")
		company := stringToArray(ctx.URLParam("company"))
		car_ := stringToArray(ctx.URLParam("car"))
		type_ := stringToArray(ctx.URLParam("type"))
		style := stringToArray(ctx.URLParam("style"))
		minCost, minCostExists := ctx.URLParamFloat64("mincost")
		maxCost, maxCostExists := ctx.URLParamFloat64("maxcost")
		dateFrom, _ := time.Parse(time.RFC3339, ctx.URLParam("dateFrom")+"T00:00:00.000Z")
		dateTo, _ := time.Parse(time.RFC3339, ctx.URLParam("dateTo")+"T00:00:00.000Z")

		fmt.Println("Getting car rental data for -> /" + country + "/" + city)

		if dateTo.Sub(dateFrom) < 0 {
			server.Response(ctx, http.StatusInternalServerError, data.Error{
				Error: "from date can not be greater than to date",
			})
			return
		}

		multiplier := dateMultiplier(dateFrom)
		var carsData []data.Car
		for _, car := range Cars {
			if reflect.ValueOf(car).FieldByName("Country").String() == capitalize(country) &&
				reflect.ValueOf(car).FieldByName("City").String() == capitalize(city) &&
				(strings.Trim(company[0], " ") == "" || contains(company, reflect.ValueOf(car).FieldByName("RentalCompany").String())) &&
				(strings.Trim(car_[0], " ") == "" || contains(car_, reflect.ValueOf(car).FieldByName("Name").String())) &&
				(strings.Trim(type_[0], " ") == "" || contains(type_, reflect.ValueOf(car).FieldByName("BodyType").String())) &&
				(strings.Trim(style[0], " ") == "" || contains(style, reflect.ValueOf(car).FieldByName("Style").String())) &&
				(minCostExists != nil || minCost <= car.Cost*multiplier) &&
				(maxCostExists != nil || car.Cost <= maxCost*multiplier) {
				carsData = append(carsData, data.Car{
					RentalCompany: car.RentalCompany,
					City:          car.City,
					BodyType:      car.BodyType,
					Cost:          car.Cost * multiplier,
					Name:          car.Name,
					Country:       car.Country,
					Image:         car.Image,
					CarId:         car.CarId,
					Id:            car.Id,
					Style:         car.Style,
				})
			}
		}
		server.Response(ctx, http.StatusOK, carsData)
	}
}

func GetFilterList() server.RequestHandler {
	return func(ctx server.RequestContext) {
		filterType := ctx.Params().Get("tag")
		fmt.Println("Getting info for " + filterType)

		var listOfFilterOptions []string
		var isDup = false
		var currentFilterOption string

		for _, car := range Cars {
			for _, filterOption := range listOfFilterOptions {
				currentFilterOption = reflect.ValueOf(car).FieldByName(strings.Title(strings.ReplaceAll(filterType, "_", ""))).String()
				if currentFilterOption == filterOption {
					isDup = true
					break
				}
			}
			if isDup == false {
				listOfFilterOptions = append(listOfFilterOptions, currentFilterOption)
			}
			isDup = false
		}
		/*err := nil
		if err != nil {
			server.Response(ctx, http.StatusForbidden, Error{
				Error: err.Error(),
			})
			return
		}*/
		server.Response(ctx, http.StatusOK, listOfFilterOptions[1:])
	}
}

func GetCarByID() server.RequestHandler {
	return func(ctx server.RequestContext) {
		id := ctx.Params().Get("id")
		dateFrom, _ := time.Parse(time.RFC3339, ctx.URLParam("dateFrom")+"T00:00:00.000Z")
		dateTo, _ := time.Parse(time.RFC3339, ctx.URLParam("dateTo")+"T00:00:00.000Z")
		if dateTo.Sub(dateFrom) < 0 {
			server.Response(ctx, http.StatusInternalServerError, data.Error{
				Error: "from date can not be greater than to date",
			})
			return
		}

		for _, car := range Cars {
			if reflect.ValueOf(car).FieldByName("Id").String() == id {
				multiplier := dateMultiplier(dateFrom)
				if multiplier == -1 {
					server.Response(ctx, http.StatusInternalServerError, data.Error{
						Error: "from date can not be greater than to date",
					})
					return
				}

				server.Response(ctx, http.StatusOK, data.Car{
					RentalCompany: car.RentalCompany,
					City:          car.City,
					BodyType:      car.BodyType,
					Cost:          car.Cost * multiplier,
					Name:          car.Name,
					Country:       car.Country,
					Image:         car.Image,
					CarId:         car.CarId,
					Id:            car.Id,
					Style:         car.Style,
					DateFrom:      dateFrom.Format("2006-01-02"),
					DateTo:        dateTo.Format("2006-01-02"),
				})
				return
			}
		}

		server.Response(ctx, http.StatusNotFound, data.Error{
			Error: "not found",
		})
		return
	}
}

func stringToArray(s string) []string {
	return strings.Split(s, ",")
}

func capitalize(str string) string {
	lcExceptions := []string{"es", "de", "au"}
	dashes := false

	splitText := strings.Split(strings.ToLower(str), "-")
	for i, word := range splitText {
		if contains(lcExceptions, word) {
			if word == lcExceptions[2] {
				dashes = true
			}
		} else {
			splitText[i] = strings.Title(word)
		}
	}

	if dashes {
		return strings.Join(splitText, "-")
	}
	return strings.Join(splitText, " ")
}

func dateMultiplier(dateFrom time.Time) float64 {
	dateNow := time.Now()
	numDays := dateFrom.Sub(dateNow).Hours() / 24
	if numDays < 0 {
		return -1
	} else if numDays < 2 {
		return 2.25
	} else if numDays < 7 {
		return 1.75
	} else if numDays < 14 {
		return 1.5
	} else if numDays < 21 {
		return 1.2
	} else if numDays < 45 {
		return 1
	} else if numDays < 90 {
		return 0.8
	} else {
		return -1
	}
}

func contains(arr []string, str string) bool {
	for _, value := range arr {
		if value == str {
			return true
		}
	}
	return false
}
