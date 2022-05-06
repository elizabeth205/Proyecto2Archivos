package structs

import "container/list"

/**************************************************************
	EStructuras para guardar los datos del proyecto
***************************************************************/
type Mbr struct {
	Tamano int64
	Fecha  [16]byte
	Fit    [2]byte
	Firma  int8
	Tabla  [4]Particion
}

type Particion struct {
	Status byte
	Type   byte
	Fit    [2]byte
	Start  int64
	Size   int64
	Name   [16]byte
}

type Ebr struct {
	Status byte
	Fit    [2]byte
	Start  int64
	Size   int64
	Next   int64
	Name   [16]byte
}

type Montadas struct {
	Path   [100]byte
	ID     int
	Estado int
	Lista  [100]ParticionMontada
}

type ParticionMontada struct {
	ID            [4]byte
	Nombre        [16]byte
	EstadoFormato byte
	EstadoMount   byte
}

type Superbloque struct {
	NombreHD               [100]byte
	ArbolVirtualCount      int64
	DetalleDirectorioCount int64
	InodosCount            int64
	BloquesCount           int64
	ArbolVirtualFree       int64
	DetalleDirectorioFree  int64
	InodosFree             int64
	BloquesFree            int64
	DateCreacion           [16]byte
	DateUltimoMontaje      [16]byte
	MontajesCount          int64
	InicioBMAV             int64
	InicioAV               int64
	InicioBMDD             int64
	InicioDD               int64
	InicioBMInodos         int64
	InicioInodos           int64
	InicioBMBloques        int64
	InicioBloques          int64
	InicioLog              int64
	TamAV                  int64
	TamDD                  int64
	TamInodo               int64
	TamBloque              int64
	PrimerLibreAV          int64
	PrimerLibreDD          int64
	PrimerLibreInodo       int64
	PrimerLibreBloque      int64
	MagicNum               int64
}

type AVDpart struct {
	AVDFechaCreacion            [16]byte
	AVDNombreDirectorio         [16]byte
	AVDApArraySubdirectorios    [6]int64
	AVDApDetalleDirectorio      int64
	AVDApArbolVirtualDirectorio int64
	AVDProper                   [16]byte
}

type efine struct {
	DDArrayFiles          [5]archivo
	DDApDetalleDirectorio int64
}

type Inodo struct {
	ICountInodo            int64
	ISizeArchivo           int64
	ICountBloquesAsignados int64
	IArrayBloques          [4]int64
	IApIndirecto           int64
	IIdProper              [16]byte
}

type Bloque struct {
	DBData [25]byte
}

type Jour struct {
	LogTipoOperacion int64
	LogTipo          int64
	LogNombre        [16]byte
	LogContenido     int64
	LogFecha         [16]byte
}

type archivo struct {
	FileNombre           [16]byte
	FileApInodo          int64
	FileDateCreacion     [16]byte
	FileDateModificacion [16]byte
}

var discos = list.New()

func Mount_disk(montado Montadas) {
	discos.PushBack(montado)
}

func Mountedisk() *list.List {
	return discos
}
