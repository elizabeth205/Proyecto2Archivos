package rep

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"structs"
	"unsafe"
)

var path = ""
var name string = ""
var id = ""
var mountdisk = structs.Montadas{}

var commandpath = regexp.MustCompile("(?i)\\s?-\\s?path\\s?=\\s?/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.jpg|.png|.gif)")
var commandname = regexp.MustCompile("(?i)\\s?-\\s?name\\s?=\\s?([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var commandid = regexp.MustCompile("(?i)\\s?-\\s?id\\s?=\\s?([0-9]{3}[A-Z])")

var route = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.jpg|.png|.gif)")
var disdirect = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.dk|.txt)")
var rutasgs = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var nombrerep = regexp.MustCompile("([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var nombreid = regexp.MustCompile("([0-9]{3}[A-Z])")

var masterBootRecord = structs.Mbr{}
var verebr = structs.Ebr{}

var pos2 = 0
var totalroute = ""

var discopath = ""

func Analizador(recibido string) {
	if commandname.MatchString(recibido) && commandpath.MatchString(recibido) {
		path = route.FindString(commandpath.FindString(recibido))
		name = commandname.FindString(commandname.FindString(recibido))
		name = strings.ReplaceAll(name, "-", "")
		name = strings.ReplaceAll(name, "name", "")
		name = strings.ReplaceAll(name, "=", "")
		name = strings.ReplaceAll(name, " ", "")
		recibido = commandpath.ReplaceAllLiteralString(recibido, "")
		recibido = commandname.ReplaceAllLiteralString(recibido, "")
		id = nombreid.FindString(commandid.FindString(recibido))
		fmt.Printf("path = %s\n", path)
		fmt.Printf("name = %s\n", name)
		fmt.Printf("id = %s\n", id)
	} else {
		fmt.Println("Parametros incorrectos")
		fmt.Println("no se reconocio path")
		fmt.Printf("%q\n", commandpath.Split(recibido, -1))
		fmt.Println("no se reconocio name")
		fmt.Printf("%q\n", commandname.Split(recibido, -1))
	}

}

func Reps() {
	if name == "mbr" {
		MBRReport()
	}

}

func MBRReport() {
	var datas string = string(id[2])
	disk, _ := strconv.Atoi(datas)
	if structs.Mountedisk().Len() > 0 {
		var newdisk = structs.Montadas{}

		for k := structs.Mountedisk().Front(); k != nil; k = k.Next() {
			itdisk := structs.Montadas(k.Value.(structs.Montadas))

			if itdisk.ID == disk {
				newdisk = itdisk

				break
			}
		}
		var name_partition = ""
		for k := range newdisk.Lista {
			var id_part_mount = string(newdisk.Lista[k].ID[:])
			if id_part_mount == id {
				name_partition = string(newdisk.Lista[k].Nombre[:])
				discopath = string(newdisk.Path[:])
				discopath = disdirect.FindString(discopath)
				fmt.Println(discopath)
				break
			}
		}
		if name_partition != "" {
			for pos, char := range path {
				if char == '/' {
					pos2 = pos
				}
			}
			totalroute = path[:pos2]
			fmt.Print(ArchivoExiste(totalroute))
			if !ArchivoExiste(totalroute) {
				var err = os.Mkdir(totalroute, 0755)
				if err != nil {
					// Aquí puedes manejar mejor el error, es un ejemplo
					panic(err)
				}
			}
			var dotruta = totalroute + "/" + "mbr" + ".dot"
			Abrir_mbr()
			var part1 = nombrerep.FindString(string(masterBootRecord.Tabla[0].Name[:]))
			var part2 = nombrerep.FindString(string(masterBootRecord.Tabla[1].Name[:]))
			var part3 = nombrerep.FindString(string(masterBootRecord.Tabla[2].Name[:]))
			var part4 = nombrerep.FindString(string(masterBootRecord.Tabla[3].Name[:]))
			var tipodato byte = 'p'
			var codigodeldot = "digraph G { \n" +
				"ordering = out \n" +
				"forcelabels=true \n" +
				"graph[ranksep=1,margin=0.3  ]; \n" +
				"node [shape = plaintext];\n " +
				"1 [ label = <<TABLE color = \"black\"> \n" +
				"<TR>\n" +
				"<td > mbr tamaño_disco= " + strconv.Itoa(int(masterBootRecord.Tamano)) + "</td>\n"
			//por particion
			if masterBootRecord.Tabla[0].Type == tipodato {
				var porcentaje = strconv.Itoa(int(masterBootRecord.Tabla[0].Size) * 100 / int(masterBootRecord.Tamano))
				codigodeldot += "<td >" + part1 + "\n " + porcentaje + "%" + "</td>\n"
			} else {
				tam, code := logiparts()
				var colspan = strconv.Itoa(tam)
				codigodeldot += "<td coslspan=\"" + colspan + "\"" + ">" + "extendida" + "\n " + "</td>\n"
				codigodeldot += code
			}
			if masterBootRecord.Tabla[1].Type == tipodato {
				var porcentaje = strconv.Itoa(int(masterBootRecord.Tabla[1].Size) * 100 / int(masterBootRecord.Tamano))
				codigodeldot += "<td >" + part2 + "\n " + porcentaje + "%" + "</td>\n"
			} else {
				tam, code := logiparts()
				var colspan = strconv.Itoa(tam)
				codigodeldot += "<td coslspan=\"" + colspan + "\"" + ">" + "extendida" + "\n " + "</td>\n"
				codigodeldot += code
			}
			if masterBootRecord.Tabla[2].Type == tipodato {
				var porcentaje = strconv.Itoa(int(masterBootRecord.Tabla[2].Size) * 100 / int(masterBootRecord.Tamano))
				codigodeldot += "<td >" + part3 + "\n " + porcentaje + "%" + "</td>\n"

			} else {
				tam, code := logiparts()
				var colspan = strconv.Itoa(tam)
				codigodeldot += "<td coslspan=\"" + colspan + "\"" + ">" + "extendida" + "\n " + "</td>\n"
				codigodeldot += code
			}
			if masterBootRecord.Tabla[3].Type == tipodato {
				var porcentaje = strconv.Itoa(int(masterBootRecord.Tabla[3].Size) * 100 / int(masterBootRecord.Tamano))
				codigodeldot += "<td >" + part4 + "\n " + porcentaje + "%" + "</td>\n"

			} else {

			}

			codigodeldot += "</TR>\n" +
				"</TABLE>> dir =none color=white style =none]\n" +
				"}"
			f, err := os.Create(dotruta)
			verify(err)
			defer f.Close()

			f.Sync()
			w := bufio.NewWriter(f)
			n4, err := w.WriteString(codigodeldot)
			verify(err)
			fmt.Printf("se escribieron %d bytes \n", n4)
			w.Flush()
			dot := "dot"
			format := "-Tjpg"
			dot_file := dotruta
			ouput := "-o"
			ab_pa := totalroute + "/" + id + ".jpg"
			cmd := exec.Command(dot, format, dot_file, ouput, ab_pa)

			stdout, err := cmd.Output()

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			fmt.Println(string(stdout))
		} else {
			fmt.Println("esta particion no esta montada")
		}
	} else {
		fmt.Println("No esta montada la particion")
	}
}
func verify(e error) {
	if e != nil {
		panic(e)
	}
}

func ArchivoExiste(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}

func Abrir_mbr() {
	for pos, char := range discopath {
		if char == '/' {
			pos2 = pos
		}
	}
	totalroute = discopath[:pos2]
	fmt.Print(totalroute)
	if ArchivoExiste(discopath) {
		file, err := os.Open(discopath)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		var tamano_masterBoot int64 = int64(unsafe.Sizeof(masterBootRecord))
		data := leerBytes(file, tamano_masterBoot)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &masterBootRecord)
		if err != nil {
			log.Fatal("leer archivobinary.Read failed", err)
		}
		/*
			fmt.Println("Mbr Tamano:", masterBootRecord.Tamano)
			fmt.Println("Mbr Fecha creacion:", string(masterBootRecord.Fecha[:]))
			fmt.Println("Mbr Disk Signarue:", masterBootRecord.Firma)
			fmt.Println("Disco Fit:", string(masterBootRecord.Fit[:]))
			for k := range masterBootRecord.Tabla {
				fmt.Println("particion:", string(masterBootRecord.Tabla[k].Name[:]))
				fmt.Println("size :", masterBootRecord.Tabla[k].Size)
			}*/
		file.Close()

	} else {
		fmt.Print("no existe disco")
	}
}

