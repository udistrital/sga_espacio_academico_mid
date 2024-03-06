// @APIVersion 1.0.0
// @Title SGA MID - Espacios Académicos
// @Description Microservicio MID del SGA MID que complementa los espacios academicos
package routers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/sga_espacio_academico_mid/controllers"
	"github.com/udistrital/utils_oas/errorhandler"
)

func init() {

	beego.ErrorController(&errorhandler.ErrorHandlerController{})

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/espacios-academicos",
			beego.NSInclude(
				&controllers.EspaciosAcademicosController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
