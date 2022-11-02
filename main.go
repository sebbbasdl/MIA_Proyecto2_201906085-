package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/rs/cors"
)

var arregloMountId [20]string
var arregloMountPath [20]string
var arregloMountPart [20]string
var contadorMount int = 0
var arregloletra [20]string
var arregloDiscos [20]string
var contadorDiscos1 int = 0
var contadorMount2 int = 0
var activa bool = false
var usuario_actual string
var path_actual string
var respuesta_exec string
var respuesta_rep string

type partition = struct {
	Part_status [100]byte
	Part_type   [100]byte
	Part_fit    [100]byte
	Part_start  [100]byte
	Part_size   [100]byte
	Part_name   [16]byte
}

type MBR = struct {
	Mbr_tamano         [100]byte
	Mbr_fecha_creacion [100]byte
	Mbr_disk_signature [100]byte
	Mbr_partition_1    partition
	Mbr_partition_2    partition
	Mbr_partition_3    partition
	Mbr_partition_4    partition
}

type ejemplo = struct {
	Id        [100]byte
	Nombre    [100]byte
	Direccion [100]byte
	Telefono  [100]byte
}

type cmdstruct struct {
	Cmd string `json:"cmd"`
}

func main() {
	fmt.Println("Proyecto 2 - MIA - 201906085 - Sebastian Alejandro de Leon Tenaz")

	mux := http.NewServeMux()

	mux.HandleFunc("/analizar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var Content cmdstruct
		respuesta := ""
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &Content)
		respuesta = split_comando(Content.Cmd)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "` + respuesta + `" }`))
	})

	mux.HandleFunc("/reportes", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-----SOY PATH ACTUAL " + path_sin_disco(path_actual))
		w.Header().Set("Content-Type", "application/json")
		bytes, _ := ioutil.ReadFile(path_sin_disco(path_actual) + "disk.png")
		var base64Encoding string
		base64Encoding += "data:image/jpg;base64,"
		base64Encoding += toBase64(bytes)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "` + base64Encoding + `" }`))
	})
	mux.HandleFunc("/reportes2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-----SOY PATH ACTUAL " + path_sin_disco(path_actual))
		w.Header().Set("Content-Type", "application/json")
		bytes, _ := ioutil.ReadFile(path_sin_disco(path_actual) + "mbr.png")
		var base64Encoding string
		base64Encoding += "data:image/jpg;base64,"
		base64Encoding += toBase64(bytes)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "` + base64Encoding + `" }`))
	})
	mux.HandleFunc("/reportes3", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-----SOY PATH ACTUAL " + path_sin_disco(path_actual))
		w.Header().Set("Content-Type", "application/json")
		bytes, _ := ioutil.ReadFile(path_sin_disco(path_actual) + "sb.png")
		var base64Encoding string
		base64Encoding += "data:image/jpg;base64,"
		base64Encoding += toBase64(bytes)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "` + base64Encoding + `" }`))

	})

	fmt.Println("Server ON in port 5000")
	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":5000", handler))

	//analizar()

	//exec -path="C:/Users/sebas/go/src/MIA_Proyecto2_201906085-/datoswindows.txt"
}

func msg_error(err error) {
	fmt.Println("Error: ", err)
}

func abecedario(conta int) string {
	var letra string
	if conta == 0 {
		letra = "a"

	} else if conta == 1 {
		letra = "b"

	} else if conta == 2 {
		letra = "c"

	} else if conta == 3 {
		letra = "d"

	} else if conta == 4 {
		letra = "e"

	} else if conta == 5 {
		letra = "f"

	} else if conta == 6 {
		letra = "g"

	} else if conta == 7 {
		letra = "h"

	} else if conta == 8 {
		letra = "i"

	} else if conta == 9 {
		letra = "j"

	} else if conta == 10 {
		letra = "k"

	} else if conta == 11 {
		letra = "l"

	} else if conta == 12 {
		letra = "m"

	} else if conta == 13 {
		letra = "n"

	} else if conta == 14 {
		letra = "ñ"

	} else if conta == 15 {
		letra = "o"

	} else if conta == 16 {
		letra = "p"

	} else if conta == 17 {
		letra = "q"

	} else if conta == 18 {
		letra = "r"

	}

	return letra
}
func existeE(mbr MBR) bool {
	//mbr1 := struct_to_bytes(mbr)
	var aux [100]byte
	copy(aux[:], "E")

	if mbr.Mbr_partition_1.Part_type == aux || mbr.Mbr_partition_2.Part_type == aux || mbr.Mbr_partition_3.Part_type == aux || mbr.Mbr_partition_4.Part_type == aux {
		return true
	}

	return false
}
func analizar() {
	finalizar := false
	fmt.Println("Proyecto 2 - MIA - 201906085 - Sebastian Alejandro de Leon Tenaz")
	reader := bufio.NewReader(os.Stdin)
	//  Ciclo para lectura de multiples comandos
	for !finalizar {
		fmt.Print("<Ingresar Comando>: ")
		comando, _ := reader.ReadString('\n')
		if strings.Contains(comando, "exit") {
			finalizar = true
		} else {
			if comando != "" && comando != "exit\n" {
				//  Separacion de comando y parametros
				split_comando(comando)
			}
		}
	}
}

func split_comando(comando string) string {
	var commandArray []string
	// Eliminacion de saltos de linea
	comando = strings.Replace(comando, "\n", "", 1)
	comando = strings.Replace(comando, "\r", "", 1)
	// Guardado de parametros
	if strings.Contains(comando, "mostrar") {
		commandArray = append(commandArray, comando)
	} else {
		commandArray = strings.Split(comando, " ")
	}
	// Ejecicion de comando leido
	return ejecucion_comando(commandArray)
}

func ejecucion_comando(commandArray []string) string {
	// Identificacion de comando y ejecucion
	respuesta := ""
	data := strings.ToLower(commandArray[0])
	if data == "mkdisk" {
		respuesta = crear_disco(commandArray)

	} else if data == "fdisk" {
		respuesta = crear_particion(commandArray)
	} else if data == "rmdisk" {
		respuesta = eliminar_disco(commandArray)
	} else if data == "mount" {
		//fmt.Println(commandArray)
		//fmt.Println(len(commandArray))
		if len(commandArray) == 1 {
			respuesta = onlymount()
		} else {
			respuesta = mount_particion(commandArray)
		}

	} else if data == "mkfs" {
		respuesta = mkfs(commandArray)
	} else if data == "login" {
		respuesta = login(commandArray)
	} else if data == "logout" {
		activa = false
		usuario_actual = ""
		path_actual = ""
		fmt.Println("Usted ha cerrado sesion")
		respuesta = "Usted ha cerrado sesion"

	} else if data == "mkgrp" {
		respuesta = mkgrp(commandArray)
	} else if data == "rmgrp" {
		respuesta = rmgrp(commandArray)
	} else if data == "mkusr" {
		respuesta = mkuser(commandArray)
	} else if data == "rmusr" {
		respuesta = rmusr(commandArray)
	} else if data == "mkfile" {
		mkfile(commandArray)
	} else if data == "mkdir" {
		mkdir(commandArray)
	} else if data == "mostrar" {
		//mostrar("C:/Users/sebas/go/src/Ejemplo7/disk.dk")
		//mostrar2()
		//crearDirectorioSiNoExiste("C:/Users/sebas/go/src/MIA_Proyecto2_201906085-/hola/")
	} else if data == "rep" {
		reportes(commandArray)
		respuesta = respuesta_rep
	} else if data == "exec" {
		respuesta = exece(commandArray)
	} else if data == "pause" {
		respuesta = "Press 'Enter' to continue..."
		pausee()
	} else {
		fmt.Println("Comando ingresado no es valido")
	}

	return respuesta
}

// crear_disco -tamaño=numero -dimensional=dimension/"dimension"
// mkdisk -size=15 -unit="m" -path="C:/Users/sebas/go/src/Ejemplo7/disk.dk"
func crear_disco(commandArray []string) string {
	respuesta := ""
	tamano := 0
	dimensional := ""
	tamano_archivo := 0
	limite := 0
	bloque := make([]byte, 1024)
	path_mkdisk := ""

	flag_size := false
	flag_path := false

	// Lectura de parametros del comando
	for i := 0; i < len(commandArray); i++ {
		data := commandArray[i]
		if strings.Contains(data, "-size=") {
			data := strings.ToLower(commandArray[i])
			strtam := strings.Replace(data, "-size=", "", 1)
			strtam = strings.Replace(strtam, "\"", "", 2)
			strtam = strings.Replace(strtam, "\r", "", 1)
			tamano2, err := strconv.Atoi(strtam)
			tamano = tamano2
			flag_size = true
			if err != nil {
				msg_error(err)
			}
		} else if strings.Contains(data, "-unit=") {
			data := strings.ToLower(commandArray[i])
			dimensional = strings.Replace(data, "-unit=", "", 1)
			dimensional = strings.Replace(dimensional, "\"", "", 2)

		} else if strings.Contains(data, "-path=") {
			flag_path = true
			path_mkdisk = strings.Replace(data, "-path=", "", 1)
			path_mkdisk = strings.Replace(path_mkdisk, "\"", "", 2)
		}
	}
	//crearDirectorioSiNoExiste(path_sin_disco(path_mkdisk))
	// Calculo de tamaño del archivo
	if strings.Contains(dimensional, "k") {
		tamano_archivo = tamano
	} else if strings.Contains(dimensional, "m") {
		tamano_archivo = tamano * 1024
	}

	// Preparacion del bloque a crear_particion en archivo
	for j := 0; j < 1024; j++ {
		bloque[j] = 0
	}

	// Creacion, escritura y cierre de archivo
	//crearDirectorioSiNoExiste(path_mkdisk)

	if flag_path == true && flag_size == true {
		disco, err := os.Create(path_mkdisk)

		if err != nil {
			msg_error(err)
			fmt.Println("h0")
		}
		for limite < tamano_archivo {

			_, err := disco.Write(bloque)
			if err != nil {
				fmt.Println("h01")
				msg_error(err)
			}
			limite++
		}

		t := time.Now()
		fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())
		mbr := MBR{}
		//ntel, _ := strconv.ParseInt(tamano_archivo, 10, 32)
		copy(mbr.Mbr_tamano[:], strconv.Itoa(tamano_archivo))
		copy(mbr.Mbr_fecha_creacion[:], fecha)
		copy(mbr.Mbr_disk_signature[:], strconv.Itoa(rand.Intn(20)))

		copy(mbr.Mbr_partition_1.Part_status[:], "e")
		copy(mbr.Mbr_partition_1.Part_fit[:], "-")
		copy(mbr.Mbr_partition_1.Part_name[:], "0")
		copy(mbr.Mbr_partition_1.Part_size[:], "0")
		copy(mbr.Mbr_partition_1.Part_start[:], "-1")
		copy(mbr.Mbr_partition_1.Part_type[:], "-")

		copy(mbr.Mbr_partition_2.Part_status[:], "e")
		copy(mbr.Mbr_partition_2.Part_fit[:], "-")
		copy(mbr.Mbr_partition_2.Part_name[:], "0")
		copy(mbr.Mbr_partition_2.Part_size[:], "0")
		copy(mbr.Mbr_partition_2.Part_start[:], "-1")
		copy(mbr.Mbr_partition_2.Part_type[:], "-")

		copy(mbr.Mbr_partition_3.Part_status[:], "e")
		copy(mbr.Mbr_partition_3.Part_fit[:], "-")
		copy(mbr.Mbr_partition_3.Part_name[:], "0")
		copy(mbr.Mbr_partition_3.Part_size[:], "0")
		copy(mbr.Mbr_partition_3.Part_start[:], "-1")
		copy(mbr.Mbr_partition_3.Part_type[:], "-")

		copy(mbr.Mbr_partition_4.Part_status[:], "e")
		copy(mbr.Mbr_partition_4.Part_fit[:], "-")
		copy(mbr.Mbr_partition_4.Part_name[:], "0")
		copy(mbr.Mbr_partition_4.Part_size[:], "0")
		copy(mbr.Mbr_partition_4.Part_start[:], "-1")
		copy(mbr.Mbr_partition_4.Part_type[:], "-")

		mbrbyte := struct_to_bytes(mbr)
		newpos, err := disco.Seek(int64(1*len(mbrbyte)), os.SEEK_SET)
		fmt.Println(newpos)
		if err != nil {
			msg_error(err)
		}

		_, err = disco.WriteAt(mbrbyte, newpos)
		if err != nil {
			msg_error(err)
		}

		mbrstruc := bytes_to_struct(mbrbyte)
		fmt.Println(string(mbrstruc.Mbr_tamano[:]))
		fmt.Println(string(mbrstruc.Mbr_fecha_creacion[:]))
		fmt.Println(string(mbrstruc.Mbr_disk_signature[:]))
		fmt.Println(string(mbrstruc.Mbr_partition_1.Part_fit[:]))
		fmt.Println(string(mbrstruc.Mbr_partition_4.Part_name[:]))

		respuesta = "Disco creado :    SIZE: " + strconv.Itoa(tamano) + "  PATH: " + path_mkdisk + "  UNIT: " + dimensional + "\\n"
		respuesta_exec += "Disco creado :    SIZE: " + strconv.Itoa(tamano) + "  PATH: " + path_mkdisk + "  UNIT: " + dimensional + "\\n"

		disco.Close()

	} else {
		respuesta = "¡Error, no tiene path o size en mkdisk!" + "\\n"
		respuesta_exec += "¡Error, no tiene path o size en mkdisk!" + "\\n"
		fmt.Println("¡Error, no tiene path o size en mkdisk!")
	}

	// Resumen de accion realizada
	/*fmt.Print("Creacion de Disco:")
	fmt.Print(" Tamaño: ")
	fmt.Print(tamano)
	fmt.Print(" Dimensional: ")
	fmt.Println(dimensional)
	fmt.Print(" Path: ")
	fmt.Println(path_mkdisk)*/
	/*path_aux := ""
	for i := 0; i < len(path_mkdisk); i++ {
		if path_mkdisk[i]== {

		}
	}*/

	return respuesta
}

