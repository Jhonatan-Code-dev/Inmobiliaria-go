package service

import (
	"context"
	"fmt"
	"rentals-go/internal/domain"
	"time"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/xuri/excelize/v2"
)

type GastoService struct {
	repo       domain.GastoRepository
	movRepo    domain.MovimientoCajaRepository
	tipoPgRepo domain.TipoPagoRepository
}

func NewGastoService(repo domain.GastoRepository, movRepo domain.MovimientoCajaRepository, tipoPgRepo domain.TipoPagoRepository) *GastoService {
	return &GastoService{
		repo:       repo,
		movRepo:    movRepo,
		tipoPgRepo: tipoPgRepo,
	}
}

func (s *GastoService) ListarTiposPago(ctx context.Context) ([]*domain.TipoPago, error) {
	return s.tipoPgRepo.Listar(ctx)
}

func (s *GastoService) Listar(ctx context.Context, filtros domain.GastoFiltros) ([]*domain.Gasto, int, error) {
	if filtros.Pagina <= 0 {
		filtros.Pagina = 1
	}
	if filtros.Limite <= 0 {
		filtros.Limite = 10
	}
	return s.repo.ListarPaginado(ctx, filtros)
}

func (s *GastoService) ObtenerGasto(ctx context.Context, id int, empresaID int) (*domain.Gasto, error) {
	g, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	if g.EmpresaID != empresaID {
		return nil, fmt.Errorf("gasto no pertenece a la empresa")
	}
	return g, nil
}

func (s *GastoService) RegistrarGasto(ctx context.Context, g *domain.Gasto) (*domain.Gasto, error) {
	// 1. Crear el gasto
	nuevoGasto, err := s.repo.Crear(ctx, g)
	if err != nil {
		return nil, err
	}

	// Fetch tipo pago name to use as method for MovimientoCaja
	metodo := "otro" // Default fallback
	tipos, err := s.tipoPgRepo.Listar(ctx)
	if err == nil {
		for _, t := range tipos {
			if t.ID == nuevoGasto.TipoPagoID {
				metodo = t.Nombre
				break
			}
		}
	}

	// 2. Crear el movimiento de caja (egreso)
	mov := &domain.MovimientoCaja{
		EmpresaID:       nuevoGasto.EmpresaID,
		GastoID:         &nuevoGasto.ID,
		Tipo:            "egreso",
		Concepto:        fmt.Sprintf("Gasto: %s", nuevoGasto.Descripcion),
		FechaMovimiento: nuevoGasto.Fecha,
		Moneda:          "PEN", // Por defecto usaremos PEN para el movimiento simplificado de ser necesario
		Monto:           nuevoGasto.Monto,
		Metodo:          metodo,
	}
	_, err = s.movRepo.Crear(ctx, mov)
	if err != nil {
		fmt.Printf("Error al crear movimiento de caja para gasto %d: %v\n", nuevoGasto.ID, err)
	}

	return nuevoGasto, nil
}

func (s *GastoService) ActualizarGasto(ctx context.Context, g *domain.Gasto) (*domain.Gasto, error) {
	// Verificar pertenencia antes de actualizar
	existente, err := s.repo.BuscarPorID(ctx, g.ID)
	if err != nil {
		return nil, err
	}
	if existente.EmpresaID != g.EmpresaID {
		return nil, fmt.Errorf("no autorizado para actualizar este gasto")
	}

	return s.repo.Actualizar(ctx, g)
}

func (s *GastoService) EliminarGasto(ctx context.Context, id int, empresaID int) error {
	existente, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return err
	}
	if existente.EmpresaID != empresaID {
		return fmt.Errorf("no autorizado para eliminar este gasto")
	}

	return s.repo.Eliminar(ctx, id)
}

