package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

/*
Comandos de prueba

fdisk –Size->300 –path->/home/thefernan/Desktop/disk1.dsk –name->Particion1
Mkdisk -Size->3000 –unit->K -path->/home/thefernan/Desktop/disk1.dsk
*/

type mbr struct {
	Mbrtamaño        int64
	Mbrfechacreacion [54]byte
	Mbrdisksignature int64
	Diskfit          [1]byte
	Mbrpartition_1   partition
	Mbrpartition_2   partition
	Mbrpartition_3   partition
	Mbrpartition_4   partition
}

type partition struct {
	Part_status [1]byte
	Part_type   [1]byte
	Part_fit    [1]byte
	Part_start  int64
	Part_size   int64
	Part_name   [16]byte
}

type ebr struct {
	Part_status [1]byte
	Part_fit    [1]byte
	Part_start  int64
	Part_size   int64
	Part_name   [16]byte
	Part_next   int64
}

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
		colorizefn(ColorSkyBlue, "[Fase1")
		colorizefn(ColorGreen, "@")
		colorizefn(ColorPurple, "Parrot]")
		colorizefn(ColorGreen, "$ ")
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
		} else if data == "pause\n" || data == "pause" {
			colorize(ColorYellow, "Lectura Pausada *Press Enter*")
			reader := bufio.NewReader(os.Stdin)
			comando, _ := reader.ReadString('\n')
			if comando == "\n" {
				colorize(ColorYellow, "Lectura Reanudada")
			} else {
				colorize(ColorYellow, "Lectura Reanudada")
			}
		} else if data == "mkdisk" {
			path := false
			size := false
			unit := false
			fit := false
			other := false

			var tam int
			var unidad string
			var ruta string
			var ajuste string

			for i := 1; i < len(commandArray); i++ {
				command := strings.ToLower(commandArray[i])
				caracteres := strings.Split(command, "")

				if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "s" && caracteres[2] == "i" && caracteres[3] == "z" && caracteres[4] == "e" && caracteres[5] == "-" && caracteres[6] == ">" {
					size = true
					parametros := strings.Split(command, "->")
					i1, err := strconv.Atoi(parametros[1])
					if err == nil {
						tam = i1
						if tam <= 0 {
							other = true
						}
					} else {
						other = true
					}
				} else if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "u" && caracteres[2] == "n" && caracteres[3] == "i" && caracteres[4] == "t" && caracteres[5] == "-" && caracteres[6] == ">" {
					parametros := strings.Split(command, "->")
					if parametros[1] == "K" || parametros[1] == "k" || parametros[1] == "M" || parametros[1] == "m" {
						unit = true
						unidad = parametros[1]

					} else {
						other = true
					}
				} else if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "p" && caracteres[2] == "a" && caracteres[3] == "t" && caracteres[4] == "h" && caracteres[5] == "-" && caracteres[6] == ">" {
					path = true
					parametros := strings.Split(commandArray[i], "->")
					ruta = parametros[1]
				} else if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "f" && caracteres[2] == "i" && caracteres[3] == "i" && caracteres[4] == "-" && caracteres[5] == ">" {
					parametros := strings.Split(command, "->")
					if parametros[1] == "BF" || parametros[1] == "FF" || parametros[1] == "WF" {
						fit = true
						ajuste = parametros[1]
					} else {
						other = true
					}
				} else {
					other = true
				}
			}

			if other == false && path == true && size == true {
				if unit == false {
					unidad = "M"
				}

				if fit == false {
					ajuste = "FF"
				}

				//se procede a crear el archivo
				file, err := os.Create("/home/thefernan/Desktop/disk1.dsk")
				defer file.Close()
				if err != nil {
					log.Fatal(err)
				}

				//se crea una variable temporal con un cero que nos ayudará a llenar nuestro archivo de ceros lógicos
				var temporal int8 = 0
				s := &temporal
				var binario bytes.Buffer
				binary.Write(&binario, binary.BigEndian, s)

				/*
					se realiza un for para llenar el archivo completamente de ceros
					NOTA: Para esta parte se recomienda tener un buffer con 1024 ceros (ya que 1024 es la medida
					mínima a escribir) para que este ciclo sea más eficiente
				*/
				var tamTotal = 0
				if unidad == "m" || unidad == "M" {
					tamTotal = tam * 1024 * 1024
				} else {
					tamTotal = tam * 1024
				}

				for i := 0; i < tamTotal; i++ {
					escribirBytes(file, binario.Bytes())
				}

				charFit := strings.Split(ajuste, "")
				fmt.Println(charFit)
				mbr := mbr{}

				mbr.Mbrtamaño = int64(tamTotal)
				copy(mbr.Mbrfechacreacion[:], time.Now().String())
				mbr.Mbrdisksignature = int64(512)
				copy(mbr.Diskfit[:], charFit[0])

				/*
					se escribira un estudiante por default para llevar el control.
					En el proyecto, el que nos ayuda a llevar el control de las
					particiones es el mbr
				*/

				//nos posicionamos al inicio del archivo usando la funcion Seek
				//Funcion Seek: https://ispycode.com/GO/Files-And-Directories/Seek-Positions-in-File
				file.Seek(0, 0)

				var bufferEstudiante bytes.Buffer
				binary.Write(&bufferEstudiante, binary.BigEndian, &mbr)
				escribirBytes(file, bufferEstudiante.Bytes())

				defer file.Close()

				colorize(ColorYellow, "Disco Creado -->"+ruta)
			} else {
				fmt.Println("Creacion De Disco Erronea Parametros Invalidos")
			}
		} else if data == "rmdisk" {
			parametro := commandArray[1]
			param := strings.ToLower(parametro)
			caracteres := strings.Split(param, "")
			if caracteres[0] == "-" && caracteres[1] == "p" && caracteres[2] == "a" && caracteres[3] == "t" && caracteres[4] == "h" && caracteres[5] == "-" && caracteres[6] == ">" {
				paramsParts := strings.Split(parametro, "->")
				path0 := paramsParts[1]
				path := path0[0 : len(path0)-1]

				err := os.Remove(path)

				if err != nil {
					colorize(ColorRed, "Error al eliminar el archivo.")
				} else {
					colorize(ColorYellow, "Disco Eliminado *"+path)
				}

			} else {
				colorize(ColorRed, "Comando No Aceptado")
			}
		} else if data == "fdisk" {
			path := false
			name := false
			unit := false
			tipe := false
			fit := false
			delete := false
			add := false
			size := false
			other := false

			var ruta string
			var nombre string
			var unidad string
			var tipo string
			var ajuste string
			var agragar int
			var tam int
			var eliminar string

			var primerComando string
			var firsComand bool

			for i := 1; i < len(commandArray); i++ {
				command := strings.ToLower(commandArray[i])
				caracteres := strings.Split(command, "")

				if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "p" && caracteres[2] == "a" && caracteres[3] == "t" && caracteres[4] == "h" && caracteres[5] == "-" && caracteres[6] == ">" {
					path = true
					parametros := strings.Split(commandArray[i], "->")
					ruta = parametros[1]
				} else if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "n" && caracteres[2] == "a" && caracteres[3] == "m" && caracteres[4] == "e" && caracteres[5] == "-" && caracteres[6] == ">" {
					name = true
					parametros := strings.Split(command, "->")
					nombre = parametros[1]
				} else if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "u" && caracteres[2] == "n" && caracteres[3] == "i" && caracteres[4] == "t" && caracteres[5] == "-" && caracteres[6] == ">" {
					parametros := strings.Split(command, "->")
					if parametros[1] == "K" || parametros[1] == "k" || parametros[1] == "M" || parametros[1] == "m" {
						unit = true
						unidad = parametros[1]
					} else {
						other = true
					}
				} else if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "t" && caracteres[2] == "y" && caracteres[3] == "p" && caracteres[4] == "e" && caracteres[5] == "-" && caracteres[6] == ">" {
					parametros := strings.Split(command, "->")
					if parametros[1] == "p" || parametros[1] == "e" || parametros[1] == "l" {
						tipe = true
						tipo = parametros[1]
					} else {
						other = true
					}
				} else if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "f" && caracteres[2] == "i" && caracteres[3] == "i" && caracteres[4] == "-" && caracteres[5] == ">" {
					parametros := strings.Split(command, "->")
					if parametros[1] == "BF" || parametros[1] == "FF" || parametros[1] == "WF" {
						fit = true
						ajuste = parametros[1]
					} else {
						other = true
					}
				} else if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "d" && caracteres[2] == "e" && caracteres[3] == "l" && caracteres[4] == "e" && caracteres[5] == "t" && caracteres[6] == "e" {
					delete = true
					parametros := strings.Split(command, "->")
					configuracionDel := strings.ToLower(parametros[1])
					if configuracionDel == "fast" || configuracionDel == "full" {
						eliminar = configuracionDel
						if firsComand == false {
							firsComand = true
							primerComando = "delete"
						}
					} else {
						colorize(ColorRed, "Error En  Los Parametros De Eliminacion")
					}
				} else if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "a" && caracteres[2] == "d" && caracteres[3] == "d" {
					add = true
					parametros := strings.Split(command, "->")
					i1, err := strconv.Atoi(parametros[1])
					if err == nil {
						agragar = i1
						if firsComand == false {
							firsComand = true
							primerComando = "add"
						}
					} else {
						other = true
					}
				} else if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "s" && caracteres[2] == "i" && caracteres[3] == "z" && caracteres[4] == "e" && caracteres[5] == "-" && caracteres[6] == ">" {
					size = true
					parametros := strings.Split(command, "->")
					i1, err := strconv.Atoi(parametros[1])
					if err == nil {
						tam = i1
						if tam <= 0 {
							other = true
						}
						if firsComand == false {
							firsComand = true
							primerComando = "create"
						}
					} else {
						other = true
					}
				} else {

				}
			}

			if other == false && path == true && name == true {
				if unit == false {
					unidad = "K"
				}

				if tipe == false {
					tipo = "P"
				}

				if fit == false {
					ajuste = "WF"
				}

				if primerComando == "delete" && delete == true {
					fmt.Println("eliminando ", nombre, "-", ruta, "-", eliminar)
				} else if primerComando == "add" && add == true {
					fmt.Println("añadiendo ", agragar, "-", unidad, "-", nombre, "-", ruta)
				} else if primerComando == "create" && size == true {
					mbrTemp := leerMBR(ruta)
					status := [1]byte{65}
					fmt.Println(mbrTemp)
					if mbrTemp.Mbrpartition_1.Part_status != status && mbrTemp.Mbrpartition_2.Part_status != status && mbrTemp.Mbrpartition_3.Part_status != status && mbrTemp.Mbrpartition_4.Part_status != status {
						var tamTot int64
						if unidad == "K" || unidad == "k" {
							tamTot = int64(tam) * 1024
						} else if unidad == "M" || unidad == "m" {
							tamTot = int64(tam) * 1024 * 1024
						}
						partition1 := partition{}
						copy(partition1.Part_status[:], "A")
						copy(partition1.Part_type[:], tipo)
						copy(partition1.Part_fit[:], ajuste)
						partition1.Part_start = int64(unsafe.Sizeof(mbrTemp))
						partition1.Part_size = tamTot
						copy(partition1.Part_name[:], nombre)
						mbrTemp.Mbrpartition_1 = partition1
						file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
						defer file.Close()
						if err != nil {
							log.Fatal(err)
						}
						file.Seek(0, 0)

						var bufferEstudiante bytes.Buffer
						binary.Write(&bufferEstudiante, binary.BigEndian, &mbrTemp)
						escribirBytes(file, bufferEstudiante.Bytes())
						defer file.Close()
					} else if mbrTemp.Mbrpartition_1.Part_status != status || mbrTemp.Mbrpartition_2.Part_status != status || mbrTemp.Mbrpartition_3.Part_status != status || mbrTemp.Mbrpartition_4.Part_status != status {
						if mbrTemp.Mbrpartition_1.Part_status != status {
							fmt.Println("Creando Particion ", tam, "-", unidad, "-", nombre, "-", ruta, "-", tipo, "-", ajuste)
							var tamTot int64
							if unidad == "K" || unidad == "k" {
								tamTot = int64(tam) * 1024
							} else if unidad == "M" || unidad == "m" {
								tamTot = int64(tam) * 1024 * 1024
							}

							////////////////////////////////////////////////////////////
							////////////////////////////////////////////////////////////
							////////////////////////////////////////////////////////////
							partition1 := partition{}
							copy(partition1.Part_status[:], "A")
							copy(partition1.Part_type[:], tipo)
							copy(partition1.Part_fit[:], ajuste)
							partition1.Part_size = tamTot
							copy(partition1.Part_name[:], nombre)

							crearParticion(mbrTemp, partition1)

							mbrTemp.Mbrpartition_1 = partition1
							file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
							defer file.Close()
							if err != nil {
								log.Fatal(err)
							}
							file.Seek(0, 0)

							var bufferEstudiante bytes.Buffer
							binary.Write(&bufferEstudiante, binary.BigEndian, &mbrTemp)
							escribirBytes(file, bufferEstudiante.Bytes())
							defer file.Close()
						} else if mbrTemp.Mbrpartition_2.Part_status != status {
							fmt.Println("Creando Particion ", tam, "-", unidad, "-", nombre, "-", ruta, "-", tipo, "-", ajuste)
							partition1 := partition{}
							copy(partition1.Part_status[:], "A")
							copy(partition1.Part_type[:], tipo)
							copy(partition1.Part_fit[:], ajuste)
							partition1.Part_start = int64(unsafe.Sizeof(mbrTemp))
							partition1.Part_size = int64(unsafe.Sizeof(partition1))
							copy(partition1.Part_name[:], nombre)
							mbrTemp.Mbrpartition_2 = partition1
							file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
							defer file.Close()
							if err != nil {
								log.Fatal(err)
							}
							file.Seek(0, 0)

							var bufferEstudiante bytes.Buffer
							binary.Write(&bufferEstudiante, binary.BigEndian, &mbrTemp)
							escribirBytes(file, bufferEstudiante.Bytes())
							defer file.Close()
						} else if mbrTemp.Mbrpartition_3.Part_status != status {
							fmt.Println("Creando Particion ", tam, "-", unidad, "-", nombre, "-", ruta, "-", tipo, "-", ajuste)
							partition1 := partition{}
							copy(partition1.Part_status[:], "A")
							copy(partition1.Part_type[:], tipo)
							copy(partition1.Part_fit[:], ajuste)
							partition1.Part_start = int64(unsafe.Sizeof(mbrTemp))
							partition1.Part_size = int64(unsafe.Sizeof(partition1))
							copy(partition1.Part_name[:], nombre)
							mbrTemp.Mbrpartition_3 = partition1
							file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
							defer file.Close()
							if err != nil {
								log.Fatal(err)
							}
							file.Seek(0, 0)

							var bufferEstudiante bytes.Buffer
							binary.Write(&bufferEstudiante, binary.BigEndian, &mbrTemp)
							escribirBytes(file, bufferEstudiante.Bytes())
							defer file.Close()
						} else if mbrTemp.Mbrpartition_4.Part_status != status {
							fmt.Println("Creando Particion ", tam, "-", unidad, "-", nombre, "-", ruta, "-", tipo, "-", ajuste)
							partition1 := partition{}
							copy(partition1.Part_status[:], "A")
							copy(partition1.Part_type[:], tipo)
							copy(partition1.Part_fit[:], ajuste)
							partition1.Part_start = int64(unsafe.Sizeof(mbrTemp))
							partition1.Part_size = int64(unsafe.Sizeof(partition1))
							copy(partition1.Part_name[:], nombre)
							mbrTemp.Mbrpartition_4 = partition1
							file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
							defer file.Close()
							if err != nil {
								log.Fatal(err)
							}
							file.Seek(0, 0)

							var bufferEstudiante bytes.Buffer
							binary.Write(&bufferEstudiante, binary.BigEndian, &mbrTemp)
							escribirBytes(file, bufferEstudiante.Bytes())
							defer file.Close()
						}
					} else {
						colorize(ColorRed, "Error Las Particiones Estan Completas")
					}
				}
			} else {
				fmt.Println("Administracion De Discos Cerrada Por Comandos Erroneos")
			}
			colorize(ColorYellow, "Administrando Disco")
		} else if data == "mount" {
			path := false
			name := false
			other := false

			var ruta string
			var nombre string

			for i := 1; i < len(commandArray); i++ {
				command := strings.ToLower(commandArray[i])
				caracteres := strings.Split(command, "")

				if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "p" && caracteres[2] == "a" && caracteres[3] == "t" && caracteres[4] == "h" && caracteres[5] == "-" && caracteres[6] == ">" {
					path = true
					parametros := strings.Split(command, "->")
					ruta = parametros[1]
				} else if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "n" && caracteres[2] == "a" && caracteres[3] == "m" && caracteres[4] == "e" && caracteres[5] == "-" && caracteres[6] == ">" {
					name = true
					parametros := strings.Split(command, "->")
					nombre = parametros[1]
				} else {
					other = true
				}
			}

			if other == false && path == true && name == true {
				fmt.Println("Se ha Montado ", nombre, "-", ruta)
			} else {
				fmt.Println("No Se Ha podido Montar El Disco Error En Los Parametros")
			}
		} else if data == "unmount" {
			id := false
			other := false

			var identificador string
			for i := 1; i < len(commandArray); i++ {
				command := strings.ToLower(commandArray[i])
				caracteres := strings.Split(command, "")
				if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "i" && caracteres[2] == "d" && caracteres[3] == "-" && caracteres[4] == ">" {
					id = true
					parametros := strings.Split(command, "->")
					identificador = parametros[1]
				} else {
					other = true
				}
			}

			if other == false && id == true {
				fmt.Println("Se ha Desmontado ", identificador)
			} else {
				fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
			}
		} else if data == "rep" {
			id := false
			name := false
			path := false
			other := false

			var identificador string
			var nombre string
			var ruta string

			for i := 1; i < len(commandArray); i++ {
				command := strings.ToLower(commandArray[i])
				caracteres := strings.Split(command, "")
				if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "p" && caracteres[2] == "a" && caracteres[3] == "t" && caracteres[4] == "h" && caracteres[5] == "-" && caracteres[6] == ">" {
					path = true
					parametros := strings.Split(command, "->")
					ruta = parametros[1]
				} else if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "n" && caracteres[2] == "a" && caracteres[3] == "m" && caracteres[4] == "e" && caracteres[5] == "-" && caracteres[6] == ">" {
					name = true
					parametros := strings.Split(command, "->")
					if parametros[1] == "mbr" || parametros[1] == "disk" {
						nombre = parametros[1]
					} else {
						other = true
					}
				} else if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "i" && caracteres[2] == "d" && caracteres[3] == "-" && caracteres[4] == ">" {
					id = true
					parametros := strings.Split(command, "->")
					identificador = parametros[1]
				} else {
					other = true
				}
			}

			if other == false && id == true && name == true && path == true {
				if nombre == "disk" {
					fmt.Println("reporte de disco ", identificador, "-", nombre, "-", ruta)
				} else {
					fmt.Println("reporte de mbr ", identificador, "-", nombre, "-", ruta)
				}
			} else {
				fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
			}
		} else {
			colorize(ColorYellow, "Comando Incorrecto")
		}
	} else {
		colorize(ColorYellow, "Comentario De Script")
	}
}

func leerMBR(path string) mbr {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	mbrTemp := mbr{}

	var size int = int(unsafe.Sizeof(mbrTemp))
	file.Seek(0, 0)
	mbrTemp = obtenerMBR(file, size, mbrTemp)

	return mbrTemp
}

func crearParticion(mbrTemp mbr, partition1 partition) {

}

func obtenerMBR(file *os.File, size int, mbrTemp mbr) mbr {
	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, size)

	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable estudianteTemporal
	err := binary.Read(buffer, binary.BigEndian, &mbrTemp)
	if err != nil {
		log.Fatal("binary.Read failed ", err)
	}

	//retornamos el estudiante
	return mbrTemp
}

func leerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal("Error De Lectura ", err)
	}

	return bytes
}

func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
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
