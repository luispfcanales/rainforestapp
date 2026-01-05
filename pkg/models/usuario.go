package models

import (
	"fmt"
	"strings"
	"time"
)

// Substructs for nested data

type DatosConyuge struct {
	ApellidosNombres string `json:"apellidos_nombres" firestore:"apellidos_nombres"`
	Genero           string `json:"genero" firestore:"genero"`
	FechaNacimiento  string `json:"fecha_nacimiento" firestore:"fecha_nacimiento"`
	Dni              string `json:"dni" firestore:"dni"`
	Direccion        string `json:"direccion" firestore:"direccion"`
	CopiaDni         string `json:"copia_dni,omitempty" firestore:"copia_dni,omitempty"`
}

type Hijo struct {
	ID               string `json:"id" firestore:"id"`
	ApellidosNombres string `json:"apellidos_nombres" firestore:"apellidos_nombres"`
	FechaNacimiento  string `json:"fecha_nacimiento" firestore:"fecha_nacimiento"`
	Direccion        string `json:"direccion" firestore:"direccion"`
	Dni              string `json:"dni" firestore:"dni"`
	Edad             int    `json:"edad" firestore:"edad"`
	CopiaDni         string `json:"copia_dni,omitempty" firestore:"copia_dni,omitempty"`
}

type Padre struct {
	ID               string `json:"id" firestore:"id"`
	ApellidosNombres string `json:"apellidos_nombres" firestore:"apellidos_nombres"`
	FechaNacimiento  string `json:"fecha_nacimiento" firestore:"fecha_nacimiento"`
	Ocupacion        string `json:"ocupacion" firestore:"ocupacion"`
	EstadoCivil      string `json:"estado_civil" firestore:"estado_civil"`
	Vive             bool   `json:"vive" firestore:"vive"`
}

type EducacionBasica struct {
	Nivel          string `json:"nivel" firestore:"nivel"`
	Completa       bool   `json:"completa" firestore:"completa"`
	CentroEstudios string `json:"centro_estudios" firestore:"centro_estudios"`
	Desde          string `json:"desde" firestore:"desde"`
	Hasta          string `json:"hasta" firestore:"hasta"`
}

type EducacionSuperior struct {
	Nivel          string `json:"nivel" firestore:"nivel"`
	Especialidad   string `json:"especialidad" firestore:"especialidad"`
	CentroEstudios string `json:"centro_estudios" firestore:"centro_estudios"`
	Desde          string `json:"desde" firestore:"desde"`
	Hasta          string `json:"hasta" firestore:"hasta"`
	Completa       bool   `json:"completa" firestore:"completa"`
	GradoAcademico string `json:"grado_academico" firestore:"grado_academico"`
}

type Capacitacion struct {
	ID          string `json:"id" firestore:"id"`
	Nombre      string `json:"nombre" firestore:"nombre"`
	Institucion string `json:"institucion" firestore:"institucion"`
	Horas       int    `json:"horas" firestore:"horas"`
}

type ExperienciaLaboral struct {
	ID                string `json:"id" firestore:"id"`
	Cargo             string `json:"cargo" firestore:"cargo"`
	Empresa           string `json:"empresa" firestore:"empresa"`
	FechaIngreso      string `json:"fecha_ingreso" firestore:"fecha_ingreso"`
	FechaCese         string `json:"fecha_cese" firestore:"fecha_cese"`
	TiempoPermanencia string `json:"tiempo_permanencia" firestore:"tiempo_permanencia"`
	MotivoCese        string `json:"motivo_cese" firestore:"motivo_cese"`
}

type Idioma struct {
	ID      string `json:"id" firestore:"id"`
	Idioma  string `json:"idioma" firestore:"idioma"`
	Lee     string `json:"lee" firestore:"lee"`
	Habla   string `json:"habla" firestore:"habla"`
	Escribe string `json:"escribe" firestore:"escribe"`
}

