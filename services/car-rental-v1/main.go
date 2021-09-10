package main

import (
	"car-rental-v1/internal/handler"
	"car-rental-v1/wrappers/server"

	instana "github.com/instana/go-sensor"
)

func main() {
	handler.GetData()

	if err := server.Start("car-rental-v1", initializeRouter); err != nil {
		panic(err)
	}
}

func initializeRouter(router server.PathRouter, _ *instana.Sensor) {
	router.Path("/api/v1/cars", func(router server.PathRouter) {
		router.Path("/info/{tag:string}", func(router server.PathRouter) {
			/**
			* GET /api/v1/cars/info/{filter}
			* @tag Car Rental
			* @summary Get filter list
			* @description Gets list of a type to filter Car Rental data by.
			* @pathParam {FilterType} filter - The name of the filter to get options for.
			* @response 200 - OK
			* @response 400 - Filter Not Found Error
			* @response 500 - Internal Server Error
			 */
			router.Get(handler.GetFilterList())
		})
		router.Path("/{id:string}", func(router server.PathRouter) {
			/**
			* GET /api/v1/cars/{id}
			* @tag Car Rental
			* @summary Get car by id
			* @description Gets data associated with a specific car ID.
			* @pathParam {string} id - id of the car
			* @queryParam {string} dateFrom - Date From
			* @queryParam {string} dateTo - Date To
			* @response 200 - OK
			* @response 404 - not found
			* @response 500 - Internal Server Error
			 */
			router.Get(handler.GetCarByID())
		})
		router.Path("/{country:string}/{city:string}", func(router server.PathRouter) {
			/**
			* GET /api/v1/cars/{country}/{city}
			* @tag Car Rental
			* @summary Get list of cars
			* @description Gets data associated with a specific city.
			* @pathParam {string} country - Country of the rental company using slug casing.
			* @pathParam {string} city - City of the rental company using slug casing.
			* @queryParam {string} dateFrom - Date From
			* @queryParam {string} dateTo - Date To
			* @queryParam {string} [company] - Rental Company name.
			* @queryParam {string} [car] - Car Name.
			* @queryParam {string} [type] - Car Type.
			* @queryParam {string} [style] - Car Style.
			* @queryParam {number} [mincost] - Min Cost.
			* @queryParam {number} [maxcost] - Max Cost.
			* @response 200 - OK
			* @response 500 - Internal Server Error
			 */
			router.Get(handler.GetCars())
		})
	})
}
