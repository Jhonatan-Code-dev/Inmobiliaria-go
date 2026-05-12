package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"rentals-go/internal/domain"
)

type AlquilerService struct {
	repo    domain.AlquilerRepository
	cliente domain.ClienteRepository
	empresa domain.EmpresaRepository
}

func NewAlquilerService(repo domain.AlquilerRepository, cliente domain.ClienteRepository, empresa domain.EmpresaRepository) *AlquilerService {
	return &AlquilerService{repo: repo, cliente: cliente, empresa: empresa}
}

func (s *AlquilerService) Listar(ctx context.Context, filtros domain.AlquilerFiltros) ([]*domain.Alquiler, int, error) {
	if filtros.Pagina <= 0 {
		filtros.Pagina = 1
	}
	if filtros.Limite <= 0 {
		filtros.Limite = 10
	}
	return s.repo.ListarPaginado(ctx, filtros)
}

func (s *AlquilerService) Obtener(ctx context.Context, id int, empresaID int) (*domain.Alquiler, error) {
	item, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}
	if item.EmpresaID != empresaID {
		return nil, fmt.Errorf("%w: alquiler no pertenece a la empresa", domain.ErrForbidden)
	}
	return item, nil
}

func (s *AlquilerService) Crear(ctx context.Context, alquiler *domain.Alquiler) (*domain.Alquiler, error) {
	if alquiler.Tipo == "" {
		alquiler.Tipo = "alquiler"
	}
	if alquiler.Moneda == "" {
		alquiler.Moneda = "PEN"
	}
	if alquiler.Estado == "" {
		alquiler.Estado = "activo"
	}
	alquiler.ActivoParaCobro = true
	alquiler.Codigo = fmt.Sprintf("ALQ-%d", time.Now().UTC().UnixNano())
	return s.repo.Crear(ctx, alquiler)
}

func (s *AlquilerService) Actualizar(ctx context.Context, id int, empresaID int, alq *domain.Alquiler) (*domain.Alquiler, error) {
	old, err := s.Obtener(ctx, id, empresaID)
	if err != nil {
		return nil, err
	}
	alq.ID = id
	alq.EmpresaID = empresaID
	// Mantener el código original
	alq.Codigo = old.Codigo
	return s.repo.Actualizar(ctx, alq)
}

func (s *AlquilerService) Eliminar(ctx context.Context, id int, empresaID int) error {
	_, err := s.Obtener(ctx, id, empresaID)
	if err != nil {
		return err
	}
	return s.repo.Eliminar(ctx, id)
}

func (s *AlquilerService) Terminar(ctx context.Context, id int, empresaID int) error {
	alq, err := s.Obtener(ctx, id, empresaID)
	if err != nil {
		return err
	}
	alq.Estado = "finalizado"
	alq.ActivoParaCobro = false
	// Al actualizar a finalizado, el repo también debería liberar la unidad o podemos hacerlo explícito.
	// En mi implementación de repo.Eliminar lo hace, pero aquí es una actualización de estado.
	// Ajustaré el repo.Actualizar o lo haré aquí.
	_, err = s.repo.Actualizar(ctx, alq)
	return err
}

// --- Plantillas ---

func (s *AlquilerService) ListarPlantillas(ctx context.Context, empresaID int) ([]*domain.PlantillaContrato, error) {
	return s.repo.ListarPlantillas(ctx, empresaID)
}

func (s *AlquilerService) ObtenerPlantilla(ctx context.Context, id int, empresaID int) (*domain.PlantillaContrato, error) {
	return s.repo.ObtenerPlantilla(ctx, id, empresaID)
}

func (s *AlquilerService) GuardarPlantilla(ctx context.Context, p *domain.PlantillaContrato) (*domain.PlantillaContrato, error) {
	if p.ID > 0 {
		return s.repo.ActualizarPlantilla(ctx, p)
	}
	return s.repo.CrearPlantilla(ctx, p)
}

