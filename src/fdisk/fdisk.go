package fdisk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"structs"
	"unsafe"
)

var size = ""
var fit = ""
var unit = ""
var path = ""
var type_ = ""
var name = ""

var commandsize = regexp.MustCompile("(?i)\\s?-\\s?size\\s?=\\s?[0-9]+")
var commandfit = regexp.MustCompile("(?i)\\s?-\\s?fit\\s?=\\s?(bf|ff|wf)")
var commandunit = regexp.MustCompile("(?i)\\s?-\\s?unit\\s?=\\s?(k|m)")
var commandpath = regexp.MustCompile("(?i)\\s?-\\s?path\\s?=\\s?/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.dk|.txt)")
var commandtype = regexp.MustCompile("(?i)\\s?-\\s?type\\s?=\\s?(P|p|E|e|L|l)")
var othertype = regexp.MustCompile("(?i)\\s?-\\s?type\\s?=\\s?")
var commandname = regexp.MustCompile("(?i)\\s?-\\s?name\\s?=\\s?([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")

var nums = regexp.MustCompile("[0-9]+")
var route1 = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.dk|.txt)")
var route2 = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var fitscals = regexp.MustCompile("(bf|ff|wf)")
var unidadvals = regexp.MustCompile("(k|m|M|K|b|B)")
var type_val = regexp.MustCompile("(P|p|E|e|L|l)")
var ids = regexp.MustCompile("([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var slash = regexp.MustCompile("/")

var masterBoot = structs.Mbr{}
var Ebr = structs.Ebr{}

func PartAnalise(recibido string) {
	if commandsize.MatchString(recibido) && commandpath.MatchString(recibido) && commandname.MatchString(recibido) {
		size = nums.FindString(commandsize.FindString(recibido))
		path = route1.FindString(commandpath.FindString(recibido))
		name = commandname.FindString(commandname.FindString(recibido))
		name = strings.ReplaceAll(name, "-", "")
		name = strings.ReplaceAll(name, "name", "")
		name = strings.ReplaceAll(name, "=", "")
		name = strings.ReplaceAll(name, " ", "")
		recibido = commandsize.ReplaceAllLiteralString(recibido, "")
		recibido = commandpath.ReplaceAllLiteralString(recibido, "")
		recibido = commandname.ReplaceAllLiteralString(recibido, "")
		fmt.Printf("size = %s\n", size)
		fmt.Printf("path = %s\n", path)
		fmt.Printf("name = %s\n", name)
		if regexp.MustCompile("(?i)fit").MatchString(recibido) {
			if commandfit.MatchString(recibido) {
				fit = fitscals.FindString(commandfit.FindString(recibido))
				recibido = commandfit.ReplaceAllLiteralString(recibido, "")
				fmt.Printf("fit = %s\n", fit)
			} else {
				fmt.Println("orden incorrecto")
				fmt.Println("no se reconocio fit ")
				fmt.Printf("%q\n", commandfit.Split(recibido, -1))
			}
		} else {
			fit = "wf"
			fmt.Printf("fit = %s\n", fit)
		}
		if regexp.MustCompile("(?i)unit").MatchString(recibido) {
			if commandunit.MatchString(recibido) {
				unit = unidadvals.FindString(commandunit.FindString(recibido))
				recibido = commandunit.ReplaceAllLiteralString(recibido, "")
				fmt.Printf("unit = %s\n", unit)
			} else {
				fmt.Println("orden incorrecto")
				fmt.Println("no se reconocion unit ")
				fmt.Printf("%q\n", commandunit.Split(recibido, -1))
			}
		} else {
			unit = "K"
			fmt.Printf("unit = %s\n", unit)
		}
		if regexp.MustCompile("(?i)type").MatchString(recibido) {
			if commandtype.MatchString(recibido) {
				recibido = othertype.ReplaceAllLiteralString(recibido, "")
				type_ = strings.ToLower(type_val.FindString(recibido))
				fmt.Printf("type = %s\n", type_)
			} else {
				fmt.Println("paramentro incorrecto")
				fmt.Println("no se reconocio type ")
				fmt.Printf("%q\n", commandtype.Split(recibido, -1))
			}
		} else {
			type_ = "p"
			fmt.Printf("type = %s\n", type_)
		}

	} else {
		fmt.Println("Faltan parametros")
		fmt.Println("no se reconocio path")
		fmt.Printf("%q\n", commandpath.Split(recibido, -1))
		fmt.Println("no se reconocio size")
		fmt.Printf("%q\n", commandsize.Split(recibido, -1))
		fmt.Println("no se reconocio name ")
		fmt.Printf("%q\n", commandname.Split(recibido, -1))
	}

}

func EscBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

func ArchivoExiste(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}

var pos2 = 0
var totpas = ""

func partinfo(size_newP int64) [3]int64 {
	datos := [3]int64{-1, -1, -1}
	var f1 = masterBoot.Tabla[0].Start + masterBoot.Tabla[0].Size
	var f2 = masterBoot.Tabla[1].Start + masterBoot.Tabla[1].Size
	var f3 = masterBoot.Tabla[2].Start + masterBoot.Tabla[2].Size

	if masterBoot.Tabla[0].Size == 0 || masterBoot.Tabla[1].Size == 0 ||
		masterBoot.Tabla[2].Size == 0 || masterBoot.Tabla[3].Size == 0 {
		if masterBoot.Tabla[0].Size == 0 && masterBoot.Tabla[0].Status == 0 {
			datos[0] = 0
			datos[1] = size_newP + int64(unsafe.Sizeof(structs.Mbr{}))
			datos[2] = int64(unsafe.Sizeof(structs.Mbr{})) + 1
		} else if masterBoot.Tabla[1].Size == 0 && masterBoot.Tabla[1].Status == 0 {
			datos[0] = 1
			datos[1] = masterBoot.Tabla[0].Size + size_newP + int64(unsafe.Sizeof(structs.Mbr{}))
			datos[2] = f1 + 1
		} else if masterBoot.Tabla[2].Size == 0 && masterBoot.Tabla[2].Status == 0 {
			datos[0] = 2
			datos[1] = masterBoot.Tabla[0].Size + masterBoot.Tabla[1].Size + size_newP + int64(unsafe.Sizeof(structs.Mbr{}))
			datos[2] = f2 + 1
		} else if masterBoot.Tabla[3].Size == 0 && masterBoot.Tabla[3].Status == 0 {
			datos[0] = 3
			datos[1] = masterBoot.Tabla[0].Size + masterBoot.Tabla[1].Size + masterBoot.Tabla[2].Size + size_newP + int64(unsafe.Sizeof(structs.Mbr{}))
			datos[2] = f3 + 1
		} else {
			fmt.Println("ESpacio insuficiente")
		}
	}

	return datos
}

func Abrir_mbr() {
	for pos, char := range path {
		if char == '/' {
			pos2 = pos
		}
	}
	totpas = path[:pos2]
	fmt.Print(totpas)
	if ArchivoExiste(path) {
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		var tammast int64 = int64(unsafe.Sizeof(masterBoot))
		data := leerBytes(file, tammast)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &masterBoot)
		if err != nil {
			log.Fatal("fallo lectura", err)
		}
		fmt.Println("Tamano:", masterBoot.Tamano)
		fmt.Println("Fecha creacion:", string(masterBoot.Fecha[:]))
		fmt.Println("disk Signarue:", masterBoot.Firma)
		fmt.Println("Fit:", string(masterBoot.Fit[:]))
		for k := range masterBoot.Tabla {
			fmt.Println("particion:", string(masterBoot.Tabla[k].Name[:]))
			fmt.Println("size :", masterBoot.Tabla[k].Size)
		}
		file.Close()
		partcreate()
	} else {
		fmt.Print("NO eciste disco")
	}
}

