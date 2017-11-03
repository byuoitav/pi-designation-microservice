package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	ac "github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

const VARIABLE_COLUMN_NAME = "value"
const VARIABLE_MAPPINGS_TABLE = "variable_mappings"
const VARIABLE_DEFINITION_COLUMN = "variable_id"
const VARIABLE_DEFINITION_TABLE = "variable_definitions"

func AddVariableMapping(context echo.Context) error {

	log.Printf("[handlers] binding new variable mapping...")

	var mapping ac.VariableMapping
	err := context.Bind(&mapping)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	id, err := ac.AddMapping(
		VARIABLE_MAPPINGS_TABLE,
		VARIABLE_DEFINITION_COLUMN,
		VARIABLE_COLUMN_NAME,
		mapping.Value,
		mapping.Variable.ID,
		mapping.Class.ID,
		mapping.Designation.ID)
	if err != nil {
		msg := fmt.Sprintf("unable to add mapping: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	var entry ac.VariableMapping
	err = ac.GetVariableMappingById(id, &entry)
	if err != nil {
		msg := fmt.Sprintf("new entry not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusInternalServerError, msg)
	}

	return context.JSON(http.StatusOK, entry)
}

//relies on MySQL for most logic
//e.g. foreign keys, duplicates, etc
func AddVariableMappings(context echo.Context) error {

	log.Printf("[handlers] binding new variable mappings...")

	var mappings ac.Batch
	err := context.Bind(&mappings)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	lastInserted, err := ac.AddMappings(
		VARIABLE_MAPPINGS_TABLE,
		VARIABLE_DEFINITION_COLUMN,
		VARIABLE_COLUMN_NAME,
		&mappings)
	if err != nil {
		msg := fmt.Sprintf("variables not added: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	entries, err := ac.GetVariableMappingsById(lastInserted)
	if err != nil {
		msg := fmt.Sprintf("new entries not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusInternalServerError, msg)
	}

	return context.JSON(http.StatusOK, entries)
}

func EditVariableMapping(context echo.Context) error {

	log.Printf("[handlers] binding variable mapping...")

	var mapping ac.VariableMapping
	err := context.Bind(&mapping)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	err = ac.EditMapping(
		VARIABLE_MAPPINGS_TABLE,
		VARIABLE_DEFINITION_COLUMN,
		VARIABLE_COLUMN_NAME,
		mapping.Value,
		mapping.Variable.ID,
		mapping.Class.ID,
		mapping.Designation.ID,
		mapping.ID)
	if err != nil {
		msg := fmt.Sprintf("variables not added: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	var entry ac.VariableMapping
	err = ac.GetVariableMappingById(mapping.ID, &entry)
	if err != nil {
		msg := fmt.Sprintf("new entries not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusInternalServerError, msg)
	}

	return context.JSON(http.StatusOK, entry)
}

func AddVariableDefinition(context echo.Context) error {

	log.Printf("[handlers] binding new variable definition...")

	var variable ac.Definition
	err := context.Bind(&variable)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	log.Printf("[handlers] adding variable definition...")

	err = ac.AddDefinition(VARIABLE_DEFINITION_TABLE, &variable)
	if err != nil {
		msg := fmt.Sprintf("variable definition failed: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, variable)
}

func EditVariableDefinition(context echo.Context) error {

	log.Printf("[handlers] binding variable definition...")

	var variable ac.Definition
	err := context.Bind(&variable)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	log.Printf("[handlers] editing variable definition...")

	err = ac.EditDefinition(VARIABLE_DEFINITION_TABLE, &variable)
	if err != nil {
		msg := fmt.Sprintf("edit failed: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, variable)
}

func GetVariableDefinition(context echo.Context) error {

	stringId := context.Param("id")

	intId, err := strconv.Atoi(stringId)
	if err != nil {
		msg := fmt.Sprintf("invalid ID: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	log.Printf("[handlers] getting variable definition with ID: %d", intId)

	variable := ac.Definition{ID: int64(intId)}
	err = ac.GetDefinitionById(VARIABLE_DEFINITION_TABLE, &variable)
	if err != nil {
		msg := fmt.Sprintf("accessor error: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, variable)
}

func GetAllVariableDefinitions(context echo.Context) error {

	log.Printf("[handlers] fetching all variable definitions...")

	var variables []ac.Definition
	err := ac.GetAllDefinitions(VARIABLE_DEFINITION_TABLE, &variables)
	if err != nil {
		msg := fmt.Sprintf("accessor error: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, variables)
}

func GetVariableMappingById(context echo.Context) error {

	stringId := context.Param("id")
	intId, err := strconv.Atoi(stringId)
	if err != nil {
		msg := fmt.Sprintf("invalid ID: %s")
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	log.Printf("[handlers] getting variable mapping with ID: %d", intId)

	var variable ac.VariableMapping
	err = ac.GetVariableMappingById(int64(intId), &variable)
	if err != nil {
		msg := fmt.Sprintf("accessor error: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, variable)
}

func GetAllVariableMappings(context echo.Context) error { return nil }
