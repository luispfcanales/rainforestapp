package pdf

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
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

// GenerateUsuarioPDF - Versión mejorada con diseño profesional
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
					// Intentar cargar el logo - búsqueda en varias rutas
					paths := []string{"rainforest.png", "./rainforest.png", "../../rainforest.png", "../rainforest.png"}
					var imgBytes []byte
					var err error
					for _, path := range paths {
						imgBytes, err = os.ReadFile(path)
						if err == nil {
							break
						}
					}

					if len(imgBytes) > 0 {
						return image.NewFromBytes(imgBytes, extension.Png, props.Rect{
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
					// Cuadro "Sin Foto"
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
	)

	m.AddRows(
		row.New(5),
		row.New(2).Add(
			col.New(12).Add(line.New(props.Line{Color: getHeaderColor(), Thickness: 2})),
		),
		row.New(5),
	)

	// Helper para pares etiqueta-valor
	addProp := func(label, value string) core.Component {
		if value == "" {
			value = "-"
		}
		return text.New(fmt.Sprintf("%s: %s", label, value), props.Text{
			Size:  9,
			Align: align.Left,
		})
	}

	formatBool := func(b bool) string {
		if b {
			return "Si"
		}
		return "No"
	}

	// -- I. DATOS PERSONALES --
	m.AddRows(row.New(10).Add(col.New(12).Add(text.New("I. DATOS PERSONALES", props.Text{
		Style: fontstyle.Bold, Size: 11, Color: getHeaderColor(),
	}))))

	m.AddRows(
		row.New(6).Add(
			col.New(4).Add(addProp("DNI", usuario.Dni)),
			col.New(8).Add(addProp("Apellidos y Nombres", fmt.Sprintf("%s %s %s", usuario.ApellidoPaterno, usuario.ApellidoMaterno, usuario.Nombres))),
		),
		row.New(6).Add(
			col.New(4).Add(addProp("Fecha Nacimiento", usuario.FechaNacimiento)),
			col.New(4).Add(addProp("Sexo", usuario.Sexo)),
			col.New(4).Add(addProp("Estado Civil", usuario.EstadoCivil)),
		),
		row.New(6).Add(
			col.New(4).Add(addProp("Licencia Conducir", usuario.LicenciaConducir)),
			col.New(4).Add(addProp("Categoría Licencia", usuario.CategoriaLicencia)),
			col.New(4).Add(addProp("Grupo Sanguíneo", usuario.GrupoSanguineo)),
		),
		row.New(6).Add(
			col.New(12).Add(addProp("Dirección", usuario.DireccionDomicilio)),
		),
		row.New(6).Add(
			col.New(12).Add(addProp("Lugar Nacimiento", fmt.Sprintf("%s - %s - %s", usuario.LugarNacimientoDepartamento, usuario.LugarNacimientoProvincia, usuario.LugarNacimientoDistrito))),
		),
	)

	m.AddRows(row.New(5))

	// -- II. DATOS DE CONTACTO Y EMERGENCIA --
	m.AddRows(row.New(10).Add(col.New(12).Add(text.New("II. CONTACTO Y EMERGENCIA", props.Text{
		Style: fontstyle.Bold, Size: 11, Color: getHeaderColor(),
	}))))

	m.AddRows(
		row.New(6).Add(
			col.New(6).Add(addProp("Teléfono", usuario.Telefono)),
			col.New(6).Add(addProp("Email", usuario.Email)),
		),
		row.New(6).Add(
			col.New(6).Add(addProp("Contacto Emergencia", usuario.ContactoNombre)),
			col.New(3).Add(addProp("Parentesco", usuario.ContactoParentesco)),
			col.New(3).Add(addProp("Celular", usuario.ContactoCelular)),
		),
		row.New(6).Add(
			col.New(12).Add(addProp("Dirección Emergencia", usuario.ContactoDireccion)),
		),
	)

	m.AddRows(row.New(5))

	// -- III. DATOS LABORALES --
	m.AddRows(row.New(10).Add(col.New(12).Add(text.New("III. DATOS LABORALES", props.Text{
		Style: fontstyle.Bold, Size: 11, Color: getHeaderColor(),
	}))))

	m.AddRows(
		row.New(6).Add(
			col.New(4).Add(addProp("Puesto Actual", usuario.PuestoActual)),
			col.New(4).Add(addProp("Lugar Trabajo", usuario.LugarTrabajo)),
			col.New(4).Add(addProp("Fecha Ingreso", usuario.FechaIngreso)),
		),
		row.New(6).Add(
			col.New(4).Add(addProp("Régimen Pensionario", usuario.RegimenPensionario)),
			col.New(4).Add(addProp("AFP/ONP", usuario.AfpNombre+" "+usuario.Cuspp)),
			col.New(4).Add(addProp("Regimen Salud", usuario.RegimenSalud)),
		),
		row.New(6).Add(
			col.New(4).Add(addProp("Situación Contractual", usuario.SituacionContractual)),
			col.New(4).Add(addProp("Autoriza BCP", formatBool(usuario.AutorizaBcp))),
			col.New(4).Add(addProp("Autoriza CTS BCP", formatBool(usuario.AutorizaCtsBcp))),
		),
		row.New(6).Add(
			col.New(12).Add(addProp("Otro Banco", fmt.Sprintf("%s (%s - CCI: %s)", usuario.OtroBancoNombre, usuario.OtroBancoCuenta, usuario.OtroBancoCci))),
		),
	)

	m.AddRows(row.New(5))

	// -- IV. EDUCACIÓN --
	m.AddRows(row.New(10).Add(col.New(12).Add(text.New("IV. EDUCACIÓN", props.Text{
		Style: fontstyle.Bold, Size: 11, Color: getHeaderColor(),
	}))))

	if len(usuario.EducacionBasica) > 0 {
		m.AddRows(row.New(6).Add(col.New(12).Add(text.New("Educación Básica:", props.Text{Style: fontstyle.Bold, Size: 9}))))
		for _, edu := range usuario.EducacionBasica {
			m.AddRows(row.New(5).Add(
				col.New(12).Add(text.New(fmt.Sprintf("%s - %s (%s - %s) Completa: %s", edu.Nivel, edu.CentroEstudios, edu.Desde, edu.Hasta, formatBool(edu.Completa)), props.Text{Size: 8})),
			))
		}
	}
	if len(usuario.EducacionSuperior) > 0 {
		m.AddRows(row.New(6).Add(col.New(12).Add(text.New("Educación Superior:", props.Text{Style: fontstyle.Bold, Size: 9}))))
		for _, edu := range usuario.EducacionSuperior {
			m.AddRows(row.New(5).Add(
				col.New(12).Add(text.New(fmt.Sprintf("%s en %s (%s - %s) - %s. Grado: %s", edu.Nivel, edu.CentroEstudios, edu.Desde, edu.Hasta, edu.Especialidad, edu.GradoAcademico), props.Text{Size: 8})),
			))
		}
	}

	m.AddRows(row.New(5))

	// -- V. CAPACITACIONES --
	if len(usuario.Capacitaciones) > 0 {
		m.AddRows(row.New(10).Add(col.New(12).Add(text.New("V. CAPACITACIONES", props.Text{
			Style: fontstyle.Bold, Size: 11, Color: getHeaderColor(),
		}))))
		for _, cap := range usuario.Capacitaciones {
			m.AddRows(row.New(5).Add(
				col.New(12).Add(text.New(fmt.Sprintf("- %s (%s) - %d Horas", cap.Nombre, cap.Institucion, cap.Horas), props.Text{Size: 8})),
			))
		}
		m.AddRows(row.New(5))
	}

	// -- VI. EXPERIENCIA LABORAL --
	if len(usuario.ExperienciaLaboral) > 0 {
		m.AddRows(row.New(10).Add(col.New(12).Add(text.New("VI. EXPERIENCIA LABORAL", props.Text{
			Style: fontstyle.Bold, Size: 11, Color: getHeaderColor(),
		}))))
		for _, exp := range usuario.ExperienciaLaboral {
			m.AddRows(row.New(5).Add(
				col.New(12).Add(text.New(fmt.Sprintf("- %s en %s (%s al %s) - %s. Motivo: %s", exp.Cargo, exp.Empresa, exp.FechaIngreso, exp.FechaCese, exp.TiempoPermanencia, exp.MotivoCese), props.Text{Size: 8})),
			))
		}
		m.AddRows(row.New(5))
	}

	// -- VII. IDIOMAS --
	if len(usuario.Idiomas) > 0 {
		m.AddRows(row.New(10).Add(col.New(12).Add(text.New("VII. IDIOMAS", props.Text{
			Style: fontstyle.Bold, Size: 11, Color: getHeaderColor(),
		}))))
		for _, idi := range usuario.Idiomas {
			m.AddRows(row.New(5).Add(
				col.New(12).Add(text.New(fmt.Sprintf("- %s: Lee(%s), Habla(%s), Escribe(%s)", idi.Idioma, idi.Lee, idi.Habla, idi.Escribe), props.Text{Size: 8})),
			))
		}
		m.AddRows(row.New(5))
	}

	// -- VIII. INFORMACIÓN FAMILIAR --
	if usuario.DatosConyuge != nil || len(usuario.Hijos) > 0 || len(usuario.Padres) > 0 {
		m.AddRows(row.New(10).Add(col.New(12).Add(text.New("VIII. INFORMACIÓN FAMILIAR", props.Text{
			Style: fontstyle.Bold, Size: 11, Color: getHeaderColor(),
		}))))

		if usuario.DatosConyuge != nil && usuario.DatosConyuge.ApellidosNombres != "" {
			m.AddRows(
				row.New(6).Add(
					col.New(12).Add(text.New("Cónyuge / Conviviente:", props.Text{Style: fontstyle.Bold, Size: 9})),
				),
				row.New(6).Add(
					col.New(6).Add(addProp("Nombre", usuario.DatosConyuge.ApellidosNombres)),
					col.New(3).Add(addProp("DNI", usuario.DatosConyuge.Dni)),
					col.New(3).Add(addProp("F. Nacimiento", usuario.DatosConyuge.FechaNacimiento)),
				),
			)
		}

		if len(usuario.Hijos) > 0 {
			m.AddRows(row.New(6).Add(col.New(12).Add(text.New("Hijos:", props.Text{Style: fontstyle.Bold, Size: 9}))))
			for i, hijo := range usuario.Hijos {
				m.AddRows(row.New(5).Add(
					col.New(12).Add(text.New(fmt.Sprintf("%d. %s (DNI: %s) - F. Nac: %s", i+1, hijo.ApellidosNombres, hijo.Dni, hijo.FechaNacimiento), props.Text{Size: 8})),
				))
			}
		}

		if len(usuario.Padres) > 0 {
			m.AddRows(row.New(6).Add(col.New(12).Add(text.New("Padres:", props.Text{Style: fontstyle.Bold, Size: 9}))))
			for i, padre := range usuario.Padres {
				vive := "No"
				if padre.Vive {
					vive = "Si"
				}
				m.AddRows(row.New(5).Add(
					col.New(12).Add(text.New(fmt.Sprintf("%d. %s - F. Nac: %s - Ocupación: %s - Vive: %s", i+1, padre.ApellidosNombres, padre.FechaNacimiento, padre.Ocupacion, vive), props.Text{Size: 8})),
				))
			}
		}
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
