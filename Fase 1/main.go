package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

/*
Comandos de prueba
/home/thefernan/Desktop/Fase 1/main.go
fdisk –Size->300 –path->/home/thefernan/Desktop/disk3.dsk –name->Particion1
Mkdisk -Size->3000 –unit->K -path->/home/thefernan/Desktop/disk1.dsk
mount -path->/home/thefernan/Desktop/disk3.dsk -name->particion3
*/

var discosMounted []discos

type discos struct {
	id                string
	path              string
	partitionsMounted []partitionMounted
}

type partitionMounted struct {
	particion   partition
	id          string
	correlativo int
}

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
	if strings.Contains(comando, "\n") {
		commandArray = strings.Split(comando, "\n")
	}
	commandArray = strings.Split(commandArray[0], " ")
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
				file, err := os.Create(ruta)
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
				path := paramsParts[1]

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
				} else if (caracteres[0] == "-" || caracteres[0] == "–") && caracteres[1] == "f" && caracteres[2] == "i" && caracteres[3] == "t" && caracteres[4] == "-" && caracteres[5] == ">" {
					parametros := strings.Split(command, "->")
					if parametros[1] == "BF" || parametros[1] == "FF" || parametros[1] == "WF" || parametros[1] == "bf" || parametros[1] == "ff" || parametros[1] == "wf" {
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
					var tamTot int64
					if unidad == "K" || unidad == "k" {
						tamTot = int64(tam) * 1024
					} else if unidad == "M" || unidad == "m" {
						tamTot = int64(tam) * 1024 * 1024
					}
					fmt.Println(ajuste)
					charfit := strings.Split(ajuste, "")
					ajuste = charfit[0]

					partition1 := partition{}
					copy(partition1.Part_status[:], "A")
					copy(partition1.Part_type[:], tipo)
					copy(partition1.Part_fit[:], ajuste)
					partition1.Part_size = tamTot
					copy(partition1.Part_name[:], nombre)

					if mbrTemp.Mbrpartition_1.Part_status != status && mbrTemp.Mbrpartition_2.Part_status != status && mbrTemp.Mbrpartition_3.Part_status != status && mbrTemp.Mbrpartition_4.Part_status != status {

						partition1.Part_start = int64(unsafe.Sizeof(mbrTemp))

						if tipo == "P" || tipo == "p" {
							mbrTemp = asignPartition(mbrTemp, partition1)
							mbrTemp = sortPartitions(mbrTemp)
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
						} else if tipo == "E" || tipo == "e" {
							tipo1 := string(mbrTemp.Mbrpartition_1.Part_type[:])
							tipo2 := string(mbrTemp.Mbrpartition_2.Part_type[:])
							tipo3 := string(mbrTemp.Mbrpartition_3.Part_type[:])
							tipo4 := string(mbrTemp.Mbrpartition_4.Part_type[:])

							if (tipo1 != "E" && tipo1 != "e") && (tipo2 != "E" && tipo2 != "e") && (tipo3 != "E" && tipo3 != "e") && (tipo4 != "E" && tipo4 != "e") {

								mbrTemp = asignPartition(mbrTemp, partition1)
								mbrTemp = sortPartitions(mbrTemp)
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

								ebrTemp := ebr{}
								copy(ebrTemp.Part_status[:], "I")
								ebrTemp.Part_fit = partition1.Part_fit
								ebrTemp.Part_start = (partition1.Part_start + int64(unsafe.Sizeof(ebrTemp)))
								ebrTemp.Part_next = -1

								file.Seek(partition1.Part_start, 0)

								var bufferEbr bytes.Buffer
								binary.Write(&bufferEbr, binary.BigEndian, &mbrTemp)
								escribirBytes(file, bufferEbr.Bytes())
								defer file.Close()
							} else {
								colorize(ColorRed, "Error: Ya Existe una Particion Extendida")
							}

						}

					} else if mbrTemp.Mbrpartition_1.Part_status != status || mbrTemp.Mbrpartition_2.Part_status != status || mbrTemp.Mbrpartition_3.Part_status != status || mbrTemp.Mbrpartition_4.Part_status != status {

						fmt.Println("Creando Particion ", tam, "-", unidad, "-", nombre, "-", ruta, "-", tipo, "-", ajuste)

						position := crearParticion(mbrTemp, partition1)
						if position != -1 {
							partition1.Part_start = position
							if tipo == "P" || tipo == "p" {
								mbrTemp = asignPartition(mbrTemp, partition1)
								mbrTemp = sortPartitions(mbrTemp)
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
							} else if tipo == "E" || tipo == "e" {
								tipo1 := string(mbrTemp.Mbrpartition_1.Part_type[:])
								tipo2 := string(mbrTemp.Mbrpartition_2.Part_type[:])
								tipo3 := string(mbrTemp.Mbrpartition_3.Part_type[:])
								tipo4 := string(mbrTemp.Mbrpartition_4.Part_type[:])
								fmt.Println(tipo1, tipo2, tipo3, tipo4)

								if tipo1 != "E" && tipo1 != "e" && tipo2 != "E" && tipo2 != "e" && tipo3 != "E" && tipo3 != "e" && tipo4 != "E" && tipo4 != "e" {
									mbrTemp = asignPartition(mbrTemp, partition1)
									mbrTemp = sortPartitions(mbrTemp)
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

									ebrTemp := ebr{}
									copy(ebrTemp.Part_status[:], "I")
									ebrTemp.Part_fit = partition1.Part_fit
									ebrTemp.Part_start = (partition1.Part_start + int64(unsafe.Sizeof(ebrTemp)))
									ebrTemp.Part_next = -1

									file.Seek(partition1.Part_start, 0)

									var bufferEbr bytes.Buffer
									binary.Write(&bufferEbr, binary.BigEndian, &mbrTemp)
									escribirBytes(file, bufferEbr.Bytes())
									defer file.Close()

									fmt.Println(mbrTemp)
								} else {
									colorize(ColorRed, "Error: Ya Existe una Particion Extendida")
								}

							}
						} else {
							colorize(ColorRed, "No Hay Espacio Para La Particion")
						}
					} else if mbrTemp.Mbrpartition_1.Part_status == status && mbrTemp.Mbrpartition_2.Part_status == status && mbrTemp.Mbrpartition_3.Part_status == status && mbrTemp.Mbrpartition_4.Part_status == status {
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
					parametros := strings.Split(commandArray[i], "->")
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
				mbrTemp := leerMBR(ruta)
				var compareName [16]byte
				copy(compareName[:], nombre)

				diskId := "error"
				disk := discos{}
				for i := 0; i < len(discosMounted); i++ {
					if discosMounted[i].path == ruta {
						diskId = discosMounted[i].id
						disk = discosMounted[i]
						break
					}
				}

				if diskId == "error" {
					disk.id = "vd" + generarIdDisco(len(discosMounted))
					disk.path = ruta
				}

				if mbrTemp.Mbrpartition_1.Part_name == compareName {
					mounted := mountVerify(disk, mbrTemp.Mbrpartition_1)
					if mounted == false {
						montarParticion(ruta, disk, mbrTemp.Mbrpartition_1)
					} else {
						colorize(ColorRed, "Error: La Particion Ya Esta Montada")
					}
				} else if mbrTemp.Mbrpartition_2.Part_name == compareName {
					mounted := mountVerify(disk, mbrTemp.Mbrpartition_2)
					if mounted == false {
						montarParticion(ruta, disk, mbrTemp.Mbrpartition_2)
					} else {
						colorize(ColorRed, "Error: La Particion Ya Esta Montada")
					}
				} else if mbrTemp.Mbrpartition_3.Part_name == compareName {
					mounted := mountVerify(disk, mbrTemp.Mbrpartition_3)
					if mounted == false {
						montarParticion(ruta, disk, mbrTemp.Mbrpartition_3)
					} else {
						colorize(ColorRed, "Error: La Particion Ya Esta Montada")
					}
				} else if mbrTemp.Mbrpartition_4.Part_name == compareName {
					mounted := mountVerify(disk, mbrTemp.Mbrpartition_4)
					if mounted == false {
						montarParticion(ruta, disk, mbrTemp.Mbrpartition_4)
					} else {
						colorize(ColorRed, "Error: La Particion Ya Esta Montada")
					}
				} else {
					colorize(ColorRed, "Error: El Nombre De La Particion Es Invalido")
				}
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
				desmontarParticion(identificador)
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
					encontrado := false
					var rutaMbr string
					for i := 0; i < len(discosMounted); i++ {
						for k := 0; k < len(discosMounted[i].partitionsMounted); k++ {
							if discosMounted[i].partitionsMounted[k].id == identificador {
								encontrado = true
								rutaMbr = discosMounted[i].path
								break
							}
						}
					}

					if encontrado == true {
						mbrTemp := leerMBR(rutaMbr)
						graficarMbr(mbrTemp, ruta)
						fmt.Println("reporte de mbr ", identificador, "-", nombre, "-", ruta)
					} else {
						colorize(ColorRed, "Error: La particion no esta montada")
					}

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

func montarParticion(ruta string, disk discos, partitionTemp partition) {
	erro := true
	for i := 0; i < len(discosMounted); i++ {
		if discosMounted[i].path == ruta {
			discosMounted[i] = disk
			erro = false
			break
		}
	}

	if erro == true {
		discosMounted = append(discosMounted, disk)
	}
	correlativo := 1
	if len(disk.partitionsMounted) != 0 {
		correlativo = disk.partitionsMounted[len(disk.partitionsMounted)-1].correlativo + 1
	}
	id := disk.id + strconv.Itoa(correlativo)
	partitionMountedTemp := partitionMounted{}
	partitionMountedTemp.particion = partitionTemp
	partitionMountedTemp.id = id
	partitionMountedTemp.correlativo = correlativo

	disk.partitionsMounted = append(disk.partitionsMounted, partitionMountedTemp)

	for i := 0; i < len(discosMounted); i++ {
		if discosMounted[i].path == ruta {
			discosMounted[i].partitionsMounted = disk.partitionsMounted
			break
		}
	}

	fmt.Println(partitionMountedTemp.id)
	for k := 0; k < len(discosMounted); k++ {
		for i := 0; i < len(discosMounted[k].partitionsMounted); i++ {
			fmt.Println(discosMounted[k].partitionsMounted[i].id, "-", string(discosMounted[k].partitionsMounted[i].particion.Part_name[:]))
		}
	}
}

func desmontarParticion(id string) {
	erro := true
	for k := 0; k < len(discosMounted); k++ {
		for i := 0; i < len(discosMounted[k].partitionsMounted); i++ {
			if discosMounted[k].partitionsMounted[i].id == id {
				erro = false
				discosMounted[k].partitionsMounted[i] = discosMounted[k].partitionsMounted[len(discosMounted[k].partitionsMounted)-1]
				discosMounted[k].partitionsMounted = discosMounted[k].partitionsMounted[:len(discosMounted[k].partitionsMounted)-1]
				break
			}
		}
	}

	if erro == true {
		colorize(ColorRed, "Error Id No Encontrado")
	} else {
		colorize(ColorBlue, "Particion Desmontada")
	}

	for k := 0; k < len(discosMounted); k++ {
		for i := 0; i < len(discosMounted[k].partitionsMounted); i++ {
			fmt.Println(discosMounted[k].partitionsMounted[i].id, "-", string(discosMounted[k].partitionsMounted[i].particion.Part_name[:]))
		}
	}
}

func generarIdDisco(longitud int) string {
	id := "error"

	if longitud == 0 {
		id = "a"
	} else if longitud == 1 {
		id = "b"
	} else if longitud == 2 {
		id = "c"
	} else if longitud == 3 {
		id = "d"
	} else if longitud == 4 {
		id = "e"
	} else if longitud == 5 {
		id = "f"
	} else if longitud == 6 {
		id = "g"
	} else if longitud == 7 {
		id = "h"
	} else if longitud == 8 {
		id = "i"
	} else if longitud == 9 {
		id = "j"
	} else if longitud == 10 {
		id = "k"
	} else if longitud == 11 {
		id = "l"
	} else if longitud == 12 {
		id = "m"
	} else if longitud == 13 {
		id = "n"
	} else if longitud == 14 {
		id = "o"
	} else if longitud == 15 {
		id = "p"
	} else if longitud == 16 {
		id = "q"
	} else if longitud == 17 {
		id = "r"
	} else if longitud == 18 {
		id = "s"
	} else if longitud == 19 {
		id = "t"
	} else if longitud == 20 {
		id = "u"
	} else if longitud == 21 {
		id = "v"
	} else if longitud == 22 {
		id = "w"
	} else if longitud == 23 {
		id = "x"
	} else if longitud == 24 {
		id = "y"
	} else if longitud == 25 {
		id = "z"
	}
	return id
}

func asignPartition(mbrTemp mbr, partition1 partition) mbr {
	status := [1]byte{65}
	if mbrTemp.Mbrpartition_1.Part_status != status {
		mbrTemp.Mbrpartition_1 = partition1
	} else if mbrTemp.Mbrpartition_2.Part_status != status {
		mbrTemp.Mbrpartition_2 = partition1
	} else if mbrTemp.Mbrpartition_3.Part_status != status {
		mbrTemp.Mbrpartition_3 = partition1
	} else if mbrTemp.Mbrpartition_4.Part_status != status {
		mbrTemp.Mbrpartition_4 = partition1
	}

	return mbrTemp
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

func mountVerify(discosM discos, partitionTemp partition) bool {
	mounted := false
	for i := 0; i < len(discosM.partitionsMounted); i++ {
		if discosM.partitionsMounted[i].particion.Part_name == partitionTemp.Part_name {
			mounted = true
			break
		}
	}
	return mounted
}

func crearParticion(mbrTemp mbr, partition1 partition) int64 {
	position := -1

	partSize := int64(partition1.Part_size)
	sizeMbr := int64(unsafe.Sizeof(mbrTemp))
	sizePart := partition1.Part_size
	sizeFree := mbrTemp.Mbrpartition_1.Part_start - sizeMbr
	ajuste := string(partition1.Part_fit[:])

	status := [1]byte{65}
	if ajuste == "F" || ajuste == "f" {
		if mbrTemp.Mbrpartition_1.Part_status == status {
			if sizeFree >= sizePart {
				position = int(sizeMbr)
			} else {
				if mbrTemp.Mbrpartition_2.Part_status == status {
					sizeFree = mbrTemp.Mbrpartition_2.Part_start - (mbrTemp.Mbrpartition_1.Part_start + partSize)
					if sizeFree >= sizePart {
						position = int(mbrTemp.Mbrpartition_1.Part_start + +partSize)
					} else {
						if mbrTemp.Mbrpartition_3.Part_status == status {
							sizeFree = mbrTemp.Mbrpartition_3.Part_start - (mbrTemp.Mbrpartition_2.Part_start + partSize)
							if sizeFree >= sizePart {
								position = int(mbrTemp.Mbrpartition_2.Part_start + partSize)
							} else {
								if mbrTemp.Mbrpartition_4.Part_status == status {
									sizeFree = mbrTemp.Mbrpartition_4.Part_start - (mbrTemp.Mbrpartition_3.Part_start + partSize)
									if sizeFree >= sizePart {
										position = int(mbrTemp.Mbrpartition_3.Part_start + partSize)
									} else {
										sizeFree = mbrTemp.Mbrtamaño - (mbrTemp.Mbrpartition_4.Part_start + partSize)
										if sizeFree >= sizePart {
											position = int(mbrTemp.Mbrpartition_4.Part_start + partSize)
										}
									}
								} else {
									if sizePart <= (mbrTemp.Mbrtamaño - sizeMbr) {
										position = int(mbrTemp.Mbrpartition_3.Part_start + partSize)
									}
								}
							}
						} else {
							if sizePart <= (mbrTemp.Mbrtamaño - sizeMbr) {
								position = int(mbrTemp.Mbrpartition_2.Part_start + partSize)
							}
						}
					}
				} else {
					if sizePart <= (mbrTemp.Mbrtamaño - sizeMbr) {
						position = int(mbrTemp.Mbrpartition_1.Part_start + partSize)
					}
				}
			}
		} else {
			if sizePart < (mbrTemp.Mbrtamaño - sizeMbr) {
				position = int(sizeMbr)
			}
		}
	} else if ajuste == "M" || ajuste == "m" {
		freeSpace := int64(0)

		if mbrTemp.Mbrpartition_1.Part_status == status {
			sizeFree = mbrTemp.Mbrpartition_1.Part_start - sizeMbr
			if sizeFree >= sizePart {
				position = int(sizeMbr)
				freeSpace = int64(int(sizeFree) - (position + int(partSize)))
			}
		}

		if mbrTemp.Mbrpartition_2.Part_status == status {
			sizeFree = mbrTemp.Mbrpartition_2.Part_start - (mbrTemp.Mbrpartition_1.Part_start + sizePart)
			if sizeFree >= sizePart {
				if freeSpace == 0 {
					position = int(mbrTemp.Mbrpartition_1.Part_start + sizePart)
					freeSpace = int64(int(sizeFree) - (position + int(partSize)))
				} else {
					espacioLibre := sizeFree - (mbrTemp.Mbrpartition_1.Part_start + sizePart + partSize)
					if espacioLibre < freeSpace {
						freeSpace = espacioLibre
						position = int(mbrTemp.Mbrpartition_1.Part_start + sizePart)
					}
				}
			}
		} else {
			if position == -1 {
				espacioLibre := sizeFree - (mbrTemp.Mbrpartition_1.Part_start + sizePart + partSize)
				freeSpace = espacioLibre
				position = int(mbrTemp.Mbrpartition_1.Part_start + sizePart)
			}
		}

		if mbrTemp.Mbrpartition_3.Part_status == status {
			sizeFree = mbrTemp.Mbrpartition_3.Part_start - (mbrTemp.Mbrpartition_2.Part_start + sizePart)
			if sizeFree >= sizePart {
				if freeSpace == 0 {
					position = int(mbrTemp.Mbrpartition_2.Part_start + sizePart)
					freeSpace = int64(int(sizeFree) - (position + int(partSize)))
				} else {
					espacioLibre := sizeFree - (mbrTemp.Mbrpartition_2.Part_start + sizePart + partSize)
					if espacioLibre < freeSpace {
						freeSpace = espacioLibre
						position = int(mbrTemp.Mbrpartition_2.Part_start + sizePart)
					}
				}
			}
		} else {
			if position == -1 {
				espacioLibre := sizeFree - (mbrTemp.Mbrpartition_2.Part_start + sizePart + partSize)
				freeSpace = espacioLibre
				position = int(mbrTemp.Mbrpartition_2.Part_start + sizePart)
			}
		}

		if mbrTemp.Mbrpartition_4.Part_status == status {
			sizeFree = mbrTemp.Mbrpartition_4.Part_start - (mbrTemp.Mbrpartition_3.Part_start + sizePart)
			if sizeFree >= sizePart {
				if freeSpace == 0 {
					position = int(mbrTemp.Mbrpartition_3.Part_start + sizePart)
					freeSpace = int64(int(sizeFree) - (position + int(partSize)))
				} else {
					espacioLibre := sizeFree - (mbrTemp.Mbrpartition_3.Part_start + sizePart + partSize)
					if espacioLibre < freeSpace {
						freeSpace = espacioLibre
						position = int(mbrTemp.Mbrpartition_3.Part_start + sizePart)
					}
				}
			}
		} else {
			if position == -1 {
				espacioLibre := sizeFree - (mbrTemp.Mbrpartition_3.Part_start + sizePart + partSize)
				freeSpace = espacioLibre
				position = int(mbrTemp.Mbrpartition_3.Part_start + sizePart)
			}
		}
	} else if ajuste == "W" || ajuste == "w" {
		freeSpace := int64(0)

		if mbrTemp.Mbrpartition_1.Part_status == status {
			sizeFree = mbrTemp.Mbrpartition_1.Part_start - sizeMbr
			if sizeFree >= sizePart {
				position = int(sizeMbr)
				freeSpace = int64(int(sizeFree) - (position + int(partSize)))
			}
		}

		if mbrTemp.Mbrpartition_2.Part_status == status {
			sizeFree = mbrTemp.Mbrpartition_2.Part_start - (mbrTemp.Mbrpartition_1.Part_start + sizePart)
			if sizeFree >= sizePart {
				if freeSpace == 0 {
					position = int(mbrTemp.Mbrpartition_1.Part_start + sizePart)
					freeSpace = int64(int(sizeFree) - (position + int(partSize)))
				} else {
					espacioLibre := sizeFree - (mbrTemp.Mbrpartition_1.Part_start + sizePart + partSize)
					if espacioLibre > freeSpace {
						freeSpace = espacioLibre
						position = int(mbrTemp.Mbrpartition_1.Part_start + sizePart)
					}
				}
			}
		} else {
			if position == -1 {
				espacioLibre := sizeFree - (mbrTemp.Mbrpartition_1.Part_start + sizePart + partSize)
				freeSpace = espacioLibre
				position = int(mbrTemp.Mbrpartition_1.Part_start + sizePart)
			}
		}

		if mbrTemp.Mbrpartition_3.Part_status == status {
			sizeFree = mbrTemp.Mbrpartition_3.Part_start - (mbrTemp.Mbrpartition_2.Part_start + sizePart)
			if sizeFree >= sizePart {
				if freeSpace == 0 {
					position = int(mbrTemp.Mbrpartition_2.Part_start + sizePart)
					freeSpace = int64(int(sizeFree) - (position + int(partSize)))
				} else {
					espacioLibre := sizeFree - (mbrTemp.Mbrpartition_2.Part_start + sizePart + partSize)
					if espacioLibre > freeSpace {
						freeSpace = espacioLibre
						position = int(mbrTemp.Mbrpartition_2.Part_start + sizePart)
					}
				}
			}
		} else {
			if position == -1 {
				espacioLibre := sizeFree - (mbrTemp.Mbrpartition_1.Part_start + sizePart + partSize)
				freeSpace = espacioLibre
				position = int(mbrTemp.Mbrpartition_2.Part_start + sizePart)
			}
		}

		if mbrTemp.Mbrpartition_4.Part_status == status {
			sizeFree = mbrTemp.Mbrpartition_4.Part_start - (mbrTemp.Mbrpartition_3.Part_start + sizePart)
			if sizeFree >= sizePart {
				if freeSpace == 0 {
					position = int(mbrTemp.Mbrpartition_3.Part_start + sizePart)
					freeSpace = int64(int(sizeFree) - (position + int(partSize)))
				} else {
					espacioLibre := sizeFree - (mbrTemp.Mbrpartition_3.Part_start + sizePart + partSize)
					if espacioLibre > freeSpace {
						freeSpace = espacioLibre
						position = int(mbrTemp.Mbrpartition_3.Part_start + sizePart)
					}
				}
			}
		} else {
			if position == -1 {
				espacioLibre := sizeFree - (mbrTemp.Mbrpartition_1.Part_start + sizePart + partSize)
				freeSpace = espacioLibre
				position = int(mbrTemp.Mbrpartition_3.Part_start + sizePart)
			}
		}
	}
	return int64(position)
}

func sortPartitions(mbrTemp mbr) mbr {
	var cambio partition
	listaPartitions := []partition{mbrTemp.Mbrpartition_1, mbrTemp.Mbrpartition_2, mbrTemp.Mbrpartition_3, mbrTemp.Mbrpartition_4}
	n := len(listaPartitions)

	status := [1]byte{65}

	for k := 0; k < n; k++ {
		if listaPartitions[k].Part_status != status {
			listaPartitions[k].Part_start = mbrTemp.Mbrtamaño
		}
	}

	for k := 1; k < n; k++ {
		for i := 0; i < (n - k); i++ {
			if listaPartitions[i].Part_start > listaPartitions[i+1].Part_start {
				cambio = listaPartitions[i]
				listaPartitions[i] = listaPartitions[i+1]
				listaPartitions[i+1] = cambio
			}
		}
	}

	mbrTemp.Mbrpartition_1 = listaPartitions[0]
	mbrTemp.Mbrpartition_2 = listaPartitions[1]
	mbrTemp.Mbrpartition_3 = listaPartitions[2]
	mbrTemp.Mbrpartition_4 = listaPartitions[3]

	return mbrTemp
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

func graficarMbr(mbrTemp mbr, ruta string) {
	status := [1]byte{65}
	str := "digraph {\n"
	str = str + "tbl [shape=plaintext\n"
	str = str + "label=<\n"
	str = str + "<table border='0' cellborder='1' color='blue' cellspacing='0'>\n"
	str = str + "	<tr><td>Nombre</td><td>Valor</td></tr>\n"
	str = str + "	<tr><td>mbr_tamaño</td><td>" + strconv.FormatInt(mbrTemp.Mbrtamaño, 10) + "</td></tr>\n"
	str = str + "	<tr><td>mbr_fecha_creacion</td><td>" + string(mbrTemp.Mbrfechacreacion[:]) + "</td></tr>\n"
	str = str + "	<tr><td>mbr_disk_signature</td><td>" + strconv.FormatInt(mbrTemp.Mbrdisksignature, 10) + "</td></tr>\n"

	if mbrTemp.Mbrpartition_1.Part_status == status {
		label := mbrTemp.Mbrpartition_1.Part_name[:]
		str = str + "	<tr><td>part_status_1</td><td>" + string(mbrTemp.Mbrpartition_1.Part_status[:]) + "</td></tr>\n"
		str = str + "	<tr><td>part_type_1</td><td>" + string(mbrTemp.Mbrpartition_1.Part_type[:]) + "</td></tr>\n"
		str = str + "	<tr><td>part_fit_1</td><td>" + string(mbrTemp.Mbrpartition_1.Part_fit[:]) + "</td></tr>\n"
		str = str + "	<tr><td>part_start_1</td><td>" + strconv.FormatInt(mbrTemp.Mbrpartition_1.Part_start, 10) + "</td></tr>\n"
		str = str + "	<tr><td>part_size_1</td><td>" + strconv.FormatInt(mbrTemp.Mbrpartition_1.Part_size, 10) + "</td></tr>\n"
		str = str + "	<tr><td>part_name_1</td><td>" + string(label[:clen(label)]) + "</td></tr>\n"
	} else {
		str = str + "	<tr><td>part_status_1</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_type_1</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_fit_1</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_start_1</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_size_1</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_name_1</td><td> -- </td></tr>\n"
	}

	if mbrTemp.Mbrpartition_2.Part_status == status {
		label := mbrTemp.Mbrpartition_2.Part_name[:]
		str = str + "	<tr><td>part_status_2</td><td>" + string(mbrTemp.Mbrpartition_2.Part_status[:]) + "</td></tr>\n"
		str = str + "	<tr><td>part_type_2</td><td>" + string(mbrTemp.Mbrpartition_2.Part_type[:]) + "</td></tr>\n"
		str = str + "	<tr><td>part_fit_2</td><td>" + string(mbrTemp.Mbrpartition_2.Part_fit[:]) + "</td></tr>\n"
		str = str + "	<tr><td>part_start_2</td><td>" + strconv.FormatInt(mbrTemp.Mbrpartition_2.Part_start, 10) + "</td></tr>\n"
		str = str + "	<tr><td>part_size_2</td><td>" + strconv.FormatInt(mbrTemp.Mbrpartition_2.Part_size, 10) + "</td></tr>\n"
		str = str + "	<tr><td>part_name_2</td><td>" + string(label[:clen(label)]) + "</td></tr>\n"
	} else {
		str = str + "	<tr><td>part_status_2</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_type_2</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_fit_2</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_start_2</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_size_2</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_name_2</td><td> -- </td></tr>\n"
	}

	if mbrTemp.Mbrpartition_3.Part_status == status {
		label := mbrTemp.Mbrpartition_3.Part_name[:]
		str = str + "	<tr><td>part_status_3</td><td>" + string(mbrTemp.Mbrpartition_3.Part_status[:]) + "</td></tr>\n"
		str = str + "	<tr><td>part_type_3</td><td>" + string(mbrTemp.Mbrpartition_3.Part_type[:]) + "</td></tr>\n"
		str = str + "	<tr><td>part_fit_3</td><td>" + string(mbrTemp.Mbrpartition_3.Part_fit[:]) + "</td></tr>\n"
		str = str + "	<tr><td>part_start_3</td><td>" + strconv.FormatInt(mbrTemp.Mbrpartition_3.Part_start, 10) + "</td></tr>\n"
		str = str + "	<tr><td>part_size_3</td><td>" + strconv.FormatInt(mbrTemp.Mbrpartition_3.Part_size, 10) + "</td></tr>\n"
		str = str + "	<tr><td>part_name_3</td><td>" + string(label[:clen(label)]) + "</td></tr>\n"
	} else {
		str = str + "	<tr><td>part_status_3</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_type_3</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_fit_3</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_start_3</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_size_3</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_name_3</td><td> -- </td></tr>\n"
	}

	if mbrTemp.Mbrpartition_4.Part_status == status {
		label := mbrTemp.Mbrpartition_4.Part_name[:]
		str = str + "	<tr><td>part_status_4</td><td>" + string(mbrTemp.Mbrpartition_4.Part_status[:]) + "</td></tr>\n"
		str = str + "	<tr><td>part_type_4</td><td>" + string(mbrTemp.Mbrpartition_4.Part_type[:]) + "</td></tr>\n"
		str = str + "	<tr><td>part_fit_4</td><td>" + string(mbrTemp.Mbrpartition_4.Part_fit[:]) + "</td></tr>\n"
		str = str + "	<tr><td>part_start_4</td><td>" + strconv.FormatInt(mbrTemp.Mbrpartition_4.Part_start, 10) + "</td></tr>\n"
		str = str + "	<tr><td>part_size_4</td><td>" + strconv.FormatInt(mbrTemp.Mbrpartition_4.Part_size, 10) + "</td></tr>\n"
		str = str + "	<tr><td>part_name_4</td><td>" + string(label[:clen(label)]) + "</td></tr>\n"
	} else {
		str = str + "	<tr><td>part_status_4</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_type_4</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_fit_4</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_start_4</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_size_4</td><td> -- </td></tr>\n"
		str = str + "	<tr><td>part_name_4</td><td> -- </td></tr>\n"
	}

	str = str + "</table>\n"
	str = str + ">];\n"
	str = str + "}\n"

	fmt.Println(str)

	b := []byte(str)
	erro := ioutil.WriteFile("reporteMbr.dot", b, 0664)
	if erro != nil {
		log.Fatal(erro)
	}

	path, _ := exec.LookPath("dot")
	cmd, _ := exec.Command(path, "-Tpng", "reporteMbr.dot").Output()
	mode := int(0777)
	ioutil.WriteFile("outfile.png", cmd, os.FileMode(mode))
}

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}

func colorizefn(color Color, message string) {
	fmt.Print(string(color), message, string(ColorReset))
}