// crear_particion -nombre=nombre/"nombre" -direccion=direccion/"direccion" -telefono=numero -veces=numero
// fdisk -size=1 -unit=m -path="C:\Users\sebas\go\src\Ejemplo7\disk.dk" -tipo=p -fit=ff -name=hola1
// fdisk -size=4 -unit=m -path="C:\Users\sebas\go\src\Ejemplo7\disk.dk" -tipo=p -fit=ff -name=hola2
// fdisk -size=5 -unit=m -path="C:\Users\sebas\go\src\Ejemplo7\disk.dk" -tipo=p -fit=ff -name=hola3
func crear_particion(commandArray []string) string {
	respuesta := ""
	straux := ""
	flagf := true
	size := 0
	unit := ""
	path := ""
	tipo := ""
	fit := ""
	name := ""
	// Lectura de parametros del comando
	for i := 0; i < len(commandArray); i++ {
		data := strings.ToLower(commandArray[i])
		data2 := commandArray[i]
		if strings.Contains(data, "-size=") {
			straux = strings.Replace(data, "-size=", "", 1)
			straux = strings.Replace(straux, "\"", "", 2)
			straux = strings.Replace(straux, "\r", "", 1)
			size2, err := strconv.Atoi(straux)
			size = size2
			if err != nil {
				msg_error(err)
			}
		} else if strings.Contains(data, "-unit=") {
			straux = strings.Replace(data, "-unit=", "", 1)
			straux = strings.Replace(straux, "\"", "", 2)
			straux = strings.Replace(straux, "\r", "", 1)
			unit = straux
		} else if strings.Contains(data2, "-path=") {
			straux = strings.Replace(data2, "-path=", "", 1)
			straux = strings.Replace(straux, "\"", "", 2)
			straux = strings.Replace(straux, "\r", "", 1)
			path = straux
		} else if strings.Contains(data, "-tipo=") {
			straux = strings.Replace(data, "-tipo=", "", 1)
			straux = strings.Replace(straux, "\"", "", 2)
			straux = strings.Replace(straux, "\r", "", 1)
			tipo = straux
		} else if strings.Contains(data, "-fit=") {
			straux = strings.Replace(data, "-fit=", "", 1)
			straux = strings.Replace(straux, "\"", "", 2)
			straux = strings.Replace(straux, "\r", "", 1)
			fit = straux
		} else if strings.Contains(data, "-name=") {
			straux = strings.Replace(data, "-name=", "", 1)
			straux = strings.Replace(straux, "\"", "", 2)
			straux = strings.Replace(straux, "\r", "", 1)
			name = straux
		}
	}

	/*	// Apertura del archivo
		disco, err := os.OpenFile("Ejemplo7.dk", os.O_RDWR, 0660)
		if err != nil {
			msg_error(err)
		}

		// Escritura en el archivo utilizando SEEK_SET
		nnombre := ""
		ndir := ""
		for k := 0; k < veces; k++ {
			index := k + 1
			nnombre = string(nombreejemplo) + " " + strconv.Itoa(index)
			ndir = string(direjemplo) + " " + strconv.Itoa(index)
			ntel, _ := strconv.ParseInt(telejemplo, 10, 32)
			ntel = ntel + int64(index)*int64(index)
			copy(ejm.Id[:], strconv.Itoa(index))
			copy(ejm.Nombre[:], nnombre)
			copy(ejm.Direccion[:], ndir)
			copy(ejm.Telefono[:], strconv.Itoa(int(ntel)))
			// Conversion de struct a bytes
			ejmbyte := struct_to_bytes(ejm)
			// Cambio de posicion de puntero dentro del archivo
			newpos, err := disco.Seek(int64(k*len(ejmbyte)), os.SEEK_SET)
			if err != nil {
				msg_error(err)
			}
			// Escritura de struct en archivo binario
			_, err = disco.WriteAt(ejmbyte, newpos)
			if err != nil {
				msg_error(err)
			}
		}
		disco.Close()*/

	disco, err := os.OpenFile(path, os.O_RDWR, 0660)
	if err != nil {
		msg_error(err)
	}
	// Calculo del tamano de struct en bytes
	mbr := MBR{}
	ejm2 := struct_to_bytes(mbr)
	sstruct := len(ejm2)
	fmt.Println(sstruct)

	// Lectrura de conjunto de bytes en archivo binario
	lectura := make([]byte, sstruct)
	_, err = disco.ReadAt(lectura, int64(sstruct))
	if err != nil && err != io.EOF {
		msg_error(err)
	}

	// Conversion de bytes a struct
	ejm := bytes_to_struct(lectura)
	sstruct = len(lectura)

	if err != nil {
		msg_error(err)
	}

	if unit == "m" {
		size = size * 1024
		fmt.Println(size)
	} else if unit == "k" {
		size = size * 1024
	} else if unit == "b" {

	}

	var aux [100]byte
	copy(aux[:], "e")
	if tipo == "p" || tipo == "e" {
		if ejm.Mbr_partition_1.Part_status != aux && ejm.Mbr_partition_2.Part_status != aux && ejm.Mbr_partition_3.Part_status != aux && ejm.Mbr_partition_4.Part_status != aux {
			fmt.Println("Todas las particiones estan ocupadas")
			respuesta = "Todas las particiones estan ocupadas" + "\\n"
			respuesta_exec += "Todas las particiones estan ocupadas" + "\\n"
			flagf = false
		}
	}

	if tipo == "e" && flagf == true {
		if existeE(ejm) == true {
			fmt.Println("Ya existe una particion extendida en " + name)
			flagf = false
			respuesta = "Ya existe una particion extendida en " + name + "\\n"
			respuesta_exec += "Ya existe una particion extendida en " + name + "\\n"

		}
	} else if tipo == "l" && flagf == true {
		if existeE(ejm) == false {
			fmt.Println("No existe un particion extendida en " + name + " entonces no se puede crear la particion logica.")
			flagf = false
			respuesta = "No existe un particion extendida en " + name + " entonces no se puede crear la particion logica." + "\\n"
			respuesta_exec += "No existe un particion extendida en " + name + " entonces no se puede crear la particion logica." + "\\n"

		}

	}

	//para poder obtener tamaño mbr
	flagnum := false
	cont := 0
	aux_tamano := ""
	for !flagnum {
		if string(ejm.Mbr_tamano[cont]) == "0" || string(ejm.Mbr_tamano[cont]) == "1" || string(ejm.Mbr_tamano[cont]) == "2" || string(ejm.Mbr_tamano[cont]) == "3" || string(ejm.Mbr_tamano[cont]) == "4" || string(ejm.Mbr_tamano[cont]) == "5" || string(ejm.Mbr_tamano[cont]) == "6" || string(ejm.Mbr_tamano[cont]) == "7" || string(ejm.Mbr_tamano[cont]) == "8" || string(ejm.Mbr_tamano[cont]) == "9" {
			fmt.Println(string(ejm.Mbr_tamano[cont]))
			aux_tamano += string(ejm.Mbr_tamano[cont])
			cont++
		} else {
			flagnum = true
		}

	}
	println(flagf)
	if tipo == "l" && flagf == true {

	} else if flagf == true {

		var espaciototal string = aux_tamano
		fmt.Print("sstruct")
		fmt.Println(sstruct)
		espaciototal2, err := strconv.Atoi(espaciototal)
		espaciototal2 = espaciototal2 - sstruct

		if err != nil {
			msg_error(err)
		}

		var arreglo_espacios [4]int
		var esp1 string = "0"
		var esp2 string = "0"
		var esp3 string = "0"
		var esp4 string = "0"

		copy(aux[:], "e")

		if ejm.Mbr_partition_1.Part_status == aux && ejm.Mbr_partition_2.Part_status == aux && ejm.Mbr_partition_3.Part_status == aux && ejm.Mbr_partition_4.Part_status == aux {
			//cout<<"entre 0"<<endl;
			if size <= espaciototal2 {
				//cout<<"entre"<<endl;
				//fmt.Println(ejm.Mbr_tamano)
				//
				copy(ejm.Mbr_partition_1.Part_status[:], "0")
				copy(ejm.Mbr_partition_1.Part_fit[:], fit)
				copy(ejm.Mbr_partition_1.Part_name[:], name)
				copy(ejm.Mbr_partition_1.Part_size[:], strconv.Itoa(size))
				copy(ejm.Mbr_partition_1.Part_start[:], strconv.Itoa(espaciototal2))
				copy(ejm.Mbr_partition_1.Part_type[:], tipo)

				fmt.Println(ejm.Mbr_partition_1.Part_status)
				auxejm := struct_to_bytes(ejm)
				disco.Seek(0, 0)
				_, err = disco.WriteAt(auxejm, int64(sstruct))
				if err != nil {
					msg_error(err)
				}
				fmt.Println(size)
				fmt.Print("START:")
				fmt.Println(espaciototal2)
				//rewind(file)
				//fwrite(&mbr, sizeof(MBR), 1, file)
				respuesta = "Particion relealizada:  " + name + "\\n"
				respuesta_exec += "Particion relealizada:  " + name + "\\n"
			}
		} else {
			fmt.Println(string(ejm.Mbr_partition_1.Part_status[:]))
			fmt.Println(string(aux[:]))
			if ejm.Mbr_partition_1.Part_status != aux {
				esp1 = byte_to_string(ejm)
				//fmt.Println("entre")
				//fmt.Println(esp1)
			}

			if ejm.Mbr_partition_2.Part_status != aux {
				esp2 = byte_to_string2(ejm)

			}

			if ejm.Mbr_partition_3.Part_status != aux {
				esp3 = byte_to_string3(ejm)

			}

			if ejm.Mbr_partition_4.Part_status != aux {
				esp4 = byte_to_string4(ejm)

			}
			fmt.Println("aqui")
			esp11, err1 := strconv.Atoi(esp1)
			fmt.Println(esp11)
			if err1 != nil {
				fmt.Println("aqui3")
				msg_error(err)
			}
			esp22, err2 := strconv.Atoi(esp2)
			if err2 != nil {
				fmt.Println("aqui3")
				msg_error(err)
			}
			esp33, err3 := strconv.Atoi(esp3)
			if err3 != nil {
				fmt.Println("aqui3")
				msg_error(err)
			}

			esp44, err4 := strconv.Atoi(esp4)
			if err4 != nil {
				fmt.Println("aqui3")
				msg_error(err)
			}
			fmt.Println("aqui2")
			fmt.Println(espaciototal2)
			fmt.Println(strconv.Itoa(esp11) + "," + strconv.Itoa(esp22) + "," + strconv.Itoa(esp33) + "," + strconv.Itoa(esp44))

			espaciototal2 = espaciototal2 - esp11 - esp22 - esp33 - esp44 + 2899

			fmt.Println(size)
			fmt.Println(espaciototal2)

			if size <= espaciototal2 {
				arreglo_espacios[0] = esp11
				arreglo_espacios[1] = esp22
				arreglo_espacios[2] = esp33
				arreglo_espacios[3] = esp44

				//var espaciodispo int
				var arreglolibre [4]int
				fmt.Println("muere")
				//fmt.Println(arreglo_start)
				fmt.Println(arreglo_espacios)
				/*for i := 0; i < 4; i++ {
					aux_tamano1, err := strconv.Atoi(aux_tamano)
					if err != nil {
						msg_error(err)
					}
					fmt.Println("muere11")

					arreglo_start[i] = arreglo_start[i-1] + arreglo_espacios[i-1]
					espaciodispo = aux_tamano1 - arreglo_start[i]

					fmt.Println("muere22")
					for j := i; j < 4; j++ {
						if arreglo_espacios[j] != 0 {
							espaciodispo = arreglo_start[j] - arreglo_start[i]
							break

						}

					}

					arreglolibre[i] = espaciodispo

				}*/
				fmt.Println("muere2")
				fmt.Println(string(ejm.Mbr_partition_1.Part_size[:]))
				fmt.Println(espaciototal2)

				fmt.Println("---arreglo LIBRE---")
				if esp11 != 0 {
					fmt.Println("estoy")
					arreglolibre[0] = 0
				} else {
					fmt.Println("esto1y")
					arreglolibre[0] = espaciototal2
				}

				if esp22 != 0 {

					arreglolibre[1] = 0
				} else {
					arreglolibre[1] = espaciototal2
				}

				if esp33 != 0 {

					arreglolibre[2] = 0
				} else {
					arreglolibre[2] = espaciototal2
				}

				if esp44 != 0 {

					arreglolibre[3] = 0
				} else {
					arreglolibre[3] = espaciototal2
				}

				fmt.Println(arreglolibre)

				fmt.Println(size)
				if fit == "ff" || fit == "bf" || fit == "wf" {
					fmt.Println("estoy aca jeje0")
					for i := 0; i < 4; i++ {
						if size <= arreglolibre[i] {
							fmt.Println("estoy aca jeje")
							if i == 0 {
								copy(ejm.Mbr_partition_1.Part_status[:], "0")
								copy(ejm.Mbr_partition_1.Part_fit[:], fit)
								copy(ejm.Mbr_partition_1.Part_name[:], name)
								copy(ejm.Mbr_partition_1.Part_size[:], strconv.Itoa(size))
								copy(ejm.Mbr_partition_1.Part_start[:], strconv.Itoa(espaciototal2))
								copy(ejm.Mbr_partition_1.Part_type[:], tipo)
								fmt.Println("entreee")
								//std::cout<<mbr.mbr_partition_1.part_status<<endl;

								auxejm := struct_to_bytes(ejm)
								disco.Seek(0, 0)
								_, err = disco.WriteAt(auxejm, int64(sstruct))
								if err != nil {
									msg_error(err)
								}
								break

							} else if i == 1 {
								copy(ejm.Mbr_partition_2.Part_status[:], "0")
								copy(ejm.Mbr_partition_2.Part_fit[:], fit)
								copy(ejm.Mbr_partition_2.Part_name[:], name)
								copy(ejm.Mbr_partition_2.Part_size[:], strconv.Itoa(size))
								copy(ejm.Mbr_partition_2.Part_start[:], strconv.Itoa(espaciototal2))
								copy(ejm.Mbr_partition_2.Part_type[:], tipo)

								//std::cout<<mbr.mbr_partition_1.part_status<<endl;
								fmt.Println("entreee")
								auxejm := struct_to_bytes(ejm)
								disco.Seek(0, 0)
								_, err = disco.WriteAt(auxejm, int64(sstruct))
								if err != nil {
									msg_error(err)
								}
								break

							} else if i == 2 {
								copy(ejm.Mbr_partition_3.Part_status[:], "0")
								copy(ejm.Mbr_partition_3.Part_fit[:], fit)
								copy(ejm.Mbr_partition_3.Part_name[:], name)
								copy(ejm.Mbr_partition_3.Part_size[:], strconv.Itoa(size))
								copy(ejm.Mbr_partition_3.Part_start[:], strconv.Itoa(espaciototal2))
								copy(ejm.Mbr_partition_3.Part_type[:], tipo)

								//std::cout<<mbr.mbr_partition_1.part_status<<endl;

								auxejm := struct_to_bytes(ejm)
								disco.Seek(0, 0)
								_, err = disco.WriteAt(auxejm, int64(sstruct))
								if err != nil {
									msg_error(err)
								}
								break

							} else if i == 3 {
								copy(ejm.Mbr_partition_4.Part_status[:], "0")
								copy(ejm.Mbr_partition_4.Part_fit[:], fit)
								copy(ejm.Mbr_partition_4.Part_name[:], name)
								copy(ejm.Mbr_partition_4.Part_size[:], strconv.Itoa(size))
								copy(ejm.Mbr_partition_4.Part_start[:], strconv.Itoa(espaciototal2))
								copy(ejm.Mbr_partition_4.Part_type[:], tipo)

								//std::cout<<mbr.mbr_partition_1.part_status<<endl;

								auxejm := struct_to_bytes(ejm)
								disco.Seek(0, 0)
								_, err = disco.WriteAt(auxejm, int64(sstruct))
								if err != nil {
									msg_error(err)
								}
								break

							}
						}
					}

				}
			}
			respuesta = "Particion relealizada:  " + name + "\\n"
			respuesta_exec += "Particion relealizada:  " + name + "\\n"
		}
	}

	disco.Close()

	// Resumen de accion realizada
	fmt.Println("Escritura en Disco de Struct con los siguientes datos :")
	fmt.Println(" size: ")
	fmt.Println(size)
	fmt.Println(" unit: ")
	fmt.Println(unit)
	fmt.Println(" path: ")
	fmt.Println(path)
	fmt.Println(" tipo: ")
	fmt.Println(tipo)
	fmt.Println(" fit: ")
	fmt.Println(fit)
	fmt.Println(" name: ")
	fmt.Println(name)

	return respuesta

}