func (s *GastoService) ExportarExcel(ctx context.Context, filtros domain.GastoFiltros) ([]byte, error) {
	// 1. Obtener datos (sin paginación para el reporte completo)
	filtros.Pagina = 1
	filtros.Limite = 10000 // Un límite alto razonable
	gastos, _, err := s.repo.ListarPaginado(ctx, filtros)
	if err != nil {
		return nil, err
	}

	// 2. Obtener tipos de pago para mapear nombres
	tipos, _ := s.tipoPgRepo.Listar(ctx)
	mapTipos := make(map[int]string)
	for _, t := range tipos {
		mapTipos[t.ID] = t.Nombre
	}

	// 3. Crear Excel
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Gastos"
	f.SetSheetName("Sheet1", sheet)

	// Encabezados
	headers := []string{"ID", "Fecha", "Descripción", "Monto", "Tipo de Pago"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Estilo para encabezado
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#CCCCCC"}, Pattern: 1},
	})
	f.SetCellStyle(sheet, "A1", "E1", style)

	// Llenar datos
	for i, g := range gastos {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), g.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), g.Fecha.Format("2006-01-02"))
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), g.Descripcion)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), g.Monto)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), mapTipos[g.TipoPagoID])
	}

	// Ajustar ancho de columnas
	f.SetColWidth(sheet, "A", "A", 10)
	f.SetColWidth(sheet, "B", "B", 15)
	f.SetColWidth(sheet, "C", "C", 40)
	f.SetColWidth(sheet, "D", "D", 15)
	f.SetColWidth(sheet, "E", "E", 20)

	// Guardar en buffer
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (s *GastoService) ExportarPDF(ctx context.Context, filtros domain.GastoFiltros) ([]byte, error) {
	// 1. Obtener datos
	filtros.Pagina = 1
	filtros.Limite = 10000
	gastos, _, err := s.repo.ListarPaginado(ctx, filtros)
	if err != nil {
		return nil, err
	}

	tipos, _ := s.tipoPgRepo.Listar(ctx)
	mapTipos := make(map[int]string)
	for _, t := range tipos {
		mapTipos[t.ID] = t.Nombre
	}

	// 2. Crear PDF con Maroto v2
	cfg := config.NewBuilder().
		WithPageNumber().
		Build()

	m := maroto.New(cfg)

	// Header
	m.AddRow(20,
		col.New(12).Add(
			text.New("REPORTE DE GASTOS", props.Text{
				Size:  16,
				Style: fontstyle.Bold,
				Align: align.Center,
			}),
		),
	)

	// Información de Filtros (opcional)
	m.AddRow(10,
		col.New(12).Add(
			text.New(fmt.Sprintf("Fecha Generación: %s", time.Now().Format("2006-01-02 15:04")), props.Text{
				Size:  10,
				Align: align.Right,
			}),
		),
	)

	// Tabla Header
	m.AddRow(10,
		col.New(1).Add(text.New("ID", props.Text{Style: fontstyle.Bold, Size: 9})),
		col.New(2).Add(text.New("Fecha", props.Text{Style: fontstyle.Bold, Size: 9})),
		col.New(5).Add(text.New("Descripción", props.Text{Style: fontstyle.Bold, Size: 9})),
		col.New(2).Add(text.New("Tipo Pago", props.Text{Style: fontstyle.Bold, Size: 9})),
		col.New(2).Add(text.New("Monto", props.Text{Style: fontstyle.Bold, Size: 9, Align: align.Right})),
	)

	var total float64
	for _, g := range gastos {
		total += g.Monto
		m.AddRow(8,
			col.New(1).Add(text.New(fmt.Sprintf("%d", g.ID), props.Text{Size: 8})),
			col.New(2).Add(text.New(g.Fecha.Format("2006-01-02"), props.Text{Size: 8})),
			col.New(5).Add(text.New(g.Descripcion, props.Text{Size: 8})),
			col.New(2).Add(text.New(mapTipos[g.TipoPagoID], props.Text{Size: 8})),
			col.New(2).Add(text.New(fmt.Sprintf("%.2f", g.Monto), props.Text{Size: 8, Align: align.Right})),
		)
	}

	// Línea de Total
	m.AddRow(10,
		col.New(10).Add(text.New("TOTAL:", props.Text{Style: fontstyle.Bold, Size: 10, Align: align.Right})),
		col.New(2).Add(text.New(fmt.Sprintf("%.2f", total), props.Text{Style: fontstyle.Bold, Size: 10, Align: align.Right})),
	)

	doc, err := m.Generate()
	if err != nil {
		return nil, err
	}

	return doc.GetBytes(), nil
}