func ObtenerSize() int64 {
	var num int64 = 0
	Size, err := strconv.ParseInt(size, 10, 64)
	fmt.Println(size, err, reflect.TypeOf(size))
	if err != nil {
		log.Fatal(err)
	}
	if strings.Compare(strings.ToLower(unit), "m") == 0 {
		num = int64(Size) * 1024 * 1024
	} else if strings.Compare(strings.ToLower(unit), "k") == 0 {
		num = int64(Size) * 1024
	}
	num = num - 1
	return num
}

func Repetido() bool {
	var existe = false
	for k := range masterBoot.Tabla {
		sigau := string(masterBoot.Tabla[k].Name[:])
		fmt.Println(sigau)
		fmt.Print(name)
		if sigau == name {
			existe = true
			break
		}
	}
	return existe
}

func verextendida() bool {
	var ext bool = false
	for k := range masterBoot.Tabla {
		tipo_ext := string(masterBoot.Tabla[k].Type)
		if tipo_ext == "e" {
			ext = true
			break
		}
	}
	return ext
}

func calpart() [3]int64 {
	data := [3]int64{-1, -1, -1}
	for k := range masterBoot.Tabla {
		if string(masterBoot.Tabla[k].Type) == "e" {
			data[0] = int64(k)
			data[1] = masterBoot.Tabla[k].Start
			data[2] = masterBoot.Tabla[k].Size
			break
		}
	}

	return data
}