// rmdisk -path="C:\Users\sebas\go\src\Ejemplo7\disk.dk"
func eliminar_disco(commandArray []string) string {
	respuesta := ""
	path_rmdisk := ""

	flag_path := false
	for i := 0; i < len(commandArray); i++ {
		data := commandArray[i]
		if strings.Contains(data, "-path=") {
			flag_path = true
			path_rmdisk = strings.Replace(data, "-path=", "", 1)
			path_rmdisk = strings.Replace(path_rmdisk, "\"", "", 2)
		}
	}

	if flag_path == true {

		fmt.Println(path_rmdisk)
		err := os.Remove(path_rmdisk)
		if err != nil {
			fmt.Printf("Error eliminando archivo: %v\n", err)
			respuesta = "Error eliminando archivo" + "\\n"
			respuesta_exec += "Error eliminando archivo" + "\\n"
		} else {
			fmt.Println("Eliminado correctamente")
			respuesta = "Eliminado correctamente"
			respuesta_exec += "Eliminado correctamente" + "\\n"
		}
	} else {
		fmt.Println("¡Error, no es posible ejecutar el comando rmdisk ya que falta el dato de path!")
		respuesta = "¡Error, no es posible ejecutar el comando rmdisk ya que falta el dato de path!" + "\\n"
		respuesta_exec += "¡Error, no es posible ejecutar el comando rmdisk ya que falta el dato de path!" + "\\n"
	}
	return respuesta
}

