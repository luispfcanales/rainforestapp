package pdf

import (
	"context"
	_ "embed"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/extension"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/luispfcanales/rainforestapp/pkg/models"
)

//go:embed rainforest.png
var logoBytes []byte

type PDFGenerator struct{}

func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

// Colores como en la documentación oficial
func getHeaderColor() *props.Color {
	return &props.Color{
		Red:   34,
		Green: 139,
		Blue:  34,
	}
}

func getTextColor() *props.Color {
	return &props.Color{
		Red:   55,
		Green: 55,
		Blue:  55,
	}
}

// decodeBase64Image decodifica el string base64 a bytes
func decodeBase64Image(base64Str string) ([]byte, error) {
	// Remover el prefijo data:image/jpeg;base64, si existe
	if strings.Contains(base64Str, "base64,") {
		parts := strings.Split(base64Str, "base64,")
		if len(parts) > 1 {
			base64Str = parts[1]
		}
	}

	// Decodificar base64
	imgBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, fmt.Errorf("error decodificando base64: %w", err)
	}

	return imgBytes, nil
}

// GenerateUsuarioPDF - Versión mejorada con diseño profesional (Cuadros)
func (g *PDFGenerator) GenerateUsuarioPDF(ctx context.Context, usuario *models.Usuario) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageSize(pagesize.A4).
		WithOrientation(orientation.Vertical).
		WithLeftMargin(15).
		WithTopMargin(15).
		WithRightMargin(15).
		WithBottomMargin(15).
		Build()

	m := maroto.New(cfg)

	// -- HEADER --
	m.AddRows(
		row.New(30).Add(
			// Logo (Izquierda)
			col.New(3).Add(
				func() core.Component {
					if len(logoBytes) > 0 {
						return image.NewFromBytes(logoBytes, extension.Png, props.Rect{
							Center:  true,
							Percent: 90,
						})
					}
					return text.New("LOGO", props.Text{
						Style: fontstyle.Bold,
						Align: align.Center,
						Top:   10,
					})
				}(),
			),
			// Título (Centro)
			col.New(6).Add(
				text.New("RAINFOREST ENTERPRISE", props.Text{
					Size:  16,
					Style: fontstyle.Bold,
					Align: align.Center,
					Color: getHeaderColor(),
					Top:   5,
				}),
				text.New("FICHA DE DATOS DEL PERSONAL", props.Text{
					Size:  12,
					Style: fontstyle.Bold,
					Align: align.Center,
					Top:   15,
				}),
			),
			// Foto (Derecha)
			col.New(3).Add(
				func() core.Component {
					if usuario.Foto != "" && usuario.Foto != "null" && usuario.Foto != "undefined" {
						imgBytes, err := decodeBase64Image(usuario.Foto)
						if err == nil {
							return image.NewFromBytes(imgBytes, extension.Jpg, props.Rect{
								Percent: 95,
								Center:  true,
							})
						}
					}
					return text.New("[ SIN FOTO ]", props.Text{
						Size:  10,
						Style: fontstyle.BoldItalic,
						Align: align.Center,
						Top:   10,
						Color: &props.Color{Red: 150, Green: 150, Blue: 150},
					})
				}(),
			),
		),
		row.New(5),
		row.New(1).Add(col.New(12).Add(line.New(props.Line{Color: getHeaderColor(), Thickness: 2}))),
		row.New(5),
	)

	// Helpers
	formatBool := func(b bool) string {
		if b {
			return "Si"
		}
		return "No"
	}

	headerTextStyle := props.Text{Style: fontstyle.Bold, Size: 9, Align: align.Center}
	cellTextStyle := props.Text{Size: 8, Align: align.Center}

	// Helper para Sección Título
	addSectionTitle := func(title string) {
		m.AddRows(
			row.New(8).Add(col.New(12).Add(text.New(title, props.Text{
				Style: fontstyle.Bold, Size: 11, Color: getHeaderColor(), Top: 1,
			}))),
			row.New(1).Add(col.New(12).Add(line.New(props.Line{Color: getHeaderColor(), Thickness: 1}))),
			row.New(3),
		)
	}

	// Helper para celda con borde simulado (Label: Value)
	addLabelValue := func(label, value string) core.Component {
		if value == "" {
			value = "-"
		}
		// En Maroto v2, para simular borde en celda, podemos usar simplemente texto por ahora
		// Ya que la implementación de bordes exactos por celda requiere layouts complejos.
		// Usaremos formato "Label: Value"
		return text.New(fmt.Sprintf("%s: %s", label, value), props.Text{Size: 8, Align: align.Left, Left: 2})
	}

	// -- I. DATOS PERSONALES --
	addSectionTitle("I. DATOS PERSONALES")

	m.AddRows(
		row.New(6).Add(
			col.New(3).Add(addLabelValue("DNI", usuario.Dni)),
			col.New(5).Add(addLabelValue("Apellidos", fmt.Sprintf("%s %s", usuario.ApellidoPaterno, usuario.ApellidoMaterno))),
			col.New(4).Add(addLabelValue("Nombres", usuario.Nombres)),
		),
		row.New(6).Add(
			col.New(3).Add(addLabelValue("F. Nacimiento", usuario.FechaNacimiento)),
			col.New(3).Add(addLabelValue("Sexo", usuario.Sexo)),
			col.New(3).Add(addLabelValue("Estado Civil", usuario.EstadoCivil)),
			col.New(3).Add(addLabelValue("G. Sanguíneo", usuario.GrupoSanguineo)),
		),
		row.New(6).Add(
			col.New(3).Add(addLabelValue("Licencia", usuario.LicenciaConducir)),
			col.New(3).Add(addLabelValue("Cat. Licencia", usuario.CategoriaLicencia)),
		),
		row.New(6).Add(
			col.New(12).Add(addLabelValue("Dirección Domicilio", usuario.DireccionDomicilio)),
		),
		row.New(6).Add(
			col.New(12).Add(addLabelValue("Lugar Nacimiento", fmt.Sprintf("%s - %s - %s", usuario.LugarNacimientoDepartamento, usuario.LugarNacimientoProvincia, usuario.LugarNacimientoDistrito))),
		),
		row.New(5),
	)

	// -- II. CONTACTO Y EMERGENCIA --
	addSectionTitle("II. CONTACTO Y EMERGENCIA")
	m.AddRows(
		row.New(6).Add(
			col.New(4).Add(addLabelValue("Teléfono", usuario.Telefono)),
			col.New(4).Add(addLabelValue("Email", usuario.Email)),
		),
		row.New(1).Add(col.New(12).Add(line.New(props.Line{Color: getTextColor(), Thickness: 0.5}))), // Separador
		row.New(6).Add(
			col.New(6).Add(addLabelValue("Contacto Emergencia", usuario.ContactoNombre)),
			col.New(3).Add(addLabelValue("Parentesco", usuario.ContactoParentesco)),
			col.New(3).Add(addLabelValue("Celular", usuario.ContactoCelular)),
		),
		row.New(6).Add(
			col.New(12).Add(addLabelValue("Dirección Emergencia", usuario.ContactoDireccion)),
		),
		row.New(5),
	)

	// -- III. DATOS LABORALES --
	addSectionTitle("III. DATOS LABORALES")
	m.AddRows(
		row.New(6).Add(
			col.New(4).Add(addLabelValue("Puesto", usuario.PuestoActual)),
			col.New(4).Add(addLabelValue("Lugar", usuario.LugarTrabajo)),
			col.New(4).Add(addLabelValue("F. Ingreso", usuario.FechaIngreso)),
		),
		row.New(6).Add(
			col.New(4).Add(addLabelValue("Régimen Pénsion", usuario.RegimenPensionario)),
			col.New(4).Add(addLabelValue("AFP/ONP", usuario.AfpNombre+" "+usuario.Cuspp)),
			col.New(4).Add(addLabelValue("Salud", usuario.RegimenSalud)),
		),
		row.New(6).Add(
			col.New(12).Add(addLabelValue("Situación Contractual", usuario.SituacionContractual)),
		),
		row.New(5),
	)

	// -- IV. APERTURA DE CUENTAS --
	addSectionTitle("IV. APERTURA DE CUENTAS")
	m.AddRows(
		row.New(6).Add(
			col.New(6).Add(addLabelValue("Autoriza Cuenta Sueldo BCP", formatBool(usuario.AutorizaBcp))),
			col.New(6).Add(addLabelValue("Autoriza CTS BCP", formatBool(usuario.AutorizaCtsBcp))),
		),
	)
	if usuario.AutorizaOtroBanco {
		m.AddRows(row.New(6).Add(col.New(12).Add(addLabelValue("Otro Banco", fmt.Sprintf("%s (Cta: %s / CCI: %s)", usuario.OtroBancoNombre, usuario.OtroBancoCuenta, usuario.OtroBancoCci)))))
	}
	m.AddRows(row.New(5))

	// -- V. DATOS FAMILIARES --
	addSectionTitle("V. DATOS FAMILIARES")

	// Cónyuge Table
	m.AddRows(row.New(6).Add(col.New(12).Add(text.New("Datos del Cónyuge / Conviviente", props.Text{Style: fontstyle.Bold, Size: 9}))))
	// Header Table
	m.AddRows(row.New(6).Add(
		col.New(6).Add(text.New("Apellidos y Nombres", headerTextStyle)),
		col.New(3).Add(text.New("DNI", headerTextStyle)),
		col.New(3).Add(text.New("F. Nacimiento", headerTextStyle)),
	))
	m.AddRows(row.New(1).Add(col.New(12).Add(line.New(props.Line{Color: getTextColor(), Thickness: 0.5}))))
	// Data or Empty
	if usuario.DatosConyuge != nil && usuario.DatosConyuge.ApellidosNombres != "" {
		m.AddRows(row.New(5).Add(
			col.New(6).Add(text.New(usuario.DatosConyuge.ApellidosNombres, cellTextStyle)),
			col.New(3).Add(text.New(usuario.DatosConyuge.Dni, cellTextStyle)),
			col.New(3).Add(text.New(usuario.DatosConyuge.FechaNacimiento, cellTextStyle)),
		))
	} else {
		m.AddRows(row.New(5).Add(col.New(12).Add(text.New("- Sin información -", cellTextStyle))))
	}
	m.AddRows(row.New(5))

	// Hijos Table
	m.AddRows(row.New(6).Add(col.New(12).Add(text.New("Hijos", props.Text{Style: fontstyle.Bold, Size: 9}))))
	m.AddRows(row.New(6).Add(
		col.New(6).Add(text.New("Apellidos y Nombres", headerTextStyle)),
		col.New(2).Add(text.New("DNI", headerTextStyle)),
		col.New(2).Add(text.New("F. Nac", headerTextStyle)),
		col.New(2).Add(text.New("Edad", headerTextStyle)),
	))
	m.AddRows(row.New(1).Add(col.New(12).Add(line.New(props.Line{Color: getTextColor(), Thickness: 0.5}))))
	if len(usuario.Hijos) > 0 {
		for _, h := range usuario.Hijos {
			m.AddRows(row.New(5).Add(
				col.New(6).Add(text.New(h.ApellidosNombres, cellTextStyle)),
				col.New(2).Add(text.New(h.Dni, cellTextStyle)),
				col.New(2).Add(text.New(h.FechaNacimiento, cellTextStyle)),
				col.New(2).Add(text.New(fmt.Sprintf("%d", h.Edad), cellTextStyle)),
			))
		}
	} else {
		m.AddRows(row.New(5).Add(col.New(12).Add(text.New("- Sin hijos registrados -", cellTextStyle))))
	}
	m.AddRows(row.New(5))

	// -- VI. EDUCACIÓN --
	addSectionTitle("VI. EDUCACIÓN")

	// Basica
	m.AddRows(row.New(6).Add(col.New(12).Add(text.New("Educación Básica", props.Text{Style: fontstyle.Bold, Size: 9}))))
	m.AddRows(row.New(6).Add(
		col.New(2).Add(text.New("Nivel", headerTextStyle)),
		col.New(5).Add(text.New("Institución", headerTextStyle)),
		col.New(2).Add(text.New("Desde", headerTextStyle)),
		col.New(2).Add(text.New("Hasta", headerTextStyle)),
		col.New(1).Add(text.New("Comp.", headerTextStyle)),
	))
	m.AddRows(row.New(1).Add(col.New(12).Add(line.New(props.Line{Color: getTextColor(), Thickness: 0.5}))))
	if len(usuario.EducacionBasica) > 0 {
		for _, e := range usuario.EducacionBasica {
			m.AddRows(row.New(5).Add(
				col.New(2).Add(text.New(e.Nivel, cellTextStyle)),
				col.New(5).Add(text.New(e.CentroEstudios, cellTextStyle)),
				col.New(2).Add(text.New(e.Desde, cellTextStyle)),
				col.New(2).Add(text.New(e.Hasta, cellTextStyle)),
				col.New(1).Add(text.New(formatBool(e.Completa), cellTextStyle)),
			))
		}
	} else {
		m.AddRows(row.New(5).Add(col.New(12).Add(text.New("- Sin información -", cellTextStyle))))
	}
	m.AddRows(row.New(5))

	// Superior
	m.AddRows(row.New(6).Add(col.New(12).Add(text.New("Educación Superior", props.Text{Style: fontstyle.Bold, Size: 9}))))
	m.AddRows(row.New(6).Add(
		col.New(2).Add(text.New("Nivel", headerTextStyle)),
		col.New(3).Add(text.New("Institución", headerTextStyle)),
		col.New(3).Add(text.New("Especialidad", headerTextStyle)),
		col.New(2).Add(text.New("Periodo", headerTextStyle)),
		col.New(2).Add(text.New("Grado", headerTextStyle)),
	))
	m.AddRows(row.New(1).Add(col.New(12).Add(line.New(props.Line{Color: getTextColor(), Thickness: 0.5}))))
	if len(usuario.EducacionSuperior) > 0 {
		for _, e := range usuario.EducacionSuperior {
			m.AddRows(row.New(5).Add(
				col.New(2).Add(text.New(e.Nivel, cellTextStyle)),
				col.New(3).Add(text.New(e.CentroEstudios, cellTextStyle)),
				col.New(3).Add(text.New(e.Especialidad, cellTextStyle)),
				col.New(2).Add(text.New(e.Desde+"-"+e.Hasta, cellTextStyle)),
				col.New(2).Add(text.New(e.GradoAcademico, cellTextStyle)),
			))
		}
	} else {
		m.AddRows(row.New(5).Add(col.New(12).Add(text.New("- Sin información -", cellTextStyle))))
	}
	m.AddRows(row.New(5))

	// -- VII. EXPERIENCIA LABORAL --
	addSectionTitle("VII. EXPERIENCIA LABORAL")
	m.AddRows(row.New(6).Add(
		col.New(3).Add(text.New("Empresa", headerTextStyle)),
		col.New(3).Add(text.New("Cargo", headerTextStyle)),
		col.New(2).Add(text.New("F. Inicio", headerTextStyle)),
		col.New(2).Add(text.New("F. Fin", headerTextStyle)),
		col.New(2).Add(text.New("Motivo", headerTextStyle)),
	))
	m.AddRows(row.New(1).Add(col.New(12).Add(line.New(props.Line{Color: getTextColor(), Thickness: 0.5}))))
	if len(usuario.ExperienciaLaboral) > 0 {
		for _, ex := range usuario.ExperienciaLaboral {
			m.AddRows(row.New(5).Add(
				col.New(3).Add(text.New(ex.Empresa, cellTextStyle)),
				col.New(3).Add(text.New(ex.Cargo, cellTextStyle)),
				col.New(2).Add(text.New(ex.FechaIngreso, cellTextStyle)),
				col.New(2).Add(text.New(ex.FechaCese, cellTextStyle)),
				col.New(2).Add(text.New(ex.MotivoCese, cellTextStyle)),
			))
		}
	} else {
		m.AddRows(row.New(5).Add(col.New(12).Add(text.New("- Sin experiencia registrada -", cellTextStyle))))
	}
	m.AddRows(row.New(5))

	// -- VIII. IDIOMAS --
	addSectionTitle("VIII. IDIOMAS")
	m.AddRows(row.New(6).Add(
		col.New(3).Add(text.New("Idioma", headerTextStyle)),
		col.New(3).Add(text.New("Lee", headerTextStyle)),
		col.New(3).Add(text.New("Habla", headerTextStyle)),
		col.New(3).Add(text.New("Escribe", headerTextStyle)),
	))
	m.AddRows(row.New(1).Add(col.New(12).Add(line.New(props.Line{Color: getTextColor(), Thickness: 0.5}))))
	if len(usuario.Idiomas) > 0 {
		for _, i := range usuario.Idiomas {
			m.AddRows(row.New(5).Add(
				col.New(3).Add(text.New(i.Idioma, cellTextStyle)),
				col.New(3).Add(text.New(i.Lee, cellTextStyle)),
				col.New(3).Add(text.New(i.Habla, cellTextStyle)),
				col.New(3).Add(text.New(i.Escribe, cellTextStyle)),
			))
		}
	} else {
		m.AddRows(row.New(5).Add(col.New(12).Add(text.New("- Sin idiomas registrados -", cellTextStyle))))
	}

	// FOOTER
	m.AddRows(
		row.New(15).Add(
			col.New(12).Add(
				line.New(props.Line{Color: getHeaderColor(), Thickness: 1}),
			),
		),
		row.New(5).Add(
			col.New(12).Add(
				text.New(fmt.Sprintf("Generado el: %s", time.Now().Format("02/01/2006 15:04:05")), props.Text{
					Size: 8, Align: align.Right, Color: getTextColor(),
				}),
			),
		),
	)

	// Generar
	document, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("error generando PDF: %w", err)
	}

	return document.GetBytes(), nil
}

// calcularTiempoEnSistema calcula cuánto tiempo lleva el usuario en el sistema
func (g *PDFGenerator) calcularTiempoEnSistema(fechaRegistro time.Time) string {
	diferencia := time.Since(fechaRegistro)

	dias := int(diferencia.Hours() / 24)
	if dias < 1 {
		return "Menos de 1 día"
	} else if dias == 1 {
		return "1 día"
	} else if dias < 30 {
		return fmt.Sprintf("%d días", dias)
	} else if dias < 365 {
		meses := dias / 30
		return fmt.Sprintf("%d meses", meses)
	} else {
		anos := dias / 365
		return fmt.Sprintf("%d años", anos)
	}
}
