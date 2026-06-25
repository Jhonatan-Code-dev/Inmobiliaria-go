package service

import (
	"context"
	"fmt"
	"time"

	"rentals-go/internal/domain"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type ReclamacionService struct {
	repo    domain.ReclamacionRepository
	empRepo domain.EmpresaRepository
}

func NewReclamacionService(repo domain.ReclamacionRepository, empRepo domain.EmpresaRepository) *ReclamacionService {
	return &ReclamacionService{
		repo:    repo,
		empRepo: empRepo,
	}
}

func (s *ReclamacionService) Listar(ctx context.Context, empresaID int, pag int) ([]*domain.Reclamacion, int, error) {
	if pag <= 0 {
		pag = 1
	}
	limite := 10
	return s.repo.ListarPaginado(ctx, empresaID, pag, limite)
}

func (s *ReclamacionService) Obtener(ctx context.Context, id int, empresaID int) (*domain.Reclamacion, error) {
	return s.repo.BuscarPorID(ctx, id, empresaID)
}

func (s *ReclamacionService) Registrar(ctx context.Context, rec *domain.Reclamacion) (*domain.Reclamacion, error) {
	// Generate correlative code
	prefix := fmt.Sprintf("REC-%d-", time.Now().Year())

	// Just get total count by paginating with 1 item, but repo list paginated returns total
	_, total, err := s.repo.ListarPaginado(ctx, rec.EmpresaID, 1, 1)
	if err != nil {
		return nil, err
	}

	rec.Codigo = fmt.Sprintf("%s%05d", prefix, total+1)
	rec.Estado = "PENDIENTE"
	rec.CreadoEn = time.Now()

	return s.repo.Crear(ctx, rec)
}

func (s *ReclamacionService) Responder(ctx context.Context, id int, empresaID int, respuesta string) (*domain.Reclamacion, error) {
	rec, err := s.repo.BuscarPorID(ctx, id, empresaID)
	if err != nil {
		return nil, err
	}

	rec.Estado = "RESUELTO"
	rec.RespuestaDetalle = respuesta
	now := time.Now()
	rec.RespondidoEn = &now

	return s.repo.Actualizar(ctx, rec)
}

func (s *ReclamacionService) Eliminar(ctx context.Context, id int, empresaID int) error {
	return s.repo.Eliminar(ctx, id, empresaID)
}

func (s *ReclamacionService) GenerarPDF(ctx context.Context, id int, empresaID int) ([]byte, error) {
	rec, err := s.repo.BuscarPorID(ctx, id, empresaID)
	if err != nil {
		return nil, err
	}

	emp, err := s.empRepo.BuscarPorID(ctx, rec.EmpresaID)
	empName := "Alquilamax Tenant"
	if err == nil && emp != nil {
		empName = emp.Nombre
	}

	cfg := config.NewBuilder().
		WithPageNumber().
		Build()

	m := maroto.New(cfg)

	// Header Title
	m.AddRow(15,
		col.New(12).Add(
			text.New("HOJA DE RECLAMACIÓN", props.Text{
				Size:  16,
				Style: fontstyle.Bold,
				Align: align.Center,
			}),
		),
	)
	m.AddRow(10,
		col.New(12).Add(
			text.New("LIBRO DE RECLAMACIONES VIRTUAL", props.Text{
				Size:  12,
				Style: fontstyle.Bold,
				Align: align.Center,
			}),
		),
	)

	// Details header
	m.AddRow(8,
		col.New(6).Add(text.New(fmt.Sprintf("Establecimiento: %s", empName), props.Text{Size: 9, Style: fontstyle.Bold})),
		col.New(6).Add(text.New(fmt.Sprintf("Hoja Nro: %s", rec.Codigo), props.Text{Size: 9, Style: fontstyle.Bold, Align: align.Right})),
	)
	m.AddRow(8,
		col.New(6).Add(text.New(fmt.Sprintf("Fecha de Registro: %s", rec.CreadoEn.Format("2006-01-02 15:04")), props.Text{Size: 9})),
		col.New(6).Add(text.New(fmt.Sprintf("Estado: %s", rec.Estado), props.Text{Size: 9, Style: fontstyle.Bold, Align: align.Right})),
	)

	m.AddRow(4, col.New(12)) // Spacer

	// 1. Consumer Ident
	m.AddRow(8,
		col.New(12).Add(text.New("1. IDENTIFICACIÓN DEL CONSUMIDOR RECLAMANTE", props.Text{
			Size:  10,
			Style: fontstyle.Bold,
		})),
	)
	m.AddRow(6,
		col.New(12).Add(text.New(fmt.Sprintf("Nombre Completo: %s %s", rec.Nombres, rec.Apellidos), props.Text{Size: 9})),
	)
	m.AddRow(6,
		col.New(6).Add(text.New(fmt.Sprintf("Documento: %s - %s", rec.TipoDocumento, rec.NumeroDocumento), props.Text{Size: 9})),
		col.New(6).Add(text.New(fmt.Sprintf("Teléfono: %s", rec.Telefono), props.Text{Size: 9})),
	)
	m.AddRow(6,
		col.New(12).Add(text.New(fmt.Sprintf("E-mail: %s", rec.Email), props.Text{Size: 9})),
	)
	m.AddRow(6,
		col.New(12).Add(text.New(fmt.Sprintf("Dirección: %s", rec.Direccion), props.Text{Size: 9})),
	)
	if rec.MenorEdad {
		m.AddRow(6,
			col.New(12).Add(text.New(fmt.Sprintf("Apoderado (Menor de edad): %s", rec.NombreApoderado), props.Text{Size: 9, Style: fontstyle.Italic})),
		)
	}

	m.AddRow(4, col.New(12)) // Spacer

	// 2. Bien Contratado
	m.AddRow(8,
		col.New(12).Add(text.New("2. IDENTIFICACIÓN DEL BIEN CONTRATADO", props.Text{
			Size:  10,
			Style: fontstyle.Bold,
		})),
	)
	m.AddRow(6,
		col.New(6).Add(text.New(fmt.Sprintf("Tipo de Bien: %s", rec.TipoBien), props.Text{Size: 9})),
		col.New(6).Add(text.New(fmt.Sprintf("Monto Reclamado: S/. %.2f", rec.MontoReclamado), props.Text{Size: 9})),
	)
	m.AddRow(8,
		col.New(12).Add(text.New(fmt.Sprintf("Descripción del Bien: %s", rec.DescripcionBien), props.Text{Size: 9})),
	)

	m.AddRow(4, col.New(12)) // Spacer

	// 3. Detalle de Reclamación
	m.AddRow(8,
		col.New(12).Add(text.New("3. DETALLE DE LA RECLAMACIÓN Y PEDIDO DEL CONSUMIDOR", props.Text{
			Size:  10,
			Style: fontstyle.Bold,
		})),
	)
	m.AddRow(6,
		col.New(12).Add(text.New(fmt.Sprintf("Tipo de Incidencia: %s", rec.TipoReclamacion), props.Text{Size: 9, Style: fontstyle.Bold})),
	)
	m.AddRow(12,
		col.New(12).Add(text.New(fmt.Sprintf("Detalle: %s", rec.DetalleReclamacion), props.Text{Size: 9})),
	)
	m.AddRow(12,
		col.New(12).Add(text.New(fmt.Sprintf("Pedido Solicitado: %s", rec.PedidoConsumidor), props.Text{Size: 9})),
	)

	m.AddRow(4, col.New(12)) // Spacer

	// 4. Respuesta Proveedor
	m.AddRow(8,
		col.New(12).Add(text.New("4. ACCIONES Y RESPUESTA DEL PROVEEDOR", props.Text{
			Size:  10,
			Style: fontstyle.Bold,
		})),
	)
	if rec.Estado == "RESUELTO" {
		m.AddRow(6,
			col.New(12).Add(text.New(fmt.Sprintf("Fecha Respuesta: %s", rec.RespondidoEn.Format("2006-01-02 15:04")), props.Text{Size: 9})),
		)
		m.AddRow(15,
			col.New(12).Add(text.New(fmt.Sprintf("Detalle de Respuesta: %s", rec.RespuestaDetalle), props.Text{Size: 9})),
		)
	} else {
		m.AddRow(8,
			col.New(12).Add(text.New("Reclamación en proceso de evaluación. El proveedor responderá dentro del plazo legal.", props.Text{
				Size:  9,
				Style: fontstyle.Italic,
			})),
		)
	}

	doc, err := m.Generate()
	if err != nil {
		return nil, err
	}

	return doc.GetBytes(), nil
}