func (s *AlquilerService) EliminarPlantilla(ctx context.Context, id int, empresaID int) error {
	return s.repo.EliminarPlantilla(ctx, id, empresaID)
}

// --- Generación de Documentos ---

func (s *AlquilerService) GenerarContrato(ctx context.Context, id int, empresaID int, plantillaID int) (string, error) {
	alq, err := s.Obtener(ctx, id, empresaID)
	if err != nil {
		return "", err
	}

	cliente, err := s.cliente.BuscarPorID(ctx, alq.ClienteID)
	if err != nil {
		return "", fmt.Errorf("no se pudo cargar datos del cliente")
	}

	var contenido string
	if plantillaID > 0 {
		plantilla, err := s.ObtenerPlantilla(ctx, plantillaID, empresaID)
		if err != nil {
			return "", fmt.Errorf("plantilla no encontrada")
		}
		contenido = plantilla.Contenido
	} else {
		contenido = `
# {{titulo}}
 
Conste por el presente documento el contrato de arrendamiento que celebran de una parte **{{empresa_nombre}}** con RUC N.º **{{empresa_ruc}}**, con domicilio en **{{empresa_direccion}}**, a quien en lo sucesivo se denominará **EL ARRENDADOR**; y de otra parte **{{cliente_nombre}} {{cliente_apellidos}}** con DNI N.º **{{cliente_documento}}**, con domicilio en **{{cliente_direccion}}**, a quien en lo sucesivo se denominará **EL ARRENDATARIO**; en los términos contenidos en las cláusulas siguientes:
 
### ANTECEDENTES:
**PRIMERA.**- EL ARRENDADOR es propietario del inmueble ubicado en **{{unidad_codigo}}**, el cual se encuentra en buen estado de funcionamiento y conservación.
**SEGUNDA.**- EL ARRENDADOR deja constancia que el bien se encuentra en buen estado de funcionamiento y conservación.
 
### OBJETO DEL CONTRATO:
**TERCERA.**- Por el presente contrato EL ARRENDADOR se obliga a ceder el uso del bien descrito en la cláusula primera en favor de EL ARRENDATARIO, a título de arrendamiento.
 
### RENTA: FORMA Y OPORTUNIDAD DE PAGO:
**CUARTA.**- Las partes acuerdan que el monto de la renta que pagará EL ARRENDATARIO asciende a la suma de **{{moneda}} {{monto_renta}}** (**{{monto_renta_letras}}**) mensuales.
**QUINTA.**- La forma de pago de la renta será por mensualidades adelantadas que EL ARRENDATARIO pagará el día **{{dia_vencimiento}}** de cada mes.
 
### PLAZO DEL CONTRATO:
**SEXTA.**- Las partes convienen fijar un plazo de duración para el presente contrato.
- **Fecha de Inicio:** {{fecha_inicio}}
- **Fecha de Finalización:** {{fecha_fin}}

**SÉTIMA.**- Cualquiera de las partes podrá dar por concluido el presente contrato cursando a la otra una comunicación por vía notarial con no menos de 30 días de anticipación.
 
### OBLIGACIONES DE LAS PARTES:
**OCTAVA.**- EL ARRENDADOR se obliga a entregar el bien objeto de la prestación a su cargo en la fecha de suscripción de este documento.
**NOVENA.**- EL ARRENDATARIO se obliga a pagar puntualmente el monto de la renta.
**DÉCIMA.**- EL ARRENDATARIO se obliga a emplear el bien arrendado única y exclusivamente para el uso a que está destinado: Casa - habitación.
 
### CLÁUSULA RESOLUTORIA EXPRESA:
**DÉCIMO QUINTA.**- El incumplimiento de las obligaciones constituirá causal de resolución del presente contrato, al amparo del artículo 1430 del Código Civil.
 
### CLÁUSULA DE GARANTÍA:
**DÉCIMO SEXTA.**- EL ARRENDATARIO entrega a EL ARRENDADOR la suma de **{{moneda}} {{monto_deposito}}** (**{{monto_deposito_letras}}**) en calidad de depósito.
 
### COMPETENCIA TERRITORIAL:
**DÉCIMO NOVENA.**- Las partes se someten a la competencia territorial de los jueces y tribunales de la ciudad de suscripción.
 
En señal de conformidad las partes suscriben este documento el día {{fecha_inicio}}.
 
<br><br><br>
__________________________&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;__________________________
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;**EL ARRENDADOR**&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;**EL ARRENDATARIO**
`
	}

	clienteApellidos := ""
	if cliente.Apellidos != nil {
		clienteApellidos = *cliente.Apellidos
	}
	clienteCorreo := ""
	if cliente.Correo != nil {
		clienteCorreo = *cliente.Correo
	}
	clienteDireccion := "...................."
	if cliente.Direccion != nil && *cliente.Direccion != "" {
		clienteDireccion = *cliente.Direccion
	}

	emp, _ := s.empresa.BuscarPorID(ctx, empresaID)
	empNombre := "LA MERCANTIL E.I.R.L."
	empRUC := "20615683346"
	empDireccion := "Lambayeque, Chiclayo, PERÚ"
	// Ignoramos el nombre de la DB para usar "La Mercantil" directo como pidió el usuario
	if emp != nil {
		if emp.RUC != nil && *emp.RUC != "" && *emp.RUC != "20615683346" {
			empRUC = *emp.RUC
		}
		if emp.Direccion != nil && *emp.Direccion != "" && *emp.Direccion != "Lambayeque, Chiclayo, PERÚ" {
			empDireccion = *emp.Direccion
		}
	}

	fechaFinStr := "Indeterminado"
	if alq.FechaFin != nil {
		fechaFinStr = alq.FechaFin.Format("2006-01-02")
	}
	observacionesStr := "Ninguna"
	if alq.Observaciones != nil && *alq.Observaciones != "" {
		observacionesStr = *alq.Observaciones
	}

	// Reemplazar placeholders
	replacer := strings.NewReplacer(
		"{{empresa_nombre}}", empNombre,
		"{{empresa_ruc}}", empRUC,
		"{{empresa_direccion}}", empDireccion,
		"{{cliente_nombre}}", alq.ClienteNombre,
		"{{cliente_apellidos}}", clienteApellidos,
		"{{cliente_correo}}", clienteCorreo,
		"{{cliente_direccion}}", clienteDireccion,
		"{{cliente_documento}}", cliente.DocumentoNumero,
		"{{unidad_codigo}}", alq.UnidadCodigo,
		"{{monto_renta}}", fmt.Sprintf("%.2f", alq.MontoRenta),
		"{{monto_renta_letras}}", s.NumeroALetras(alq.MontoRenta, alq.Moneda),
		"{{monto_deposito}}", fmt.Sprintf("%.2f", alq.MontoDeposito),
		"{{monto_deposito_letras}}", s.NumeroALetras(alq.MontoDeposito, alq.Moneda),
		"{{moneda}}", alq.Moneda,
		"{{fecha_inicio}}", alq.FechaInicio.Format("2006-01-02"),
		"{{fecha_fin}}", fechaFinStr,
		"{{dia_vencimiento}}", fmt.Sprintf("%d", alq.DiaVencimiento),
		"{{observaciones}}", observacionesStr,
		"{{titulo}}", fmt.Sprintf("CONTRATO DE ARRENDAMIENTO (ALQUILER) DE BIEN INMUEBLE A %s", strings.ToUpper(func() string {
			if alq.FechaFin != nil {
				return "PLAZO DETERMINADO"
			}
			return "PLAZO INDETERMINADO"
		}())),
	)

	return replacer.Replace(contenido), nil
}