func partcreate() {
	if emptypart() {
		if type_ != "l" {
			masterBoot.Tabla[0].Size = ObtenerSize()
			if masterBoot.Tabla[0].Size > 0 && masterBoot.Tabla[0].Size < masterBoot.Tamano {
				masterBoot.Tabla[0].Status = 1
				copy(masterBoot.Tabla[0].Fit[:], fit)
				if type_ == "p" {
					masterBoot.Tabla[0].Type = 'p'
				} else if type_ == "e" {
					masterBoot.Tabla[0].Type = 'e'
					var size_mbr = int64(unsafe.Sizeof(masterBoot))
					recebr(size_mbr)
				}
				masterBoot.Tabla[0].Start = int64(unsafe.Sizeof(structs.Mbr{})) + 1
				copy(masterBoot.Tabla[0].Name[:], name)
				insertar_mbr()
			} else {
				fmt.Println("error tamanio incorrecto")
			}
		} else {
			fmt.Println("(error no eciste extendida")
		}
	} else {
		if type_ != "l" {
			var size_parti int64 = ObtenerSize()
			if size_parti < masterBoot.Tamano {
				var resultado [3]int64 = partinfo(size_parti)
				var selectpart = resultado[0]
				var using = resultado[1]
				var inipart = resultado[2]
				if selectpart != -1 && using != -1 && inipart != -1 {
					if !Repetido() {
						if using < masterBoot.Tamano {
							if type_ == "p" {
								masterBoot.Tabla[selectpart].Status = 1
								copy(masterBoot.Tabla[selectpart].Fit[:], fit)
								masterBoot.Tabla[selectpart].Type = 'p'
								masterBoot.Tabla[selectpart].Start = inipart
								masterBoot.Tabla[selectpart].Size = ObtenerSize()
								copy(masterBoot.Tabla[selectpart].Name[:], name)
								insertar_mbr()
								fmt.Println("se creo la particion")
							} else if !verextendida() && type_ == "e" {
								masterBoot.Tabla[selectpart].Status = 1
								copy(masterBoot.Tabla[selectpart].Fit[:], fit)
								masterBoot.Tabla[selectpart].Type = 'e'
								masterBoot.Tabla[selectpart].Start = inipart
								masterBoot.Tabla[selectpart].Size = ObtenerSize()
								copy(masterBoot.Tabla[selectpart].Name[:], name)
								insertar_mbr()
								recebr(inipart + 1)
								fmt.Println("se creo correctamente")
							} else {
								fmt.Println("error ya hay una extendida ")
							}
						} else {
							fmt.Println("error exceso de tamanio ")
						}
					} else {
						fmt.Println("ya existe nombre")
					}
				} else {
					fmt.Println("sin particiones")
				}

			}
		} else {
			fmt.Println("se creara la logica")
			if verextendida() {
				var result = calpart()
				var indice = result[0]
				var inicio = result[1]
				var tamano_ext = result[2]
				if indice != -1 && inicio != -1 && tamano_ext != -1 && masterBoot.Tabla[indice].Status != 0 {
					Abrir_ebr(inicio)
					if Ebr.Status == 1 {
						var repeted = false
						var name_ebr = ""
						for Ebr.Next != -1 {
							Abrir_ebr(Ebr.Next)
							name_ebr = string(Ebr.Name[:])
							if strings.Compare(name_ebr, name) == 1 {
								repeted = true
								break
							}
						}
						if !repeted {
							name_ebr = string(Ebr.Name[:])
							if Ebr.Next == -1 && strings.Compare(name_ebr, name) != 1 {
								Ebr.Next = Ebr.Size + Ebr.Start + 1
								var inicio_nueva = Ebr.Next
								if (inicio_nueva + ObtenerSize() + 1) < (inicio + tamano_ext) {
									recebr(Ebr.Start)
									copy(Ebr.Fit[:], fit)
									copy(Ebr.Name[:], name)
									Ebr.Start = inicio_nueva
									Ebr.Size = ObtenerSize() + int64(unsafe.Sizeof(Ebr))
									Ebr.Next = -1
									Ebr.Status = 1
									recebr(inicio_nueva)
								} else {
									fmt.Println("sin espacio")
								}
							} else {
								fmt.Println("nombre repetido")
							}
						} else {
							fmt.Println("nombre repetido, error")
						}
					} else {
						copy(Ebr.Fit[:], fit)
						copy(Ebr.Name[:], name)
						Ebr.Start = inicio
						Ebr.Size = ObtenerSize()
						Ebr.Next = -1
						Ebr.Status = 1
						recebr(inicio)
					}
				} else {
					fmt.Println("no hay ectendida ")
				}
			} else {
				fmt.Println("no hay extendida")
			}

		}
	}
}

