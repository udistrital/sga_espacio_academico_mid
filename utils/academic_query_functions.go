// Academic query functions.
// Funciones generalizadas para consultar los servicios de
// espacios académicos o proyeto académico y obtener los
// regisros resultantes
package utils

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_espacio_academico_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

func GetAcademicSpacesByQuery(query string) (any, error) {
	var resSpaces interface{}
	urlAcademicSpaces := "http://" + beego.AppConfig.String("EspaciosAcademicosService") +
		"espacio-academico?" + query
	if errSpace := request.GetJson(urlAcademicSpaces, &resSpaces); errSpace == nil {
		if resSpaces.(map[string]interface{})["Data"] != nil {
			return resSpaces.(map[string]interface{})["Data"], nil
		} else {
			return nil, fmt.Errorf("EspaciosAcademicosService No se encuentran espacios académicos")
		}
	} else {
		return nil, errSpace
	}
}

// UpdateAcademicSpace update the academic space and return data of updated space
func UpdateAcademicSpace(id string, spaceData interface{}) (map[string]interface{}, error) {
	var updatedSpace map[string]interface{}
	urlAcademicSpaces := "http://" + beego.AppConfig.String("EspaciosAcademicosService") +
		"espacio-academico/" + fmt.Sprintf("%v", id)
	if errUpdate := helpers.SendJson(urlAcademicSpaces, "PUT", &updatedSpace, spaceData); errUpdate == nil {
		return updatedSpace["Data"].(map[string]interface{}), nil
	} else {
		return nil, fmt.Errorf(fmt.Sprintf("Error actualizando espacio académico %v", id))
	}
}
