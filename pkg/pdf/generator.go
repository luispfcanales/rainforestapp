package pdf

import (
	"context"
	"fmt"
	"time"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
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

func getLightColor() *props.Color {
	return &props.Color{
		Red:   240,
		Green: 240,
		Blue:  240,
	}
}

// GenerateUsuarioPDFSimple - Versión simple basada en la documentación
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

	// Header
	m.AddRows(
		row.New(20).Add(
			col.New(12).Add(
				text.New("RAINFOREST ENTERPRISE", props.Text{
					Size:  18,
					Style: fontstyle.Bold,
					Align: align.Center,
					Color: getHeaderColor(),
				}),
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
				text.New("FICHA DE USUARIO", props.Text{
					Size:  14,
					Style: fontstyle.Bold,
					Align: align.Center,
				}),
			),
		),
	)

	// User Info Header
	m.AddRows(
		row.New(8).Add(
			col.New(12).Add(
				text.New("Información del Usuario", props.Text{
					Size:  12,
					Style: fontstyle.Bold,
					Align: align.Center,
					Color: &props.WhiteColor,
				}),
			),
		).WithStyle(&props.Cell{
			BackgroundColor: getHeaderColor(),
		}),
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
				text.New(usuario.Nombre, props.Text{
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
				text.New(usuario.Apellido, props.Text{
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
