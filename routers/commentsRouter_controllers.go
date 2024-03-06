package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/sga_espacio_academico_mid/controllers:EspaciosAcademicosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_espacio_academico_mid/controllers:EspaciosAcademicosController"],
        beego.ControllerComments{
            Method: "PostAcademicSpacesBySon",
            Router: "/hijos",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_espacio_academico_mid/controllers:EspaciosAcademicosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_espacio_academico_mid/controllers:EspaciosAcademicosController"],
        beego.ControllerComments{
            Method: "PutAcademicSpaceAssignPeriod",
            Router: "/hijos/asignar-periodo",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_espacio_academico_mid/controllers:EspaciosAcademicosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_espacio_academico_mid/controllers:EspaciosAcademicosController"],
        beego.ControllerComments{
            Method: "GetAcademicSpacesByProject",
            Router: "/proyectos/:proyecto_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
