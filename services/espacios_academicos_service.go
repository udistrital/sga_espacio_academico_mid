package services

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_espacio_academico_mid/helpers"
	"github.com/udistrital/sga_espacio_academico_mid/utils"
	requestmanager "github.com/udistrital/sga_espacio_academico_mid/utils/requestManager"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

// FUNCIONES PRINCIPALES

func GetAcademicSpacesByProject(idProyecto int64) requestresponse.APIResponse {
	/*
		consulta espacios academicos por proyecto
	*/
	var Espacios_academicos_1 map[string]interface{}
	Espacios_academicos_1Err := request.GetJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+
		fmt.Sprintf("espacio-academico?query=activo:true,proyecto_academico_id:%v,espacio_academico_padre&limit=0", idProyecto), &Espacios_academicos_1)
	if Espacios_academicos_1Err != nil || Espacios_academicos_1["Success"] == false || Espacios_academicos_1["Status"] != "200" {
		if Espacios_academicos_1Err == nil {
			Espacios_academicos_1Err = fmt.Errorf("EspaciosAcademicosService: %v", Espacios_academicos_1["Message"])
		}
		logs.Error(Espacios_academicos_1Err.Error())
		return requestresponse.APIResponseDTO(false, 404, nil, Espacios_academicos_1Err.Error())
	}
	/*
		consulta parametros, clase, enfoque
	*/
	idClase := 51
	var ClaseEspacio map[string]interface{}
	ClaseEspacioErr := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+
		fmt.Sprintf("parametro?query=TipoParametroId:%v&limit=0&fields=Id,Nombre,CodigoAbreviacion", idClase), &ClaseEspacio)
	if ClaseEspacioErr != nil || fmt.Sprintf("%v", ClaseEspacio) == "[map[]]" {
		if ClaseEspacioErr == nil {
			ClaseEspacioErr = fmt.Errorf("ParametroService: query for clases is empty")
		}
		logs.Error(ClaseEspacioErr.Error())
		return requestresponse.APIResponseDTO(false, 404, nil, ClaseEspacioErr.Error())
	}
	clases := ClaseEspacio["Data"].([]interface{})
	idEnfoque := 68
	var EnfoqueEspacio map[string]interface{}
	EnfoqueEspacioErr := request.GetJson("http://"+beego.AppConfig.String("ParametroService")+
		fmt.Sprintf("parametro?query=TipoParametroId:%v&limit=0&fields=Id,CodigoAbreviacion", idEnfoque), &EnfoqueEspacio)
	if EnfoqueEspacioErr != nil || fmt.Sprintf("%v", EnfoqueEspacio) == "[map[]]" {
		if EnfoqueEspacioErr == nil {
			EnfoqueEspacioErr = fmt.Errorf("ParametroService: query for enfoques is empty")
		}
		logs.Error(EnfoqueEspacioErr.Error())
		return requestresponse.APIResponseDTO(false, 404, nil, EnfoqueEspacioErr.Error())
	}
	enfoques := EnfoqueEspacio["Data"].([]interface{})
	/*
		Construcción información requerida
	*/
	var EspaciosAcademicos []interface{}
	for _, espacio := range Espacios_academicos_1["Data"].([]interface{}) {
		var nombresEspacios []map[string]interface{}
		var nombresEspaciosStr string = ""
		if reflect.TypeOf(espacio.(map[string]interface{})["espacios_requeridos"]).Kind() == reflect.Slice {
			for _, requerido := range espacio.(map[string]interface{})["espacios_requeridos"].([]interface{}) {
				nombreEspacio, err := getLocalEspacioAcademico(requerido.(string), Espacios_academicos_1["Data"].([]interface{}))
				if err != nil {
					nombreEspacio, err = getLineaEspacioAcademico(requerido.(string))
					if err != nil {
						nombreEspacio = "No encontrado..."
					}
				}
				nombresEspacios = append(nombresEspacios, map[string]interface{}{
					"_id":    requerido.(string),
					"nombre": nombreEspacio,
				})
				nombresEspaciosStr += nombreEspacio + ", "
			}
		}
		nombreClase, err := getClase(espacio.(map[string]interface{})["clasificacion_espacio_id"].(float64), clases)
		if err != nil {
			nombreClase = "No encontrado..."
		}
		formatoEspacio := map[string]interface{}{
			"_id":               espacio.(map[string]interface{})["_id"].(string),
			"nombre":            espacio.(map[string]interface{})["nombre"].(string),
			"prerequisitos":     nombresEspacios,
			"prerequisitos_str": nombresEspaciosStr,
			"clase":             nombreClase,
			"creditos":          espacio.(map[string]interface{})["creditos"].(float64),
			"htd":               espacio.(map[string]interface{})["distribucion_horas"].(map[string]interface{})["HTD"].(float64),
			"htc":               espacio.(map[string]interface{})["distribucion_horas"].(map[string]interface{})["HTC"].(float64),
			"hta":               espacio.(map[string]interface{})["distribucion_horas"].(map[string]interface{})["HTA"].(float64),
		}
		for _, clase := range clases {
			code := clase.(map[string]interface{})["CodigoAbreviacion"].(string)
			value := 0
			if clase.(map[string]interface{})["Id"].(float64) == espacio.(map[string]interface{})["clasificacion_espacio_id"].(float64) {
				value = 1
			}
			formatoEspacio[code] = value
		}
		for _, enfoque := range enfoques {
			code := enfoque.(map[string]interface{})["CodigoAbreviacion"].(string)
			code = strings.Replace(code, "-", "_", -1)
			value := 0
			if enfoque.(map[string]interface{})["Id"].(float64) == espacio.(map[string]interface{})["enfoque_id"].(float64) {
				value = 1
			}
			formatoEspacio[code] = value
		}
		EspaciosAcademicos = append(EspaciosAcademicos, formatoEspacio)
	}
	/*
		entrega de respuesta existosa :)
	*/
	return requestresponse.APIResponseDTO(true, 200, EspaciosAcademicos)
}

