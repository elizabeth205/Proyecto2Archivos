package main

import (
	"bufio"
	"fdisk"
	"fmt"
	"log"
	"mkdisk"
	mkfs2 "mkfs"
	"mount"
	"os"
	"regexp"
	"rep"
	"strings"
)

var commkdisk = regexp.MustCompile("(?i)mkdisk")
var comrmdisk = regexp.MustCompile("(?i)rmdisk")
var read = regexp.MustCompile("(?i)read")
var comfdisk = regexp.MustCompile("(?i)fdisk")
var commount = regexp.MustCompile("(?i)mount")
var comexec = regexp.MustCompile("(?i)exec")
var comrep = regexp.MustCompile("(?i)rep")
var commkfs = regexp.MustCompile("(?i)commkfs")

var rutaexec = ""

//para leer el path en exec
var ruta1 = regexp.MustCompile("(?i)\\s?-\\s?path\\s?=\\s?/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.sh|.txt|.script)")
var ruta2 = regexp.MustCompile("/([a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)*(/[a-zA-Z]+([a-zA-Z]+|[0-9]+|_)*)+(.sh|.txt|.script)")

func main() {

	menu :=
		`
	------------------------------Escriba su comando o Exec seguido de la ruta del archivo que sea ejecutar------------------------------
	
>`

	reader := bufio.NewReader(os.Stdin)

	// Read para cada linea

	for {
		fmt.Print(menu)
		linea, _ := reader.ReadString('\n')
		command := strings.TrimRight(linea, "\n") // Le quito el salto de linea a la linea

		if command == "ESC" {
			fmt.Println("Adios!")
			break
		} else if comexec.MatchString(command) {
			ComandoExex(command)
			if ExistenciArchivo(rutaexec) {
				file, err := os.Open(rutaexec)
				if err != nil {
					log.Fatalf("No se puede abrir el archivo: %s", err)
				}
				filenuevo := bufio.NewScanner(file)
				for filenuevo.Scan() {
					if filenuevo.Text() == "pause" {
						fmt.Println("Esta es una pausa, Presione una tecla para continuar")
						lector2 := bufio.NewReader(os.Stdin)
						linea2, _ := lector2.ReadString('\n')
						command2 := strings.TrimRight(linea2, "\n")
						fmt.Println(command2)

					}
					sep := strings.TrimRight(filenuevo.Text(), "\n")
					AnalizadorComandos(sep[0:])
				}
				if err := filenuevo.Err(); err != nil {
					log.Fatalf("Error al abrir el archivo: %s", err)
				}
				file.Close()
			} else {
				fmt.Println("el archivo no existe!")
			}

		} else {
			fmt.Printf("%q\n", commkdisk.FindString(command))
			AnalizadorComandos(command)
		}

	}
}

func ComandoExex(input string) {
	if ruta1.MatchString(input) {
		rutaexec = ruta2.FindString(ruta1.FindString(input))
		input = ruta1.ReplaceAllLiteralString(input, "")
		fmt.Printf("ruta = %s\n", rutaexec)
		fmt.Printf("ruta = %s\n", rutaexec)
	} else {
		fmt.Println("Faltan elementos para este comando ")
		fmt.Println("No se reconocen estos valores ")
		fmt.Printf("%q\n", ruta1.Split(input, -1))
	}
}

func AnalizadorComandos(selected_command string) {
	fmt.Println(selected_command)
	if commkdisk.MatchString(selected_command) {
		fmt.Println("se creara un disco")
		selected_command = commkdisk.ReplaceAllLiteralString(selected_command, "")
		mkdisk.Filtro(selected_command)
		mkdisk.MKDISK()
	} else if comrmdisk.MatchString(selected_command) {
		fmt.Println("Se eliminara el disco")
		selected_command = comrmdisk.ReplaceAllLiteralString(selected_command, "")
		mkdisk.Filtro2(selected_command)
		mkdisk.DeleteDisk()
	} else if comfdisk.MatchString(selected_command) {
		fmt.Println("Creacion de particiones ")
		selected_command = comfdisk.ReplaceAllLiteralString(selected_command, "")
		fdisk.PartAnalise(selected_command)
		fdisk.Abrir_mbr()
	} else if commount.MatchString(selected_command) {
		fmt.Println("Se montara la particion seleccionada")
		selected_command = commount.ReplaceAllLiteralString(selected_command, "")
		mount.CommandAnaysis(selected_command)
		mount.OrganizacionMount()
	} else if comrep.MatchString(selected_command) {
		fmt.Println("Creando reportes")
		selected_command = commount.ReplaceAllLiteralString(selected_command, "")
		rep.Analizador(selected_command)
		rep.Reps()
	} else if commkdisk.MatchString(selected_command) {
		fmt.Println("Intentando formatear a EXT2")
		selected_command = commkdisk.ReplaceAllLiteralString(selected_command, "")
		mkfs2.Analizador(selected_command)
	}
}

func ExistenciArchivo(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}