// Usuario representa un usuario en el sistema con todos sus datos detallados
type Usuario struct {
	ID        string    `json:"id,omitempty" firestore:"-"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" firestore:"updated_at,omitempty"`

	// Datos Personales
	ApellidoPaterno   string `json:"apellido_paterno" firestore:"apellido_paterno"`
	ApellidoMaterno   string `json:"apellido_materno" firestore:"apellido_materno"`
	Nombres           string `json:"nombres" firestore:"nombres"`
	Sexo              string `json:"sexo" firestore:"sexo"`
	Dni               string `json:"dni" firestore:"dni"`
	LicenciaConducir  string `json:"licencia_conducir,omitempty" firestore:"licencia_conducir,omitempty"`
	CategoriaLicencia string `json:"categoria_licencia,omitempty" firestore:"categoria_licencia,omitempty"`
	FechaNacimiento   string `json:"fecha_nacimiento" firestore:"fecha_nacimiento"`

	// Lugar de Nacimiento
	LugarNacimientoDistrito     string `json:"lugar_nacimiento_distrito" firestore:"lugar_nacimiento_distrito"`
	LugarNacimientoProvincia    string `json:"lugar_nacimiento_provincia" firestore:"lugar_nacimiento_provincia"`
	LugarNacimientoDepartamento string `json:"lugar_nacimiento_departamento" firestore:"lugar_nacimiento_departamento"`

	// Domicilio
	DireccionDomicilio string `json:"direccion_domicilio" firestore:"direccion_domicilio"`

	// Datos Laborales
	FechaIngreso         string `json:"fecha_ingreso" firestore:"fecha_ingreso"`
	LugarTrabajo         string `json:"lugar_trabajo" firestore:"lugar_trabajo"`
	PuestoActual         string `json:"puesto_actual" firestore:"puesto_actual"`
	Telefono             string `json:"telefono" firestore:"telefono"`
	Email                string `json:"email" firestore:"email"`
	SituacionContractual string `json:"situacion_contractual" firestore:"situacion_contractual"`

	// Régimen
	RegimenPensionario string `json:"regimen_pensionario" firestore:"regimen_pensionario"`
	AfpNombre          string `json:"afp_nombre,omitempty" firestore:"afp_nombre,omitempty"`
	Cuspp              string `json:"cuspp,omitempty" firestore:"cuspp,omitempty"`
	RegimenSalud       string `json:"regimen_salud" firestore:"regimen_salud"`

	// Información de Emergencia
	ContactoNombre       string `json:"contacto_nombre" firestore:"contacto_nombre"`
	ContactoParentesco   string `json:"contacto_parentesco" firestore:"contacto_parentesco"`
	ContactoCelular      string `json:"contacto_celular" firestore:"contacto_celular"`
	ContactoTelefonoFijo string `json:"contacto_telefono_fijo,omitempty" firestore:"contacto_telefono_fijo,omitempty"`
	ContactoDireccion    string `json:"contacto_direccion" firestore:"contacto_direccion"`
	GrupoSanguineo       string `json:"grupo_sanguineo" firestore:"grupo_sanguineo"`

	EstadoCivil           string `json:"estado_civil" firestore:"estado_civil"`
	ConstanciaEstadoCivil string `json:"constancia_estado_civil,omitempty" firestore:"constancia_estado_civil,omitempty"`

	// Datos Familiares
	DatosConyuge *DatosConyuge `json:"datos_conyuge,omitempty" firestore:"datos_conyuge,omitempty"`
	Hijos        []Hijo        `json:"hijos,omitempty" firestore:"hijos,omitempty"`
	Padres       []Padre       `json:"padres,omitempty" firestore:"padres,omitempty"`

	// Datos de Instrucción
	EducacionBasica   []EducacionBasica   `json:"educacion_basica,omitempty" firestore:"educacion_basica,omitempty"`
	EducacionSuperior []EducacionSuperior `json:"educacion_superior,omitempty" firestore:"educacion_superior,omitempty"`

	Capacitaciones []Capacitacion `json:"capacitaciones,omitempty" firestore:"capacitaciones,omitempty"`

	// Experiencia Laboral
	ExperienciaLaboral []ExperienciaLaboral `json:"experiencia_laboral,omitempty" firestore:"experiencia_laboral,omitempty"`

	// Idiomas
	Idiomas []Idioma `json:"idiomas,omitempty" firestore:"idiomas,omitempty"`

	// Apertura de Cuenta Sueldo
	AutorizaBcp       bool   `json:"autoriza_bcp" firestore:"autoriza_bcp"`
	AutorizaOtroBanco bool   `json:"autoriza_otro_banco" firestore:"autoriza_otro_banco"`
	OtroBancoNombre   string `json:"otro_banco_nombre,omitempty" firestore:"otro_banco_nombre,omitempty"`
	OtroBancoCuenta   string `json:"otro_banco_cuenta,omitempty" firestore:"otro_banco_cuenta,omitempty"`
	OtroBancoCci      string `json:"otro_banco_cci,omitempty" firestore:"otro_banco_cci,omitempty"`

	// Apertura de Cuenta CTS
	AutorizaCtsBcp bool `json:"autoriza_cts_bcp" firestore:"autoriza_cts_bcp"`

	Foto string `json:"foto,omitempty" firestore:"foto,omitempty"`
}