// mount -path="C:\Users\sebas\go\src\Ejemplo7\disk.dk" -name=hola1
func mount_particion(commandArray []string) string {
	respuesta := ""
	path_mount := ""
	name := ""

	flag_path := false
	flag_name := false
	for i := 0; i < len(commandArray); i++ {
		data := commandArray[i]
		if strings.Contains(data, "-path=") {
			flag_path = true
			path_mount = strings.Replace(data, "-path=", "", 1)
			path_mount = strings.Replace(path_mount, "\"", "", 2)
		} else if strings.Contains(data, "-name=") {
			flag_name = true
			name = strings.Replace(data, "-name=", "", 1)
			name = strings.Replace(name, "\"", "", 2)
		}
	}

	if flag_path == true && flag_name == true {
		fmt.Println("mount")
		fmt.Println(path_mount)
		fmt.Println(name)
		//asignacion mount
		//var auxx [2]string
		//auxx[0] = name
		//auxx[1] = path_mount

		var auxx0 [16]byte
		copy(auxx0[:], name)

		var auxx1 [100]byte
		copy(auxx1[:], path_mount)

		var nombredisco string
		poss := 0
		poss2 := 0

		for i := 0; i < len(path_mount); i++ {
			if path_mount[i] == '/' {
				poss = i
			}

			if path_mount[i] == '.' {
				poss2 = i
			}
		}

		for i := poss + 1; i < poss2; i++ {
			nombredisco += string(path_mount[i])
		}
		fmt.Println(nombredisco)
		fmt.Println("path_mount: ")
		fmt.Println(path_mount)

		disco, err := os.OpenFile(path_mount, os.O_RDWR, 0660)
		if err != nil {
			fmt.Println("es aca")
			msg_error(err)
		}

		// Calculo del tamano de struct en bytes
		mbr := MBR{}
		ejm2 := struct_to_bytes(mbr)
		sstruct := len(ejm2)
		fmt.Println(sstruct)

		// Lectrura de conjunto de bytes en archivo binario
		lectura := make([]byte, sstruct)
		_, err = disco.ReadAt(lectura, int64(sstruct))
		if err != nil && err != io.EOF {
			fmt.Println("es aca2")
			msg_error(err)
		}

		// Conversion de bytes a struct
		ejm := bytes_to_struct(lectura)
		sstruct = len(lectura)

		if err != nil {
			msg_error(err)
		}

		//disco.Seek(0, 0)
		if contadorDiscos1 == 0 {
			if ejm.Mbr_partition_1.Part_name == auxx0 {

				//cout<<"entre"<<endl;
				copy(ejm.Mbr_partition_1.Part_status[:], "1")
				//cout<<"-------------------"<<contadorDiscos1<<endl;
				arregloMountId[contadorMount] = "85" + strconv.Itoa(contadorMount+1) + nombredisco
				arregloMountPart[contadorMount] = string(auxx0[:])
				arregloMountPath[contadorMount] = path_mount

				contadorMount += 1
				arregloletra[contadorDiscos1] = "a"
				contadorDiscos1 += 1
				auxejm := struct_to_bytes(ejm)
				disco.Seek(0, 0)
				_, err = disco.WriteAt(auxejm, int64(sstruct))
				if err != nil {
					msg_error(err)
				}

			} else if ejm.Mbr_partition_2.Part_name == auxx0 {
				copy(ejm.Mbr_partition_2.Part_status[:], "1")
				arregloMountId[contadorMount] = "85" + strconv.Itoa(contadorMount+1) + nombredisco
				arregloMountPart[contadorMount] = string(auxx0[:])
				arregloMountPath[contadorMount] = path_mount
				contadorMount += 1
				arregloletra[contadorDiscos1] = "a"
				contadorDiscos1 += 1
				auxejm := struct_to_bytes(ejm)
				disco.Seek(0, 0)
				_, err = disco.WriteAt(auxejm, int64(sstruct))
				if err != nil {
					msg_error(err)
				}

			} else if ejm.Mbr_partition_3.Part_name == auxx0 {
				copy(ejm.Mbr_partition_3.Part_status[:], "1")
				arregloMountId[contadorMount] = "85" + strconv.Itoa(contadorMount+1) + nombredisco
				arregloMountPart[contadorMount] = string(auxx0[:])
				arregloMountPath[contadorMount] = path_mount
				contadorMount += 1
				arregloletra[contadorDiscos1] = "a"
				contadorDiscos1 += 1
				auxejm := struct_to_bytes(ejm)
				disco.Seek(0, 0)
				_, err = disco.WriteAt(auxejm, int64(sstruct))
				if err != nil {
					msg_error(err)
				}

			} else if ejm.Mbr_partition_4.Part_name == auxx0 {
				copy(ejm.Mbr_partition_4.Part_status[:], "1")
				arregloMountId[contadorMount] = "85" + strconv.Itoa(contadorMount+1) + nombredisco
				arregloMountPart[contadorMount] = string(auxx0[:])
				arregloMountPath[contadorMount] = path_mount
				contadorMount += 1
				arregloletra[contadorDiscos1] = "a"
				contadorDiscos1 += 1
				auxejm := struct_to_bytes(ejm)
				disco.Seek(0, 0)
				_, err = disco.WriteAt(auxejm, int64(sstruct))
				if err != nil {
					msg_error(err)
				}

			} else {
				fmt.Println("Esta particion no existe en el disco con la ruta: " + string(auxx1[:]))
				respuesta = "Esta particion no existe en el disco con la ruta: " + string(auxx1[:]) + "\\n"
				respuesta_exec += "Esta particion no existe en el disco con la ruta: " + string(auxx1[:]) + "\\n"
			}
			respuesta += name + "ha sido montada." + "\\n"
			respuesta_exec += name + "ha sido montada." + "\\n"
			//cout << "------STATUS----" << endl

		} else {
			var contt int = 1
			var auxlugar int = 0

			for i := 0; i < len(arregloMountPath); i++ {
				if path_mount == arregloMountPath[i] {
					contt += 1
					fmt.Println("cont: " + strconv.Itoa(contt))
					auxlugar = i
					fmt.Println("lugar: " + strconv.Itoa(auxlugar))
				}
			}

			for i := 0; i < len(arregloletra); i++ {
				fmt.Println(arregloletra[i])
			}

			for i := 0; i < len(arregloMountPath); i++ {
				if path_mount == arregloMountPath[i] || contt > 1 {
					//cout << "mori2" << endl
					if ejm.Mbr_partition_1.Part_name == auxx0 {
						//cout << "mori3" << endl
						copy(ejm.Mbr_partition_1.Part_status[:], "1")

						arregloMountId[contadorMount] = "85" + strconv.Itoa(contadorMount+1) + nombredisco
						arregloMountPart[contadorMount] = string(auxx0[:])
						arregloMountPath[contadorMount] = path_mount
						arregloletra[contadorMount] = arregloletra[auxlugar]
						contadorMount += 1

						auxejm := struct_to_bytes(ejm)
						disco.Seek(0, 0)
						_, err = disco.WriteAt(auxejm, int64(sstruct))
						if err != nil {
							msg_error(err)
						}
						break

					} else if ejm.Mbr_partition_2.Part_name == auxx0 {
						//cout << "mori3" << endl
						copy(ejm.Mbr_partition_2.Part_status[:], "1")

						arregloMountId[contadorMount] = "85" + strconv.Itoa(contadorMount+1) + nombredisco
						arregloMountPart[contadorMount] = string(auxx0[:])
						arregloMountPath[contadorMount] = path_mount
						arregloletra[contadorMount] = arregloletra[auxlugar]
						contadorMount += 1

						auxejm := struct_to_bytes(ejm)
						disco.Seek(0, 0)
						_, err = disco.WriteAt(auxejm, int64(sstruct))
						if err != nil {
							msg_error(err)
						}
						break

					} else if ejm.Mbr_partition_3.Part_name == auxx0 {
						//cout << "mori3" << endl
						copy(ejm.Mbr_partition_3.Part_status[:], "1")

						arregloMountId[contadorMount] = "85" + strconv.Itoa(contadorMount+1) + nombredisco
						arregloMountPart[contadorMount] = string(auxx0[:])
						arregloMountPath[contadorMount] = path_mount
						arregloletra[contadorMount] = arregloletra[auxlugar]
						contadorMount += 1

						auxejm := struct_to_bytes(ejm)
						disco.Seek(0, 0)
						_, err = disco.WriteAt(auxejm, int64(sstruct))
						if err != nil {
							msg_error(err)
						}
						break

					} else if ejm.Mbr_partition_4.Part_name == auxx0 {
						//cout << "mori3" << endl
						copy(ejm.Mbr_partition_4.Part_status[:], "1")

						arregloMountId[contadorMount] = "85" + strconv.Itoa(contadorMount+1) + nombredisco
						arregloMountPart[contadorMount] = string(auxx0[:])
						arregloMountPath[contadorMount] = path_mount
						arregloletra[contadorMount] = arregloletra[auxlugar]
						contadorMount += 1

						auxejm := struct_to_bytes(ejm)
						disco.Seek(0, 0)
						_, err = disco.WriteAt(auxejm, int64(sstruct))
						if err != nil {
							msg_error(err)
						}
						break

					}

				} else {

					if ejm.Mbr_partition_1.Part_name == auxx0 {
						//cout << "CONTADOR DISCOSSSS : " << contadorDiscos1 << endl

						copy(ejm.Mbr_partition_1.Part_status[:], "1")

						arregloMountId[contadorMount] = "85" + strconv.Itoa(contadorMount+1) + nombredisco
						arregloMountPart[contadorMount] = string(auxx0[:])
						arregloMountPath[contadorMount] = path_mount

						arregloletra[contadorMount] = abecedario(contadorDiscos1)
						contadorMount += 1
						contadorDiscos1 += 1
						auxejm := struct_to_bytes(ejm)
						disco.Seek(0, 0)
						_, err = disco.WriteAt(auxejm, int64(sstruct))
						if err != nil {
							msg_error(err)
						}

						contadorDiscos1 += 1
						break
					} else if ejm.Mbr_partition_2.Part_name == auxx0 {
						copy(ejm.Mbr_partition_2.Part_status[:], "1")

						arregloMountId[contadorMount] = "85" + strconv.Itoa(contadorMount+1) + nombredisco
						arregloMountPart[contadorMount] = string(auxx0[:])
						arregloMountPath[contadorMount] = path_mount

						arregloletra[contadorMount] = abecedario(contadorDiscos1)
						contadorMount += 1
						contadorDiscos1 += 1
						auxejm := struct_to_bytes(ejm)
						disco.Seek(0, 0)
						_, err = disco.WriteAt(auxejm, int64(sstruct))
						if err != nil {
							msg_error(err)
						}

						contadorDiscos1 += 1
						break
					} else if ejm.Mbr_partition_3.Part_name == auxx0 {
						copy(ejm.Mbr_partition_3.Part_status[:], "1")

						arregloMountId[contadorMount] = "85" + strconv.Itoa(contadorMount+1) + nombredisco
						arregloMountPart[contadorMount] = string(auxx0[:])
						arregloMountPath[contadorMount] = path_mount

						arregloletra[contadorMount] = abecedario(contadorDiscos1)
						contadorMount += 1
						contadorDiscos1 += 1
						auxejm := struct_to_bytes(ejm)
						disco.Seek(0, 0)
						_, err = disco.WriteAt(auxejm, int64(sstruct))
						if err != nil {
							msg_error(err)
						}

						contadorDiscos1 += 1
						break
					} else if ejm.Mbr_partition_4.Part_name == auxx0 {
						copy(ejm.Mbr_partition_4.Part_status[:], "1")

						arregloMountId[contadorMount] = "85" + strconv.Itoa(contadorMount+1) + nombredisco
						arregloMountPart[contadorMount] = string(auxx0[:])
						arregloMountPath[contadorMount] = path_mount

						arregloletra[contadorMount] = abecedario(contadorDiscos1)
						contadorMount += 1
						contadorDiscos1 += 1
						auxejm := struct_to_bytes(ejm)
						disco.Seek(0, 0)
						_, err = disco.WriteAt(auxejm, int64(sstruct))
						if err != nil {
							msg_error(err)
						}

						contadorDiscos1 += 1
						break
					}

				}
			}
			respuesta = name + "ha sido montada." + "\\n"
			respuesta_exec += name + "ha sido montada." + "\\n"
		}
		fmt.Println(string(ejm.Mbr_partition_1.Part_status[:]))
		fmt.Println(string(ejm.Mbr_partition_2.Part_status[:]))
		fmt.Println(string(ejm.Mbr_partition_3.Part_status[:]))
		fmt.Println(string(ejm.Mbr_partition_4.Part_status[:]))
		disco.Close()

		/*disco1, err := os.OpenFile(path_mount, os.O_RDWR, 0660)
		if err != nil {
			fmt.Println("hoy si")
			msg_error(err)
		}

		// Calculo del tamano de struct en bytes
		mbr1 := MBR{}
		ejm21 := struct_to_bytes(mbr1)
		sstruct1 := len(ejm21)
		fmt.Println(sstruct1)

		// Lectrura de conjunto de bytes en archivo binario
		lectura1 := make([]byte, sstruct1)
		_, err = disco.ReadAt(lectura1, int64(sstruct1))
		if err != nil && err != io.EOF {
			msg_error(err)
			fmt.Println("hoy si2")
		}

		// Conversion de bytes a struct
		ejm1 := bytes_to_struct(lectura1)
		sstruct1 = len(lectura1)

		if err != nil {
			msg_error(err)
		}

		fmt.Println(string(ejm1.Mbr_partition_1.Part_status[:]))
		fmt.Println(string(ejm1.Mbr_partition_2.Part_status[:]))
		fmt.Println(string(ejm1.Mbr_partition_3.Part_status[:]))
		fmt.Println(string(ejm1.Mbr_partition_4.Part_status[:]))

		disco1.Close()*/
	} else {
		respuesta = "¡Error, no es posible ejecutar el comando mount hizo falta un dato!" + "\\n"
		respuesta_exec += "¡Error, no es posible ejecutar el comando mount hizo falta un dato!" + "\\n"
		fmt.Println("¡Error, no es posible ejecutar el comando mount hizo falta un dato!")
	}
	return respuesta
}

func onlymount() string {
	respuesta := ""
	//cout << contadorMount << endl
	for i := 0; i < contadorMount; i++ {
		fmt.Println(strconv.Itoa(i) + " " + arregloMountId[i] + " " + arregloMountPath[i] + " " + arregloMountPart[i])
		respuesta = strconv.Itoa(i) + " " + arregloMountId[i] + "\\n"
		respuesta_exec += strconv.Itoa(i) + " " + arregloMountId[i] + "\\n"

	}
	return respuesta

}

