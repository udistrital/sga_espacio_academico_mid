package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/sga_espacio_academico_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// EspaciosAcademicosController operations for Espacios_academicos
type EspaciosAcademicosController struct {
	beego.Controller
}

// URLMapping ...
func (c *EspaciosAcademicosController) URLMapping() {
	c.Mapping("GetAcademicSpacesByProject", c.GetAcademicSpacesByProject)
	c.Mapping("PostAcademicSpacesBySon", c.PostAcademicSpacesBySon)
	c.Mapping("PutAcademicSpaceAssignPeriod", c.PutAcademicSpaceAssignPeriod)
}

// GetAcademicSpacesByProject ...
// @Title GetAcademicSpacesByProject
// @Description get Espacios_academicos for Plan Estudios
// @Param	id_proyecto		path	int	true	"Id del proyecto"
// @Success 200 {}
// @Failure 404 not found resource
// @router /proyectos/:proyecto_id [get]
func (c *EspaciosAcademicosController) GetAcademicSpacesByProject() {
	defer errorhandler.HandlePanic(&c.Controller)
	/*
		check validez de id proyecto
	*/
	idProyectoStr := c.Ctx.Input.Param(":proyecto_id")
	idProyecto, errId := strconv.ParseInt(idProyectoStr, 10, 64)
	if errId != nil || idProyecto <= 0 {
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Error en el formato del id del proyecto")
		c.Ctx.Output.SetStatus(400)
	} else {
		resultado := services.GetAcademicSpacesByProject(idProyecto)
		c.Data["json"] = resultado
		c.Ctx.Output.SetStatus(resultado.Status)
	}
	c.ServeJSON()
}

// PostAcademicSpacesBySon ...
// @Title PostAcademicSpacesBySon
// @Description post EspaciosAcademicos for Plan Estudios
// @Param   body        body    {}  true        "body crear espacio academico content"
// @Success 200 {}
// @Failure 403 :body is empty
// @router /hijos [post]
func (c *EspaciosAcademicosController) PostAcademicSpacesBySon() {
	defer errorhandler.HandlePanic(&c.Controller)
	dataBody := c.Ctx.Input.RequestBody
	resultado := services.PostAcademicSpacesBySon(dataBody)
	c.Data["json"] = resultado
	c.Ctx.Output.SetStatus(resultado.Status)
	c.ServeJSON()
}

// PutAcademicSpaceAssignPeriod ...
// @Title PutAcademicSpaceAssignPeriod
// @Description Asigna el periodo a los grupos/espacios académicos indicados
// @Param   body        body    {}  true        "Asignar periodo a los espacios académicos"
// @Success 200 {}
// @Failure 400 the request contains incorrect syntaxis
// @router /hijos/asignar-periodo [put]
func (c *EspaciosAcademicosController) PutAcademicSpaceAssignPeriod() {
	defer errorhandler.HandlePanic(&c.Controller)
	/*
		{
			"grupo": ["Grupo 1", "Grupo 3"],
			"periodo_id": 36,
			"padre": "649cf98ecf8adba537ca9052"
		}
	*/
	dataBody := c.Ctx.Input.RequestBody
	resultado := services.PutAcademicSpaceAssignPeriod(dataBody)
	c.Data["json"] = resultado
	c.Ctx.Output.SetStatus(resultado.Status)
	c.ServeJSON()
}