func PostAcademicSpacesBySon(data []byte) requestresponse.APIResponse {
	var espacioAcademicoRequest map[string]interface{}
	var EspacioPadrePost map[string]interface{}
	var EspacioPadrePostTempo map[string]interface{}

	if err := json.Unmarshal(data, &espacioAcademicoRequest); err == nil {

		grupos_espacios := espacioAcademicoRequest["grupo"]
		strGrupos := fmt.Sprintf("%v", grupos_espacios)
		_, GrupoIn := contarYSepararGrupos(strGrupos)

		if err := helpers.SendJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico", "POST", &EspacioPadrePost, espacioAcademicoRequest); err != nil {
			panic(map[string]interface{}{"funcion": "FuncionPostHijosEspacio", "err": "Error al generar el espacio padre  ", "status": "400", "log": err})
		}

		responseEspacioPadre := EspacioPadrePost["Data"].(map[string]interface{})
		IdEspacioAcademicoPadre := fmt.Sprintf("%v", responseEspacioPadre["_id"])
		EspacioAcademicoHijoTemporal := espacioAcademicoRequest

		EspacioAcademicoHijoTemporal["espacio_academico_padre"] = IdEspacioAcademicoPadre

		//fmt.Println(".---------------------------Espacio temporal--------------------------")
		for _, grupo := range GrupoIn {
			EspacioAcademicoHijoTemporal["grupo"] = grupo
			if err := helpers.SendJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+"espacio-academico", "POST", &EspacioPadrePostTempo, EspacioAcademicoHijoTemporal); err != nil {
				panic(map[string]interface{}{"funcion": "VersionarPlan", "err": "Error al generar el espacio padre  ", "status": "400", "log": err})
			}
		}
		return requestresponse.APIResponseDTO(true, 201, responseEspacioPadre)
	}
	return requestresponse.APIResponseDTO(false, 400, nil, "Error en el formato del cuerpo de la petición")
}