// CreateUsuarioRequest DTO para crear usuario
type CreateUsuarioRequest struct {
	ApellidoPaterno             string               `json:"apellido_paterno"`
	ApellidoMaterno             string               `json:"apellido_materno"`
	Nombres                     string               `json:"nombres"`
	Sexo                        string               `json:"sexo"`
	Dni                         string               `json:"dni"`
	LicenciaConducir            string               `json:"licencia_conducir,omitempty"`
	CategoriaLicencia           string               `json:"categoria_licencia,omitempty"`
	FechaNacimiento             string               `json:"fecha_nacimiento"`
	LugarNacimientoDistrito     string               `json:"lugar_nacimiento_distrito"`
	LugarNacimientoProvincia    string               `json:"lugar_nacimiento_provincia"`
	LugarNacimientoDepartamento string               `json:"lugar_nacimiento_departamento"`
	DireccionDomicilio          string               `json:"direccion_domicilio"`
	FechaIngreso                string               `json:"fecha_ingreso"`
	LugarTrabajo                string               `json:"lugar_trabajo"`
	PuestoActual                string               `json:"puesto_actual"`
	Telefono                    string               `json:"telefono"`
	Email                       string               `json:"email"`
	SituacionContractual        string               `json:"situacion_contractual"`
	RegimenPensionario          string               `json:"regimen_pensionario"`
	AfpNombre                   string               `json:"afp_nombre,omitempty"`
	Cuspp                       string               `json:"cuspp,omitempty"`
	RegimenSalud                string               `json:"regimen_salud"`
	ContactoNombre              string               `json:"contacto_nombre"`
	ContactoParentesco          string               `json:"contacto_parentesco"`
	ContactoCelular             string               `json:"contacto_celular"`
	ContactoTelefonoFijo        string               `json:"contacto_telefono_fijo,omitempty"`
	ContactoDireccion           string               `json:"contacto_direccion"`
	GrupoSanguineo              string               `json:"grupo_sanguineo"`
	EstadoCivil                 string               `json:"estado_civil"`
	ConstanciaEstadoCivil       string               `json:"constancia_estado_civil,omitempty"`
	DatosConyuge                *DatosConyuge        `json:"datos_conyuge,omitempty"`
	Hijos                       []Hijo               `json:"hijos,omitempty"`
	Padres                      []Padre              `json:"padres,omitempty"`
	EducacionBasica             []EducacionBasica    `json:"educacion_basica,omitempty"`
	EducacionSuperior           []EducacionSuperior  `json:"educacion_superior,omitempty"`
	Capacitaciones              []Capacitacion       `json:"capacitaciones,omitempty"`
	ExperienciaLaboral          []ExperienciaLaboral `json:"experiencia_laboral,omitempty"`
	Idiomas                     []Idioma             `json:"idiomas,omitempty"`
	AutorizaBcp                 bool                 `json:"autoriza_bcp"`
	AutorizaOtroBanco           bool                 `json:"autoriza_otro_banco"`
	OtroBancoNombre             string               `json:"otro_banco_nombre,omitempty"`
	OtroBancoCuenta             string               `json:"otro_banco_cuenta,omitempty"`
	OtroBancoCci                string               `json:"otro_banco_cci,omitempty"`
	AutorizaCtsBcp              bool                 `json:"autoriza_cts_bcp"`
	Foto                        string               `json:"foto,omitempty"`
}

