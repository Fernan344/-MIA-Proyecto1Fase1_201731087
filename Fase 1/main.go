package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math"
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
fdisk –Size->1000 –path->/home/thefernan/Desktop/disk3.dsk –name->Particion4
Mkdisk -Size->3000 –unit->K -path->/home/thefernan/Desktop/disk3.dsk
mount -path->/home/thefernan/Desktop/disk3.dsk -name->particion3
rep -id->vda1 -Path->/home/user/reports/reporte2.pdf -name->disk

Mkdisk -Size->3000 –unit->K -path->/home/thefernan/Desktop/disk3.dsk
fdisk –Size->700 –path->/home/thefernan/Desktop/disk3.dsk –name->Particion4 -fit->wf -type->p
fdisk –Size->200 –path->/home/thefernan/Desktop/disk3.dsk –name->Particion2 -fit->wf -type->e
fdisk –Size->400 –path->/home/thefernan/Desktop/disk3.dsk –name->Particion3 -fit->wf
fdisk –Size->1000 –path->/home/thefernan/Desktop/disk3.dsk –name->Particion1 -fit->wf
rep -id->vda1 -Path->/home/user/reports/reporte2.pdf -name->disk
*/

var discosMounted []discos
var session []string

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

type superbloque struct {
	S_filesystem_type   int64
	S_inodes_count      int64
	S_blocks_count      int64
	S_free_blocks_count int64
	S_free_inodes_count int64
	S_mtime             [54]byte
	S_umtime            [54]byte
	S_mnt_count         int64
	S_magic             int64
	S_inode_size        int64
	S_block_size        int64
	S_first_ino         int64
	S_first_blo         int64
	S_bm_inode_start    int64
	S_bm_block_start    int64
	S_inode_start       int64
	S_block_start       int64
}

type journaling struct {
	Journal_tipo_operacion [1]byte
	Journal_tipo           [1]byte
	Journal_nombre         [12]byte
	Journal_contenido      [1]byte
	Journal_fecha          [54]byte
	Journal_propietario    int64
	Journal_permisos       int64
}

type inodo struct {
	I_uid   int64
	I_gid   int64
	I_size  int64
	I_atime [54]byte
	I_ctime [54]byte
	I_mtime [54]byte
	I_block [15]int64
	I_type  [1]byte
	I_perm  int64
}

type carpeta struct {
	B_content [4]content
}

type content struct {
	B_name  [12]byte
	B_inodo int32
}

type archivo struct {
	B_content [64]byte
}

type apuntador struct {
	B_content [16]int32
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
		commandArray = strings.Split(commandArray[0], " -")
	} else {
		commandArray = strings.Split(comando, " -")
	}
	if commandArray[0] != "" {
		executeComand(commandArray)
	}
}