func PutAcademicSpaceAssignPeriod(data []byte) requestresponse.APIResponse {
	var periodRequestBody map[string]interface{}
	var response []map[string]interface{}

	if err := json.Unmarshal(data, &periodRequestBody); err == nil {
		parentId := fmt.Sprintf("%v", periodRequestBody["padre"])
		queryParams := "query=activo:true,espacio_academico_padre:" +
			parentId + "&fields=_id,grupo,periodo_id"
		groups := utils.Slice2SliceString(periodRequestBody["grupo"].([]interface{}))
		periodIdReq := int(periodRequestBody["periodo_id"].(float64))

		if resSpaces, errSpace := utils.GetAcademicSpacesByQuery(queryParams); errSpace == nil {
			if resSpaces != nil {
				spaces := resSpaces.([]any)
				if assignedSpaces, errAssign := assignExistingPeriod(spaces, &groups, periodIdReq); errAssign == nil {
					response = append(response, assignedSpaces...)
					if len(groups) > 0 {
						if newSpaces, errNewSpaces := createAcademicSpaceChild(parentId, groups, periodIdReq); errNewSpaces == nil {
							response = append(response, newSpaces...)
							return requestresponse.APIResponseDTO(true, 200, response)
						} else {
							if newSpaces != nil {
								response = append(response, newSpaces...)
							}
							return requestresponse.APIResponseDTO(false, 400, response, "No fue posible asignar todos los espacios académicos")
						}
					} else {
						return requestresponse.APIResponseDTO(true, 200, response)
					}
				} else {
					if assignedSpaces != nil {
						return requestresponse.APIResponseDTO(false, 400, assignedSpaces, "No fue posible asignar todos los espacios académicos")
					} else {
						return requestresponse.APIResponseDTO(false, 400, nil, "Espacios académicos no encontrados")
					}
				}
			} else {
				return requestresponse.APIResponseDTO(false, 400, nil, "Espacios académicos no encontrados")
			}
		} else {
			if newSpaces, errNewSpaces := createAcademicSpaceChild(parentId, groups, periodIdReq); errNewSpaces == nil {
				response = append(response, newSpaces...)
				return requestresponse.APIResponseDTO(true, 200, response)
			} else {
				if newSpaces != nil {
					response = append(response, newSpaces...)
				}
				return requestresponse.APIResponseDTO(false, 400, response, "No fue posible asignar todos los espacios académicos")
			}
		}
	} else {
		errResponse, statusCode := requestmanager.MidResponseFormat(
			"AsignarPeriodoEspacioAcadémico", "PUT", false, err.Error())
		return requestresponse.APIResponseDTO(false, statusCode, errResponse, err.Error())
	}
}

// FUNCIONES AXILIARES

func getLocalEspacioAcademico(_id string, espacios []interface{}) (string, error) {
	for _, espacio := range espacios {
		if _id == espacio.(map[string]interface{})["_id"] {
			return espacio.(map[string]interface{})["nombre"].(string), nil
		}
	}
	return "", fmt.Errorf("not found")
}

func getLineaEspacioAcademico(_id string) (string, error) {
	var nombreEspacio map[string]interface{}
	nombreEspacioErr := request.GetJson("http://"+beego.AppConfig.String("EspaciosAcademicosService")+
		fmt.Sprintf("espacio-academico/%v", _id), &nombreEspacio)
	if nombreEspacioErr != nil || nombreEspacio["Success"] == false || nombreEspacio["Status"] != "200" {
		if nombreEspacioErr == nil {
			nombreEspacioErr = fmt.Errorf("EspaciosAcademicosService: %v", nombreEspacio["Message"])
		}
		return "", nombreEspacioErr
	}
	return nombreEspacio["Data"].(map[string]interface{})["nombre"].(string), nil
}

func getClase(id float64, clases []interface{}) (string, error) {
	for _, clase := range clases {
		if id == clase.(map[string]interface{})["Id"].(float64) {
			return clase.(map[string]interface{})["Nombre"].(string), nil
		}
	}
	return "", fmt.Errorf("not found")
}

func contarYSepararGrupos(cadena string) (int, []string) {
	// Dividir la cadena en Grupos utilizando la coma como separador
	grupos := strings.Split(cadena, ",")

	// Eliminar espacios en blanco alrededor de cada Grupos
	for i := 0; i < len(grupos); i++ {
		grupos[i] = strings.TrimSpace(grupos[i])
	}

	// Devolver la cantidad de Grupos y el slice de Grupos
	return len(grupos), grupos
}