func (s *AlquilerService) GenerarContratoWord(ctx context.Context, id int, empresaID int, plantillaID int) ([]byte, error) {
	textoMarkdown, err := s.GenerarContrato(ctx, id, empresaID, plantillaID)
	if err != nil {
		return nil, err
	}

	// Convertir Markdown a HTML optimizado para Word
	// 1. Primero convertir títulos de sección (###) - el CSS ya aplica bold, no hace falta <b> extra
	htmlBody := strings.ReplaceAll(textoMarkdown, "### ", "<h3>")

	// 2. Título principal (#) - centrado y en mayúsculas por CSS
	htmlBody = strings.ReplaceAll(htmlBody, "# ", "<h1>")

	// 3. Convertir **texto** correctamente a <b>texto</b>
	//    Alternamos entre abrir <b> y cerrar </b> para cada par de **
	htmlBody = convertirNegritas(htmlBody)

	// 4. Convertir saltos de línea en párrafos reales para que Word respete el "Justify"
	parrafos := strings.Split(htmlBody, "\n")
	var bodyFinal string
	for _, p := range parrafos {
		trimP := strings.TrimSpace(p)
		if trimP == "" {
			bodyFinal += "<p class='MsoNormal'>&nbsp;</p>"
			continue
		}
		if strings.HasPrefix(trimP, "<h1>") {
			bodyFinal += trimP + "</h1>"
		} else if strings.HasPrefix(trimP, "<h3>") {
			bodyFinal += trimP + "</h3>"
		} else {
			bodyFinal += fmt.Sprintf("<p class='MsoNormal'>%s</p>", trimP)
		}
	}

	// Estructura XML/HTML de alta fidelidad para Microsoft Word
	wordXml := fmt.Sprintf(`
<html xmlns:o='urn:schemas-microsoft-com:office:office' 
      xmlns:w='urn:schemas-microsoft-com:office:word' 
      xmlns='http://www.w3.org/TR/REC-html40'>
<head>
<meta charset='utf-8'>
<!--[if gte mso 9]>
<xml>
 <w:WordDocument>
  <w:View>Print</w:View>
  <w:Zoom>100</w:Zoom>
  <w:DoNotOptimizeForBrowser/>
 </w:WordDocument>
</xml>
<![endif]-->
<style>
 <!--
  @page {
    size: 8.5in 11in;
    margin: 1.25in 1.0in 1.25in 1.0in;
    mso-header-margin: .5in;
    mso-footer-margin: .5in;
    mso-paper-source: 0;
  }
  body {
    font-family: "Times New Roman", serif;
    mso-ascii-font-family: "Times New Roman";
    mso-hansi-font-family: "Times New Roman";
  }
  p.MsoNormal, li.MsoNormal, div.MsoNormal {
    mso-style-unhide: no;
    mso-style-qformat: yes;
    mso-style-parent: "";
    margin: 0in;
    margin-bottom: .0001pt;
    mso-pagination: widow-orphan;
    font-size: 12.0pt;
    font-family: "Times New Roman", serif;
    mso-fareast-font-family: "Times New Roman";
    text-align: justify;
    line-height: 150%%;
  }
  h1 {
    mso-style-label: "Título 1";
    font-size: 16.0pt;
    font-family: "Times New Roman", serif;
    font-weight: bold;
    text-align: center;
    text-transform: uppercase;
    margin-bottom: 12pt;
  }
  h3 {
    mso-style-label: "Título 3";
    font-size: 12.0pt;
    font-family: "Times New Roman", serif;
    font-weight: bold;
    margin-top: 12pt;
    margin-bottom: 6pt;
    text-decoration: underline;
  }
  b { font-weight: bold; }
 -->
</style>
</head>
<body lang=ES-PE style='tab-interval:.5in'>
 <div class=Section1>
  %s
 </div>
</body>
</html>`, bodyFinal)

	return []byte(wordXml), nil
}

