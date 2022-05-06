package mkdisk

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"structs"
	"time"
)

var tam = ""
var fit = ""
var unit = ""
var path = ""

//variables utilizadas para analizar las entradas
var elementoruta = regexp.MustCompile("(?i)\\s?-\\s?path\\s?=\\s?/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.dk|.txt)")
var elementosize = regexp.MustCompile("(?i)\\s?-\\s?size\\s?=\\s?[0-9]+")
var elementofit = regexp.MustCompile("(?i)\\s?-\\s?fit\\s?=\\s?(bf|ff|wf)")
var elementounit = regexp.MustCompile("(?i)\\s?-\\s?unit\\s?=\\s?(k|m|K|M)")
var nums = regexp.MustCompile("[0-9]+")
var direc = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.dk|.txt)")
var rutasg = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)")
var ajustes = regexp.MustCompile("(bf|ff|wf|BF|FF|WF)")
var unidades = regexp.MustCompile("(k|m|M|K)")
var slash = regexp.MustCompile("/")

func Filtro2(input string) {
	if elementoruta.MatchString(input) {
		path = direc.FindString(elementoruta.FindString(input))
		fmt.Printf("path = %s\n", path)
	} else {
		fmt.Println("Orden incorrecto")
		fmt.Println("no se reconocio ningun valor")
		fmt.Printf("%q\n", elementoruta.Split(input, -1))
	}
}
func Filtro(enter string) {
	if elementosize.MatchString(enter) && elementoruta.MatchString(enter) {
		tam = nums.FindString(elementosize.FindString(enter))
		path = direc.FindString(elementoruta.FindString(enter))
		enter = elementosize.ReplaceAllLiteralString(enter, "")
		enter = elementoruta.ReplaceAllLiteralString(enter, "")
		fmt.Printf("size = %s\n", tam)
		fmt.Printf("path = %s\n", path)
		if regexp.MustCompile("(?i)fit").MatchString(enter) {
			if elementofit.MatchString(enter) {
				fit = ajustes.FindString(elementofit.FindString(enter))
				fmt.Printf("fit = %s\n", fit)
			} else {
				fmt.Println("Orden icorrecto")
				fmt.Println("NO se reconocio ningun valor")
				fmt.Printf("%q\n", elementofit.Split(enter, -1))
			}
		} else {
			fit = "wf"
			fmt.Printf("fit = %s\n", fit)
		}
		if regexp.MustCompile("(?i)unit").MatchString(enter) {
			if elementounit.MatchString(enter) {
				unit = unidades.FindString(elementounit.FindString(enter))
				fmt.Printf("unit = %s\n", unit)
			} else {
				fmt.Println("orden incorrecto")
				fmt.Println("No se reconocio la unidad ")
				fmt.Printf("%q\n", elementounit.Split(enter, -1))
			}
		} else {
			unit = "m"
			fmt.Printf("unit = %s\n", unit)
		}
	} else {
		fmt.Println("orden incorrecto")
		fmt.Println("ruta no reconocida ")
		fmt.Printf("%q\n", elementoruta.Split(enter, -1))
		fmt.Println("tamanio no reconocido")
		fmt.Printf("%q\n", elementosize.Split(enter, -1))

	}
}

var pos2 = 0
var path_rou = ""

func MKDISK() {

	for pos, char := range path {
		if char == '/' {
			pos2 = pos
		}
	}
	path_rou = path[:pos2]
	fmt.Print(ArchivoExiste(path_rou))

	if !ArchivoExiste(path_rou) {
		var fail = os.Mkdir(path_rou, 0755)
		if fail != nil {
			// Aquí puedes manejar mejor el error, es un ejemplo
			panic(fail)
		}
	}
	arch, err := os.Create(path)

	size, err := strconv.ParseInt(tam, 10, 64)
	fmt.Println(size, err, reflect.TypeOf(size))
	defer arch.Close()
	if err != nil {
		log.Fatal(err)
	}
	var empty int8 = 0
	s := &empty
	var num int64 = 0

	if strings.Compare(strings.ToLower(unit), "m") == 0 {
		num = int64(size) * 1024 * 1024
	} else if strings.Compare(strings.ToLower(unit), "k") == 0 {
		num = int64(size) * 1024
	}
	num = num - 1

	var binfile bytes.Buffer
	binary.Write(&binfile, binary.BigEndian, s)
	EscribirA(arch, binfile.Bytes())

	//situando el cursor en la ultima posicion
	arch.Seek(num, 0)

	//colocando el ultimo byte para rellenar
	var binfile2 bytes.Buffer
	binary.Write(&binfile2, binary.BigEndian, s)
	EscribirA(arch, binfile2.Bytes())

	arch.Seek(0, 0)

	disco := structs.Mbr{}
	disco.Tamano = num + 1

	datetime := time.Now()
	datetimearrange := strings.Split(datetime.String(), "")
	caddatetime := ""
	for i := 0; i < 16; i++ {
		caddatetime = caddatetime + datetimearrange[i]
	}
	copy(disco.Fecha[:], caddatetime)
	copy(disco.Fit[:], fit)
	var signature int8
	binary.Read(rand.Reader, binary.LittleEndian, &signature)
	if signature < 0 {
		signature = signature * -1
	}
	disco.Firma = signature

	var binfile3 bytes.Buffer
	binary.Write(&binfile3, binary.BigEndian, disco)
	EscribirA(arch, binfile3.Bytes())
	//path := path
	//graficarDISCO(path)
	//graficarMBR(path)
}

func EscribirA(file *os.File, bytes []byte) {
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

func DeleteDisk() {
	if ArchivoExiste(path) {
		fmt.Println("desea eliminar el archivo  (y/n) ")

		reader := bufio.NewReader(os.Stdin)
		entrada, _ := reader.ReadString('\n')
		eleccion := strings.TrimRight(entrada, "\n")
		if eleccion == "y" {
			err := os.Remove(path)
			if err != nil {
				fmt.Printf("Error eliminando disco: %v\n", err)
			} else {
				fmt.Println("Se elimino el disco correctamente")
			}
		} else {
			fmt.Println("Operacion no ejecutada")
		}
	} else {
		fmt.Println("Disco inexistente")
	}
}