// mkfs -id=851disk
func mkfs(commandArray []string) string {
	respuesta := ""
	id := ""
	tipo := ""

	flag_id := false
	for i := 0; i < len(commandArray); i++ {
		data := commandArray[i]
		if strings.Contains(data, "-id=") {
			flag_id = true
			id = strings.Replace(data, "-id=", "", 1)
			id = strings.Replace(id, "\"", "", 2)
		} else if strings.Contains(data, "-type=") {

			tipo = strings.Replace(data, "-type=", "", 1)
			tipo = strings.Replace(tipo, "\"", "", 2)
		}
	}

	if flag_id == true {
		fmt.Println("mkfs")
		fmt.Println(id)
		fmt.Println(tipo)

		var path string
		var path2 string
		var name string
		var flag_mkfs bool = false
		var poss int
		name = name + ""
		poss = poss + 0
		for i := 0; i < contadorMount; i++ {
			if id == arregloMountId[i] {
				path = arregloMountPath[i]
				name = arregloMountPart[i]
				poss = i

				flag_mkfs = true
				break
			}
		}

		if flag_mkfs == true {
			path2 = path
			var auxpath string = "users.txt"
			path2 = path_sin_disco(path2)
			path2 += auxpath
			fmt.Println(path2)
			crearArchivo_UyG(path2)
			respuesta = "mkfs realizado correctamente" + "\\n"
			respuesta_exec += "mkfs realizado correctamente" + "\\n"
		}

	} else {
		fmt.Println("¡Error, no es posible ejecutar el comando mkfs")
		respuesta = "¡Error, no es posible ejecutar el comando mkfs " + "\\n"
		respuesta_exec += "¡Error, no es posible ejecutar el comando mkfs " + "\\n"
	}
	return respuesta
}

//login -usuario=root -password=123 -id=851disk

func login(commandArray []string) string {
	respuesta := ""
	id := ""
	usuario := ""
	password := ""

	flag_id := false
	flag_user := false
	flag_password := false
	for i := 0; i < len(commandArray); i++ {
		data := commandArray[i]
		if strings.Contains(data, "-id=") {
			flag_id = true
			id = strings.Replace(data, "-id=", "", 1)
			id = strings.Replace(id, "\"", "", 2)

		} else if strings.Contains(data, "-usuario=") {
			flag_user = true
			usuario = strings.Replace(data, "-usuario=", "", 1)
			usuario = strings.Replace(usuario, "\"", "", 2)

		} else if strings.Contains(data, "-password=") {
			flag_password = true
			password = strings.Replace(data, "-password=", "", 1)
			password = strings.Replace(password, "\"", "", 2)
		}
	}

	if flag_id == true && flag_password == true && flag_user == true {
		fmt.Println("login")
		fmt.Println(id)
		fmt.Println(usuario)
		fmt.Println(password)

		var path string
		//int contadorlog=0;
		var poss int
		var datos [5]string
		var cont_datos int
		var aux_dato string
		var flag bool = false
		poss = poss + 0
		for i := 0; i < contadorMount; i++ {
			if id == arregloMountId[i] {
				path = arregloMountPath[i]
				poss = i

				flag = true
				break
			}
		}

		path = path_sin_disco(path)
		pathenviar := path
		path += "users.txt"

		if activa == false && flag == true {

			file, err := os.OpenFile(path, os.O_RDWR, 0660)
			if err != nil {
				fmt.Println("es aca")
				msg_error(err)
			}
			var texto string

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				//contadorlog+=1;
				texto = scanner.Text()
				//cout<<contadorlog<<texto<<endl;
				//cout<<texto.length()<<endl;
				//cout<<texto[8]<<endl;
				cont_datos = 0

				for i := 0; i < len(texto); i++ {
					if texto[i] == ',' {
						fmt.Println(aux_dato)
						datos[cont_datos] = aux_dato
						cont_datos += 1
						aux_dato = ""

					} else {
						if i == len(texto)-1 {
							aux_dato += string(texto[i])
							fmt.Println(aux_dato)
							datos[cont_datos] = aux_dato
							cont_datos += 1
							aux_dato = ""
						} else {
							aux_dato += string(texto[i])
						}

					}

				}
				//cout<<cont_datos<<endl;

				if cont_datos == 5 {
					if datos[3] == usuario && datos[4] == password {
						activa = true
						fmt.Println("Inicio de sesion correcto, bienvenido: " + usuario)
						respuesta = pathenviar
						respuesta_exec += "Inicio de sesion correcto, bienvenido: " + usuario + "\\n"
						path_actual = path
						if usuario == "root" {
							fmt.Println("Usuario root activo")
							respuesta += "\\n Usuario root activo"
							respuesta_exec += "\\n" + "Usuario root activo" + "\\n"
						}

						usuario_actual = usuario

						break
					} else {
						respuesta += "Usuario o contraseña son incorrectos\n"
						respuesta_exec += "Usuario o contraseña son incorrectos" + "\\n"
						fmt.Println("Usuario o contraseña son incorrectos")
					}
				}

			}
			file.Close()

		} else {
			respuesta += "No se puede iniciar sesion ya que hay un usuario activo o el id no existe\\n"
			respuesta_exec += "No se puede iniciar sesion ya que hay un usuario activo o el id no existe" + "\\n"
			fmt.Println("No se puede iniciar sesion ya que hay un usuario activo o el id no existe")
		}

	} else {
		respuesta += "¡Error, no es posible ejecutar el comando Login ya que falta algun dato " + "\\n"
		respuesta_exec += "¡Error, no es posible ejecutar el comando Login ya que falta algun dato " + "\\n"
		fmt.Println("¡Error, no es posible ejecutar el comando Login ya que falta algun dato ")
	}
	return respuesta
}

func mkgrp(commandArray []string) string {
	respuesta := ""
	name := ""

	flag_name := false
	for i := 0; i < len(commandArray); i++ {
		data := commandArray[i]
		if strings.Contains(data, "-name=") {
			flag_name = true
			name = strings.Replace(data, "-name=", "", 1)
			name = strings.Replace(name, "\"", "", 2)
		}
	}

	if flag_name == true {
		fmt.Println("mkgrp")
		fmt.Println(name)

		var datos [5]string
		var cont_datos int
		var aux_dato string
		var id_grupo int
		var flag bool = false
		path := path_actual

		if activa == true && usuario_actual == "root" {

			file, err := os.OpenFile(path, os.O_RDWR, 0660)
			if err != nil {
				fmt.Println("es aca")
				msg_error(err)
			}
			var texto string
			var aux_texto string

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				//contadorlog+=1;
				texto = scanner.Text()
				//cout<<contadorlog<<texto<<endl;
				//cout<<texto.length()<<endl;
				//cout<<texto[8]<<endl;
				cont_datos = 0
				aux_texto += texto + "\n"

				for i := 0; i < len(texto); i++ {
					if texto[i] == ',' {
						fmt.Println(aux_dato)
						datos[cont_datos] = aux_dato
						cont_datos += 1
						aux_dato = ""

					} else {
						if i == len(texto)-1 {
							aux_dato += string(texto[i])
							fmt.Println(aux_dato)
							datos[cont_datos] = aux_dato
							cont_datos += 1
							aux_dato = ""
						} else {
							aux_dato += string(texto[i])
						}

					}

				}
				//cout<<cont_datos<<endl;

				if cont_datos == 3 {
					if datos[2] == name {
						fmt.Println("Este grupo ya existe")
						flag = false
					} else {
						flag = true
						//auxgrp := string(datos[0])
						id_grupo, err = strconv.Atoi(datos[0])
						if err != nil {
							//	fmt.Println("es aca")
							msg_error(err)
						}
						//cout<<"xxxxxxxxx "<<stoi(datos[0])+1<<endl;
					}

				}

			}

			if flag == true {
				respuesta = "Creacion de grupo correcta" + name

				_, err = file.WriteString("\n" + strconv.Itoa(id_grupo+1) + ",G," + name)
				if err != nil {
					//	fmt.Println("es aca")
					msg_error(err)
				}

			}

			file.Close()

		} else {
			respuesta += "No se puede crear un grupo ya que la sesion no es del usuario root" + "\\n"
			respuesta_exec += "No se puede crear un grupo ya que la sesion no es del usuario root" + "\\n"
			fmt.Println("No se puede crear un grupo ya que la sesion no es del usuario root")
		}

	} else {
		respuesta += "¡Error, no es posible ejecutar el comando rmdisk ya que falta el dato de path!" + "\\n"
		respuesta_exec += "¡Error, no es posible ejecutar el comando rmdisk ya que falta el dato de path!" + "\\n"
		fmt.Println("¡Error, no es posible ejecutar el comando rmdisk ya que falta el dato de path!")
	}
	return respuesta
}

//rmgrp -name=grupo1
func rmgrp(commandArray []string) string {
	respuesta := ""
	name := ""

	flag_name := false
	for i := 0; i < len(commandArray); i++ {
		data := commandArray[i]
		if strings.Contains(data, "-name=") {
			flag_name = true
			name = strings.Replace(data, "-name=", "", 1)
			name = strings.Replace(name, "\"", "", 2)
		}
	}

	if flag_name == true {
		fmt.Println("mkgrp")
		fmt.Println(name)

		var datos [5]string
		var cont_datos int
		var aux_dato string
		var id_grupo int
		id_grupo = id_grupo + 0
		var flag bool = false
		path := path_actual
		if activa == true && usuario_actual == "root" {

			file, err := os.OpenFile(path, os.O_RDWR, 0660)
			if err != nil {
				fmt.Println("es aca")
				msg_error(err)
			}
			var texto string
			var aux_texto string

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				//contadorlog+=1;
				texto = scanner.Text()
				//cout<<contadorlog<<texto<<endl;
				//cout<<texto.length()<<endl;
				//cout<<texto[8]<<endl;
				cont_datos = 0

				for i := 0; i < len(texto); i++ {
					if texto[i] == ',' {
						fmt.Println(aux_dato)
						datos[cont_datos] = aux_dato
						cont_datos += 1
						aux_dato = ""

					} else {
						if i == len(texto)-1 {
							aux_dato += string(texto[i])
							fmt.Println(aux_dato)
							datos[cont_datos] = aux_dato
							cont_datos += 1
							aux_dato = ""
						} else {
							aux_dato += string(texto[i])
						}

					}

				}
				//cout<<cont_datos<<endl;

				if cont_datos == 3 {
					if datos[2] == name {
						fmt.Println("Existe el grupo")
						texto = ""
						for i := 0; i < cont_datos; i++ {
							datos[0] = "0"
							if i == cont_datos-1 {
								texto += datos[i]
							} else {
								texto += datos[i] + ","
							}

						}

						flag = true
					} else {
						flag = false
						id_grupo, err = strconv.Atoi(datos[0])
						if err != nil {
							//	fmt.Println("es aca")
							msg_error(err)
						}
						//cout<<"xxxxxxxxx "<<stoi(datos[0])+1<<endl;
					}

				}
				aux_texto += texto + "\n"

			}
			file.Close()

			file2, err := os.Create(path)

			if err != nil {
				msg_error(err)
				fmt.Println("h0")
			}

			if flag == true {

				_, err = file2.WriteString(aux_texto)
				if err != nil {
					//	fmt.Println("es aca")
					msg_error(err)
				}
				respuesta += "Grupo eliminado" + name + "\n"

			}
			file2.Close()

		} else {
			respuesta += "No se puede crear un grupo ya que la sesion no es del usuario root\n"
			respuesta_exec += "No se puede crear un grupo ya que la sesion no es del usuario root\n"
			fmt.Println("No se puede crear un grupo ya que la sesion no es del usuario root")
		}
	} else {
		respuesta += "¡Error, no es posible ejecutar el comando rmgrp!\\n"
		respuesta_exec += "¡Error, no es posible ejecutar el comando rmgrp!" + "\\n"
		fmt.Println("¡Error, no es posible ejecutar el comando rmgrp!")
	}
	return respuesta
}

//mkusr -usuario=usuario1 -pwd=321 -grp=grupo1

