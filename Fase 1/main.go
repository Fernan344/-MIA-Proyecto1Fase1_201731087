package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	colorize(ColorRed, "                uuuuuuu")
	colorize(ColorRed, "            uu$$$$$$$$$$$uu")
	colorize(ColorRed, "          uu$$$$$$$$$$$$$$$$$uu")
	colorize(ColorRed, "         u$$$$$$$$$$$$$$$$$$$$$u")
	colorize(ColorRed, "        u$$$$$$$$$$$$$$$$$$$$$$$u")
	colorize(ColorRed, "       u$$$$$$$$$$$$$$$$$$$$$$$$$u")
	colorize(ColorRed, "       u$$$$$$$$$$$$$$$$$$$$$$$$$u")
	colorize(ColorRed, "       u$$$$$$\"   \"$$$\"   \"$$$$$$u")
	colorize(ColorRed, "       \"$$$$\"      u$u       $$$$\"")
	colorize(ColorRed, "        $$$u       u$u       u$$$")
	colorize(ColorRed, "        $$$u      u$$$u      u$$$")
	colorize(ColorRed, "         \"$$$$uu$$$   $$$uu$$$$\"")
	colorize(ColorRed, "          \"$$$$$$$\"   \"$$$$$$$\"")
	colorize(ColorRed, "            u$$$$$$$u$$$$$$$u")
	colorize(ColorRed, "	     u$\"$\"$\"$\"$\"$\"$u")
	colorize(ColorRed, "  uuu        $$u$ $ $ $ $u$$       uuu")
	colorize(ColorRed, " u$$$$        $$$$$u$u$u$$$       u$$$$")
	colorize(ColorRed, "  $$$$$uu      \"$$$$$$$$$\"     uu$$$$$$")
	colorize(ColorRed, "u$$$$$$$$$$$uu    \"\"\"\"\"    uuuu$$$$$$$$$$")
	colorize(ColorRed, "$$$$\"\"\"$$$$$$$$$$uuu   uu$$$$$$$$$\"\"\"$$$\"")
	colorize(ColorRed, " \"\"\"      \"\"$$$$$$$$$$$uu \"\"$\"\"\"")
	colorize(ColorRed, "                uuuu \"\"$$$$$$$$$$uuu")
	colorize(ColorRed, "  u$$$uuu$$$$$$$$$uu \"\"$$$$$$$$$$$uuu$$$")
	colorize(ColorRed, "  $$$$$$$$$$\"\"\"\"           \"\"$$$$$$$$$$$\"")
	colorize(ColorRed, "   \"$$$$$\"                      \"\"$$$$\"\"")
	colorize(ColorRed, "      $$$\"                         $$$$")
	colorize(ColorWhite, "*****************************************")
	colorize(ColorWhite, "*         CONSOLA DE COMANDOS           *")
	colorize(ColorWhite, "*****************************************")
	interpretar()
}

func interpretar() {
	for {
		colorizefn(ColorSkyBlue, "Fase1@Parrot$: ")
		reader := bufio.NewReader(os.Stdin)
		comando, _ := reader.ReadString('\n')
		if comando == "exit\n" {
			break
		} else {
			lineaDeComandos(comando)
		}
	}
}

func lineaDeComandos(comando string) {
	var commandArray []string
	commandArray = strings.Split(comando, " ")
	executeComand(commandArray)
}

func executeComand(commandArray []string) {
	data := strings.ToLower(commandArray[0])
	comment := strings.Split(data, "")
	if comment[0] != "#" {
		if data == "exec" {
			parametro := commandArray[1]
			param := strings.ToLower(parametro)
			caracteres := strings.Split(param, "")

			if caracteres[0] == "-" && caracteres[1] == "p" && caracteres[2] == "a" && caracteres[3] == "t" && caracteres[4] == "h" && caracteres[5] == "-" && caracteres[6] == ">" {
				paramsParts := strings.Split(parametro, "->")
				path0 := paramsParts[1]
				path := path0[0 : len(path0)-1]
				b, err := ioutil.ReadFile(path)
				if err != nil {
					fmt.Print(err)
				}
				str := string(b) // convert content to a 'string'
				lineas := strings.Split(str, "\n")
				colorize(ColorYellow, "Corriendo Scripts")
				for i := 0; i < len(lineas)-1; i++ {
					colorize(ColorReset, lineas[i])
					lineaDeComandos(lineas[i])
				}
			} else {
				colorize(ColorRed, "Comando No Aceptado")
			}
		} else if data == "pause" {
			colorize(ColorYellow, "Lectura Pausado")
		} else if data == "mkdisk" {
			colorize(ColorYellow, "Creando Disco")
		} else if data == "rmdisk" {
			colorize(ColorYellow, "Eliminando Disco")
		} else if data == "fdisk" {
			colorize(ColorYellow, "Administrando Disco")
		} else if data == "mount" {
			colorize(ColorYellow, "Montando Disco")
		} else if data == "unmount" {
			colorize(ColorYellow, "Desmontando Disco")
		} else if data == "rep" {
			colorize(ColorYellow, "Creando Reportes")
		} else {
			colorize(ColorYellow, "Comando Incorrecto")
		}
	} else {
		colorize(ColorYellow, "Comentario De Script")
	}
}

type Color string

const (
	ColorBlack    Color = "\u001b[30m"
	ColorRed            = "\u001b[31m"
	ColorGreen          = "\u001b[32m"
	ColorYellow         = "\u001b[33m"
	ColorBlue           = "\u001b[34m"
	ColorReset          = "\u001b[0m"
	ColorPurple         = "\u001b[35m"
	ColorSkyBlue        = "\u001b[36m"
	ColorWhite          = "\u001b[37m"
	ColorSurprise       = "\u001b[41m"
)

func colorize(color Color, message string) {
	fmt.Println(string(color), message, string(ColorReset))
}

func colorizefn(color Color, message string) {
	fmt.Print(string(color), message, string(ColorReset))
}
