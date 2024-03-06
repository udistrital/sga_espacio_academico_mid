# Espacios Académicos - SGA MID

API MID intermediaria entre el cliente SGA y la API de espacios académicos con los endpoints requeridos para la gestión de la información necesaria en los módulos del SGA cliente.

## Especificaciones Técnicas

### Tecnologías Implementadas y Versiones
* [Golang](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/golang.md)
* [BeeGo](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/beego.md)
* [Docker](https://docs.docker.com/engine/install/ubuntu/)
* [Docker Compose](https://docs.docker.com/compose/)

### Variables de Entorno
```shell
SGA_ESPACIO_ACADEMICO_MID_HTTP_PORT: [Puerto de ejecución API]
ESPACIOS_ACADEMICOS_SERVICE: ${Servicio API ESPACIOS ACADEMICOS}
PARAMETRO_SERVICE: [Servicio API PARAMETROS]
```
**NOTA:** Las variables se pueden ver en el fichero conf/app.conf.

### Ejecución del Proyecto
```shell
#1. Obtener el repositorio con Go
go get github.com/udistrital/sga_mid_espacios_academicos

#2. Moverse a la carpeta del repositorio
cd $GOPATH/src/github.com/udistrital/sga_mid_espacios_academicos

# 3. Moverse a la rama **develop**
git pull origin develop && git checkout develop

# 4. alimentar todas las variables de entorno que utiliza el proyecto.
SGA_MID_HTTP_PORT=8080 RULER_URL=127.0.0.1:27017 SGA_MID_SOME_VARIABLE=some_value bee run

# 5. Ejecutar comandos para descargar dependencias
go mod init && go mod tidy

# 6. Ejecutar proyecto
bee run -downdoc=true -gendoc=true
```

### Ejecución Dockerfile
```shell
# docker build --tag=sga_mid_espacios_academicos . --no-cache
# docker run -p 80:80 sga_mid_espacios_academicos
```

### Ejecución docker-compose
```shell
#1. Clonar el repositorio
git clone -b develop https://github.com/udistrital/sga_mid_espacios_academicos

#2. Moverse a la carpeta del repositorio
cd sga_mid_espacios_academicos

#3. Crear un fichero con el nombre **custom.env**
# En windows ejecutar:* ` ni custom.env`
touch custom.env

#4. Crear la network **back_end** para los contenedores
docker network create back_end

#5. Ejecutar el compose del contenedor
docker-compose up --build

#6. Comprobar que los contenedores estén en ejecución
docker ps
```

### Ejecución Pruebas

Pruebas unitarias
```shell
# En Proceso
```

## Estado CI

| Develop | Relese 0.0.1 | Master |
| -- | -- | -- |
| [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/sga_mid_espacios_academicos/status.svg?ref=refs/heads/develop)](https://hubci.portaloas.udistrital.edu.co/udistrital/sga_mid_espacios_academicos) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/sga_mid_espacios_academicos/status.svg?ref=refs/heads/release/0.0.1)](https://hubci.portaloas.udistrital.edu.co/udistrital/sga_mid_espacios_academicos) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/sga_mid_espacios_academicos/status.svg)](https://hubci.portaloas.udistrital.edu.co/udistrital/sga_mid_espacios_academicos) |

## Licencia

This file is part of sga_mid_espacios_academicos.

sga_mid_espacios_academicos is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

sga_mid is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with sga_mid. If not, see https://www.gnu.org/licenses/.