func mkuser(commandArray []string) string {
	respuesta := ""
	grp := ""
	usuario := ""
	password := ""

	flag_grp := false
	flag_user := false
	flag_password := false
	for i := 0; i < len(commandArray); i++ {
		data := commandArray[i]
		if strings.Contains(data, "-grp=") {
			flag_grp = true
			grp = strings.Replace(data, "-grp=", "", 1)
			grp = strings.Replace(grp, "\"", "", 2)

		} else if strings.Contains(data, "-usuario=") {
			flag_user = true
			usuario = strings.Replace(data, "-usuario=", "", 1)
			usuario = strings.Replace(usuario, "\"", "", 2)

		} else if strings.Contains(data, "-pwd=") {
			flag_password = true
			password = strings.Replace(data, "-pwd=", "", 1)
			password = strings.Replace(password, "\"", "", 2)
		}
	}

	if flag_grp == true && flag_password == true && flag_user == true {
		fmt.Println("login")
		fmt.Println(grp)
		fmt.Println(usuario)
		fmt.Println(password)

		var datos [5]string
		var cont_datos int
		var aux_dato string
		var id_grupo int
		var cont_grupos int = 0
		var flaguser bool = false
		var flaggroup bool = false
		var arr_grupos [30]string
		if activa == true && usuario_actual == "root" {

			file, err := os.OpenFile(path_actual, os.O_RDWR, 0660)
			if err != nil {
				fmt.Println("es aca")
				msg_error(err)
			}
			var texto string
			var aux_texto string

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				//contadorlog+=1;
				texto = scanner.Text()
				//cout<<contadorlog<<texto<<endl;
				//cout<<texto.length()<<endl;
				//cout<<texto[8]<<endl;
				cont_datos = 0
				aux_texto += texto + "\n"

				for i := 0; i < len(texto); i++ {
					if texto[i] == ',' {
						fmt.Println(aux_dato)
						datos[cont_datos] = aux_dato
						cont_datos += 1
						aux_dato = ""

					} else {
						if i == len(texto)-1 {
							aux_dato += string(texto[i])
							fmt.Println(aux_dato)
							datos[cont_datos] = aux_dato
							cont_datos += 1
							aux_dato = ""
						} else {
							aux_dato += string(texto[i])

						}

					}

				}
				fmt.Println("---" + strconv.Itoa(cont_datos))

				if cont_datos == 3 {
					//cout<<"entre"<<cont_grupos<<endl;
					arr_grupos[cont_grupos] = datos[2]
					//cout<<"entre"<<endl;
					cont_grupos += 1
				}

				if cont_datos == 5 {
					if datos[3] == usuario {
						flaguser = false
						fmt.Println("Ya existe el usuario: " + usuario + " por lo tanto no se podrá crear.")
						respuesta += "Ya existe el usuario: " + usuario + " por lo tanto no se podrá crear." + "\\n"
						respuesta_exec += "Ya existe el usuario: " + usuario + " por lo tanto no se podrá crear." + "\\n"
					} else {
						id_grupo, err = strconv.Atoi(datos[0])
						if err != nil {
							//	fmt.Println("es aca")
							msg_error(err)
						}
						flaguser = true

					}
				}
				/*if(cont_datos==2){
					if(datos[2]==name){
						cout<<"Este grupo ya existe"<<endl;
						flag=false;
					}else{
						flag=true;
						id_grupo=stoi(datos[0]);
						cout<<"xxxxxxxxx "<<stoi(datos[0])+1<<endl;
					}



				}*/

			}

			for i := 0; i < cont_grupos; i++ {
				if grp == arr_grupos[i] {
					flaggroup = true
					break
				} else {
					flaggroup = false
				}
			}

			if flaggroup == false {
				fmt.Println("El grupo no existe: " + grp)
				respuesta += "El grupo no existe: " + grp + "\\n"
				respuesta_exec += "El grupo no existe: " + grp + "\\n"
			}

			if flaguser == true && flaggroup == true {

				_, err = file.WriteString("\n" + strconv.Itoa(id_grupo+1) + ",U," + grp + "," + usuario + "," + password)
				if err != nil {
					//	fmt.Println("es aca")
					msg_error(err)
				}
				respuesta += "Creacion correcta de" + usuario + "\\n"
				respuesta_exec += "Creacion correcta de" + usuario + "\\n"

			}

			file.Close()

		} else {
			respuesta += "No se puede crear usuario ya que no ha iniciado sesion en el usuario root!!!" + "\\n"
			respuesta_exec += "No se puede crear usuario ya que no ha iniciado sesion en el usuario root!!!" + "\\n"
			fmt.Println("No se puede crear usuario ya que no ha iniciado sesion en el usuario root!!!")
		}
	} else {
		respuesta += "¡Error, no es posible ejecutar el comando mkusr ya que falta el dato de path!" + "\\n"
		respuesta_exec += "¡Error, no es posible ejecutar el comando mkusr ya que falta el dato de path!" + "\\n"
		fmt.Println("¡Error, no es posible ejecutar el comando mkusr ya que falta el dato de path!")
	}
	return respuesta
}

func rmusr(commandArray []string) string {
	respuesta := ""
	usuario := ""

	flag_usuario := false
	for i := 0; i < len(commandArray); i++ {
		data := commandArray[i]
		if strings.Contains(data, "-usuario=") {
			flag_usuario = true
			usuario = strings.Replace(data, "-usuario=", "", 1)
			usuario = strings.Replace(usuario, "\"", "", 2)
		}
	}

	if flag_usuario == true {
		fmt.Println("rmusr")
		fmt.Println(usuario)

		var datos [5]string
		var cont_datos int
		var aux_dato string
		var id_grupo int
		id_grupo = id_grupo + 0
		var flag bool = false
		path := path_actual
		if activa == true && usuario_actual == "root" {

			file, err := os.OpenFile(path, os.O_RDWR, 0660)
			if err != nil {
				fmt.Println("es aca")
				msg_error(err)
			}
			var texto string
			var aux_texto string

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				//contadorlog+=1;
				texto = scanner.Text()
				//cout<<contadorlog<<texto<<endl;
				//cout<<texto.length()<<endl;
				//cout<<texto[8]<<endl;
				cont_datos = 0

				for i := 0; i < len(texto); i++ {
					if texto[i] == ',' {
						fmt.Println(aux_dato)
						datos[cont_datos] = aux_dato
						cont_datos += 1
						aux_dato = ""

					} else {
						if i == len(texto)-1 {
							aux_dato += string(texto[i])
							fmt.Println(aux_dato)
							datos[cont_datos] = aux_dato
							cont_datos += 1
							aux_dato = ""
						} else {
							aux_dato += string(texto[i])
						}

					}

				}
				//cout<<cont_datos<<endl;

				if cont_datos == 5 {
					if datos[3] == usuario {
						fmt.Println("Existe el usuario")
						texto = ""
						for i := 0; i < cont_datos; i++ {
							datos[0] = "0"
							if i == cont_datos-1 {
								texto += datos[i]
							} else {
								texto += datos[i] + ","
							}

						}

						flag = true
					} else {
						flag = false
						id_grupo, err = strconv.Atoi(datos[0])
						if err != nil {
							//	fmt.Println("es aca")
							msg_error(err)
						}
						//cout<<"xxxxxxxxx "<<stoi(datos[0])+1<<endl;
					}

				}
				aux_texto += texto + "\n"

			}
			file.Close()

			file2, err := os.Create(path)

			if err != nil {
				msg_error(err)
				fmt.Println("h0")
			}

			if flag == true {

				_, err = file2.WriteString(aux_texto)
				if err != nil {
					//	fmt.Println("es aca")
					msg_error(err)
				}
				respuesta += "Usuario eliminado" + usuario + "\n"

			}
			file2.Close()

		} else {
			respuesta += "No se puede crear un grupo ya que la sesion no es del usuario root\\n"
			respuesta_exec += "No se puede crear un grupo ya que la sesion no es del usuario root" + "\\n"
			fmt.Println("No se puede crear un grupo ya que la sesion no es del usuario root")
		}
	} else {
		respuesta = "¡Error, no es posible ejecutar el comando rmusr!\\n"
		respuesta_exec = "¡Error, no es posible ejecutar el comando rmusr!" + "\\n"
		fmt.Println("¡Error, no es posible ejecutar el comando rmusr!")
	}

	return respuesta
}

func mkfile(commandArray []string) {

	path := ""
	r := ""
	size := ""
	cont := ""
	flag_path := false

	for i := 0; i < len(commandArray); i++ {
		data := commandArray[i]
		if strings.Contains(data, "-path=") {
			flag_path = true
			path = strings.Replace(data, "-path=", "", 1)
			path = strings.Replace(path, "\"", "", 2)

		} else if strings.Contains(data, "-r=") {

			r = strings.Replace(data, "-r=", "", 1)
			r = strings.Replace(r, "\"", "", 2)

		} else if strings.Contains(data, "-size=") {

			size = strings.Replace(data, "-size=", "", 1)
			size = strings.Replace(size, "\"", "", 2)

		} else if strings.Contains(data, "-cont=") {

			cont = strings.Replace(data, "-cont=", "", 1)
			cont = strings.Replace(size, "\"", "", 2)
		}
	}

	if flag_path == true {
		fmt.Println("mkfile")
		fmt.Println(path)
		fmt.Println(r)
		fmt.Println(size)
		fmt.Println(cont)
	} else {
		fmt.Println("¡Error, no es posible ejecutar el comando rmdisk ya que falta el dato de path!")
	}
}

func mkdir(commandArray []string) {

	path := ""
	p := ""

	flag_path := false

	for i := 0; i < len(commandArray); i++ {
		data := commandArray[i]
		if strings.Contains(data, "-path=") {
			flag_path = true
			path = strings.Replace(data, "-path=", "", 1)
			path = strings.Replace(path, "\"", "", 2)

		} else if strings.Contains(data, "-r=") {

			p = strings.Replace(data, "-r=", "", 1)
			p = strings.Replace(p, "\"", "", 2)

		}
	}

	if flag_path == true {
		fmt.Println("mkfile")
		fmt.Println(path)
		fmt.Println(p)

	} else {
		fmt.Println("¡Error, no es posible ejecutar el comando rmdisk ya que falta el dato de path!")
	}
}

