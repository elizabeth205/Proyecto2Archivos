package mount

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"structs"
	"unsafe"
)

var path = ""
var name = ""
var montardisco = structs.Montadas{}

var commandpath = regexp.MustCompile("(?i)\\s?-\\s?path\\s?=\\s?/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.dk|.txt)")
var commandname = regexp.MustCompile("(?i)\\s?-\\s?name\\s?=\\s?([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")

var rutaprin = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.dk|.txt)")
var rutasig = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var idnombre = regexp.MustCompile("([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")

var masterBootRec = structs.Mbr{}
var SaveEbr = structs.Ebr{}

var pos2 = 0
var rutacom = ""

var alph = [26]string{"A", "B", "C", "D", "E", "F", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func CommandAnaysis(recibida string) {
	if commandname.MatchString(recibida) && commandpath.MatchString(recibida) {
		path = rutaprin.FindString(commandpath.FindString(recibida))
		name = commandname.FindString(commandname.FindString(recibida))
		name = strings.ReplaceAll(name, "-", "")
		name = strings.ReplaceAll(name, "name", "")
		name = strings.ReplaceAll(name, "=", "")
		name = strings.ReplaceAll(name, " ", "")
		recibida = commandpath.ReplaceAllLiteralString(recibida, "")
		recibida = commandname.ReplaceAllLiteralString(recibida, "")
		fmt.Printf("path = %s\n", path)
		fmt.Printf("name = %s\n", name)
	} else {
		fmt.Println("Error faltan parametros ")
		fmt.Println("no se reconocio path")
		fmt.Printf("%q\n", commandpath.Split(recibida, -1))
		fmt.Println("no se reconocio name")
		fmt.Printf("%q\n", commandname.Split(recibida, -1))
	}

}

func RecuperarDisco() structs.Montadas {
	var newdisk = structs.Montadas{}
	newdisk.ID = 1

	if structs.Mountedisk().Len() > 0 {
		for k := structs.Mountedisk().Front(); k != nil; k = k.Next() {
			itdisk := structs.Montadas(k.Value.(structs.Montadas))
			var rutarecibida = string(itdisk.Path[:])
			if strings.Compare(rutarecibida, path) == 1 {
				return itdisk
			}
			var id = itdisk.ID + 1
			newdisk.ID = id
		}

		copy(newdisk.Path[:], path)
		return newdisk
	}
	copy(newdisk.Path[:], path)
	return newdisk
}

func OrganizacionMount() {
	var mountdisk = RecuperarDisco()
	Abrir_mbr()
	var extpart = -1
	var mountpart = structs.ParticionMontada{}
	for k := range masterBootRec.Tabla {
		var nomparti = string(masterBootRec.Tabla[k].Name[:])
		if string(masterBootRec.Tabla[k].Type) == "e" {

			if strings.Compare(nomparti, name) == 1 {
				fmt.Println("NO es posible el montaje")
				return
			}
		}

		if strings.Compare(nomparti, name) == 1 {
			mountpart.EstadoMount = 1
			copy(mountpart.Nombre[:], name)
			for l := range mountdisk.Lista {
				if mountdisk.Lista[l].EstadoMount == 0 {
					id_disk := strconv.Itoa(mountdisk.ID)
					var idcarnet = "19" + id_disk + alph[l]
					copy(mountpart.ID[:], idcarnet)
					mountdisk.Lista[l] = mountpart
					structs.Mount_disk(mountdisk)
					fmt.Printf("particion montada id:  %s \n", idcarnet)
					return
				}
			}

		}
	}
	for s := range masterBootRec.Tabla {
		if string(masterBootRec.Tabla[s].Type) == "e" {
			var nombreparti = string(masterBootRec.Tabla[s].Name[:])
			if strings.Compare(nombreparti, name) == 1 {
				extpart = s
				break
			}
		}
	}

	if extpart != -1 {
		var init = masterBootRec.Tabla[extpart].Start
		LeerEBR(init)
		if SaveEbr.Status == 1 {
			var name_ebr = ""
			for SaveEbr.Next != -1 {
				LeerEBR(SaveEbr.Next)
				name_ebr = string(SaveEbr.Name[:])
				if strings.Compare(name_ebr, name) == 1 {
					mountpart.EstadoMount = 1
					copy(mountpart.Nombre[:], name)
					for l := range mountdisk.Lista {
						if mountdisk.Lista[l].EstadoMount == 0 {
							id_disk := strconv.Itoa(mountdisk.ID)
							var id = "19" + id_disk + alph[l]
							copy(mountpart.ID[:], id)
							mountdisk.Lista[l] = mountpart
							structs.Mount_disk(mountdisk)
							fmt.Printf("particion montada id:%s \n", id)
							return
						}
					}
				}
			}
			name_ebr = string(SaveEbr.Name[:])
			if strings.Compare(name_ebr, name) == 1 {
				mountpart.EstadoMount = 1
				copy(mountpart.Nombre[:], name)
				for l := range mountdisk.Lista {
					if mountdisk.Lista[l].EstadoMount == 0 {
						id_disk := strconv.Itoa(mountdisk.ID)
						var id = "19" + id_disk + alph[l]
						copy(mountpart.ID[:], id)
						mountdisk.Lista[l] = mountpart
						structs.Mount_disk(mountdisk)
						fmt.Printf("particion montada id:  %s \n", id)
						return
					}
				}
			}
		}
	}

}

func LeerEBR(initebr int64) {
	for pos, char := range path {
		if char == '/' {
			pos2 = pos
		}
	}
	rutacom = path[:pos2]

	if ArchivoExiste(path) {
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}

		var ebrtam int64 = int64(unsafe.Sizeof(SaveEbr))
		file.Seek(initebr, 0)
		data := leerBytes(file, ebrtam)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &SaveEbr)
		if err != nil {
			log.Fatal("leer archivobinary.Read failed", err)
		}
		fmt.Printf("nombre: %s \n", SaveEbr.Name[:])
		fmt.Printf("status: %d \n", SaveEbr.Status)
		fmt.Printf("fit: %s \n", SaveEbr.Fit)
		fmt.Printf("siguiente: %d \n", SaveEbr.Next)
		fmt.Printf("size : %d \n", SaveEbr.Size)
		fmt.Printf("inicio: %d \n", SaveEbr.Start)
		file.Close()

	} else {
		fmt.Print("NO eiste disco")
	}
}

func Abrir_mbr() {
	for pos, char := range path {
		if char == '/' {
			pos2 = pos
		}
	}
	rutacom = path[:pos2]
	fmt.Print(rutacom)
	if ArchivoExiste(path) {
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		var tammaster int64 = int64(unsafe.Sizeof(masterBootRec))
		data := leerBytes(file, tammaster)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &masterBootRec)
		if err != nil {
			log.Fatal("No se pudo leer", err)
		}
		/*
			fmt.Println("Mbr Tamano:", masterBootRec.Tamano)
			fmt.Println("Mbr Fecha creacion:", string(masterBootRec.Fecha[:]))
			fmt.Println("Mbr Disk Signarue:", masterBootRec.Firma)
			fmt.Println("Disco Fit:", string(masterBootRec.Fit[:]))
			for k := range masterBootRec.Tabla {
				fmt.Println("particion:", string(masterBootRec.Tabla[k].Name[:]))
				fmt.Println("size :", masterBootRec.Tabla[k].Size)
			}*/
		file.Close()

	} else {
		fmt.Print("NO existe el disco")
	}
}

func ArchivoExiste(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}

func leerBytes(file *os.File, number int64) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal("ERROR", err)
	}

	return bytes
}