func assignExistingPeriod(academicSpaces []interface{}, groups *[]string, periodIdReq int) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	spaceBody := map[string]interface{}{"periodo_id": periodIdReq}

	for _, space := range academicSpaces {
		spaceMap := space.(map[string]interface{})

		// unassigned period
		periodId := spaceMap["periodo_id"]
		if periodId == nil {
			validSpace, errValidation := validateGroup(groups, fmt.Sprintf("%v", spaceMap["grupo"]))
			if validSpace {
				// partial update
				if responseSpace, errSpace := utils.UpdateAcademicSpace(fmt.Sprintf("%v", spaceMap["_id"]), spaceBody); errSpace == nil {
					result = append(result, responseSpace)
				} else {
					return result, errValidation
				}
			} else if errValidation != nil {
				return result, errValidation
			}
		} else if reflect.TypeOf(periodId).Kind() == reflect.Int || reflect.TypeOf(periodId).Kind() == reflect.Float64 {
			if int(periodId.(float64)) == 0 {
				validSpace, errValidation := validateGroup(groups, fmt.Sprintf("%v", spaceMap["grupo"]))
				if validSpace {
					// partial update
					if responseSpace, errSpace := utils.UpdateAcademicSpace(fmt.Sprintf("%v", spaceMap["_id"]), spaceBody); errSpace != nil {
						result = append(result, responseSpace)
					} else {
						return result, errValidation
					}
				} else if errValidation != nil {
					return result, errValidation
				}
			}
		} else if reflect.TypeOf(periodId).Kind() == reflect.String {
			validSpace, errValidation := validateGroup(groups, fmt.Sprintf("%v", spaceMap["grupo"]))
			if validSpace {
				// partial update
				if responseSpace, errSpace := utils.UpdateAcademicSpace(fmt.Sprintf("%v", spaceMap["_id"]), spaceBody); errSpace != nil {
					result = append(result, responseSpace)
				} else {
					return result, errValidation
				}
			} else if errValidation != nil {
				return result, errValidation
			}
		}

		if len(*groups) == 0 {
			return result, nil
		}
	}
	return result, nil
}

func createAcademicSpaceChild(parent string, groups []string, periodIdReq int) ([]map[string]interface{}, error) {
	var newSpace map[string]interface{}
	var result []map[string]interface{}
	queryParams := "query=_id:" + fmt.Sprintf("%v", parent)
	urlAcademicSpaces := "http://" + beego.AppConfig.String("EspaciosAcademicosService") + "espacio-academico"

	if resSpaces, errSpace := utils.GetAcademicSpacesByQuery(queryParams); errSpace == nil {
		if space := resSpaces.([]any); space != nil {
			spaceBody := space[0].(map[string]any)
			spaceBody["espacio_academico_padre"] = spaceBody["_id"]
			delete(spaceBody, "_id")
			delete(spaceBody, "fecha_creacion")
			delete(spaceBody, "fecha_modificacion")

			for _, group := range groups {
				spaceBody["grupo"] = group
				spaceBody["periodo_id"] = periodIdReq
				if errNewSpace := helpers.SendJson(urlAcademicSpaces, "POST", &newSpace, spaceBody); errNewSpace == nil {
					result = append(result, newSpace["Data"].(map[string]interface{}))
				} else {
					return result, fmt.Errorf("EspaciosAcademicosService Error creando espacios académicos")
				}
			}
			return result, nil
		} else {
			return nil, fmt.Errorf("Espacio académico padre no encontrado")
		}
	} else {
		return nil, errSpace
	}
}

func validateGroup(groups *[]string, group string) (bool, error) {
	var errRemove error
	contains, idx := utils.ContainsStringIndex(*groups, group)
	if contains {
		*groups, errRemove = utils.RemoveIndexString(*groups, idx)
		if errRemove == nil {
			return true, nil
		} else {
			return false, errRemove
		}
	}
	return false, nil
}