// Validate valida los datos del usuario
func (u *CreateUsuarioRequest) Validate() error {
	if strings.TrimSpace(u.Nombres) == "" {
		return fmt.Errorf("los nombres son requeridos")
	}
	if strings.TrimSpace(u.ApellidoPaterno) == "" {
		return fmt.Errorf("el apellido paterno es requerido")
	}
	if strings.TrimSpace(u.Dni) == "" {
		return fmt.Errorf("el DNI es requerido")
	}
	if len(u.Nombres) < 2 {
		return fmt.Errorf("los nombres deben tener al menos 2 caracteres")
	}
	return nil
}

// ToUsuario convierte el request a un modelo Usuario
func (u *CreateUsuarioRequest) ToUsuario() *Usuario {
	return &Usuario{
		ApellidoPaterno:             strings.TrimSpace(u.ApellidoPaterno),
		ApellidoMaterno:             strings.TrimSpace(u.ApellidoMaterno),
		Nombres:                     strings.TrimSpace(u.Nombres),
		Sexo:                        strings.TrimSpace(u.Sexo),
		Dni:                         strings.TrimSpace(u.Dni),
		LicenciaConducir:            strings.TrimSpace(u.LicenciaConducir),
		CategoriaLicencia:           strings.TrimSpace(u.CategoriaLicencia),
		FechaNacimiento:             strings.TrimSpace(u.FechaNacimiento),
		LugarNacimientoDistrito:     strings.TrimSpace(u.LugarNacimientoDistrito),
		LugarNacimientoProvincia:    strings.TrimSpace(u.LugarNacimientoProvincia),
		LugarNacimientoDepartamento: strings.TrimSpace(u.LugarNacimientoDepartamento),
		DireccionDomicilio:          strings.TrimSpace(u.DireccionDomicilio),
		FechaIngreso:                strings.TrimSpace(u.FechaIngreso),
		LugarTrabajo:                strings.TrimSpace(u.LugarTrabajo),
		PuestoActual:                strings.TrimSpace(u.PuestoActual),
		Telefono:                    strings.TrimSpace(u.Telefono),
		Email:                       strings.TrimSpace(u.Email),
		SituacionContractual:        strings.TrimSpace(u.SituacionContractual),
		RegimenPensionario:          strings.TrimSpace(u.RegimenPensionario),
		AfpNombre:                   strings.TrimSpace(u.AfpNombre),
		Cuspp:                       strings.TrimSpace(u.Cuspp),
		RegimenSalud:                strings.TrimSpace(u.RegimenSalud),
		ContactoNombre:              strings.TrimSpace(u.ContactoNombre),
		ContactoParentesco:          strings.TrimSpace(u.ContactoParentesco),
		ContactoCelular:             strings.TrimSpace(u.ContactoCelular),
		ContactoTelefonoFijo:        strings.TrimSpace(u.ContactoTelefonoFijo),
		ContactoDireccion:           strings.TrimSpace(u.ContactoDireccion),
		GrupoSanguineo:              strings.TrimSpace(u.GrupoSanguineo),
		EstadoCivil:                 strings.TrimSpace(u.EstadoCivil),
		ConstanciaEstadoCivil:       strings.TrimSpace(u.ConstanciaEstadoCivil),
		DatosConyuge:                u.DatosConyuge,
		Hijos:                       u.Hijos,
		Padres:                      u.Padres,
		EducacionBasica:             u.EducacionBasica,
		EducacionSuperior:           u.EducacionSuperior,
		Capacitaciones:              u.Capacitaciones,
		ExperienciaLaboral:          u.ExperienciaLaboral,
		Idiomas:                     u.Idiomas,
		AutorizaBcp:                 u.AutorizaBcp,
		AutorizaOtroBanco:           u.AutorizaOtroBanco,
		OtroBancoNombre:             strings.TrimSpace(u.OtroBancoNombre),
		OtroBancoCuenta:             strings.TrimSpace(u.OtroBancoCuenta),
		OtroBancoCci:                strings.TrimSpace(u.OtroBancoCci),
		AutorizaCtsBcp:              u.AutorizaCtsBcp,
		Foto:                        strings.TrimSpace(u.Foto),
		CreatedAt:                   time.Now(),
	}
}