func reportes(commandArray []string) {
	respuesta_rep = ""

	name := ""
	path := ""
	id := ""

	flag_id := false
	flag_path := false
	flag_name := false
	for i := 0; i < len(commandArray); i++ {
		data := commandArray[i]
		if strings.Contains(data, "-id=") {
			flag_id = true
			id = strings.Replace(data, "-id=", "", 1)
			id = strings.Replace(id, "\"", "", 2)

		} else if strings.Contains(data, "-path=") {
			flag_path = true
			path = strings.Replace(data, "-path=", "", 1)
			path = strings.Replace(path, "\"", "", 2)

		} else if strings.Contains(data, "-name=") {
			flag_name = true
			name = strings.Replace(data, "-name=", "", 1)
			name = strings.Replace(name, "\"", "", 2)
		}
	}

	if flag_id == true && flag_name == true && flag_path == true {
		//crearDirectorioSiNoExiste(path_sin_disco(path))
		fmt.Println("Rep:")
		fmt.Println(id)
		fmt.Println(path)
		fmt.Println(name)
		var flag bool = false
		lugar := 0
		auxcont := 0

		lugar = lugar + 0
		auxcont = auxcont + 0

		for x := 0; x < contadorMount; x++ {
			if id == /*arregloMountId[x]*/ arregloMountId[x] {

				//if(a_rep[0]==mbr.mbr_partition_1.part_name || a_rep[0]==mbr.mbr_partition_2.part_name || a_rep[0]==mbr.mbr_partition_3.part_name || a_rep[0]==mbr.mbr_partition_4.part_name) {
				flag = true
				lugar = x
				break
				//}
				//auxcont += 1

			}
		}

		disco, err := os.OpenFile(arregloMountPath[lugar], os.O_RDWR, 0660)
		if err != nil {
			msg_error(err)
		}
		disco.Seek(0, 0)

		mbr := MBR{}

		ejm2 := struct_to_bytes(mbr)
		sstruct := len(ejm2)
		fmt.Println(sstruct)

		// Lectrura de conjunto de bytes en archivo binario
		lectura := make([]byte, sstruct)
		_, err = disco.ReadAt(lectura, int64(sstruct))
		if err != nil && err != io.EOF {
			msg_error(err)
		}

		// Conversion de bytes a struct
		ejm := bytes_to_struct(lectura)
		sstruct = len(lectura)

		if err != nil {
			msg_error(err)
		}

		//byte_string(string(ejm.Mbr_tamano[:]))

		if flag == true && name == "mbr" {
			var info string = ""
			info += "digraph H{\n node [shape=plaintext];\n B [ label=< <TABLE BORDER =\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n<TR PORT=\"header\"><TD COLSPAN=\"2\">MBR</TD></TR>\n"
			info += "<TR><TD PORT=\"1\">MBR_TAMANO</TD><TD> " + byte_string(ejm.Mbr_tamano) + " bytes</TD></TR>\n"
			info += "<TR><TD PORT=\"2\">MBR_FECHA_CREACION</TD><TD>" + byte_string(ejm.Mbr_fecha_creacion) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"3\">MBR_DISK_SIGNATURE</TD><TD>" + byte_string(ejm.Mbr_disk_signature) + "</TD></TR>\n"

			info += "<TR><TD  bgcolor='red' PORT=\"4\">PARTITION</TD><TD bgcolor='red' >1</TD></TR>"
			info += "<TR><TD PORT=\"5\">PARTITION1_NAME</TD><TD>" + byte_string2(ejm.Mbr_partition_1.Part_name) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"6\">PARTITION1_SIZE</TD><TD>" + byte_string(ejm.Mbr_partition_1.Part_size) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"7\">PARTITION1_START</TD><TD>" + byte_string(ejm.Mbr_partition_1.Part_start) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"8\">PARTITION1_STATUS</TD><TD>" + byte_string(ejm.Mbr_partition_1.Part_status) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"9\">PARTITION1_TYPE</TD><TD>" + byte_string(ejm.Mbr_partition_1.Part_type) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"10\">PARTITION1_FIT</TD><TD>" + byte_string(ejm.Mbr_partition_1.Part_fit) + "</TD></TR>\n"

			info += "<TR><TD  bgcolor='red' PORT=\"11\">PARTITION</TD><TD bgcolor='red' >2</TD></TR>\n"
			info += "<TR><TD PORT=\"12\">PARTITION2_NAME</TD><TD>" + byte_string2(ejm.Mbr_partition_2.Part_name) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"13\">PARTITION2_SIZE</TD><TD>" + byte_string(ejm.Mbr_partition_2.Part_size) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"14\">PARTITION2_START</TD><TD>" + byte_string(ejm.Mbr_partition_2.Part_start) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"15\">PARTITION2_STATUS</TD><TD>" + byte_string(ejm.Mbr_partition_2.Part_status) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"16\">PARTITION2_TYPE</TD><TD>" + byte_string(ejm.Mbr_partition_2.Part_type) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"17\">PARTITION2_FIT</TD><TD>" + byte_string(ejm.Mbr_partition_2.Part_fit) + "</TD></TR>\n"

			info += "<TR><TD  bgcolor='red' PORT=\"18\">PARTITION</TD><TD bgcolor='red' >3</TD></TR>\n"
			info += "<TR><TD PORT=\"19\">PARTITION3_NAME</TD><TD>" + byte_string2(ejm.Mbr_partition_3.Part_name) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"20\">PARTITION3_SIZE</TD><TD>" + byte_string(ejm.Mbr_partition_3.Part_size) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"21\">PARTITION3_START</TD><TD>" + byte_string(ejm.Mbr_partition_3.Part_start) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"22\">PARTITION3_STATUS</TD><TD>" + byte_string(ejm.Mbr_partition_3.Part_status) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"23\">PARTITION3_TYPE</TD><TD>" + byte_string(ejm.Mbr_partition_3.Part_type) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"24\">PARTITION3_FIT</TD><TD>" + byte_string(ejm.Mbr_partition_3.Part_fit) + "</TD></TR>\n"

			info += "<TR><TD  bgcolor='red' PORT=\"25\">PARTITION</TD><TD bgcolor='red' >4</TD></TR>\n"
			info += "<TR><TD PORT=\"26\">PARTITION4_NAME</TD><TD>" + byte_string2(ejm.Mbr_partition_4.Part_name) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"27\">PARTITION4_SIZE</TD><TD>" + byte_string(ejm.Mbr_partition_4.Part_size) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"28\">PARTITION4_START</TD><TD>" + byte_string(ejm.Mbr_partition_4.Part_start) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"29\">PARTITION4_STATUS</TD><TD>" + byte_string(ejm.Mbr_partition_4.Part_status) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"30\">PARTITION4_TYPE</TD><TD>" + byte_string(ejm.Mbr_partition_4.Part_type) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"31\">PARTITION4_FIT</TD><TD>" + byte_string(ejm.Mbr_partition_4.Part_fit) + "</TD></TR>\n"
			info += "</TABLE> >];}\n"
			disco.Close()

			f, err := os.Create(path_sin_disco(path_actual) + "mbr.txt")
			defer f.Close()
			if err != nil {
				fmt.Println(">> Error drawing graph!")
			}
			file, err := os.Open(path_sin_disco(path_actual) + "mbr.txt")

			defer file.Close()
			if err != nil {
				fmt.Println(">> Error reading the file. Try again.")
				return
			}

			f.WriteString(info)

			e := exec.Command("dot", "-Tpng", "mbr.txt", "-o", "mbr.png")
			if er := e.Run(); er != nil {
				fmt.Println(">> Error", er)
				return
			}
			fmt.Println("Reporte MBR realizado")
			respuesta_rep = "Reporte MBR realizado\\n"
			respuesta_exec += "Reporte MBR realizado\\n"

		} else if flag == true && name == "disk" {
			text := ""
			text += "digraph H {\nparent [shape=plaintext \nlabel=<\n<table border='1' cellborder='1'>\n <tr>\n"

			size1, err := strconv.Atoi(byte_string(ejm.Mbr_partition_1.Part_size))
			if err != nil && err != io.EOF {
				msg_error(err)
			}
			size2, err := strconv.Atoi(byte_string(ejm.Mbr_partition_2.Part_size))
			if err != nil && err != io.EOF {
				msg_error(err)
			}
			size3, err := strconv.Atoi(byte_string(ejm.Mbr_partition_3.Part_size))
			if err != nil && err != io.EOF {
				msg_error(err)
			}
			size4, err := strconv.Atoi(byte_string(ejm.Mbr_partition_4.Part_size))
			if err != nil && err != io.EOF {
				msg_error(err)
			}

			tamanio, err := strconv.Atoi(byte_string(ejm.Mbr_tamano))
			if err != nil && err != io.EOF {
				msg_error(err)
			}

			/*var xd int64 = 10
			var xd2 int64 = 99

			var valor float64 = float64(float64(xd) / float64(xd2))
			fmt.Println("---")
			fmt.Println(tamanio)
			fmt.Println(size1)
			fmt.Println(size2)
			fmt.Println(size3)
			fmt.Println(size4)
			fmt.Println(valor)*/
			s1 := fmt.Sprintf("%.1f", ((float64(size1) / float64(tamanio)) * float64(100)))
			s2 := fmt.Sprintf("%.1f", ((float64(size2) / float64(tamanio)) * float64(100)))
			s3 := fmt.Sprintf("%.1f", ((float64(size3) / float64(tamanio)) * float64(100)))
			s4 := fmt.Sprintf("%.1f", ((float64(size4) / float64(tamanio)) * float64(100)))
			if size1 == 0 {
				text += "<td port='port_one'>LIBRE</td>\n"
			} else {
				if byte_string(ejm.Mbr_partition_1.Part_type) == "p" {
					text += "<td port='port_one'>PRIMARIA %" + s1 + "</td>\n"
				} else if byte_string(ejm.Mbr_partition_1.Part_type) == "e" {
					text += "<td port='port_one'>EXTENDIDA %" + s1 + "</td>\n"
				}
			}

			if size1 == 0 {
				text += "<td port='port_two'>LIBRE</td>\n"
			} else {
				if byte_string(ejm.Mbr_partition_2.Part_type) == "p" {
					text += "<td port='port_two'>PRIMARIA %" + s2 + "</td>\n"
				} else if byte_string(ejm.Mbr_partition_2.Part_type) == "e" {
					text += "<td port='port_two'>EXTENDIDA %" + s2 + "</td>\n"
				}
			}

			if size1 == 0 {
				text += "<td port='port_three'>LIBRE</td>\n"
			} else {
				if byte_string(ejm.Mbr_partition_3.Part_type) == "p" {
					text += "<td port='port_three'>PRIMARIA %" + s3 + "</td>\n"
				} else if byte_string(ejm.Mbr_partition_3.Part_type) == "e" {
					text += "<td port='port_three'>EXTENDIDA %" + s3 + "</td>\n"
				}
			}

			if size4 == 0 {
				text += "<td port='port_four'>LIBRE</td>\n"
			} else {
				if byte_string(ejm.Mbr_partition_4.Part_type) == "p" {
					text += "<td port='port_four'>PRIMARIA %" + s4 + "</td>\n"
				} else if byte_string(ejm.Mbr_partition_4.Part_type) == "e" {
					text += "<td port='port_four'>EXTENDIDA %" + s4 + "</td>\n"
				}
			}

			text += "</tr>\n</table>\n>];\n}"
			disco.Close()

			f, err := os.Create(path_sin_disco(path_actual) + "disk.txt")
			defer f.Close()
			if err != nil {
				fmt.Println(">> Error drawing graph!")
			}
			file, err := os.Open(path_sin_disco(path_actual) + "disk.txt")
			defer file.Close()
			if err != nil {
				fmt.Println(">> Error reading the file. Try again.")
				return
			}

			f.WriteString(text)

			e := exec.Command("dot", "-Tpng", "disk.txt", "-o", "disk.png")
			if er := e.Run(); er != nil {
				fmt.Println(">> Error", er)
				return
			}
			respuesta_rep = "Reporte DISK realizado\\n"
			respuesta_exec += "Reporte DISK realizado" + "\\n"

		} else if flag == true && name == "sb" {
			poss := 0
			poss = poss + 0
			for i := 0; i < contadorMount; i++ {
				if id == arregloMountId[i] {
					path = arregloMountPath[i]
					poss = i

					flag = true
					break
				}
			}

			path_users := path_sin_disco(path) + "users.txt"

			fmt.Println("soy actyal" + path_actual)
			arch, err := os.OpenFile(path_users, os.O_RDWR, 0660)
			if err != nil {
				fmt.Println("es aca")
				msg_error(err)
			}
			var texto string
			var aux_texto string
			var cont_datos int
			cont_datos = cont_datos + 0

			scanner := bufio.NewScanner(arch)
			for scanner.Scan() {
				//contadorlog+=1;
				texto = scanner.Text()
				//cout<<contadorlog<<texto<<endl;
				//cout<<texto.length()<<endl;
				//cout<<texto[8]<<endl;
				cont_datos = 0
				aux_texto += texto + "\n"
			}

			fmt.Println(len(aux_texto))
			arch.Close()

			bytess := len(aux_texto)
			var sizebloc float64 = float64(bytess) / float64(64)
			var entero int = int(math.Round(sizebloc))

			if sizebloc < float64(entero) {
				entero += 1
			}
			fmt.Println(entero)
			t := time.Now()
			fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
				t.Year(), t.Month(), t.Day(),
				t.Hour(), t.Minute(), t.Second())
			var info string = ""
			info += "digraph H{\n node [shape=plaintext];\n B [ label=< <TABLE BORDER =\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n<TR PORT=\"header\"><TD COLSPAN=\"2\">SUPERBLOQUE</TD></TR>\n"
			info += "<TR><TD PORT=\"1\">S_FILESYSTEM_TYPE</TD><TD> " + "ext2" + "</TD></TR>\n"
			info += "<TR><TD PORT=\"2\">S_INODES_COUNT</TD><TD>" + strconv.Itoa(2) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"3\">S_BLOCKS_COUNT</TD><TD>" + strconv.Itoa(entero) + "</TD></TR>\n"

			info += "<TR><TD  PORT=\"4\">S_FREE_BLOCKS_COUNT</TD><TD >" + strconv.Itoa(250-entero) + "</TD></TR>"
			info += "<TR><TD PORT=\"5\">S_FREE_INODES_COUNT</TD><TD>" + strconv.Itoa(250-2) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"6\">S_MTIME</TD><TD>" + fecha + "</TD></TR>\n"
			info += "<TR><TD PORT=\"7\">S_MNT_COUNT</TD><TD>" + strconv.Itoa(1) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"8\">S_MAGIC</TD><TD>" + "0xEF53" + "</TD></TR>\n"
			info += "<TR><TD PORT=\"9\">S_INODE_SIZE</TD><TD>" + strconv.Itoa(250) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"10\">S_BLOCK_SIZE</TD><TD>" + strconv.Itoa(250) + "</TD></TR>\n"

			info += "<TR><TD   PORT=\"11\">S_FIRST_INO</TD><TD  >" + strconv.Itoa(3) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"12\">S_FIRST_BLO</TD><TD>" + strconv.Itoa(entero+1) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"13\">S_BM_INODE_START</TD><TD>" + strconv.Itoa(1) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"14\">S_BM_BLOCK_START</TD><TD>" + strconv.Itoa(1) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"15\">S_INODE_START</TD><TD>" + strconv.Itoa(1) + "</TD></TR>\n"
			info += "<TR><TD PORT=\"16\">S_BLOCK_START</TD><TD>" + strconv.Itoa(1) + "</TD></TR>\n"

			info += "</TABLE> >];}\n"
			disco.Close()

			f, err := os.Create(path_sin_disco(path_actual) + "sb.txt")
			defer f.Close()
			if err != nil {
				fmt.Println(">> Error drawing graph!")
			}
			file, err := os.Open(path_sin_disco(path_actual) + "sb.txt")
			defer file.Close()
			if err != nil {
				fmt.Println(">> Error reading the file. Try again.")
				return
			}

			f.WriteString(info)

			e := exec.Command("dot", "-Tpng", "sb.txt", "-o", "sb.png")
			if er := e.Run(); er != nil {
				fmt.Println(">> Error", er)
				return
			}
			respuesta_rep = "Reporte SB realizado\\n"
			respuesta_exec += "Reporte SB realizado\\n"
		}

	} else {
		respuesta_rep = "Falta algun parametro en reportes\\n"
		respuesta_exec += "Falta algun parametro en reportes" + "\\n"
		fmt.Println("Falta algun parametro en reportes")
	}

}

