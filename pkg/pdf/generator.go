package pdf

import (
	"context"
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

// GenerateUsuarioPDFSimple - Versión simple con imagen en cabecera
func (g *PDFGenerator) GenerateUsuarioPDFSimple(ctx context.Context, usuario *models.Usuario) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageSize(pagesize.A4).
		WithOrientation(orientation.Horizontal).
		WithLeftMargin(20).
		WithTopMargin(15).
		WithRightMargin(20).
		WithBottomMargin(15).
		Build()

	m := maroto.New(cfg)

	// Header con imagen a la derecha
	m.AddRows(
		row.New(25).Add(
			// Columna para el título (izquierda - 9 columnas)
			col.New(9).Add(
				text.New("RAINFOREST ENTERPRISE", props.Text{
					Size:  18,
					Style: fontstyle.Bold,
					Align: align.Left,
					Color: getHeaderColor(),
					Top:   5,
				}),
				text.New("FICHA DE USUARIO", props.Text{
					Size:  12,
					Align: align.Left,
					Top:   12,
				}),
			),
			// Columna para la foto carnet (derecha - 3 columnas)
			col.New(3).Add(
				func() core.Component {
					if usuario.Foto != "" && usuario.Foto != "null" && usuario.Foto != "undefined" {
						imgBytes, err := decodeBase64Image(usuario.Foto)
						if err == nil {
							// Especificar el tipo de extensión (JPEG o PNG)
							return image.NewFromBytes(imgBytes, extension.Jpg, props.Rect{
								Percent: 80,
								Center:  true,
							})
						}
					}

					// Si no hay foto, mostrar placeholder
					return text.New("Sin Foto", props.Text{
						Size:  8,
						Align: align.Center,
						Color: getTextColor(),
					})
				}(),
			),
		),
	)

	// Separator
	m.AddRows(
		row.New(5).Add(
			col.New(12).Add(
				line.New(props.Line{
					Color:     getHeaderColor(),
					Thickness: 1,
				}),
			),
		),
	)

	// Title
	m.AddRows(
		row.New(15).Add(
			col.New(12).Add(
				text.New("INFORMACIÓN DEL USUARIO", props.Text{
					Size:  14,
					Style: fontstyle.Bold,
					Align: align.Center,
				}),
			),
		),
	)

	// User Data
	m.AddRows(
		row.New(7).Add(
			col.New(4).Add(
				text.New("ID:", props.Text{
					Size:  9,
					Style: fontstyle.Bold,
				}),
			),
			col.New(8).Add(
				text.New(usuario.ID, props.Text{
					Size: 9,
				}),
			),
		),
		row.New(7).Add(
			col.New(4).Add(
				text.New("Nombre:", props.Text{
					Size:  9,
					Style: fontstyle.Bold,
				}),
			),
			col.New(8).Add(
				text.New(usuario.Nombres, props.Text{
					Size: 9,
				}),
			),
		),
		row.New(7).Add(
			col.New(4).Add(
				text.New("Apellido:", props.Text{
					Size:  9,
					Style: fontstyle.Bold,
				}),
			),
			col.New(8).Add(
				text.New(fmt.Sprintf("%s %s", usuario.ApellidoPaterno, usuario.ApellidoMaterno), props.Text{
					Size: 9,
				}),
			),
		),
		row.New(7).Add(
			col.New(4).Add(
				text.New("Email:", props.Text{
					Size:  9,
					Style: fontstyle.Bold,
				}),
			),
			col.New(8).Add(
				text.New(usuario.Email, props.Text{
					Size: 9,
				}),
			),
		),
		row.New(7).Add(
			col.New(4).Add(
				text.New("Teléfono:", props.Text{
					Size:  9,
					Style: fontstyle.Bold,
				}),
			),
			col.New(8).Add(
				text.New(usuario.Telefono, props.Text{
					Size: 9,
				}),
			),
		),
		row.New(7).Add(
			col.New(4).Add(
				text.New("Registrado:", props.Text{
					Size:  9,
					Style: fontstyle.Bold,
				}),
			),
			col.New(8).Add(
				text.New(usuario.CreatedAt.Format("02/01/2006 15:04"), props.Text{
					Size: 9,
				}),
			),
		),
	)

	// Footer
	m.AddRows(
		row.New(20).Add(col.New(12)),
		row.New(8).Add(
			col.New(12).Add(
				text.New(fmt.Sprintf("Generado el %s", time.Now().Format("02/01/2006 15:04")), props.Text{
					Size:  8,
					Align: align.Center,
					Color: getTextColor(),
				}),
			),
		),
		row.New(6).Add(
			col.New(12).Add(
				text.New("Rainforest Enterprise - Sistema de Gestión", props.Text{
					Size:  7,
					Align: align.Center,
					Color: getTextColor(),
				}),
			),
		),
	)

	// Generar PDF
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