func executeComand(commandArray []string) {
	data := strings.ToLower(commandArray[0])
	comment := strings.Split(data, "")
	if comment[0] != "#" {
		if data == "exec" {
			parametro := commandArray[1]
			param := strings.ToLower(parametro)
			caracteres := strings.Split(param, "")

			if caracteres[0] == "p" && caracteres[1] == "a" && caracteres[2] == "t" && caracteres[3] == "h" && caracteres[4] == "-" && caracteres[5] == ">" {
				paramsParts := strings.Split(parametro, "->")
				path := paramsParts[1]
				b, err := ioutil.ReadFile(path)
				if err != nil {
					colorize(ColorRed, "Error Archivo No Encontrado")
				} else {
					str := string(b) // convert content to a 'string'
					lineas := strings.Split(str, "\n")
					colorize(ColorYellow, "Corriendo Scripts")
					for i := 0; i < len(lineas)-1; i++ {
						colorize(ColorReset, lineas[i])
						lineaDeComandos(lineas[i])
					}
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

				if caracteres[0] == "s" && caracteres[1] == "i" && caracteres[2] == "z" && caracteres[3] == "e" && caracteres[4] == "-" && caracteres[5] == ">" {
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
				} else if caracteres[0] == "u" && caracteres[1] == "n" && caracteres[2] == "i" && caracteres[3] == "t" && caracteres[4] == "-" && caracteres[5] == ">" {
					parametros := strings.Split(command, "->")
					if parametros[1] == "K" || parametros[1] == "k" || parametros[1] == "M" || parametros[1] == "m" {
						unit = true
						unidad = parametros[1]

					} else {
						other = true
					}
				} else if caracteres[0] == "p" && caracteres[1] == "a" && caracteres[2] == "t" && caracteres[3] == "h" && caracteres[4] == "-" && caracteres[5] == ">" {
					path = true
					parametros := strings.Split(commandArray[i], "->")
					ruta = parametros[1]
				} else if caracteres[0] == "f" && caracteres[1] == "i" && caracteres[2] == "t" && caracteres[3] == "-" && caracteres[4] == ">" {
					parametros := strings.Split(command, "->")
					if parametros[1] == "BF" || parametros[1] == "FF" || parametros[1] == "WF" || parametros[1] == "wf" || parametros[1] == "ff" || parametros[1] == "bf" {
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
					createPath(ruta)
					file, err = os.Create(ruta)
					defer file.Close()
					if err != nil {
						colorize(ColorRed, "Error En La Creacion De La ruta")
					}
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
			if caracteres[0] == "p" && caracteres[1] == "a" && caracteres[2] == "t" && caracteres[3] == "h" && caracteres[4] == "-" && caracteres[5] == ">" {
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

				if caracteres[0] == "p" && caracteres[1] == "a" && caracteres[2] == "t" && caracteres[3] == "h" && caracteres[4] == "-" && caracteres[5] == ">" {
					path = true
					parametros := strings.Split(commandArray[i], "->")
					ruta = parametros[1]
				} else if caracteres[0] == "n" && caracteres[1] == "a" && caracteres[2] == "m" && caracteres[3] == "e" && caracteres[4] == "-" && caracteres[5] == ">" {
					name = true
					parametros := strings.Split(command, "->")
					nombre = parametros[1]
				} else if caracteres[0] == "u" && caracteres[1] == "n" && caracteres[2] == "i" && caracteres[3] == "t" && caracteres[4] == "-" && caracteres[5] == ">" {
					parametros := strings.Split(command, "->")
					if parametros[1] == "K" || parametros[1] == "k" || parametros[1] == "M" || parametros[1] == "m" {
						unit = true
						unidad = parametros[1]
					} else {
						other = true
					}
				} else if caracteres[0] == "t" && caracteres[1] == "y" && caracteres[2] == "p" && caracteres[3] == "e" && caracteres[4] == "-" && caracteres[5] == ">" {
					parametros := strings.Split(command, "->")
					if parametros[1] == "p" || parametros[1] == "e" || parametros[1] == "l" {
						tipe = true
						tipo = parametros[1]
					} else {
						other = true
					}
				} else if caracteres[0] == "f" && caracteres[1] == "i" && caracteres[2] == "t" && caracteres[3] == "-" && caracteres[4] == ">" {
					parametros := strings.Split(command, "->")
					if parametros[1] == "BF" || parametros[1] == "FF" || parametros[1] == "WF" || parametros[1] == "bf" || parametros[1] == "ff" || parametros[1] == "wf" {
						fit = true
						ajuste = parametros[1]
					} else {
						other = true
					}
				} else if caracteres[0] == "d" && caracteres[1] == "e" && caracteres[2] == "l" && caracteres[3] == "e" && caracteres[4] == "t" && caracteres[5] == "e" {
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
				} else if caracteres[0] == "a" && caracteres[1] == "d" && caracteres[2] == "d" {
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
				} else if caracteres[0] == "s" && caracteres[1] == "i" && caracteres[2] == "z" && caracteres[3] == "e" && caracteres[4] == "-" && caracteres[5] == ">" {
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
					colorize(ColorRed, "Error Comando Incorrecto")
					other = true
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
					newPart := partition{}
					mbrTemp := leerMBR(ruta)

					partName := strings.ToLower(string(mbrTemp.Mbrpartition_1.Part_name[:clen(mbrTemp.Mbrpartition_1.Part_name[:])]))

					if partName == nombre {
						mbrTemp.Mbrpartition_1 = newPart
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
					} else {
						partName := strings.ToLower(string(mbrTemp.Mbrpartition_2.Part_name[:clen(mbrTemp.Mbrpartition_2.Part_name[:])]))

						if partName == nombre {
							mbrTemp.Mbrpartition_2 = newPart
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
						} else {
							partName := strings.ToLower(string(mbrTemp.Mbrpartition_3.Part_name[:clen(mbrTemp.Mbrpartition_3.Part_name[:])]))
							if partName == nombre {
								mbrTemp.Mbrpartition_3 = newPart
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
							} else {
								partName := strings.ToLower(string(mbrTemp.Mbrpartition_4.Part_name[:clen(mbrTemp.Mbrpartition_4.Part_name[:])]))

								if partName == nombre {
									mbrTemp.Mbrpartition_4 = newPart
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
								} else {
									colorize(ColorRed, "Error No Existe La Particion")
								}
							}
						}
					}

					fmt.Println("eliminando ", nombre, "-", ruta, "-", eliminar)
				} else if primerComando == "add" && add == true {
					mbrTemp := leerMBR(ruta)
					status := [1]byte{65}

					if unidad == "m" || unidad == "M" {
						agragar = agragar * 1024 * 1024
					} else if unidad == "k" || unidad == "K" {
						agragar = agragar * 1024
					}
					partName := strings.ToLower(string(mbrTemp.Mbrpartition_1.Part_name[:clen(mbrTemp.Mbrpartition_1.Part_name[:])]))
					freeSize := int64(0)
					if partName == nombre {
						if mbrTemp.Mbrpartition_2.Part_status == status {
							freeSize = mbrTemp.Mbrpartition_2.Part_start - (mbrTemp.Mbrpartition_1.Part_start + mbrTemp.Mbrpartition_1.Part_size)
						} else {
							freeSize = mbrTemp.Mbrtamaño - (mbrTemp.Mbrpartition_1.Part_start + mbrTemp.Mbrpartition_1.Part_size)
						}
						spaceOcuped := (mbrTemp.Mbrpartition_1.Part_size) + (int64(agragar))

						if spaceOcuped > 0 && int64(agragar) <= freeSize {
							mbrTemp.Mbrpartition_1.Part_size = (mbrTemp.Mbrpartition_1.Part_size + int64(agragar))
							actualizarMbr(ruta, mbrTemp)
						} else {
							colorize(ColorRed, "Error En La Expansion/Reduccion parametro de adicion invalido")
						}
					} else {
						partName := strings.ToLower(string(mbrTemp.Mbrpartition_2.Part_name[:clen(mbrTemp.Mbrpartition_2.Part_name[:])]))
						freeSize := int64(0)
						if partName == nombre {
							if mbrTemp.Mbrpartition_3.Part_status == status {
								freeSize = mbrTemp.Mbrpartition_3.Part_start - (mbrTemp.Mbrpartition_2.Part_start + mbrTemp.Mbrpartition_2.Part_size)
							} else {
								freeSize = mbrTemp.Mbrtamaño - (mbrTemp.Mbrpartition_2.Part_start + mbrTemp.Mbrpartition_2.Part_size)
							}
							spaceOcuped := (mbrTemp.Mbrpartition_2.Part_size) + (int64(agragar))

							if spaceOcuped > 0 && int64(agragar) <= freeSize {
								mbrTemp.Mbrpartition_2.Part_size = (mbrTemp.Mbrpartition_2.Part_size + int64(agragar))
								actualizarMbr(ruta, mbrTemp)
							} else {
								colorize(ColorRed, "Error En La Expansion/Reduccion parametro de adicion invalido")
							}

						} else {
							partName := strings.ToLower(string(mbrTemp.Mbrpartition_3.Part_name[:clen(mbrTemp.Mbrpartition_3.Part_name[:])]))
							freeSize := int64(0)
							if partName == nombre {
								if mbrTemp.Mbrpartition_4.Part_status == status {
									freeSize = mbrTemp.Mbrpartition_4.Part_start - (mbrTemp.Mbrpartition_3.Part_start + mbrTemp.Mbrpartition_3.Part_size)
								} else {
									freeSize = mbrTemp.Mbrtamaño - (mbrTemp.Mbrpartition_3.Part_start + mbrTemp.Mbrpartition_3.Part_size)
								}
								spaceOcuped := (mbrTemp.Mbrpartition_3.Part_size) + (int64(agragar))

								if spaceOcuped > 0 && int64(agragar) <= freeSize {
									mbrTemp.Mbrpartition_3.Part_size = (mbrTemp.Mbrpartition_3.Part_size + int64(agragar))
									actualizarMbr(ruta, mbrTemp)
								} else {
									colorize(ColorRed, "Error En La Expansion/Reduccion parametro de adicion invalido")
								}

							} else {
								partName := strings.ToLower(string(mbrTemp.Mbrpartition_4.Part_name[:clen(mbrTemp.Mbrpartition_4.Part_name[:])]))
								freeSize := int64(0)
								if partName == nombre {
									freeSize = mbrTemp.Mbrtamaño - (mbrTemp.Mbrpartition_4.Part_start + mbrTemp.Mbrpartition_4.Part_size)
									spaceOcuped := (mbrTemp.Mbrpartition_4.Part_size) + (int64(agragar))

									if spaceOcuped > 0 && int64(agragar) <= freeSize {
										mbrTemp.Mbrpartition_4.Part_size = (mbrTemp.Mbrpartition_4.Part_size + int64(agragar))
										actualizarMbr(ruta, mbrTemp)
									} else {
										colorize(ColorRed, "Error En La Expansion/Reduccion parametro de adicion invalido")
									}

								} else {
									colorize(ColorRed, "Error En La Expansion/Reduccion Particion No Encontrada")
								}
							}
						}
					}
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
									copy(ebrTemp.Part_name[:], "Default")
									ebrTemp.Part_size = 0

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

				if caracteres[0] == "p" && caracteres[1] == "a" && caracteres[2] == "t" && caracteres[3] == "h" && caracteres[4] == "-" && caracteres[5] == ">" {
					path = true
					parametros := strings.Split(commandArray[i], "->")
					ruta = parametros[1]
				} else if caracteres[0] == "n" && caracteres[1] == "a" && caracteres[2] == "m" && caracteres[3] == "e" && caracteres[4] == "-" && caracteres[5] == ">" {
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
						particion := leerParticion(ruta, mbrTemp.Mbrpartition_1.Part_start)
						if particion == true {
							fmt.Println("sihay")
							sb := leerSB(ruta, mbrTemp.Mbrpartition_1.Part_start)
							sb = updateSuperBlock(sb, "montarParticion")
							escribirSuperBloque(ruta, sb, mbrTemp.Mbrpartition_1.Part_start)
						}
					} else {
						colorize(ColorRed, "Error: La Particion Ya Esta Montada")
					}
				} else if mbrTemp.Mbrpartition_2.Part_name == compareName {
					mounted := mountVerify(disk, mbrTemp.Mbrpartition_2)
					if mounted == false {
						montarParticion(ruta, disk, mbrTemp.Mbrpartition_2)
						particion := leerParticion(ruta, mbrTemp.Mbrpartition_2.Part_start)
						if particion == true {
							sb := leerSB(ruta, mbrTemp.Mbrpartition_2.Part_start)
							sb = updateSuperBlock(sb, "montarParticion")
							escribirSuperBloque(ruta, sb, mbrTemp.Mbrpartition_2.Part_start)
						}
					} else {
						colorize(ColorRed, "Error: La Particion Ya Esta Montada")
					}
				} else if mbrTemp.Mbrpartition_3.Part_name == compareName {
					mounted := mountVerify(disk, mbrTemp.Mbrpartition_3)
					if mounted == false {
						montarParticion(ruta, disk, mbrTemp.Mbrpartition_3)
						particion := leerParticion(ruta, mbrTemp.Mbrpartition_3.Part_start)
						if particion == true {
							sb := leerSB(ruta, mbrTemp.Mbrpartition_3.Part_start)
							sb = updateSuperBlock(sb, "montarParticion")
							escribirSuperBloque(ruta, sb, mbrTemp.Mbrpartition_3.Part_start)
						}
					} else {
						colorize(ColorRed, "Error: La Particion Ya Esta Montada")
					}
				} else if mbrTemp.Mbrpartition_4.Part_name == compareName {
					mounted := mountVerify(disk, mbrTemp.Mbrpartition_4)
					if mounted == false {
						montarParticion(ruta, disk, mbrTemp.Mbrpartition_4)
						particion := leerParticion(ruta, mbrTemp.Mbrpartition_4.Part_start)
						if particion == true {
							sb := leerSB(ruta, mbrTemp.Mbrpartition_4.Part_start)
							sb = updateSuperBlock(sb, "montarParticion")
							escribirSuperBloque(ruta, sb, mbrTemp.Mbrpartition_4.Part_start)
						}
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
				if caracteres[0] == "i" && caracteres[1] == "d" && caracteres[2] == "-" && caracteres[3] == ">" {
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
				if caracteres[0] == "p" && caracteres[1] == "a" && caracteres[2] == "t" && caracteres[3] == "h" && caracteres[4] == "-" && caracteres[5] == ">" {
					path = true
					parametros := strings.Split(command, "->")
					ruta = parametros[1]
				} else if caracteres[0] == "n" && caracteres[1] == "a" && caracteres[2] == "m" && caracteres[3] == "e" && caracteres[4] == "-" && caracteres[5] == ">" {
					name = true
					parametros := strings.Split(command, "->")
					if parametros[1] == "mbr" || parametros[1] == "disk" || parametros[1] == "inode" || parametros[1] == "block" || parametros[1] == "bm_inode" {
						nombre = parametros[1]
					} else {
						other = true
					}
				} else if caracteres[0] == "i" && caracteres[1] == "d" && caracteres[2] == "-" && caracteres[3] == ">" {
					id = true
					parametros := strings.Split(command, "->")
					identificador = parametros[1]
				} else {
					other = true
				}
			}

			if other == false && id == true && name == true && path == true {
				encontrado := false
				var rutaMbr string
				var id string
				particion := partition{}
				for i := 0; i < len(discosMounted); i++ {
					for k := 0; k < len(discosMounted[i].partitionsMounted); k++ {
						if discosMounted[i].partitionsMounted[k].id == identificador {
							encontrado = true
							rutaMbr = discosMounted[i].path
							id = discosMounted[i].id
							particion = discosMounted[i].partitionsMounted[k].particion
							break
						}
					}
				}
				if nombre == "disk" {
					if encontrado == true {
						mbrTemp := leerMBR(rutaMbr)
						graficarDisco(mbrTemp, ruta, id)
						fmt.Println("reporte de disco ", identificador, "-", nombre, "-", ruta)
					} else {
						colorize(ColorRed, "Error: La particion no esta montada")
					}
				} else if nombre == "mbr" {
					if encontrado == true {
						mbrTemp := leerMBR(rutaMbr)
						graficarMbr(mbrTemp, ruta)
						fmt.Println("reporte de mbr ", identificador, "-", nombre, "-", ruta)
					} else {
						colorize(ColorRed, "Error: La particion no esta montada")
					}

				} else if nombre == "inode" {
					if encontrado == true {
						sb := leerSB(rutaMbr, particion.Part_start)
						reporteInodos(sb, ruta, rutaMbr)
						fmt.Println("reporte de Inodos ", identificador, "-", nombre, "-", ruta)
					} else {
						colorize(ColorRed, "Error: La particion no esta montada")
					}
				} else if nombre == "block" {
					if encontrado == true {
						sb := leerSB(rutaMbr, particion.Part_start)
						reporteBloques(sb, ruta, rutaMbr)
						fmt.Println("reporte de Bloques ", identificador, "-", nombre, "-", ruta)
					} else {
						colorize(ColorRed, "Error: La particion no esta montada")
					}
				} else if nombre == "bm_inode" {
					if encontrado == true {
						sb := leerSB(rutaMbr, particion.Part_start)
						reporteDeBipmapDeInodos(sb, ruta, rutaMbr)
						fmt.Println("reporte de BitMap Inodos ", identificador, "-", nombre, "-", ruta)
					} else {
						colorize(ColorRed, "Error: La particion no esta montada")
					}
				}

			} else {
				fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
			}
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************
			/*********************************************************************************/
		} else if data == "mkfs" {
			id := false
			tipe := false
			other := false
			var identificador string
			var tipo string

			for i := 1; i < len(commandArray); i++ {
				command := strings.ToLower(commandArray[i])
				caracteres := strings.Split(command, "")
				if caracteres[0] == "t" && caracteres[1] == "y" && caracteres[2] == "p" && caracteres[3] == "e" && caracteres[4] == "-" && caracteres[5] == ">" {
					tipe = true
					parametros := strings.Split(command, "->")
					if parametros[1] == "fast" || parametros[1] == "full" {
						tipo = parametros[1]
					} else {
						other = true
					}
				} else if caracteres[0] == "i" && caracteres[1] == "d" && caracteres[2] == "-" && caracteres[3] == ">" {
					id = true
					parametros := strings.Split(command, "->")
					identificador = parametros[1]
				} else {
					other = true
				}
			}
			if other == false && id == true {

				encontrado := false
				var rutaMbr string
				var id string
				var particion string
				for i := 0; i < len(discosMounted); i++ {
					for k := 0; k < len(discosMounted[i].partitionsMounted); k++ {
						if discosMounted[i].partitionsMounted[k].id == identificador {
							encontrado = true
							rutaMbr = discosMounted[i].path
							id = discosMounted[i].id
							particion = string(discosMounted[i].partitionsMounted[k].particion.Part_name[:clen(discosMounted[i].partitionsMounted[k].particion.Part_name[:])])
							break
						}
					}
				}
				mbrTemp := leerMBR(rutaMbr)
				var tamañoParticion int64
				var inicio int64

				if string(mbrTemp.Mbrpartition_1.Part_name[:clen(mbrTemp.Mbrpartition_1.Part_name[:])]) == particion {
					tamañoParticion = mbrTemp.Mbrpartition_1.Part_size
					inicio = mbrTemp.Mbrpartition_1.Part_start
				} else if string(mbrTemp.Mbrpartition_2.Part_name[:clen(mbrTemp.Mbrpartition_2.Part_name[:])]) == particion {
					tamañoParticion = mbrTemp.Mbrpartition_2.Part_size
					inicio = mbrTemp.Mbrpartition_2.Part_start
				} else if string(mbrTemp.Mbrpartition_3.Part_name[:clen(mbrTemp.Mbrpartition_3.Part_name[:])]) == particion {
					tamañoParticion = mbrTemp.Mbrpartition_3.Part_size
					inicio = mbrTemp.Mbrpartition_3.Part_start
				} else if string(mbrTemp.Mbrpartition_4.Part_name[:clen(mbrTemp.Mbrpartition_4.Part_name[:])]) == particion {
					tamañoParticion = mbrTemp.Mbrpartition_4.Part_size
					inicio = mbrTemp.Mbrpartition_4.Part_start
				}
				journal := journaling{}
				inodoX := inodo{}
				bloqueCarpeta := carpeta{}
				superBlock := superbloque{}

				tamJournaling := int64(unsafe.Sizeof(journal))
				tamInodo := int64(unsafe.Sizeof(inodoX))
				tamBlock := int64(unsafe.Sizeof(bloqueCarpeta))
				tamsuperBlock := int64(unsafe.Sizeof(superBlock))

				n := (float64(tamañoParticion) - float64(tamsuperBlock)) / (float64(tamJournaling) + float64(tamInodo) + (3 * float64(tamBlock)) + 4)

				totalEstructuras := math.Floor(n)

				if tipe == false {
					tipo = "full"
				}
				if encontrado == true {

					journalingInicio := inicio + tamsuperBlock
					inodosInicio := journalingInicio + (int64(totalEstructuras) * tamJournaling) + int64(totalEstructuras) + (3 * int64(totalEstructuras))
					bloquesInicio := inodosInicio + (int64(totalEstructuras) * tamInodo)
					bitmapInoIni := journalingInicio + (int64(totalEstructuras) * tamJournaling)
					bitmapBloIni := bitmapInoIni + int64(totalEstructuras)

					superBlock.S_filesystem_type = 3
					superBlock.S_inodes_count = 2
					superBlock.S_blocks_count = 2
					superBlock.S_free_blocks_count = (3 * int64(totalEstructuras)) - 2
					superBlock.S_free_inodes_count = int64(totalEstructuras) - 2
					copy(superBlock.S_mtime[:], time.Now().String())
					copy(superBlock.S_umtime[:], time.Now().String())
					superBlock.S_mnt_count = 1
					superBlock.S_magic = 512
					superBlock.S_inode_size = tamInodo
					superBlock.S_block_size = tamBlock
					superBlock.S_first_ino = 2
					superBlock.S_first_blo = 2
					superBlock.S_bm_inode_start = bitmapInoIni
					superBlock.S_bm_block_start = bitmapBloIni
					superBlock.S_inode_start = inodosInicio
					superBlock.S_block_start = bloquesInicio

					superBlock = updateSuperBlock(superBlock, "fullFormat")

					inodoHome := crearInodo(1, 1, 0, "0", 000)
					inodoHome.I_block[0] = 0
					inodoArchivos := crearInodo(1, 1, 0, "1", 000)
					inodoArchivos.I_block[0] = 1

					var arr []byte
					arr = append(arr, 1)
					arr = append(arr, 1)
					arr = llenarArreglo(2, int(totalEstructuras), arr)

					escribirSuperBloque(rutaMbr, superBlock, inicio)
					escribirJournaling(rutaMbr, journal, journalingInicio)
					escribirBitmap(rutaMbr, arr, bitmapInoIni)

					var arrTemp []byte
					arrTemp = append(arr, 1)
					arrTemp = append(arr, 1)
					arrTemp = llenarArreglo(2, (3 * int(totalEstructuras)), arrTemp)

					escribirBitmap(rutaMbr, arrTemp, bitmapBloIni)
					escribirInodo(rutaMbr, inodoHome, inodosInicio)
					escribirInodo(rutaMbr, inodoArchivos, inodosInicio+tamInodo)

					bloqueHome := carpeta{}

					contenido := content{}
					copy(contenido.B_name[:], ".")
					contenido.B_inodo = 0
					bloqueHome.B_content[0] = contenido

					contenido = content{}
					copy(contenido.B_name[:], "..")
					contenido.B_inodo = 0
					bloqueHome.B_content[1] = contenido

					contenido = content{}
					copy(contenido.B_name[:], "usuarios.txt")
					contenido.B_inodo = 1
					bloqueHome.B_content[2] = contenido

					contenido = content{}
					contenido.B_inodo = -1
					bloqueHome.B_content[3] = contenido

					escribirBloqueCarpeta(rutaMbr, bloqueHome, bloquesInicio)

					bloqueArch := archivo{}

					copy(bloqueArch.B_content[:], "1, G, root\n1, U, root, root, 123\n")

					escribirBloqueArc(rutaMbr, bloqueArch, bloquesInicio+tamBlock)

					if tipo == "full" {
						fmt.Println("Formateo Full-", rutaMbr, "-", id)
					} else {
						fmt.Println("Formateo Fast-", rutaMbr, "-", id)
					}
				} else {
					colorize(ColorRed, "Error La Particion No Esta Montada")
				}
			} else {
				fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
			}

		} else if data == "login" {
			if len(session) == 0 {
				id := false
				user := false
				pwd := false
				other := false
				var identificador string
				var usuario string
				var pass string

				for i := 1; i < len(commandArray); i++ {
					command := strings.ToLower(commandArray[i])
					caracteres := strings.Split(command, "")
					if caracteres[0] == "u" && caracteres[1] == "s" && caracteres[2] == "r" && caracteres[3] == "-" && caracteres[4] == ">" {
						user = true
						parametros := strings.Split(command, "->")
						usuario = parametros[1]
					} else if caracteres[0] == "i" && caracteres[1] == "d" && caracteres[2] == "-" && caracteres[3] == ">" {
						id = true
						parametros := strings.Split(command, "->")
						identificador = parametros[1]
					} else if caracteres[0] == "p" && caracteres[1] == "w" && caracteres[2] == "d" && caracteres[3] == "-" && caracteres[4] == ">" {
						pwd = true
						parametros := strings.Split(command, "->")
						pass = parametros[1]
					} else {
						other = true
					}
				}
				if other == false && id == true && pwd == true && user == true {
					encontrado := false
					var rutaMbr string
					var ide string
					var particion string
					for i := 0; i < len(discosMounted); i++ {
						for k := 0; k < len(discosMounted[i].partitionsMounted); k++ {
							if discosMounted[i].partitionsMounted[k].id == identificador {
								encontrado = true
								rutaMbr = discosMounted[i].path
								ide = discosMounted[i].partitionsMounted[k].id
								particion = string(discosMounted[i].partitionsMounted[k].particion.Part_name[:clen(discosMounted[i].partitionsMounted[k].particion.Part_name[:])])
								break
							}
						}
					}
					mbrTemp := leerMBR(rutaMbr)
					var inicio int64
					particionLogin := partition{}

					if string(mbrTemp.Mbrpartition_1.Part_name[:clen(mbrTemp.Mbrpartition_1.Part_name[:])]) == particion {
						inicio = mbrTemp.Mbrpartition_1.Part_start
						particionLogin = mbrTemp.Mbrpartition_1
					} else if string(mbrTemp.Mbrpartition_2.Part_name[:clen(mbrTemp.Mbrpartition_2.Part_name[:])]) == particion {
						inicio = mbrTemp.Mbrpartition_2.Part_start
						particionLogin = mbrTemp.Mbrpartition_2
					} else if string(mbrTemp.Mbrpartition_3.Part_name[:clen(mbrTemp.Mbrpartition_3.Part_name[:])]) == particion {
						inicio = mbrTemp.Mbrpartition_3.Part_start
						particionLogin = mbrTemp.Mbrpartition_3
					} else if string(mbrTemp.Mbrpartition_4.Part_name[:clen(mbrTemp.Mbrpartition_4.Part_name[:])]) == particion {
						inicio = mbrTemp.Mbrpartition_4.Part_start
						particionLogin = mbrTemp.Mbrpartition_4
					}

					if encontrado == true {
						sb := leerSB(rutaMbr, inicio)
						usuarios := leerArchivo("/usuarios.txt", rutaMbr, sb.S_inode_start, sb)
						login(usuarios, particionLogin, ide, usuario, pass)
					}
				} else {
					fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
				}
			} else {
				colorize(ColorRed, "Error Hay Una Session Iniciada")
			}

		} else if data == "logout" {
			if len(session) == 0 {
				colorize(ColorRed, "Error No Hay Una Session Iniciada")
			} else {
				var sessionLogout []string
				session = sessionLogout
				colorize(ColorYellow, "Session Terminada")
			}
		} else if data == "mkgrp" {
			if len(session) == 0 {
				colorize(ColorRed, "Error No Hay Una Session Iniciada")
			} else {
				if session[3] == "root" {
					name := false
					other := false
					var nombre string

					for i := 1; i < len(commandArray); i++ {
						command := strings.ToLower(commandArray[i])
						caracteres := strings.Split(command, "")
						if caracteres[0] == "n" && caracteres[1] == "a" && caracteres[2] == "m" && caracteres[3] == "e" && caracteres[4] == "-" && caracteres[5] == ">" {
							name = true
							parametros := strings.Split(command, "->")
							nombre = parametros[1]
						} else {
							other = true
						}
					}
					if other == false && name == true {
						crearGrupo(nombre)
					} else {
						fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
					}
				} else {
					colorize(ColorRed, "Error Solo El Usuario Root Puede Ejecutar Este Comando")
				}
			}
		} else if data == "mkusr" {
			if len(session) == 0 {
				colorize(ColorRed, "Error No Hay Una Session Iniciada")
			} else {
				if session[3] == "root" {
					id := false
					user := false
					pwd := false
					other := false
					var identificador string
					var usuario string
					var pass string

					for i := 1; i < len(commandArray); i++ {
						command := strings.ToLower(commandArray[i])
						caracteres := strings.Split(command, "")
						if caracteres[0] == "u" && caracteres[1] == "s" && caracteres[2] == "r" && caracteres[3] == "-" && caracteres[4] == ">" {
							user = true
							parametros := strings.Split(command, "->")
							usuario = parametros[1]
						} else if caracteres[0] == "i" && caracteres[1] == "d" && caracteres[2] == "-" && caracteres[3] == ">" {
							id = true
							parametros := strings.Split(command, "->")
							identificador = parametros[1]
						} else if caracteres[0] == "p" && caracteres[1] == "w" && caracteres[2] == "d" && caracteres[3] == "-" && caracteres[4] == ">" {
							pwd = true
							parametros := strings.Split(command, "->")
							pass = parametros[1]
						} else {
							other = true
						}
					}
					if other == false && id == true && pwd == true && user == true {
						crearUsuario(usuario, pass, identificador)
					} else {
						fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
					}
				} else {
					colorize(ColorRed, "Error Solo El Usuario Root Puede Ejecutar Este Comando")
				}
			}
		} else if data == "chmod" {
			path := false
			ugo := false
			r := false
			other := false
			var ruta string
			var permisos string

			for i := 1; i < len(commandArray); i++ {
				command := strings.ToLower(commandArray[i])
				caracteres := strings.Split(command, "")
				if len(caracteres) > 1 {
					if caracteres[0] == "p" && caracteres[1] == "a" && caracteres[2] == "t" && caracteres[3] == "h" && caracteres[4] == "-" && caracteres[5] == ">" {
						path = true
						parametros := strings.Split(command, "->")
						ruta = parametros[1]
					} else if caracteres[0] == "u" && caracteres[1] == "g" && caracteres[2] == "o" && caracteres[3] == "-" && caracteres[4] == ">" {
						ugo = true
						parametros := strings.Split(command, "->")
						permisos = parametros[1]
					} else {
						other = true
					}
				} else {
					if caracteres[0] == "r" {
						r = true
					} else {
						other = true
					}
				}

			}
			if other == false && path == true && ugo == true {

				if r == false {
					r = false
				}

				cambiarPermisos(ruta, permisos, r)
			} else {
				fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
			}

		} else if data == "mkfile" {
			path := false
			size := false
			p := false
			other := false
			var ruta string
			var tam string
			var create string

			for i := 1; i < len(commandArray); i++ {
				command := strings.ToLower(commandArray[i])
				caracteres := strings.Split(command, "")
				if len(caracteres) > 1 {
					if caracteres[0] == "p" && caracteres[1] == "a" && caracteres[2] == "t" && caracteres[3] == "h" && caracteres[4] == "-" && caracteres[5] == ">" {
						path = true
						parametros := strings.Split(command, "->")
						ruta = parametros[1]
					} else if caracteres[0] == "s" && caracteres[1] == "i" && caracteres[2] == "z" && caracteres[3] == "e" && caracteres[4] == "-" && caracteres[5] == ">" {
						size = true
						parametros := strings.Split(command, "->")
						tam = parametros[1]
					} else {
						other = true
					}
				} else {
					if caracteres[0] == "p" {
						p = true
					} else {
						other = true
					}
				}
			}
			if other == false && path == true {

				if p == false {
					p = false
				}

				if size == false {
					tam = "0"
				}
				tamaña, _ := strconv.Atoi(tam)
				crearFichero(ruta, p, int64(tamaña))

				colorize(ColorWhite, "Creando Archivo"+"-"+ruta+"-"+create+"-"+tam)
			} else {
				fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
			}

		} else if data == "cat" {
			file := false
			other := false
			var filen []string

			for i := 1; i < len(commandArray); i++ {
				command := strings.ToLower(commandArray[i])
				caracteres := strings.Split(command, "")
				if caracteres[0] == "f" && caracteres[1] == "i" && caracteres[2] == "l" && caracteres[3] == "e" {
					if strings.Contains(command, "->") {
						file = true
						parametros := strings.Split(command, "->")
						filen = append(filen, parametros[1])
					} else {
						other = true
					}
				} else {
					other = true
				}
			}
			if other == false && file == true {
				particion, encontrado, path := buscarParticionMontada(session[5])
				if encontrado == true {
					sb := leerSB(path, particion.Part_start)

					for i := 0; i < len(filen); i++ {
						colorize(ColorWhite, "Leyendo Archivo"+"-"+filen[i])
						colorize(ColorBlue, leerArchivo(filen[i], path, sb.S_inode_start, sb))
						fmt.Println(leerArchivo(filen[i], path, sb.S_inode_start, sb))
					}
				}
			} else {
				fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
			}

		} else if data == "rem" {
			path := false
			other := false
			var ruta string

			for i := 1; i < len(commandArray); i++ {
				command := strings.ToLower(commandArray[i])
				caracteres := strings.Split(command, "")
				if caracteres[0] == "p" && caracteres[1] == "a" && caracteres[2] == "t" && caracteres[3] == "h" && caracteres[4] == "-" && caracteres[5] == ">" {
					path = true
					parametros := strings.Split(command, "->")
					ruta = parametros[1]
				} else {
					other = true
				}
			}
			if other == false && path == true {

				colorize(ColorWhite, "Eliminando"+"-"+ruta)
			} else {
				fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
			}

		} else if data == "ren" {
			path := false
			name := false
			other := false
			var ruta string
			var nombre string

			for i := 1; i < len(commandArray); i++ {
				command := strings.ToLower(commandArray[i])
				caracteres := strings.Split(command, "")
				if caracteres[0] == "p" && caracteres[1] == "a" && caracteres[2] == "t" && caracteres[3] == "h" && caracteres[4] == "-" && caracteres[5] == ">" {
					path = true
					parametros := strings.Split(command, "->")
					ruta = parametros[1]
				} else if caracteres[0] == "n" && caracteres[1] == "a" && caracteres[2] == "m" && caracteres[3] == "e" && caracteres[4] == "-" && caracteres[5] == ">" {
					name = true
					parametros := strings.Split(command, "->")
					nombre = parametros[1]
				} else {
					other = true
				}
			}
			if other == false && path == true && name == true {

				colorize(ColorWhite, "Cambiando Nombre"+"-"+nombre+"-"+ruta)
			} else {
				fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
			}

		} else if data == "mkdir" {
			path := false
			p := false
			other := false
			var ruta string

			for i := 1; i < len(commandArray); i++ {
				command := strings.ToLower(commandArray[i])
				caracteres := strings.Split(command, "")
				if len(caracteres) > 1 {
					if caracteres[0] == "p" && caracteres[1] == "a" && caracteres[2] == "t" && caracteres[3] == "h" && caracteres[4] == "-" && caracteres[5] == ">" {
						path = true
						parametros := strings.Split(command, "->")
						ruta = parametros[1]
					} else {
						other = true
					}
				} else {
					if caracteres[0] == "p" {
						p = true
					} else {
						other = true
					}
				}

			}
			if other == false && path == true {

				if p == false {
					p = false
				}
				crearDirectorio(ruta, p)
				colorize(ColorBlue, "Directorio Creado")
			} else {
				fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
			}

		} else if data == "mv" {
			path := false
			dest := false
			other := false
			var ruta string
			var destino string

			for i := 1; i < len(commandArray); i++ {
				command := strings.ToLower(commandArray[i])
				caracteres := strings.Split(command, "")
				if caracteres[0] == "p" && caracteres[1] == "a" && caracteres[2] == "t" && caracteres[3] == "h" && caracteres[4] == "-" && caracteres[5] == ">" {
					path = true
					parametros := strings.Split(command, "->")
					ruta = parametros[1]
				} else if caracteres[0] == "d" && caracteres[1] == "e" && caracteres[2] == "s" && caracteres[3] == "t" && caracteres[4] == "-" && caracteres[5] == ">" {
					dest = true
					parametros := strings.Split(command, "->")
					destino = parametros[1]
				} else {
					other = true
				}
			}
			if other == false && path == true && dest == true {
				colorize(ColorWhite, "Moviendo Archivo"+"-"+ruta+"-"+destino)
			} else {
				fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
			}

		} else if data == "chown" {
			path := false
			r := false
			usr := false
			other := false
			var ruta string
			var rec string
			var user string

			for i := 1; i < len(commandArray); i++ {
				command := strings.ToLower(commandArray[i])
				caracteres := strings.Split(command, "")
				if caracteres[0] == "p" && caracteres[1] == "a" && caracteres[2] == "t" && caracteres[3] == "h" && caracteres[4] == "-" && caracteres[5] == ">" {
					path = true
					parametros := strings.Split(command, "->")
					ruta = parametros[1]
				} else if caracteres[0] == "u" && caracteres[1] == "s" && caracteres[2] == "r" && caracteres[3] == "-" && caracteres[4] == ">" {
					usr = true
					parametros := strings.Split(command, "->")
					user = parametros[1]
				} else if caracteres[0] == "r" {
					usr = true
					parametros := strings.Split(command, "->")
					user = parametros[1]
				} else {
					other = true
				}
			}
			if other == false && path == true && usr == true {
				if r == false {

				}

				colorize(ColorWhite, "Moviendo Archivo"+"-"+ruta+"-"+user+"-"+rec)
			} else {
				fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
			}

		} else if data == "chgrp" {
			path := false
			grp := false
			other := false
			var ruta string
			var group string

			for i := 1; i < len(commandArray); i++ {
				command := strings.ToLower(commandArray[i])
				caracteres := strings.Split(command, "")
				if caracteres[0] == "p" && caracteres[1] == "a" && caracteres[2] == "t" && caracteres[3] == "h" && caracteres[4] == "-" && caracteres[5] == ">" {
					path = true
					parametros := strings.Split(command, "->")
					ruta = parametros[1]
				} else if caracteres[0] == "g" && caracteres[1] == "r" && caracteres[2] == "p" && caracteres[3] == "-" && caracteres[4] == ">" {
					grp = true
					parametros := strings.Split(command, "->")
					group = parametros[1]
				} else {
					other = true
				}
			}
			if other == false && path == true && grp == true {

				colorize(ColorWhite, "Cambiando Archivo"+"-"+ruta+"-"+group)
			} else {
				fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
			}

		} else if data == "recovery" {
			id := false
			other := false
			var identificador string

			for i := 1; i < len(commandArray); i++ {
				command := strings.ToLower(commandArray[i])
				caracteres := strings.Split(command, "")
				if caracteres[0] == "i" && caracteres[1] == "d" && caracteres[2] == "-" && caracteres[3] == ">" {
					id = true
					parametros := strings.Split(command, "->")
					identificador = parametros[1]
				} else {
					other = true
				}
			}
			if other == false && id == true {

				colorize(ColorWhite, "Recuperando"+"-"+identificador)
			} else {
				fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
			}

		} else if data == "loss" {
			id := false
			other := false
			var identificador string

			for i := 1; i < len(commandArray); i++ {
				command := strings.ToLower(commandArray[i])
				caracteres := strings.Split(command, "")
				if caracteres[0] == "i" && caracteres[1] == "d" && caracteres[2] == "-" && caracteres[3] == ">" {
					id = true
					parametros := strings.Split(command, "->")
					identificador = parametros[1]
				} else {
					other = true
				}
			}
			if other == false && id == true {
				colorize(ColorWhite, "Perdiendo"+"-"+identificador)
			} else {
				fmt.Println("No Se Ha podido Desmontar El Disco Error En Los Parametros")
			}

		} else { //aqui los comandos
			colorize(ColorYellow, "Comando Incorrecto")
		}
	} else {
		colorize(ColorYellow, "Comentario De Script")
	}
}

func actualizarMbr(ruta string, mbrTemp mbr) {
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

func createPath(ruta string) {
	str := strings.Split(ruta, "/")

	path := ""
	for i := 1; i < len(str)-1; i++ {
		path = path + "/" + str[i]
	}

	fmt.Println(path)

	err := os.MkdirAll(path, 0777)
	if err != nil {
		colorize(ColorRed, "Archivo No Creado")
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
	particion := partition{}
	var ruta string
	for k := 0; k < len(discosMounted); k++ {
		for i := 0; i < len(discosMounted[k].partitionsMounted); i++ {
			if discosMounted[k].partitionsMounted[i].id == id {
				erro = false
				particion = discosMounted[k].partitionsMounted[i].particion
				ruta = discosMounted[k].path
				discosMounted[k].partitionsMounted[i] = discosMounted[k].partitionsMounted[len(discosMounted[k].partitionsMounted)-1]
				discosMounted[k].partitionsMounted = discosMounted[k].partitionsMounted[:len(discosMounted[k].partitionsMounted)-1]

				break
			}
		}
	}

	if erro == true {
		colorize(ColorRed, "Error Id No Encontrado")
	} else {
		particionB := leerParticion(ruta, particion.Part_start)
		if particionB == true {
			sb := leerSB(ruta, particion.Part_start)
			sb = updateSuperBlock(sb, "desmontarParticion")
			escribirSuperBloque(ruta, sb, particion.Part_start)
			colorize(ColorBlue, "Particion Desmontada SB Actualizado")
		}
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

func leerEBR(path string, position int64) ebr {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	mbrTemp := ebr{}

	var size int = int(unsafe.Sizeof(mbrTemp))
	file.Seek(position, 0)
	mbrTemp = obtenerEBR(file, size, mbrTemp)

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
			sizeFree = mbrTemp.Mbrpartition_1.Part_start - (sizeMbr)
			if sizeFree >= sizePart {
				position = int(sizeMbr)
			} else {
				if mbrTemp.Mbrpartition_2.Part_status == status {
					sizeFree = mbrTemp.Mbrpartition_2.Part_start - (mbrTemp.Mbrpartition_1.Part_start + mbrTemp.Mbrpartition_1.Part_size)
					if sizeFree >= sizePart {
						position = int(mbrTemp.Mbrpartition_1.Part_start + mbrTemp.Mbrpartition_1.Part_size)
					} else {
						if mbrTemp.Mbrpartition_3.Part_status == status {
							sizeFree = mbrTemp.Mbrpartition_3.Part_start - (mbrTemp.Mbrpartition_2.Part_start + mbrTemp.Mbrpartition_2.Part_size)
							if sizeFree >= sizePart {
								position = int(mbrTemp.Mbrpartition_2.Part_start + mbrTemp.Mbrpartition_2.Part_size)
							} else {
								if mbrTemp.Mbrpartition_4.Part_status == status {
									sizeFree = mbrTemp.Mbrpartition_4.Part_start - (mbrTemp.Mbrpartition_3.Part_start + mbrTemp.Mbrpartition_3.Part_size)
									if sizeFree >= sizePart {
										position = int(mbrTemp.Mbrpartition_3.Part_start + mbrTemp.Mbrpartition_3.Part_size)
									} else {
										sizeFree = mbrTemp.Mbrtamaño - (mbrTemp.Mbrpartition_4.Part_start + mbrTemp.Mbrpartition_4.Part_size)
										if sizeFree >= sizePart {
											position = int(mbrTemp.Mbrpartition_4.Part_start + mbrTemp.Mbrpartition_4.Part_size)
										}
									}
								} else {
									sizeFree = mbrTemp.Mbrtamaño - (mbrTemp.Mbrpartition_3.Part_start + mbrTemp.Mbrpartition_3.Part_size)
									if sizeFree >= sizePart {
										position = int(mbrTemp.Mbrpartition_3.Part_start + mbrTemp.Mbrpartition_3.Part_size)
									}
								}
							}
						} else {
							sizeFree = mbrTemp.Mbrtamaño - (mbrTemp.Mbrpartition_2.Part_start + mbrTemp.Mbrpartition_2.Part_size)
							if sizeFree >= sizePart {
								position = int(mbrTemp.Mbrpartition_2.Part_start + mbrTemp.Mbrpartition_2.Part_size)
							}
						}
					}
				} else {
					sizeFree = mbrTemp.Mbrtamaño - (mbrTemp.Mbrpartition_1.Part_start + mbrTemp.Mbrpartition_1.Part_size)
					if sizeFree >= sizePart {
						position = int(mbrTemp.Mbrpartition_1.Part_start + mbrTemp.Mbrpartition_1.Part_size)
					}
				}
			}
		} else {
			if sizePart < (mbrTemp.Mbrtamaño - sizeMbr) {
				position = int(sizeMbr)
			}
		}
	} else if ajuste == "B" || ajuste == "b" {
		freeSpace := int64(-1)

		if mbrTemp.Mbrpartition_1.Part_status == status {
			sizeFree = mbrTemp.Mbrpartition_1.Part_start - sizeMbr
			if sizeFree >= sizePart {
				position = int(sizeMbr)
				freeSpace = int64(int(sizeFree) - (position + int(partSize)))
			}
		}

		if mbrTemp.Mbrpartition_2.Part_status == status {
			sizeFree = mbrTemp.Mbrpartition_2.Part_start - (mbrTemp.Mbrpartition_1.Part_start + mbrTemp.Mbrpartition_1.Part_size)
			if sizeFree >= sizePart {
				if freeSpace == -1 {
					position = int(mbrTemp.Mbrpartition_1.Part_start + mbrTemp.Mbrpartition_1.Part_size)
					freeSpace = int64(int(sizeFree) - (position + int(partSize)))
				} else {
					espacioLibre := sizeFree - (mbrTemp.Mbrpartition_1.Part_start + mbrTemp.Mbrpartition_1.Part_size + sizePart)
					if espacioLibre < freeSpace {
						freeSpace = espacioLibre
						position = int(mbrTemp.Mbrpartition_1.Part_start + mbrTemp.Mbrpartition_1.Part_size)
					}
				}
			}
		} else {
			if position == -1 {
				espacioLibre := sizeFree - (mbrTemp.Mbrpartition_1.Part_start + sizePart + mbrTemp.Mbrpartition_1.Part_size)
				freeSpace = espacioLibre
				position = int(mbrTemp.Mbrpartition_1.Part_start + mbrTemp.Mbrpartition_1.Part_size)
			}
		}

		if mbrTemp.Mbrpartition_3.Part_status == status {
			sizeFree = mbrTemp.Mbrpartition_3.Part_start - (mbrTemp.Mbrpartition_2.Part_start + mbrTemp.Mbrpartition_2.Part_size)
			if sizeFree >= sizePart {
				if freeSpace == -1 {
					position = int(mbrTemp.Mbrpartition_2.Part_start + mbrTemp.Mbrpartition_2.Part_size)
					freeSpace = int64(int(sizeFree) - (position + int(partSize)))
				} else {
					espacioLibre := sizeFree - (mbrTemp.Mbrpartition_2.Part_start + sizePart + mbrTemp.Mbrpartition_2.Part_size)
					if espacioLibre < freeSpace {
						freeSpace = espacioLibre
						position = int(mbrTemp.Mbrpartition_2.Part_start + mbrTemp.Mbrpartition_2.Part_size)
					}
				}
			}
		} else {
			if position == -1 {
				espacioLibre := sizeFree - (mbrTemp.Mbrpartition_2.Part_start + sizePart + mbrTemp.Mbrpartition_2.Part_size)
				freeSpace = espacioLibre
				position = int(mbrTemp.Mbrpartition_2.Part_start + mbrTemp.Mbrpartition_2.Part_size)
			}
		}

		if mbrTemp.Mbrpartition_4.Part_status == status {
			sizeFree = mbrTemp.Mbrpartition_4.Part_start - (mbrTemp.Mbrpartition_3.Part_start + mbrTemp.Mbrpartition_3.Part_size)
			if sizeFree >= sizePart {
				if freeSpace == -1 {
					position = int(mbrTemp.Mbrpartition_3.Part_start + mbrTemp.Mbrpartition_3.Part_start)
					freeSpace = int64(int(sizeFree) - (position + int(partSize)))
				} else {
					espacioLibre := sizeFree - (mbrTemp.Mbrpartition_3.Part_start + sizePart + mbrTemp.Mbrpartition_3.Part_size)
					if espacioLibre < freeSpace {
						freeSpace = espacioLibre
						position = int(mbrTemp.Mbrpartition_3.Part_start + mbrTemp.Mbrpartition_3.Part_size)
					}
				}
			}
		} else {
			if position == -1 {
				espacioLibre := sizeFree - (mbrTemp.Mbrpartition_3.Part_start + sizePart + mbrTemp.Mbrpartition_3.Part_size)
				freeSpace = espacioLibre
				position = int(mbrTemp.Mbrpartition_3.Part_start + mbrTemp.Mbrpartition_3.Part_size)
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
			sizeFree = mbrTemp.Mbrpartition_2.Part_start - (mbrTemp.Mbrpartition_1.Part_start + mbrTemp.Mbrpartition_1.Part_size)
			if sizeFree >= sizePart {
				if freeSpace == 0 {
					position = int(mbrTemp.Mbrpartition_1.Part_start + mbrTemp.Mbrpartition_1.Part_size)
					freeSpace = int64(int(sizeFree) - (position + int(partSize)))
				} else {
					espacioLibre := sizeFree - (mbrTemp.Mbrpartition_1.Part_start + sizePart + mbrTemp.Mbrpartition_1.Part_size)
					if espacioLibre > freeSpace {
						freeSpace = espacioLibre
						position = int(mbrTemp.Mbrpartition_1.Part_start + mbrTemp.Mbrpartition_1.Part_size)
					}
				}
			}
		} else {
			if position == -1 {
				espacioLibre := sizeFree - (mbrTemp.Mbrpartition_1.Part_start + sizePart + mbrTemp.Mbrpartition_1.Part_size)
				freeSpace = espacioLibre
				position = int(mbrTemp.Mbrpartition_1.Part_start + mbrTemp.Mbrpartition_1.Part_size)
			}
		}

		if mbrTemp.Mbrpartition_3.Part_status == status {
			sizeFree = mbrTemp.Mbrpartition_3.Part_start - (mbrTemp.Mbrpartition_2.Part_start + mbrTemp.Mbrpartition_2.Part_size)
			if sizeFree >= sizePart {
				if freeSpace == 0 {
					position = int(mbrTemp.Mbrpartition_2.Part_start + mbrTemp.Mbrpartition_2.Part_size)
					freeSpace = int64(int(sizeFree) - (position + int(partSize)))
				} else {
					espacioLibre := sizeFree - (mbrTemp.Mbrpartition_2.Part_start + sizePart + mbrTemp.Mbrpartition_2.Part_size)
					if espacioLibre > freeSpace {
						freeSpace = espacioLibre
						position = int(mbrTemp.Mbrpartition_2.Part_start + mbrTemp.Mbrpartition_2.Part_size)
					}
				}
			}
		} else {
			if position == -1 {
				espacioLibre := sizeFree - (mbrTemp.Mbrpartition_2.Part_start + sizePart + mbrTemp.Mbrpartition_2.Part_size)
				freeSpace = espacioLibre
				position = int(mbrTemp.Mbrpartition_2.Part_start + mbrTemp.Mbrpartition_2.Part_size)
			}
		}

		if mbrTemp.Mbrpartition_4.Part_status == status {
			sizeFree = mbrTemp.Mbrpartition_4.Part_start - (mbrTemp.Mbrpartition_3.Part_start + mbrTemp.Mbrpartition_3.Part_size)
			if sizeFree >= sizePart {
				if freeSpace == 0 {
					position = int(mbrTemp.Mbrpartition_3.Part_start + mbrTemp.Mbrpartition_3.Part_size)
					freeSpace = int64(int(sizeFree) - (position + int(partSize)))
				} else {
					espacioLibre := sizeFree - (mbrTemp.Mbrpartition_3.Part_start + sizePart + mbrTemp.Mbrpartition_3.Part_size)
					if espacioLibre > freeSpace {
						freeSpace = espacioLibre
						position = int(mbrTemp.Mbrpartition_3.Part_start + mbrTemp.Mbrpartition_3.Part_size)
					}
				}
			}
		} else {
			if position == -1 {
				espacioLibre := sizeFree - (mbrTemp.Mbrpartition_1.Part_start + sizePart + mbrTemp.Mbrpartition_3.Part_size)
				freeSpace = espacioLibre
				position = int(mbrTemp.Mbrpartition_3.Part_start + mbrTemp.Mbrpartition_3.Part_size)
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

func obtenerEBR(file *os.File, size int, ebrTemp ebr) ebr {
	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, size)

	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable estudianteTemporal
	err := binary.Read(buffer, binary.BigEndian, &ebrTemp)
	if err != nil {
		log.Fatal("binary.Read failed ", err)
	}

	//retornamos el estudiante
	return ebrTemp
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
	err := ioutil.WriteFile(ruta, cmd, os.FileMode(mode))

	if err != nil {
		colorize(ColorRed, "Error De Creacion De Reportes")
		createPath(ruta)
		ioutil.WriteFile(ruta, cmd, os.FileMode(mode))
	}
}

func graficarDisco(mbrTemp mbr, ruta string, id string) {

	sizeMbr := int64(unsafe.Sizeof(mbrTemp))
	disksize := mbrTemp.Mbrtamaño

	var porcentMbr float64 = (float64(sizeMbr) / float64(mbrTemp.Mbrtamaño)) * 100
	str := "digraph D {\n"
	str = str + "	subgraph cluster_p {\n"
	str = str + "		label = \"" + id + "\";\n"

	part5 := partition{}
	str = graphPartition(mbrTemp.Mbrpartition_4, part5, str, disksize, 4, id)
	str = graphPartition(mbrTemp.Mbrpartition_3, mbrTemp.Mbrpartition_4, str, disksize, 3, id)
	str = graphPartition(mbrTemp.Mbrpartition_2, mbrTemp.Mbrpartition_3, str, disksize, 2, id)
	str = graphPartition(mbrTemp.Mbrpartition_1, mbrTemp.Mbrpartition_2, str, disksize, 1, id)
	status := [1]byte{65}

	if mbrTemp.Mbrpartition_1.Part_status == status {
		start := sizeMbr
		tam := mbrTemp.Mbrpartition_1.Part_start
		fs := tam - start
		porcentSF3 := (float64(fs) / float64(disksize)) * 100

		if fs != 0 {
			str = str + "		subgraph cluster_cmbr_part1{\n"
			str = str + "			label = \"" + fmt.Sprintf("%f", porcentSF3) + "%\";\n"
			str = str + "			FREE00;\n"
			str = str + "		}\n"
		}
	} else {
		start := sizeMbr
		fs := disksize - (start)
		porcentSF3 := (float64(fs) / float64(disksize)) * 100

		if fs != 0 {
			str = str + "		subgraph cluster_c00 {\n"
			str = str + "			label = \"" + fmt.Sprintf("%f", porcentSF3) + "\";\n"
			str = str + "			FREE2;\n"
			str = str + "		}\n"
		}
	}

	str = str + "		subgraph cluster_cmbr {\n"
	str = str + "			label = \"" + fmt.Sprintf("%f", porcentMbr) + "\";\n"
	str = str + "			MBR;\n"
	str = str + "		}\n"
	str = str + "	}\n"
	str = str + "}\n"

	fmt.Println(str)

	b := []byte(str)
	erro := ioutil.WriteFile("reporteDisco.dot", b, 0664)
	if erro != nil {
		log.Fatal(erro)
	}

	path, _ := exec.LookPath("dot")
	cmd, _ := exec.Command(path, "-Tpng", "reporteDisco.dot").Output()
	mode := int(0777)
	err := ioutil.WriteFile(ruta, cmd, os.FileMode(mode))

	if err != nil {
		colorize(ColorRed, "Error De Creacion De Reportes")
		createPath(ruta)
		ioutil.WriteFile(ruta, cmd, os.FileMode(mode))
	}
}

func graphPartition(Mbrpartition_1 partition, Mbrpartition_2 partition, str string, disksize int64, cluster int, id string) string {

	status := [1]byte{65}
	var porcentP4 float64
	var porcentSF3 float64
	var fs int64
	if Mbrpartition_1.Part_status == status {
		if Mbrpartition_2.Part_status == status {
			start2 := Mbrpartition_1.Part_start
			start3 := Mbrpartition_2.Part_start
			tam2 := Mbrpartition_1.Part_size
			//tam3 := mbrTemp.Mbrpartition_2.Part_size
			fs = start3 - (start2 + tam2)
			porcentSF3 = (float64(fs) / float64(disksize)) * 100

			porcentP4 = (float64(tam2) / float64(disksize)) * 100
		} else {
			start := Mbrpartition_1.Part_start
			tam := Mbrpartition_1.Part_size
			fs = disksize - (start + tam)
			porcentSF3 = (float64(fs) / float64(disksize)) * 100

			porcentP4 = (float64(tam) / float64(disksize)) * 100
		}

		tipo := string(Mbrpartition_1.Part_type[:clen(Mbrpartition_1.Part_type[:])])

		if tipo == "e" || tipo == "E" {
			tipo = "Extendida"
		} else if tipo == "p" || tipo == "P" {
			tipo = "Primaria"
		} else if tipo == "l" || tipo == "L" {
			tipo = "Logica"
		}

		if fs != 0 {
			str = str + "		subgraph cluster_c" + strconv.Itoa(cluster) + "{\n"
			str = str + "			label = \"" + fmt.Sprintf("%f", porcentSF3) + "%\";\n"
			str = str + "			FREE" + strconv.Itoa(cluster) + ";\n"
			str = str + "		}\n"
		}

		label := Mbrpartition_1.Part_name[:]
		str = str + "		subgraph cluster_c" + strconv.Itoa(cluster) + strconv.Itoa(cluster) + " {\n"
		str = str + "			label = \"" + tipo + "-" + fmt.Sprintf("%f", porcentP4) + "%\";\n"
		str = str + "			" + string(label[:clen(label)]) + ";\n"
		if tipo == "Extendida" {
			var path string

			for i := 0; i < len(discosMounted); i++ {
				if discosMounted[i].id == id {
					path = discosMounted[i].path
					break
				}
			}
			fmt.Println(path)
			ebrTemp := leerEBR(path, Mbrpartition_1.Part_start)
			sizeEbr := int64(unsafe.Sizeof(ebrTemp))
			porcentP4 = (float64(sizeEbr) / float64(disksize)) * 100
			str = str + "		subgraph cluster_c" + strconv.Itoa(cluster+1) + strconv.Itoa(cluster+1) + " {\n"
			str = str + "			label = \"" + fmt.Sprintf("%f", porcentP4) + "%\";\n"
			str = str + "			EBR_" + string(ebrTemp.Part_name[:clen(ebrTemp.Part_name[:])]) + ";\n"
			str = str + "		}\n"
		}
		str = str + "		}\n"
	}
	return str
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

/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/
func reporteDeBipmapDeInodos(sb superbloque, pathG string, rutaDisco string) {
	inodos := leerBitMapInodos(sb, rutaDisco)

	str := ""

	contador := 0
	for i := 0; i < len(inodos); i++ {
		if inodos[i] == 0 {
			str = str + "0	"
		} else {
			str = str + "1	"
		}
		contador++
		if contador == 10 {
			str = str + "\n"
			contador = 0
		}
	}

	b := []byte(str)
	err := ioutil.WriteFile(pathG, b, 0644)
	if err != nil {
		err = nil
		createPath(pathG)
		err = ioutil.WriteFile(pathG, b, 0644)
		if err != nil {
			colorize(ColorRed, "Error En La Creacion Del Archivo")
		}
	}

	fmt.Println(str)
}

func reporteBloques(sb superbloque, pathG string, rutaDisco string) {

	str := "digraph H {\n"
	str = str + "	graph [\n"
	str = str + "		rankdir = \"LR\"\n"
	str = str + "	];\n"

	var bit [1]byte
	copy(bit[:], "0")

	inodos := leerBitMapInodos(sb, rutaDisco)

	for i := 0; i < len(inodos); i++ {
		if inodos[i] == 1 {
			inode := inodo{}
			pos := sb.S_inode_start + (int64(unsafe.Sizeof(inode)) * int64(i))
			inode = leerInodo(rutaDisco, pos)

			if inode.I_type == bit {

				for j := 0; j < len(inode.I_block); j++ {
					if inode.I_block[j] != -1 {
						bCarpeta := carpeta{}
						pos := sb.S_block_start + (int64(unsafe.Sizeof(bCarpeta)) * inode.I_block[j])
						bCarpeta = leerBloqueCarpeta(rutaDisco, pos)

						str = str + "bloque_" + strconv.Itoa(int(inode.I_block[j])) + " [\n"
						str = str + "	 shape=plaintext\n"
						str = str + "	 label=<\n"
						str = str + "	   <table border='1' cellborder='1'>\n"
						str = str + "		 <tr><td colspan=\"2\"> Bloque Carpeta " + strconv.Itoa(i) + "</td></tr>\n"

						for k := 0; k < len(bCarpeta.B_content); k++ {
							str = str + "		 <tr><td>" + string(bCarpeta.B_content[k].B_name[:clen(bCarpeta.B_content[k].B_name[:])]) + "</td><td port='port_" + strconv.Itoa(k+0) + "'>" + strconv.Itoa(int(bCarpeta.B_content[k].B_inodo)) + "</td></tr>\n"
						}

						str = str + "	   </table>\n"
						str = str + "	>];\n"
					}
				}

			} else {
				for j := 0; j < len(inode.I_block); j++ {
					if inode.I_block[j] != -1 {

						bArchivo := archivo{}
						pos := sb.S_block_start + (int64(unsafe.Sizeof(bArchivo)) * inode.I_block[j])
						bArchivo = leerBloqueArchivo(rutaDisco, pos)

						str = str + "bloque_" + strconv.Itoa(i) + " [\n"
						str = str + "	 shape=plaintext\n"
						str = str + "	 label=<\n"
						str = str + "	   <table border='1' cellborder='1'>\n"
						str = str + "		 <tr><td colspan=\"1\"> Bloque Archivo " + strconv.Itoa(i) + "</td></tr>\n"
						str = str + "		 <tr><td port='port_0'>" + string(bArchivo.B_content[:clen(bArchivo.B_content[:])]) + "</td></tr>\n"
						str = str + "	   </table>\n"
						str = str + "	>];\n"
					}
				}
			}
		}
	}

	bloques := leerBitMapBloques(sb, rutaDisco)

	anterior := -1
	for i := 0; i < len(bloques); i++ {
		if bloques[i] == 1 {

			if anterior != -1 {
				str = str + "	bloque_" + strconv.Itoa(anterior) + ":port_0 -> bloque_" + strconv.Itoa(i) + ";\n"
			}
			anterior = i

		}
	}
	str = str + "  }"

	fmt.Println(str)

	b := []byte(str)
	erro := ioutil.WriteFile("reporteBloques.dot", b, 0664)
	if erro != nil {
		log.Fatal(erro)
	}

	path, _ := exec.LookPath("dot")
	cmd, _ := exec.Command(path, "-Tpdf", "reporteBloques.dot").Output()
	mode := int(0777)
	err := ioutil.WriteFile(pathG, cmd, os.FileMode(mode))

	if err != nil {
		colorize(ColorRed, "Error De Creacion De Reportes")
		createPath(pathG)
		ioutil.WriteFile(pathG, cmd, os.FileMode(mode))
	}
}

func reporteInodos(sb superbloque, pathG string, rutaDisco string) {

	str := "digraph H {\n"
	str = str + "	graph [\n"
	str = str + "		rankdir = \"LR\"\n"
	str = str + "	];\n"

	inodos := leerBitMapInodos(sb, rutaDisco)

	anterior := -1
	for i := 0; i < len(inodos); i++ {
		if inodos[i] == 1 {
			inode := inodo{}
			pos := sb.S_inode_start + (int64(unsafe.Sizeof(inode)) * int64(i))
			inode = leerInodo(rutaDisco, pos)

			str = str + "inodo_" + strconv.Itoa(i) + " [\n"
			str = str + "	 shape=plaintext\n"
			str = str + "	 label=<\n"
			str = str + "	   <table border='1' cellborder='1'>\n"
			str = str + "		 <tr><td colspan=\"2\"> inodo" + strconv.Itoa(i) + "</td></tr>\n"
			str = str + "		 <tr><td>I_uid</td><td port='port_0'>" + strconv.Itoa(int(inode.I_uid)) + "</td></tr>\n"
			str = str + "		 <tr><td>I_gid</td><td port='port_1'>" + strconv.Itoa(int(inode.I_gid)) + "</td></tr>\n"
			str = str + "		 <tr><td>I_size</td><td port='port_2'>" + strconv.Itoa(int(inode.I_size)) + "</td></tr>\n"
			str = str + "		 <tr><td>I_atime</td><td port='port_3'>" + string(inode.I_atime[:clen(inode.I_atime[:])]) + "</td></tr>\n"
			str = str + "		 <tr><td>I_ctime</td><td port='port_4'>" + string(inode.I_ctime[:clen(inode.I_ctime[:])]) + "</td></tr>\n"
			str = str + "		 <tr><td>I_mtime</td><td port='port_5'>" + string(inode.I_mtime[:clen(inode.I_mtime[:])]) + "</td></tr>\n"

			for k := 0; k < len(inode.I_block); k++ {
				str = str + "		 <tr><td>bloque" + strconv.Itoa(k) + "</td><td port='port_" + strconv.Itoa(k+6) + "'>" + strconv.Itoa(int(inode.I_block[k])) + "</td></tr>\n"
			}

			str = str + "		 <tr><td>I_type</td><td port='port_22'>" + string(inode.I_type[:clen(inode.I_type[:])]) + "</td></tr>\n"
			str = str + "		 <tr><td>I_perm</td><td port='port_23'>" + strconv.Itoa(int(inode.I_perm)) + "</td></tr>\n"
			str = str + "	   </table>\n"
			str = str + "	>];\n"

			if anterior != -1 {
				str = str + "	inodo_" + strconv.Itoa(anterior) + ":port_10 -> inodo_" + strconv.Itoa(i) + ";\n"
			}
			anterior = i
		}
	}
	str = str + "  }"

	fmt.Println(str)

	b := []byte(str)
	erro := ioutil.WriteFile("reporteInodos.dot", b, 0664)
	if erro != nil {
		log.Fatal(erro)
	}

	path, _ := exec.LookPath("dot")
	cmd, _ := exec.Command(path, "-Tpdf", "reporteInodos.dot").Output()
	mode := int(0777)
	err := ioutil.WriteFile(pathG, cmd, os.FileMode(mode))

	if err != nil {
		colorize(ColorRed, "Error De Creacion De Reportes")
		createPath(pathG)
		ioutil.WriteFile(pathG, cmd, os.FileMode(mode))
	}
}

func isPermited(nodoActual inodo) bool {
	permited := false

	perms := strconv.Itoa(int(nodoActual.I_perm))
	permsIndividuales := strings.Split(perms, "")

	if len(permsIndividuales) == 1 {
		fmt.Println("Longitud De Uno")
	} else if len(permsIndividuales) == 2 {
		fmt.Println("Longitud De Dos")
	} else if len(permsIndividuales) == 3 {
		fmt.Println("Longitud De Tres")
	}

	return permited
}

func cambiarPermisos(rutaDir string, permisos string, r bool) {
	particion, encontrado, path := buscarParticionMontada(session[5])

	var bit [1]byte
	copy(bit[:], "0")

	if encontrado == true {
		sb := leerSB(path, particion.Part_start)
		lugares := strings.Split(rutaDir, "/")
		inodoRaiz := leerInodo(path, sb.S_inode_start)

		inodoB, posicion, err := obtenerInodoBuscar(lugares, inodoRaiz, sb, path, 0)

		if err == false {

			newPerms, _ := strconv.Atoi(permisos)

			inodoB.I_perm = int64(newPerms)

			pos := sb.S_inode_start + (int64(unsafe.Sizeof(inodoB)) * posicion)
			escribirInodo(path, inodoB, pos)

			if r == true {
				for i := 0; i < len(inodoB.I_block); i++ {
					if inodoB.I_type == bit {
						if inodoB.I_block[i] != -1 {
							bloque := carpeta{}
							pos := sb.S_block_start + (int64(unsafe.Sizeof(bloque)) * inodoB.I_block[i])
							bloque = leerBloqueCarpeta(path, pos)
							for k := 0; k < len(bloque.B_content); k++ {
								if bloque.B_content[k].B_inodo != -1 {
									if string(bloque.B_content[k].B_name[:clen(bloque.B_content[k].B_name[:])]) != "." && string(bloque.B_content[k].B_name[:clen(bloque.B_content[k].B_name[:])]) != ".." {
										newRuta := rutaDir + "/" + string(bloque.B_content[k].B_name[:clen(bloque.B_content[k].B_name[:])])
										cambiarPermisos(newRuta, permisos, r)
									}
								}
							}
						}
					}
				}
			}

			colorize(ColorWhite, "Dando Permisos"+"-"+rutaDir+"-"+permisos+"-")
		} else {
			colorize(ColorRed, "Error La Direccion de Archivo/Carpeta No Existe")
		}
	}
}

func obtenerInodoBuscar(ruta []string, inodoR inodo, sb superbloque, rutaDisco string, padre int64) (inodo, int64, bool) {
	bloque := carpeta{}
	inode := inodo{}

	var newRuta []string
	for i := 1; i < len(ruta); i++ {
		newRuta = append(newRuta, ruta[i])
	}

	existe := false

	inodoRetorno := inodo{}
	posicion := int64(0)
	err := true

	if len(newRuta) != 0 {
		fmt.Println("---___---___---___---___---___---___---___---")
		for i := 0; i < len(inodoR.I_block); i++ {
			fmt.Println("Bloque-", inodoR.I_block[i])
			if inodoR.I_block[i] != -1 {
				pos := sb.S_block_start + (int64(unsafe.Sizeof(bloque)) * inodoR.I_block[i])
				bloque = leerBloqueCarpeta(rutaDisco, pos)
				for k := 0; k < len(bloque.B_content); k++ {
					fmt.Println("inodo-", bloque.B_content[k].B_inodo, string(bloque.B_content[k].B_name[:clen(bloque.B_content[k].B_name[:])]))
					if bloque.B_content[k].B_inodo != -1 {
						if string(bloque.B_content[k].B_name[:clen(bloque.B_content[k].B_name[:])]) == newRuta[0] {
							fmt.Println(string(bloque.B_content[k].B_name[:clen(bloque.B_content[k].B_name[:])]))
							existe = true
							iPos := sb.S_inode_start + (int64(unsafe.Sizeof(inode)) * int64(bloque.B_content[k].B_inodo))
							inode = leerInodo(rutaDisco, iPos)
							inodoRetorno, posicion, err = obtenerInodoBuscar(newRuta, inode, sb, rutaDisco, int64(bloque.B_content[k].B_inodo))
							break
						}
					}
				}
			}
			if existe == true {
				break
			}
		}
		fmt.Println("---___---___---___---___---___---___---___---")
		if existe == false {
			inodoRetorno = inodo{}
			posicion = 0
			err = true
		}
	} else {
		inodoRetorno = inodoR
		posicion = padre
		err = false
	}
	return inodoRetorno, posicion, err
}

func crearFichero(rutaDir string, p bool, tam int64) {
	particion, encontrado, path := buscarParticionMontada(session[5])

	if encontrado == true {
		sb := leerSB(path, particion.Part_start)
		lugares := strings.Split(rutaDir, "/")
		inodoRaiz := leerInodo(path, sb.S_inode_start)
		var newRuta []string
		for i := 0; i < len(lugares)-1; i++ {
			newRuta = append(newRuta, lugares[i])
		}

		buscarArchivo(newRuta, lugares, sb, inodoRaiz, path, p, 0, tam)
	}
}

func buscarArchivo(ruta []string, rutaArchivo []string, sb superbloque, inodoR inodo, rutaDisco string, p bool, padre int64, size int64) {
	bloque := carpeta{}
	inode := inodo{}
	var bit [1]byte
	copy(bit[:], "0")

	var newRuta []string
	for i := 1; i < len(ruta); i++ {
		newRuta = append(newRuta, ruta[i])
	}

	var newRutaArchivo []string
	for i := 1; i < len(rutaArchivo); i++ {
		newRutaArchivo = append(newRutaArchivo, rutaArchivo[i])
	}

	existe := false

	if len(newRuta) != 0 {
		fmt.Println("**************|||||||||||**************")
		for i := 0; i < len(inodoR.I_block); i++ {
			fmt.Println(inodoR.I_block[i], "-Inodo")
			if inodoR.I_block[i] != -1 {
				pos := sb.S_block_start + (int64(unsafe.Sizeof(bloque)) * inodoR.I_block[i])
				if inodoR.I_type == bit {
					bloque = leerBloqueCarpeta(rutaDisco, pos)
					for k := 0; k < len(bloque.B_content); k++ {
						fmt.Println(bloque.B_content[k].B_inodo, "-", string(bloque.B_content[k].B_name[:clen(bloque.B_content[k].B_name[:])]))
						if bloque.B_content[k].B_inodo != -1 {
							if string(bloque.B_content[k].B_name[:clen(bloque.B_content[k].B_name[:])]) == newRutaArchivo[0] {
								existe = true
								iPos := sb.S_inode_start + (int64(unsafe.Sizeof(inode)) * int64(bloque.B_content[k].B_inodo))
								inode = leerInodo(rutaDisco, iPos)
								buscarArchivo(newRuta, newRutaArchivo, sb, inode, rutaDisco, p, int64(bloque.B_content[k].B_inodo), size)
								break
							}
						}
					}
					if existe == true {
						break
					}
				} else {
					break
				}
			}
		}
		if existe == false {
			if p == true {
				crearCarpeta(ruta, sb, inodoR, rutaDisco, p, padre)
				buscarArchivo(ruta, rutaArchivo, sb, inodoR, rutaDisco, p, padre, size)
			} else {
				colorize(ColorRed, "Error: No Existe El Directorio")
			}
		}
	} else {
		if len(newRutaArchivo) == 1 {
			fmt.Println("llegamos hasta aca")
			guardarArchivo(inodoR, newRutaArchivo[0], sb, rutaDisco, padre, size)
		}
	}
}

func guardarArchivo(inodoR inodo, nombreArchivo string, sb superbloque, rutaDisco string, padre int64, size int64) {
	fmt.Println("*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*")
	var bit [1]byte
	copy(bit[:], "0")
	bloque := carpeta{}

	for i := 0; i < len(inodoR.I_block); i++ {
		fmt.Println(inodoR.I_block[i], "-Inodo")
		if inodoR.I_block[i] != -1 {
			pos := sb.S_block_start + (int64(unsafe.Sizeof(bloque)) * inodoR.I_block[i])
			if inodoR.I_type == bit {
				creado := false
				bloque = leerBloqueCarpeta(rutaDisco, pos)
				for k := 0; k < len(bloque.B_content); k++ {
					if bloque.B_content[k].B_inodo == -1 {
						bloque.B_content[k].B_inodo = int32(sb.S_first_ino)
						copy(bloque.B_content[k].B_name[:], nombreArchivo)
						fmt.Println(bloque.B_content[k].B_inodo, "-", string(bloque.B_content[k].B_name[:clen(bloque.B_content[k].B_name[:])]))

						bitmap := leerBitMapInodos(sb, rutaDisco)
						bitmap[sb.S_first_ino] = 1

						sb.S_inodes_count = sb.S_inodes_count + 1
						sb.S_free_inodes_count = sb.S_free_inodes_count - 1
						sb.S_first_ino = getFirstFree(bitmap)

						idUser, _ := strconv.Atoi(session[0])
						inodoNew := crearInodo(int64(idUser), getIdGroup(rutaDisco, sb), size, "1", 664)
						particion, _, _ := buscarParticionMontada(session[5])

						posI := sb.S_inode_start + (int64(unsafe.Sizeof(inodoR)) * int64(bloque.B_content[k].B_inodo))

						var bloques []archivo

						limite := 0
						for j := int64(0); j < size; j++ {
							limite++
							if limite == 64 {
								limite = 0
								bloqueArchivo := archivo{}
								bloques = append(bloques, bloqueArchivo)
							}
						}

						if limite != 0 {
							bloqueArchivo := archivo{}
							copy(bloqueArchivo.B_content[:], "Jose preuba Guerra")
							bloques = append(bloques, bloqueArchivo)
						}

						bitmapB := leerBitMapBloques(sb, rutaDisco)

						for j := 0; j < len(bloques); j++ {
							inodoNew.I_block[i] = sb.S_first_blo
							bitmapB[sb.S_first_blo] = 1
							bPos := sb.S_block_start + (int64(unsafe.Sizeof(bloques[i])) * int64(sb.S_first_blo))
							sb.S_blocks_count = sb.S_blocks_count + 1
							sb.S_free_blocks_count = sb.S_free_blocks_count - 1
							sb.S_first_blo = getFirstFree(bitmapB)
							escribirBloqueArc(rutaDisco, bloques[j], bPos)
						}

						escribirBitmap(rutaDisco, bitmap, sb.S_bm_inode_start)
						escribirBitmap(rutaDisco, bitmapB, sb.S_bm_block_start)
						escribirSuperBloque(rutaDisco, sb, particion.Part_start)
						escribirInodo(rutaDisco, inodoNew, posI)
						escribirBloqueCarpeta(rutaDisco, bloque, pos)
						creado = true
						break
					}
				}
				if creado == true {
					break
				}
			} else {
				break
			}
		} else {
			bloqueHome := carpeta{}
			contenido := content{}
			contenido.B_inodo = -1
			bloqueHome.B_content[0] = contenido

			contenido = content{}
			contenido.B_inodo = -1
			bloqueHome.B_content[1] = contenido

			contenido = content{}
			contenido.B_inodo = -1
			bloqueHome.B_content[2] = contenido

			contenido = content{}
			contenido.B_inodo = -1
			bloqueHome.B_content[3] = contenido

			inodoR.I_block[i] = sb.S_first_blo

			bitmapB := leerBitMapBloques(sb, rutaDisco)
			bitmapB[sb.S_first_blo] = 1

			bPos := sb.S_block_start + (int64(unsafe.Sizeof(bloqueHome)) * int64(sb.S_first_blo))
			posI := sb.S_inode_start + (int64(unsafe.Sizeof(inodoR)) * int64(padre))

			sb.S_blocks_count = sb.S_blocks_count + 1
			sb.S_free_blocks_count = sb.S_free_blocks_count - 1
			sb.S_first_blo = getFirstFree(bitmapB)

			particion, _, _ := buscarParticionMontada(session[5])

			escribirBitmap(rutaDisco, bitmapB, sb.S_bm_block_start)
			escribirSuperBloque(rutaDisco, sb, particion.Part_start)
			escribirInodo(rutaDisco, inodoR, posI)
			escribirBloqueCarpeta(rutaDisco, bloqueHome, bPos)
			guardarArchivo(inodoR, nombreArchivo, sb, rutaDisco, padre, size)
			break
		}
	}
}

func crearDirectorio(rutaDir string, p bool) {
	particion, encontrado, path := buscarParticionMontada(session[5])

	if encontrado == true {
		sb := leerSB(path, particion.Part_start)
		lugares := strings.Split(rutaDir, "/")
		inodoRaiz := leerInodo(path, sb.S_inode_start)
		crearCarpeta(lugares, sb, inodoRaiz, path, p, 0)
	}
}

func crearCarpeta(ruta []string, sb superbloque, inodoR inodo, rutaDisco string, p bool, padre int64) {
	bloque := carpeta{}
	inode := inodo{}
	var bit [1]byte
	copy(bit[:], "0")

	var newRuta []string
	for i := 1; i < len(ruta); i++ {
		newRuta = append(newRuta, ruta[i])
	}

	existe := false

	if len(newRuta) != 0 {
		fmt.Println("************************************")
		for i := 0; i < len(inodoR.I_block); i++ {
			fmt.Println(inodoR.I_block[i], "-Inodo")
			if inodoR.I_block[i] != -1 {
				pos := sb.S_block_start + (int64(unsafe.Sizeof(bloque)) * inodoR.I_block[i])
				if inodoR.I_type == bit {
					bloque = leerBloqueCarpeta(rutaDisco, pos)
					for k := 0; k < len(bloque.B_content); k++ {
						fmt.Println(bloque.B_content[k].B_inodo, "-", string(bloque.B_content[k].B_name[:clen(bloque.B_content[k].B_name[:])]))
						if bloque.B_content[k].B_inodo != -1 {
							if string(bloque.B_content[k].B_name[:clen(bloque.B_content[k].B_name[:])]) == newRuta[0] {
								existe = true
								iPos := sb.S_inode_start + (int64(unsafe.Sizeof(inode)) * int64(bloque.B_content[k].B_inodo))
								inode = leerInodo(rutaDisco, iPos)
								crearCarpeta(newRuta, sb, inode, rutaDisco, p, int64(bloque.B_content[k].B_inodo))
								break
							}
						}
					}
					if existe == true {
						break
					}
				} else {
					break
				}
			}
		}
		if existe == false {
			guardarCarpeta(inodoR, sb, rutaDisco, ruta, padre, p)
		}
	}
}

func guardarCarpeta(inodoR inodo, sb superbloque, rutaDisco string, ruta []string, padre int64, p bool) {
	var newRuta []string
	for i := 1; i < len(ruta); i++ {
		newRuta = append(newRuta, ruta[i])
	}
	fmt.Println("*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*")
	var bit [1]byte
	copy(bit[:], "0")
	bloque := carpeta{}
	if len(newRuta) != 0 {
		for i := 0; i < len(inodoR.I_block); i++ {
			fmt.Println(inodoR.I_block[i], "-Inodo")
			if inodoR.I_block[i] != -1 {
				pos := sb.S_block_start + (int64(unsafe.Sizeof(bloque)) * inodoR.I_block[i])
				if inodoR.I_type == bit {
					creado := false
					bloque = leerBloqueCarpeta(rutaDisco, pos)
					for k := 0; k < len(bloque.B_content); k++ {
						if bloque.B_content[k].B_inodo == -1 {

							bloque.B_content[k].B_inodo = int32(sb.S_first_ino)
							copy(bloque.B_content[k].B_name[:], newRuta[0])
							fmt.Println(bloque.B_content[k].B_inodo, "-", string(bloque.B_content[k].B_name[:clen(bloque.B_content[k].B_name[:])]))

							bitmap := leerBitMapInodos(sb, rutaDisco)
							bitmap[sb.S_first_ino] = 1

							sb.S_inodes_count = sb.S_inodes_count + 1
							sb.S_free_inodes_count = sb.S_free_inodes_count - 1
							sb.S_first_ino = getFirstFree(bitmap)

							fmt.Println(sb.S_first_ino, "-primer inodo libre")

							idUser, _ := strconv.Atoi(session[0])
							inodoNew := crearInodo(int64(idUser), getIdGroup(rutaDisco, sb), 0, "0", 664)
							particion, _, _ := buscarParticionMontada(session[5])

							posI := sb.S_inode_start + (int64(unsafe.Sizeof(inodoR)) * int64(bloque.B_content[k].B_inodo))

							bloqueHome := carpeta{}

							contenido := content{}
							copy(contenido.B_name[:], ".")
							contenido.B_inodo = int32(bloque.B_content[k].B_inodo)
							bloqueHome.B_content[0] = contenido

							contenido = content{}
							copy(contenido.B_name[:], "..")
							contenido.B_inodo = int32(padre)
							bloqueHome.B_content[1] = contenido

							contenido = content{}
							contenido.B_inodo = -1
							bloqueHome.B_content[2] = contenido

							contenido = content{}
							contenido.B_inodo = -1
							bloqueHome.B_content[3] = contenido

							inodoNew.I_block[0] = sb.S_first_blo

							bitmapB := leerBitMapBloques(sb, rutaDisco)
							bitmapB[sb.S_first_blo] = 1

							bPos := sb.S_block_start + (int64(unsafe.Sizeof(bloqueHome)) * int64(sb.S_first_blo))
							sb.S_blocks_count = sb.S_blocks_count + 1
							sb.S_free_blocks_count = sb.S_free_blocks_count - 1
							sb.S_first_blo = getFirstFree(bitmapB)

							if p == true {
								escribirBitmap(rutaDisco, bitmap, sb.S_bm_inode_start)
								escribirBitmap(rutaDisco, bitmapB, sb.S_bm_block_start)
								escribirSuperBloque(rutaDisco, sb, particion.Part_start)
								escribirInodo(rutaDisco, inodoNew, posI)
								escribirBloqueCarpeta(rutaDisco, bloqueHome, bPos)
								escribirBloqueCarpeta(rutaDisco, bloque, pos)
								guardarCarpeta(inodoNew, sb, rutaDisco, newRuta, int64(bloque.B_content[k].B_inodo), p)
								creado = true
							} else {
								if len(newRuta) == 1 {
									escribirBitmap(rutaDisco, bitmap, sb.S_bm_inode_start)
									escribirBitmap(rutaDisco, bitmapB, sb.S_bm_block_start)
									escribirSuperBloque(rutaDisco, sb, particion.Part_start)
									escribirInodo(rutaDisco, inodoNew, posI)
									escribirBloqueCarpeta(rutaDisco, bloqueHome, bPos)
									escribirBloqueCarpeta(rutaDisco, bloque, pos)
									creado = true
								} else {
									colorize(ColorRed, "Error En La Escritura")
								}
							}
							break
						}
					}
					if creado == true {
						break
					}
				} else {
					break
				}
			} else {
				bloqueHome := carpeta{}
				contenido := content{}
				contenido.B_inodo = -1
				bloqueHome.B_content[0] = contenido

				contenido = content{}
				contenido.B_inodo = -1
				bloqueHome.B_content[1] = contenido

				contenido = content{}
				contenido.B_inodo = -1
				bloqueHome.B_content[2] = contenido

				contenido = content{}
				contenido.B_inodo = -1
				bloqueHome.B_content[3] = contenido

				inodoR.I_block[i] = sb.S_first_blo
				fmt.Println(padre)

				bitmapB := leerBitMapBloques(sb, rutaDisco)
				bitmapB[sb.S_first_blo] = 1

				bPos := sb.S_block_start + (int64(unsafe.Sizeof(bloqueHome)) * int64(sb.S_first_blo))
				posI := sb.S_inode_start + (int64(unsafe.Sizeof(inodoR)) * int64(padre))

				sb.S_blocks_count = sb.S_blocks_count + 1
				sb.S_free_blocks_count = sb.S_free_blocks_count - 1
				sb.S_first_blo = getFirstFree(bitmapB)

				particion, _, _ := buscarParticionMontada(session[5])

				if p == true {
					escribirBitmap(rutaDisco, bitmapB, sb.S_bm_block_start)
					escribirSuperBloque(rutaDisco, sb, particion.Part_start)
					escribirInodo(rutaDisco, inodoR, posI)
					escribirBloqueCarpeta(rutaDisco, bloqueHome, bPos)
					guardarCarpeta(inodoR, sb, rutaDisco, ruta, padre, p)
				} else {
					if len(ruta) == 1 {
						escribirBitmap(rutaDisco, bitmapB, sb.S_bm_block_start)
						escribirSuperBloque(rutaDisco, sb, particion.Part_start)
						escribirInodo(rutaDisco, inodoR, posI)
						escribirBloqueCarpeta(rutaDisco, bloqueHome, bPos)
						guardarCarpeta(inodoR, sb, rutaDisco, ruta, padre, p)
					} else {
						colorize(ColorRed, "Error En La Escritura")
					}
				}
				break
			}
		}
	}
}

func getIdGroup(rutaDisco string, sb superbloque) int64 {
	users := leerArchivo("/usuarios.txt", rutaDisco, sb.S_inode_start, sb)
	registros := strings.Split(users, "\n")
	id := -1
	for i := 0; i < len(registros)-1; i++ {
		data := strings.Split(registros[i], ", ")
		if data[1] == "G" {
			if data[2] == session[3] {
				id, _ = strconv.Atoi(data[0])
			}
		}
	}
	return int64(id)
}

func generarBloque() {
	bloqueHome := carpeta{}

	contenido := content{}
	copy(contenido.B_name[:], ".")
	contenido.B_inodo = 0
	bloqueHome.B_content[0] = contenido

	contenido = content{}
	copy(contenido.B_name[:], "..")
	contenido.B_inodo = 0
	bloqueHome.B_content[1] = contenido

	contenido = content{}
	copy(contenido.B_name[:], "usuarios.txt")
	contenido.B_inodo = 1
	bloqueHome.B_content[2] = contenido

	contenido = content{}
	contenido.B_inodo = -1
	bloqueHome.B_content[3] = contenido
}

func crearUsuario(name string, pass string, group string) {
	encontrado := false
	var rutaMbr string
	var inicio int64
	for i := 0; i < len(discosMounted); i++ {
		for k := 0; k < len(discosMounted[i].partitionsMounted); k++ {
			if discosMounted[i].partitionsMounted[k].id == session[5] {
				encontrado = true
				rutaMbr = discosMounted[i].path
				inicio = discosMounted[i].partitionsMounted[k].particion.Part_start
				break
			}
		}
	}
	if encontrado == true {
		sb := leerSB(rutaMbr, inicio)
		archivo := leerArchivo("/usuarios.txt", rutaMbr, sb.S_inode_start, sb)
		registros := strings.Split(archivo, "\n")
		registrado := false
		var id string

		grupo := false
		idG := ""
		for i := 0; i < len(registros)-1; i++ {
			data := strings.Split(registros[i], ", ")
			if data[1] == "U" {
				if data[3] == name {
					registrado = true
				}
				id = data[0]
			} else if data[1] == "G" {
				if data[0] == group {
					grupo = true
					idG = data[2]
				}
			}
		}

		if registrado == false && grupo == true {
			convert, _ := strconv.Atoi(id)
			newId := convert + 1
			archivo = archivo + strconv.Itoa(newId) + ", U" + ", " + idG + ", " + name + ", " + pass + "\n"
			escribirEnUsuariosTxt(archivo, sb, rutaMbr)
		} else {
			colorize(ColorRed, "Error El Usuario Esta Registrado O el Grupo No Existe")
		}
	}
}

func crearGrupo(name string) {
	encontrado := false
	var rutaMbr string
	var inicio int64
	for i := 0; i < len(discosMounted); i++ {
		for k := 0; k < len(discosMounted[i].partitionsMounted); k++ {
			if discosMounted[i].partitionsMounted[k].id == session[5] {
				encontrado = true
				rutaMbr = discosMounted[i].path
				inicio = discosMounted[i].partitionsMounted[k].particion.Part_start
				break
			}
		}
	}
	if encontrado == true {
		sb := leerSB(rutaMbr, inicio)
		archivo := leerArchivo("/usuarios.txt", rutaMbr, sb.S_inode_start, sb)
		registros := strings.Split(archivo, "\n")
		registrado := false
		var id string
		for i := 0; i < len(registros)-1; i++ {
			data := strings.Split(registros[i], ", ")
			if data[1] == "G" {
				if data[2] == name {
					registrado = true
				}
				id = data[0]
			}
		}
		if registrado == false {
			convert, _ := strconv.Atoi(id)
			newId := convert + 1
			archivo = archivo + strconv.Itoa(newId) + ", G" + ", " + name + "\n"
			escribirEnUsuariosTxt(archivo, sb, rutaMbr)
		}
	}
}

func escribirEnUsuariosTxt(archivos string, sb superbloque, ruta string) {
	inodoI := inodo{}
	bloque := archivo{}
	posI := sb.S_inode_start + (1 * int64(unsafe.Sizeof(inodoI)))
	inodoUsuarios := leerInodo(ruta, posI)
	arreglos := separarCaracteresBA(archivos)

	for i := 0; i < len(arreglos); i++ {
		if inodoUsuarios.I_block[i] != -1 {
			pos := sb.S_block_start + (inodoUsuarios.I_block[i] * int64(unsafe.Sizeof(bloque)))
			bloque := archivo{}
			bloque.B_content = arreglos[i]
			escribirBloqueArc(ruta, bloque, pos)
		} else {
			if i != 14 && i != 13 {
				bloqueArch := archivo{}
				bloqueArch.B_content = arreglos[i]
				inodoUsuarios.I_block[i] = sb.S_first_blo
				fmt.Println(inodoUsuarios.I_block[i])
				pos := sb.S_block_start + (inodoUsuarios.I_block[i] * int64(unsafe.Sizeof(bloque)))

				bitmap := leerBitMapBloques(sb, ruta)
				bitmap[sb.S_first_blo] = 1

				sb.S_blocks_count = sb.S_blocks_count + 1
				sb.S_free_blocks_count = sb.S_free_blocks_count - 1
				sb.S_first_blo = getFirstFree(bitmap)
				particion, encontrado, _ := buscarParticionMontada(session[5])

				copy(inodoUsuarios.I_mtime[:], time.Now().String())
				if encontrado == true {
					escribirBloqueArc(ruta, bloqueArch, pos)
					escribirBitmap(ruta, bitmap, sb.S_bm_block_start)
					escribirSuperBloque(ruta, sb, particion.Part_start)
					escribirInodo(ruta, inodoUsuarios, posI)
				}
			}
		}
	}
}

func buscarParticionMontada(identificador string) (partition, bool, string) {
	var particion partition
	encontrado := false
	var path string
	for i := 0; i < len(discosMounted); i++ {
		for k := 0; k < len(discosMounted[i].partitionsMounted); k++ {
			if discosMounted[i].partitionsMounted[k].id == identificador {
				particion = discosMounted[i].partitionsMounted[k].particion
				encontrado = true
				path = discosMounted[i].path
				break
			}
		}
	}
	return particion, encontrado, path
}

func getFirstFree(bitmap []byte) int64 {
	pos := -1
	for i := 0; i < len(bitmap); i++ {
		if bitmap[i] == 0 {
			pos = i
			break
		}
	}
	return int64(pos)
}

func leerBitMapBloques(sb superbloque, path string) []byte {
	tamaño := (sb.S_blocks_count + sb.S_free_blocks_count)

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(sb.S_bm_block_start, 0)

	bitmap := leerBytes(file, int(tamaño))
	return bitmap
}

func leerBitMapInodos(sb superbloque, path string) []byte {
	tamaño := (sb.S_inodes_count + sb.S_free_inodes_count)

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(sb.S_bm_inode_start, 0)

	bitmap := leerBytes(file, int(tamaño))
	return bitmap
}

func separarCaracteresBA(archivo string) [][64]byte {
	arregloCompleto := []byte(archivo)

	var arreglos [][64]byte
	var array [64]byte

	count := 0
	for i := 0; i < len(arregloCompleto); i++ {
		array[count] = arregloCompleto[i]
		count++
		if count == 64 {
			arreglos = append(arreglos, array)
			var arrayTemp [64]byte
			array = arrayTemp
			count = 0
		}
	}
	if len(array) != 0 {
		arreglos = append(arreglos, array)
	}

	return arreglos
}

func login(users string, particion partition, id string, usuario string, pass string) {
	registros := strings.Split(users, "\n")
	var uyg []string
	for i := 0; i < len(registros)-1; i++ {
		data := strings.Split(registros[i], ", ")
		if data[1] == "U" {
			if data[3] == usuario && data[4] == pass {
				uyg = data
			}
		}
	}

	if len(uyg) == 0 {
		colorize(ColorRed, "Error En Las Credenciales")
	} else {
		uyg = append(uyg, id)
		session = uyg
		colorizefn(ColorWhite, "Sesion Iniciada Como ")
		colorize(ColorBlue, usuario)
		fmt.Println(users)
	}
}

func leerArchivo(ruta string, rutaDisco string, pos int64, sb superbloque) string {
	var archivo string
	lugares := strings.Split(ruta, "/")

	inodoRaiz := leerInodo(rutaDisco, pos)
	archivo = leerInodos(lugares, rutaDisco, pos, sb, inodoRaiz)
	return archivo
}

func leerInodos(ruta []string, rutaDisco string, pos int64, sb superbloque, inodoR inodo) string {
	var archivo string
	bloque := carpeta{}
	var bit [1]byte
	copy(bit[:], "0")

	var newRuta []string

	for i := 1; i < len(ruta); i++ {
		newRuta = append(newRuta, ruta[i])
	}

	for i := 0; i < len(inodoR.I_block); i++ {
		if inodoR.I_block[i] != -1 {
			pos = sb.S_block_start + (int64(unsafe.Sizeof(bloque)) * inodoR.I_block[i])
			if inodoR.I_type == bit {
				bloqueCarpeta := leerBloqueCarpeta(rutaDisco, pos)
				archivo = archivo + leerBloqueCarpetas(newRuta, rutaDisco, pos, sb, bloqueCarpeta)
			} else {
				bloqueArchivo := leerBloqueArchivo(rutaDisco, pos)
				archivo = archivo + leerBloqueArchivos(newRuta, rutaDisco, pos, sb, bloqueArchivo)
			}
		}
	}
	return archivo
}

func leerBloqueCarpetas(ruta []string, rutaDisco string, pos int64, sb superbloque, bloque carpeta) string {

	var archivo string
	inodoI := inodo{}

	var newRuta []string

	for i := 1; i < len(ruta); i++ {
		newRuta = append(newRuta, ruta[i])
	}

	for i := 0; i < len(bloque.B_content); i++ {
		if bloque.B_content[i].B_inodo != -1 {
			if string(bloque.B_content[i].B_name[:clen(bloque.B_content[i].B_name[:])]) == ruta[0] {
				pos = sb.S_inode_start + (int64(unsafe.Sizeof(inodoI)) * int64(bloque.B_content[i].B_inodo))
				inodoI = leerInodo(rutaDisco, pos)
				archivo = archivo + leerInodos(newRuta, rutaDisco, pos, sb, inodoI)
			}
		}
	}

	return archivo
}

func leerBloqueArchivos(ruta []string, rutaDisco string, pos int64, sb superbloque, bloque archivo) string {
	archivo := string(bloque.B_content[:clen(bloque.B_content[:])])
	return archivo
}

func escribirSuperBloque(ruta string, list superbloque, pos int64) {
	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(pos, 0)

	var bufferEstudiante bytes.Buffer
	binary.Write(&bufferEstudiante, binary.BigEndian, &list)
	escribirBytes(file, bufferEstudiante.Bytes())
	defer file.Close()
}

func escribirJournaling(ruta string, list journaling, pos int64) {
	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(pos, 0)

	var bufferEstudiante bytes.Buffer
	binary.Write(&bufferEstudiante, binary.BigEndian, &list)
	escribirBytes(file, bufferEstudiante.Bytes())
	defer file.Close()
}

func escribirBitmap(ruta string, list []byte, pos int64) {
	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(pos, 0)

	escribirBytes(file, list)
	defer file.Close()
}

func escribirInodo(ruta string, list inodo, pos int64) {
	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(pos, 0)

	var bufferEstudiante bytes.Buffer
	binary.Write(&bufferEstudiante, binary.BigEndian, &list)
	escribirBytes(file, bufferEstudiante.Bytes())
	defer file.Close()
}

func escribirBloqueCarpeta(ruta string, list carpeta, pos int64) {
	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(pos, 0)

	var bufferEstudiante bytes.Buffer
	binary.Write(&bufferEstudiante, binary.BigEndian, &list)
	escribirBytes(file, bufferEstudiante.Bytes())
	defer file.Close()
}

func escribirBloqueArc(ruta string, list archivo, pos int64) {
	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	file.Seek(pos, 0)

	var bufferEstudiante bytes.Buffer
	binary.Write(&bufferEstudiante, binary.BigEndian, &list)
	escribirBytes(file, bufferEstudiante.Bytes())
	defer file.Close()
}

func updateSuperBlock(superBloqueTem superbloque, tipo string) superbloque {
	if tipo == "montarParticion" {
		superBloqueTem.S_mnt_count = superBloqueTem.S_mnt_count + 1
		copy(superBloqueTem.S_mtime[:], time.Now().String())
	} else if tipo == "desmontarParticion" {
		copy(superBloqueTem.S_mtime[:], time.Now().String())
	}
	return superBloqueTem
}

func llenarArreglo(inicio int, fin int, arreglo []byte) []byte {

	for i := inicio; i < fin; i++ {
		arreglo = append(arreglo, 0)
	}

	return arreglo
}

func crearInodo(idUser int64, idGroup int64, archSize int64, tipo string, permisos int64) inodo {
	inodoTemp := inodo{}
	inodoTemp.I_uid = idUser
	inodoTemp.I_gid = idGroup
	inodoTemp.I_size = archSize
	copy(inodoTemp.I_atime[:], time.Now().String())
	copy(inodoTemp.I_ctime[:], time.Now().String())
	copy(inodoTemp.I_mtime[:], time.Now().String())

	for i := 0; i < len(inodoTemp.I_block); i++ {
		inodoTemp.I_block[i] = -1
	}

	copy(inodoTemp.I_type[:], tipo)
	inodoTemp.I_perm = permisos
	return inodoTemp
}

func leerSB(path string, pos int64) superbloque {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	mbrTemp := superbloque{}

	var size int = int(unsafe.Sizeof(mbrTemp))
	file.Seek(pos, 0)
	mbrTemp = obtenerSB(file, size, mbrTemp)

	return mbrTemp
}

func obtenerSB(file *os.File, size int, mbrTemp superbloque) superbloque {
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

func leerInodo(path string, pos int64) inodo {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	mbrTemp := inodo{}

	var size int = int(unsafe.Sizeof(mbrTemp))
	file.Seek(pos, 0)
	mbrTemp = obtenerInodo(file, size, mbrTemp)

	return mbrTemp
}

func obtenerInodo(file *os.File, size int, mbrTemp inodo) inodo {
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

func leerBloqueCarpeta(path string, pos int64) carpeta {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	mbrTemp := carpeta{}

	var size int = int(unsafe.Sizeof(mbrTemp))
	file.Seek(pos, 0)
	mbrTemp = obtenerBloqueCarpeta(file, size, mbrTemp)

	return mbrTemp
}

func obtenerBloqueCarpeta(file *os.File, size int, mbrTemp carpeta) carpeta {
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

func leerBloqueArchivo(path string, pos int64) archivo {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	mbrTemp := archivo{}

	var size int = int(unsafe.Sizeof(mbrTemp))
	file.Seek(pos, 0)
	mbrTemp = obtenerBloqueArchivo(file, size, mbrTemp)

	return mbrTemp
}

func obtenerBloqueArchivo(file *os.File, size int, mbrTemp archivo) archivo {
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

func leerParticion(path string, pos int64) bool {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	mbrTemp := partition{}

	var size int = int(unsafe.Sizeof(mbrTemp))
	file.Seek(pos, 0)
	mbrTempbyte := obtenerParticion(file, size)
	fmt.Println(mbrTempbyte)

	datos := false
	for i := 0; i < len(mbrTempbyte); i++ {
		if mbrTempbyte[i] != 0 {
			datos = true
		}
	}

	return datos
}

func obtenerParticion(file *os.File, size int) []byte {
	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, size)

	//Convierte la data en un buffer,necesario para
	//decodificar binario

	//retornamos el estudiante
	return data
}
