package handlers

import (
	"fmt"
	"log"
	"net/http"

	ac "github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

func AddMicroserviceDefinition(context echo.Context) error {

	log.Printf("[handlers] binding new microservice definition...")

	var microservice ac.MicroserviceDefinition
	err := context.Bind(&microservice)
	if err != nil {
		msg := fmt.Sprintf("unable to JSON to struct", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	err = ac.AddMicroserviceDefinition(&microservice)
	if err != nil {
		msg := fmt.Sprintf("unable to add microservice %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	log.Printf("%s", color.HiGreenString("[handlers] successuflly added new microservice: %s", microservice.Name))

	return context.JSON(http.StatusOK, microservice)
}

func AddMicroserviceMappings(context echo.Context) error {

	log.Printf("[handlers] unmarshalling new microservice mappping...")

	var mappings ac.MicroserviceBatch
	err := context.Bind(&mappings)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	lastInserted, err := ac.AddMicroserviceMappings(&mappings)
	if err != nil {
		msg := fmt.Sprintf("variables not added: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, lastInserted)
}