func leerBytes(file *os.File, number int64) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal("ERROR", err)
	}

	return bytes
}

func Abrir_ebr(inicio_ebr int64) {
	for pos, char := range discopath {
		if char == '/' {
			pos2 = pos
		}
	}
	totalroute = discopath[:pos2]

	if ArchivoExiste(discopath) {
		file, err := os.Open(discopath)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}

		var ebrtam int64 = int64(unsafe.Sizeof(verebr))
		file.Seek(inicio_ebr, 0)
		data := leerBytes(file, ebrtam)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &verebr)
		if err != nil {
			log.Fatal("Fallo lectura", err)
		}
		fmt.Printf("nombre: %s \n", verebr.Name[:])
		fmt.Printf("status: %d \n", verebr.Status)
		fmt.Printf("fit: %s \n", verebr.Fit)
		fmt.Printf("siguiente : %d \n", verebr.Next)
		fmt.Printf("size : %d \n", verebr.Size)
		fmt.Printf("inicio: %d \n", verebr.Start)
		file.Close()

	} else {
		fmt.Print("no existe disco")
	}
}

func logiparts() (int, string) {
	var size int = 0
	var info = ""
	var extpart1 int64 = -1
	for s := range masterBootRecord.Tabla {
		if string(masterBootRecord.Tabla[s].Type) == "e" {
			extpart1 = masterBootRecord.Tabla[s].Start
			break

		}
	}

	Abrir_ebr(extpart1)

	for verebr.Next != -1 {
		Abrir_ebr(verebr.Next)
		var name_ebr = nombrerep.FindString(string(verebr.Name[:]))
		var porcentaje = strconv.Itoa(int(verebr.Size) * 100 / int(masterBootRecord.Tamano))
		info += "<td >" + name_ebr + "\n " + porcentaje + "%" + "</td>\n"
		size++
	}
	size++
	var name_ebr = nombrerep.FindString(string(verebr.Name[:]))
	var porcentaje = strconv.Itoa(int(verebr.Size) * 100 / int(masterBootRecord.Tamano))
	info += "<td >" + name_ebr + "\n " + porcentaje + "%" + "</td>\n"

	return size, info
}