func (s *AlquilerService) NumeroALetras(n float64, moneda string) string {
	entero := int64(n)
	centavos := int64((n - float64(entero)) * 100)

	unidades := []string{"", "UN", "DOS", "TRES", "CUATRO", "CINCO", "SEIS", "SIETE", "OCHO", "NUEVE"}
	decenas := []string{"", "DIEZ", "VEINTE", "TREINTA", "CUARENTA", "CINCUENTA", "SESENTA", "SETENTA", "OCHENTA", "NOVENTA"}
	especiales := []string{"DIEZ", "ONCE", "DOCE", "TRECE", "CATORCE", "QUINCE", "DIECISEIS", "DIECISIETE", "DIECIOCHO", "DIECINUEVE"}
	centenas := []string{"", "CIENTO", "DOSCIENTOS", "TRESCIENTOS", "CUATROCIENTOS", "QUINIENTOS", "SEISCIENTOS", "SETECIENTOS", "OCHOCIENTOS", "NOVECIENTOS"}

	var convertir func(num int64) string
	convertir = func(num int64) string {
		if num == 0 {
			return ""
		}
		if num < 10 {
			return unidades[num]
		}
		if num < 20 {
			return especiales[num-10]
		}
		if num < 100 {
			d := num / 10
			u := num % 10
			if u == 0 {
				return decenas[d]
			}
			return decenas[d] + " Y " + unidades[u]
		}
		if num == 100 {
			return "CIEN"
		}
		c := num / 100
		resto := num % 100
		if c == 1 && resto > 0 {
			return "CIENTO " + convertir(resto)
		}
		return centenas[c] + " " + convertir(resto)
	}

	res := ""
	if entero >= 1000 {
		m := entero / 1000
		resto := entero % 1000
		if m == 1 {
			res = "MIL "
		} else {
			res = convertir(m) + " MIL "
		}
		res += convertir(resto)
	} else {
		res = convertir(entero)
	}

	if res == "" {
		res = "CERO"
	}

	monedaNombre := "SOLES"
	if moneda == "USD" {
		monedaNombre = "DOLARES AMERICANOS"
	}

	return fmt.Sprintf("%s Y %02d/100 %s", strings.TrimSpace(res), centavos, monedaNombre)
}