// mostrar
func mostrar(path string) {
	//fin_archivo := false
	//var emptyid [100]byte
	ejm_empty := MBR{}
	//cont := 0
	// Apertura de archivo
	disco, err := os.OpenFile(path, os.O_RDWR, 0660)
	if err != nil {
		msg_error(err)
	}
	// Calculo del tamano de struct en bytes

	ejm2 := struct_to_bytes(ejm_empty)
	sstruct := len(ejm2)
	fmt.Println(sstruct)

	// Lectrura de conjunto de bytes en archivo binario
	lectura := make([]byte, sstruct)
	_, err = disco.ReadAt(lectura, int64(sstruct))
	if err != nil && err != io.EOF {
		msg_error(err)
	}

	// Conversion de bytes a struct
	ejm := bytes_to_struct(lectura)
	sstruct = len(lectura)

	if err != nil {
		msg_error(err)
	}

	fmt.Print(" Size: ")
	fmt.Print(string(ejm.Mbr_tamano[:]))
	fmt.Print(" Signature: ")
	fmt.Print(string(ejm.Mbr_disk_signature[:]))
	fmt.Print(" Direccion: ")
	fmt.Print(string(ejm.Mbr_fecha_creacion[:]))
	fmt.Print(" fit1: ")
	fmt.Println(string(ejm.Mbr_partition_1.Part_fit[:]))

	/*copy(ejm.Mbr_tamano[:], "bicho7")
	aux := struct_to_bytes(ejm)
	_, err = disco.WriteAt(aux, int64(sstruct))
	if err != nil {
		msg_error(err)
	}*/

	//fmt.Println(cont)
	//cont++
	/*var dato [100]byte
	copy(dato[:], "15360")
	if ejm.Mbr_tamano == dato {
		fmt.Println("HOLA SERENEXOOOOOOOOOOOOOOOOOOO")
	}*/
	disco.Close()
}

func struct_to_bytes(p interface{}) []byte {
	// Codificacion de Struct a []Bytes
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil && err != io.EOF {
		fmt.Println("error aqui")
		msg_error(err)
	}
	return buf.Bytes()
}

func bytes_to_struct(s []byte) MBR {
	// Decodificacion de [] Bytes a Struct ejemplo
	p := MBR{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil && err != io.EOF {
		msg_error(err)
	}
	return p
}

func byte_to_string(ejm MBR) string {
	flagnum := false
	cont := 0
	aux_tamano := ""
	for !flagnum {
		if string(ejm.Mbr_partition_1.Part_size[cont]) == "0" || string(ejm.Mbr_partition_1.Part_size[cont]) == "1" || string(ejm.Mbr_partition_1.Part_size[cont]) == "2" || string(ejm.Mbr_partition_1.Part_size[cont]) == "3" || string(ejm.Mbr_partition_1.Part_size[cont]) == "4" || string(ejm.Mbr_partition_1.Part_size[cont]) == "5" || string(ejm.Mbr_partition_1.Part_size[cont]) == "6" || string(ejm.Mbr_partition_1.Part_size[cont]) == "7" || string(ejm.Mbr_partition_1.Part_size[cont]) == "8" || string(ejm.Mbr_partition_1.Part_size[cont]) == "9" {
			fmt.Println(string(ejm.Mbr_partition_1.Part_size[cont]))
			aux_tamano += string(ejm.Mbr_partition_1.Part_size[cont])
			cont++
		} else {
			flagnum = true
		}

	}

	return aux_tamano
}

func byte_to_string2(ejm MBR) string {
	flagnum := false
	cont := 0
	aux_tamano := ""
	for !flagnum {
		if string(ejm.Mbr_partition_2.Part_size[cont]) == "0" || string(ejm.Mbr_partition_2.Part_size[cont]) == "1" || string(ejm.Mbr_partition_2.Part_size[cont]) == "2" || string(ejm.Mbr_partition_2.Part_size[cont]) == "3" || string(ejm.Mbr_partition_2.Part_size[cont]) == "4" || string(ejm.Mbr_partition_2.Part_size[cont]) == "5" || string(ejm.Mbr_partition_2.Part_size[cont]) == "6" || string(ejm.Mbr_partition_2.Part_size[cont]) == "7" || string(ejm.Mbr_partition_2.Part_size[cont]) == "8" || string(ejm.Mbr_partition_2.Part_size[cont]) == "9" {
			fmt.Println(string(ejm.Mbr_partition_2.Part_size[cont]))
			aux_tamano += string(ejm.Mbr_partition_2.Part_size[cont])
			cont++
		} else {
			flagnum = true
		}

	}

	return aux_tamano
}
func byte_to_string3(ejm MBR) string {
	flagnum := false
	cont := 0
	aux_tamano := ""
	for !flagnum {
		if string(ejm.Mbr_partition_3.Part_size[cont]) == "0" || string(ejm.Mbr_partition_3.Part_size[cont]) == "1" || string(ejm.Mbr_partition_3.Part_size[cont]) == "2" || string(ejm.Mbr_partition_3.Part_size[cont]) == "3" || string(ejm.Mbr_partition_3.Part_size[cont]) == "4" || string(ejm.Mbr_partition_3.Part_size[cont]) == "5" || string(ejm.Mbr_partition_3.Part_size[cont]) == "6" || string(ejm.Mbr_partition_3.Part_size[cont]) == "7" || string(ejm.Mbr_partition_3.Part_size[cont]) == "8" || string(ejm.Mbr_partition_3.Part_size[cont]) == "9" {
			fmt.Println(string(ejm.Mbr_partition_3.Part_size[cont]))
			aux_tamano += string(ejm.Mbr_partition_3.Part_size[cont])
			cont++
		} else {
			flagnum = true
		}

	}

	return aux_tamano
}

func byte_to_string4(ejm MBR) string {
	flagnum := false
	cont := 0
	aux_tamano := ""
	for !flagnum {
		if string(ejm.Mbr_partition_4.Part_size[cont]) == "0" || string(ejm.Mbr_partition_4.Part_size[cont]) == "1" || string(ejm.Mbr_partition_4.Part_size[cont]) == "2" || string(ejm.Mbr_partition_4.Part_size[cont]) == "3" || string(ejm.Mbr_partition_4.Part_size[cont]) == "4" || string(ejm.Mbr_partition_4.Part_size[cont]) == "5" || string(ejm.Mbr_partition_4.Part_size[cont]) == "6" || string(ejm.Mbr_partition_4.Part_size[cont]) == "7" || string(ejm.Mbr_partition_4.Part_size[cont]) == "8" || string(ejm.Mbr_partition_4.Part_size[cont]) == "9" {
			fmt.Println(string(ejm.Mbr_partition_4.Part_size[cont]))
			aux_tamano += string(ejm.Mbr_partition_4.Part_size[cont])
			cont++
		} else {
			flagnum = true
		}

	}

	return aux_tamano
}

func path_sin_disco(path string) string {
	var lugar int = 0
	var path2 string = ""
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			lugar = i
		}

	}

	for i := 0; i <= lugar; i++ {
		path2 += string(path[i])

	}

	return path2
}

func crearArchivo_UyG(path string) {
	file, err := os.Create(path)

	if err != nil {
		msg_error(err)
		fmt.Println("h0")
	}
	_, err = file.WriteString("1,G,root\n1,U,root,root,123")

	file.Close()
}

func mostrar2() {

	var datos string = "digraph G {\na1 -> b3;\n}"
	f, err := os.Create("mbr.txt")
	defer f.Close()
	if err != nil {
		fmt.Println(">> Error drawing graph!")
	}
	file, err := os.Open("C:/Users/sebas/go/src/MIA_Proyecto2_201906085-/main.go")
	defer file.Close()
	if err != nil {
		fmt.Println(">> Error reading the file. Try again.")
		return
	}

	f.WriteString(datos)

	e := exec.Command("dot", "-Tpng", "mbr.txt", "-o", "mbr.png")
	if er := e.Run(); er != nil {
		fmt.Println(">> Error", er)
		return
	}
}

func byte_string(dato [100]byte) string {
	var aux byte
	var sstring string
	for i := 0; i < len(dato); i++ {
		if dato[i] != aux {
			sstring += string(dato[i])
		}
	}

	return sstring
}

func byte_string2(dato [16]byte) string {
	var aux byte
	var sstring string
	for i := 0; i < len(dato); i++ {
		if dato[i] != aux {
			sstring += string(dato[i])
		}
	}

	return sstring
}

func crearDirectorioSiNoExiste(directorio string) {
	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		err = os.Mkdir(directorio, 0755)
		if err != nil {
			// Aquí puedes manejar mejor el error, es un ejemplo
			panic(err)
		}
	}
}

func exece(commandArray []string) string {
	respuesta_exec = ""
	path := ""

	for i := 0; i < len(commandArray); i++ {
		data := commandArray[i]
		if strings.Contains(data, "-path=") {

			path = strings.Replace(data, "-path=", "", 1)
			path = strings.Replace(path, "\"", "", 2)
		}
	}
	//respuesta := ""
	file, err := os.OpenFile(path, os.O_RDWR, 0660)
	if err != nil {
		fmt.Println("es aca")
		msg_error(err)
	}
	var texto string
	var aux_texto string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//contadorlog+=1;
		texto = scanner.Text()
		aux_texto += texto
		split_comando(texto)

		//cout<<contadorlog<<texto<<endl;
		//cout<<texto.length()<<endl;
		//cout<<texto[8]<<endl;
		//cont_datos = 0
	}
	//fmt.Println(aux_texto)
	respuesta_exec += "Ejecucion realizada con exito" + "\\n"

	return respuesta_exec
}

func pausee() {
	//respuesta := ""
	//fmt.Println("Hello world!")
	fmt.Print("Press 'Enter' to continue...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')

	//	return respuesta
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