func Abrir_ebr(inicio_ebr int64) {
	for pos, char := range path {
		if char == '/' {
			pos2 = pos
		}
	}
	totpas = path[:pos2]

	if ArchivoExiste(path) {
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}

		var tamebr int64 = int64(unsafe.Sizeof(Ebr))
		file.Seek(inicio_ebr, 0)
		data := leerBytes(file, tamebr)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &Ebr)
		if err != nil {
			log.Fatal("fallo lectura", err)
		}
		fmt.Printf("nombre  : %s \n", Ebr.Name[:])
		fmt.Printf("statu: %d \n", Ebr.Status)
		fmt.Printf("fit  : %s \n", Ebr.Fit)
		fmt.Printf("siguiente  : %d \n", Ebr.Next)
		fmt.Printf("size  : %d \n", Ebr.Size)
		fmt.Printf("inicio  : %d \n", Ebr.Start)
		file.Close()

	} else {
		fmt.Print("no existe el disco")
	}
}

func emptypart() bool {
	var resultado = true
	for _, s := range masterBoot.Tabla {
		if s.Status == 1 {
			resultado = false
			break
		}
	}
	return resultado
}

func insertar_mbr() {
	file, err := os.OpenFile(path, os.O_WRONLY, 0)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	mbr := masterBoot

	file.Seek(0, 0)
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, mbr)
	EscBytes(file, binario3.Bytes())
	file.Close()
}

func recebr(size int64) {
	file, err := os.OpenFile(path, os.O_WRONLY, 0)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(size, 0)
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, Ebr)
	EscBytes(file, binario3.Bytes())
	file.Close()
	fmt.Print("ebr creado ")
}

func leerBytes(file *os.File, number int64) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal("ERROR", err)
	}

	return bytes
}