// convertirNegritas convierte **texto** a <b>texto</b> correctamente,
// alternando entre etiqueta de apertura y cierre para cada par de asteriscos.
func convertirNegritas(s string) string {
	partes := strings.Split(s, "**")
	var resultado strings.Builder
	for i, parte := range partes {
		resultado.WriteString(parte)
		if i < len(partes)-1 {
			if i%2 == 0 {
				resultado.WriteString("<b>")
			} else {
				resultado.WriteString("</b>")
			}
		}
	}
	return resultado.String()
}

func (s *AlquilerService) GenerarContratoBorrador(ctx context.Context, empresaID int, req domain.GenerarBorradorRequest) ([]byte, error) {
	// Intentar autocompletar con datos de base de datos si existe el cliente
	if req.ClienteDocumento != "" {
		// Buscamos si hay un cliente con ese documento para esta empresa
		filtros := domain.ClienteFiltros{
			EmpresaID: empresaID,
			Busqueda:  req.ClienteDocumento,
			Pagina:    1,
			Limite:    1,
		}
		clientes, _, _ := s.cliente.ListarPaginado(ctx, filtros)
		if len(clientes) > 0 {
			// Usar datos del cliente encontrado
			c := clientes[0]
			req.ClienteNombre = c.Nombres
			if c.Apellidos != nil {
				req.ClienteApellidos = *c.Apellidos
			}
			if c.Correo != nil {
				req.ClienteCorreo = *c.Correo
			}
			if c.Direccion != nil {
				req.ClienteDireccion = *c.Direccion
			}
		}
	}

	var contenido string
	if req.PlantillaID > 0 {
		plantilla, err := s.ObtenerPlantilla(ctx, req.PlantillaID, empresaID)
		if err != nil {
			return nil, fmt.Errorf("plantilla no encontrada")
		}
		contenido = plantilla.Contenido
	} else {
		// Usar la misma plantilla formal por defecto que nos dio el usuario
		contenido = `
# {{titulo}}

Conste por el presente documento el contrato de arrendamiento que celebran de una parte **{{empresa_nombre}}** con RUC N.º **{{empresa_ruc}}**, con domicilio en **{{empresa_direccion}}**, a quien en lo sucesivo se denominará **EL ARRENDADOR**; y de otra parte **{{cliente_nombre}} {{cliente_apellidos}}** con DNI N.º **{{cliente_documento}}**, con domicilio en **{{cliente_direccion}}**, a quien en lo sucesivo se denominará **EL ARRENDATARIO**; en los términos contenidos en las cláusulas siguientes:

### ANTECEDENTES:
**PRIMERA.**- EL ARRENDADOR es propietario del inmueble ubicado en **{{unidad_codigo}}**, el cual se encuentra en buen estado de funcionamiento y conservación.
**SEGUNDA.**- EL ARRENDADOR deja constancia que el bien se encuentra en buen estado de funcionamiento y conservación.

### OBJETO DEL CONTRATO:
**TERCERA.**- Por el presente contrato EL ARRENDADOR se obliga a ceder el uso del bien descrito en la cláusula primera en favor de EL ARRENDATARIO, a título de arrendamiento.

### RENTA: FORMA Y OPORTUNIDAD DE PAGO:
**CUARTA.**- Las partes acuerdan que el monto de la renta que pagará EL ARRENDATARIO asciende a la suma de **{{moneda}} {{monto_renta}}** (**{{monto_renta_letras}}**) mensuales.
**QUINTA.**- La forma de pago de la renta será por mensualidades adelantadas que EL ARRENDATARIO pagará el día **{{dia_vencimiento}}** de cada mes.

### CLÁUSULA CUARTA: PLAZO DEL CONTRATO
Las partes convienen fijar un plazo de duración para el presente contrato.
- **Fecha de Inicio:** {{fecha_inicio}}
- **Fecha de Finalización:** {{fecha_fin}}

Cualquiera de las partes podrá dar por concluido el presente contrato cursando a la otra una comunicación por vía notarial con no menos de 30 días de anticipación.

### OBLIGACIONES DE LAS PARTES:
**OCTAVA.**- EL ARRENDADOR se obliga a entregar el bien objeto de la prestación a su cargo en la fecha de suscripción de este documento.
**NOVENA.**- EL ARRENDATARIO se obliga a pagar puntualmente el monto de la renta.
**DÉCIMA.**- EL ARRENDATARIO se obliga a emplear el bien arrendado única y exclusivamente para el uso a que está destinado: Casa - habitación.

### CLÁUSULA RESOLUTORIA EXPRESA:
**DÉCIMO QUINTA.**- El incumplimiento de las obligaciones constituirá causal de resolución del presente contrato, al amparo del artículo 1430 del Código Civil.

### CLÁUSULA DE GARANTÍA:
**DÉCIMO SEXTA.**- EL ARRENDATARIO entrega a EL ARRENDADOR la suma de **{{moneda}} {{monto_deposito}}** (**{{monto_deposito_letras}}**) en calidad de depósito.

### COMPETENCIA TERRITORIAL:
**DÉCIMO NOVENA.**- Las partes se someten a la competencia territorial de los jueces y tribunales de la ciudad de suscripción.

En señal de conformidad las partes suscriben este documento el día ` + time.Now().Format("2006-01-02") + `.

<br><br><br>
__________________________&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;__________________________
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;**EL ARRENDADOR**&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;**EL ARRENDATARIO**
`
	}

	emp, _ := s.empresa.BuscarPorID(ctx, empresaID)
	empNombre := "LA MERCANTIL E.I.R.L."
	empRUC := "20615683346"
	empDireccion := "Lambayeque, Chiclayo, PERÚ"
	// Ignoramos el nombre de la DB para usar "La Mercantil" directo como pidió el usuario
	if emp != nil {
		if emp.RUC != nil && *emp.RUC != "" && *emp.RUC != "20615683346" {
			empRUC = *emp.RUC
		}
		if emp.Direccion != nil && *emp.Direccion != "" && *emp.Direccion != "Lambayeque, Chiclayo, PERÚ" {
			empDireccion = *emp.Direccion
		}
	}

	// Reemplazar placeholders
	replacer := strings.NewReplacer(
		"{{empresa_nombre}}", empNombre,
		"{{empresa_ruc}}", empRUC,
		"{{empresa_direccion}}", empDireccion,
		"{{cliente_nombre}}", req.ClienteNombre,
		"{{cliente_apellidos}}", req.ClienteApellidos,
		"{{cliente_correo}}", req.ClienteCorreo,
		"{{cliente_direccion}}", req.ClienteDireccion,
		"{{cliente_documento}}", req.ClienteDocumento,
		"{{unidad_codigo}}", req.UnidadCodigo,
		"{{monto_renta}}", fmt.Sprintf("%.2f", req.MontoRenta),
		"{{monto_renta_letras}}", s.NumeroALetras(req.MontoRenta, "PEN"),
		"{{monto_deposito}}", fmt.Sprintf("%.2f", req.MontoDeposito),
		"{{monto_deposito_letras}}", s.NumeroALetras(req.MontoDeposito, "PEN"),
		"{{moneda}}", "S/.",
		"{{fecha_inicio}}", req.FechaInicio,
		"{{fecha_fin}}", req.FechaFin,
		"{{dia_vencimiento}}", fmt.Sprintf("%d", req.DiaVencimiento),
		"{{titulo}}", fmt.Sprintf("CONTRATO DE ARRENDAMIENTO (ALQUILER) DE BIEN INMUEBLE A %s", strings.ToUpper(func() string {
			if req.FechaFin != "" && req.FechaFin != "Indeterminado" {
				return "PLAZO DETERMINADO"
			}
			return "PLAZO INDETERMINADO"
		}())),
	)

	textoFinal := replacer.Replace(contenido)

	// Convertir a HTML optimizado para Word
	// 1. Títulos de sección (###) - el CSS ya aplica bold
	htmlBody := strings.ReplaceAll(textoFinal, "### ", "<h3>")
	// 2. Título principal (#)
	htmlBody = strings.ReplaceAll(htmlBody, "# ", "<h1>")
	// 3. Convertir **texto** correctamente a <b>texto</b>
	htmlBody = convertirNegritas(htmlBody)

	parrafos := strings.Split(htmlBody, "\n")
	var bodyFinal string
	for _, p := range parrafos {
		trimP := strings.TrimSpace(p)
		if trimP == "" {
			bodyFinal += "<p class='MsoNormal'>&nbsp;</p>"
			continue
		}
		if strings.HasPrefix(trimP, "<h1>") {
			bodyFinal += trimP + "</h1>"
		} else if strings.HasPrefix(trimP, "<h3>") {
			bodyFinal += trimP + "</h3>"
		} else {
			bodyFinal += fmt.Sprintf("<p class='MsoNormal'>%s</p>", trimP)
		}
	}

	wordXml := fmt.Sprintf(`
<html xmlns:o='urn:schemas-microsoft-com:office:office' 
      xmlns:w='urn:schemas-microsoft-com:office:word' 
      xmlns='http://www.w3.org/TR/REC-html40'>
<head>
<meta charset='utf-8'>
<!--[if gte mso 9]>
<xml>
 <w:WordDocument>
  <w:View>Print</w:View>
  <w:Zoom>100</w:Zoom>
  <w:DoNotOptimizeForBrowser/>
 </w:WordDocument>
</xml>
<![endif]-->
<style>
 <!--
  @page {
    size: 8.5in 11in;
    margin: 1.25in 1.0in 1.25in 1.0in;
    mso-header-margin: .5in;
    mso-footer-margin: .5in;
    mso-paper-source: 0;
  }
  body {
    font-family: "Times New Roman", serif;
    mso-ascii-font-family: "Times New Roman";
    mso-hansi-font-family: "Times New Roman";
  }
  p.MsoNormal, li.MsoNormal, div.MsoNormal {
    mso-style-unhide: no;
    mso-style-qformat: yes;
    mso-style-parent: "";
    margin: 0in;
    margin-bottom: .0001pt;
    mso-pagination: widow-orphan;
    font-size: 12.0pt;
    font-family: "Times New Roman", serif;
    mso-fareast-font-family: "Times New Roman";
    text-align: justify;
    line-height: 150%%;
  }
  h1 {
    mso-style-label: "Título 1";
    font-size: 16.0pt;
    font-family: "Times New Roman", serif;
    font-weight: bold;
    text-align: center;
    text-transform: uppercase;
    margin-bottom: 12pt;
  }
  h3 {
    mso-style-label: "Título 3";
    font-size: 12.0pt;
    font-family: "Times New Roman", serif;
    font-weight: bold;
    margin-top: 12pt;
    margin-bottom: 6pt;
    text-decoration: underline;
  }
  b { font-weight: bold; }
 -->
</style>
</head>
<body lang=ES-PE style='tab-interval:.5in'>
 <div class=Section1>
  %s
 </div>
</body>
</html>`, bodyFinal)

	return []byte(wordXml), nil
}


type PagoAlquilerService struct {
	repo domain.PagoAlquilerRepository
}

func NewPagoAlquilerService(repo domain.PagoAlquilerRepository) *PagoAlquilerService {
	return &PagoAlquilerService{repo: repo}
}

func (s *PagoAlquilerService) Registrar(ctx context.Context, pago *domain.RegistroPagoAlquiler) (*domain.PagoAlquiler, error) {
	if pago.MesCorrespondiente < 1 || pago.MesCorrespondiente > 12 {
		return nil, fmt.Errorf("mes_correspondiente debe estar entre 1 y 12")
	}
	if pago.FechaPago.IsZero() {
		return nil, fmt.Errorf("fecha_pago es obligatoria")
	}
	return s.repo.Registrar(ctx, pago)
}

func (s *PagoAlquilerService) ListarPendientesMesActual(ctx context.Context, empresaID int) ([]*domain.PagoPendiente, error) {
	return s.repo.ListarPendientesMesActual(ctx, empresaID, time.Now().UTC())
}

func (s *PagoAlquilerService) ListarHistorial(ctx context.Context, filtros domain.PagoFiltros) ([]*domain.PagoAlquiler, int, error) {
	if filtros.Pagina <= 0 {
		filtros.Pagina = 1
	}
	if filtros.Limite <= 0 {
		filtros.Limite = 10
	}
	return s.repo.Listar(ctx, filtros)
}

func (s *PagoAlquilerService) Obtener(ctx context.Context, id int, empresaID int) (*domain.PagoAlquiler, error) {
	return s.repo.BuscarPorID(ctx, id, empresaID)
}

func (s *PagoAlquilerService) Anular(ctx context.Context, id int, empresaID int) error {
	return s.repo.Eliminar(ctx, id, empresaID)
}

func (s *PagoAlquilerService) Actualizar(ctx context.Context, id int, empresaID int, notas *string, metodoPago string) (*domain.PagoAlquiler, error) {
	pago, err := s.Obtener(ctx, id, empresaID)
	if err != nil {
		return nil, err
	}
	if notas != nil {
		pago.Nota = notas
	}
	if metodoPago != "" {
		pago.MetodoPago = metodoPago
	}
	return s.repo.Actualizar(ctx, pago)
}
