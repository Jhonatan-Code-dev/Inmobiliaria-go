# Manual de Integración de Reportes Gerenciales en Angular (ng-apexcharts)

Este documento detalla cómo consumir los 5 nuevos endpoints de reportes gerenciales en una aplicación frontend desarrollada en **Angular** (compatible con versiones modernas v15+ y Componentes Standalone). 

Para lograr gráficos interactivos de alto nivel, con animaciones fluidas, diseño premium y adaptabilidad responsive, utilizaremos **ng-apexcharts** (la integración oficial para Angular de **ApexCharts**).

---

## 📦 Instalación y Configuración en Angular

Para comenzar, instala ApexCharts y su adaptador para Angular en tu proyecto:

```bash
npm install apexcharts ng-apexcharts --save
```

### 1. Registrar en el Componente o Módulo
Si utilizas **Componentes Standalone**, impórtalo directamente en el array de `imports` de tu componente:

```typescript
import { Component } from '@angular/core';
import { NgApexchartsModule } from 'ng-apexcharts';

@Component({
  selector: 'app-reportes-dashboard',
  standalone: true,
  imports: [NgApexchartsModule],
  templateUrl: './reportes-dashboard.component.html',
  styleUrls: ['./reportes-dashboard.component.css']
})
export class ReportesDashboardComponent {}
```

---

## 🔒 Seguridad y Autenticación (API Backend)
* **Ruta Base**: `/api/user/reportes`
* **Parámetros de Consulta**:
  - `empresa_id` (Query, Obligatorio): ID del tenant del usuario.
  - `desde` (Query, Opcional, `YYYY-MM-DD`): Fecha inicio.
  - `hasta` (Query, Opcional, `YYYY-MM-DD`): Fecha fin.

---

## 📡 1. Creación del Servicio en Angular (`ReportesService`)

Implementaremos un servicio en Angular que consuma los 5 endpoints gerenciales de manera eficiente. Utilizaremos **RxJS `forkJoin`** para disparar todas las peticiones concurrentemente con una única barra de carga o spinner.

```typescript
import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable, forkJoin } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ReportesService {
  private baseUrl = '/api/user/reportes';

  constructor(private http: HttpClient) {}

  private getParams(empresaId: number, desde?: string, hasta?: string): HttpParams {
    let params = new HttpParams().set('empresa_id', empresaId.toString());
    if (desde) params = params.set('desde', desde);
    if (hasta) params = params.set('hasta', hasta);
    return params;
  }

  // 1. Ingresos vs Gastos
  getIngresosGastos(empresaId: number, desde?: string, hasta?: string): Observable<any> {
    return this.http.get(`${this.baseUrl}/ingresos-gastos`, { params: this.getParams(empresaId, desde, hasta) });
  }

  // 2. Distribución de Métodos de Pago
  getMetodosPago(empresaId: number, desde?: string, hasta?: string): Observable<any[]> {
    return this.http.get<any[]>(`${this.baseUrl}/metodos-pago`, { params: this.getParams(empresaId, desde, hasta) });
  }

  // 3. Distribución de Categorías de Gastos
  getCategoriasGastos(empresaId: number, desde?: string, hasta?: string): Observable<any[]> {
    return this.http.get<any[]>(`${this.baseUrl}/categorias-gastos`, { params: this.getParams(empresaId, desde, hasta) });
  }

  // 4. Rentabilidad de Propiedades
  getRentabilidadPropiedades(empresaId: number, desde?: string, hasta?: string): Observable<any[]> {
    return this.http.get<any[]>(`${this.baseUrl}/rentabilidad-propiedades`, { params: this.getParams(empresaId, desde, hasta) });
  }

  // 5. Soporte y Tickets de Mantenimiento
  getTicketsMantenimiento(empresaId: number, desde?: string, hasta?: string): Observable<any> {
    return this.http.get(`${this.baseUrl}/tickets-mantenimiento`, { params: this.getParams(empresaId, desde, hasta) });
  }

  // Carga Masiva y Paralela (Dashboard Inicial)
  cargarDashboardCompleto(empresaId: number, desde?: string, hasta?: string): Observable<any> {
    return forkJoin({
      financiero: this.getIngresosGastos(empresaId, desde, hasta),
      metodosPago: this.getMetodosPago(empresaId, desde, hasta),
      categoriasGastos: this.getCategoriasGastos(empresaId, desde, hasta),
      rentabilidad: this.getRentabilidadPropiedades(empresaId, desde, hasta),
      tickets: this.getTicketsMantenimiento(empresaId, desde, hasta)
    });
  }
}
```

---

## 🎨 2. Diseños de Gráficos Gerenciales de Alto Nivel

### A. Tendencia Mensual de Ingresos vs. Gastos (Gráfico Combinado Line & Column)
Muestra barras verticales para ingresos/gastos y una línea para el balance neto. Aporta gran visibilidad a la gerencia sobre los flujos de caja mensuales.

* **Recomendación Estética**:
  - Ingresos: Color Esmeralda (`#10B981`)
  - Gastos: Color Coral/Rojo (`#EF4444`)
  - Balance Neto: Línea Azul Eléctrico (`#3B82F6`) con curvas suaves (spline).

#### Configuración del Componente Angular:
```typescript
import { Component, OnInit, ViewChild } from '@angular/core';
import { ReportesService } from './reportes.service';
import {
  ChartComponent,
  ApexAxisChartSeries,
  ApexChart,
  ApexXAxis,
  ApexDataLabels,
  ApexStroke,
  ApexYAxis,
  ApexTitleSubtitle,
  ApexLegend,
  ApexFill,
  ApexMarkers
} from 'ng-apexcharts';

export type ChartOptions = {
  series: ApexAxisChartSeries;
  chart: ApexChart;
  xaxis: ApexXAxis;
  yaxis: ApexYAxis | ApexYAxis[];
  dataLabels: ApexDataLabels;
  stroke: ApexStroke;
  legend: ApexLegend;
  fill: ApexFill;
  colors: string[];
  markers: ApexMarkers;
};

@Component({
  selector: 'app-reporte-financiero',
  standalone: true,
  imports: [NgApexchartsModule],
  template: `
    <div class="card-dashboard">
      <h3>Tendencia Financiera Mensual</h3>
      <apx-chart
        [series]="chartOptions.series"
        [chart]="chartOptions.chart"
        [xaxis]="chartOptions.xaxis"
        [yaxis]="chartOptions.yaxis"
        [stroke]="chartOptions.stroke"
        [colors]="chartOptions.colors"
        [fill]="chartOptions.fill"
        [legend]="chartOptions.legend"
        [markers]="chartOptions.markers"
      ></apx-chart>
    </div>
  `
})
export class ReporteFinancieroComponent implements OnInit {
  @ViewChild('chart') chart!: ChartComponent;
  public chartOptions!: Partial<ChartOptions>;

  constructor(private service: ReportesService) {}

  ngOnInit() {
    this.service.getIngresosGastos(1).subscribe(res => {
      const meses = res.serie_mensual.map((s: any) => s.periodo);
      const ingresos = res.serie_mensual.map((s: any) => s.ingresos);
      const gastos = res.serie_mensual.map((s: any) => s.gastos);
      const balance = res.serie_mensual.map((s: any) => s.balance);

      this.chartOptions = {
        series: [
          { name: 'Ingresos', type: 'column', data: ingresos },
          { name: 'Gastos', type: 'column', data: gastos },
          { name: 'Balance Neto', type: 'line', data: balance }
        ],
        chart: {
          height: 380,
          type: 'line',
          toolbar: { show: false },
          animations: { enabled: true, easing: 'easeinout', speed: 800 }
        },
        colors: ['#10B981', '#EF4444', '#3B82F6'],
        stroke: { width: [0, 0, 4], curve: 'smooth' },
        fill: { opacity: [0.85, 0.85, 1] },
        xaxis: { categories: meses },
        yaxis: [
          { title: { text: 'Monto (S/.)' } }
        ],
        markers: { size: 5 },
        legend: { position: 'top' }
      };
    });
  }
}
```

---

### B. Métodos de Pago y Categorías de Gasto (Doughnut Chart Premium)
Perfecto para gráficos circulares modernos con esquinas redondeadas y sombras suaves.

#### Configuración del Componente para Métodos de Pago:
```typescript
@Component({
  selector: 'app-reporte-pagos',
  standalone: true,
  imports: [NgApexchartsModule],
  template: `
    <div class="card-donut">
      <h3>Distribución por Métodos de Pago</h3>
      <apx-chart
        [series]="donutOptions.series"
        [chart]="donutOptions.chart"
        [labels]="donutOptions.labels"
        [colors]="donutOptions.colors"
        [legend]="donutOptions.legend"
        [responsive]="donutOptions.responsive"
      ></apx-chart>
    </div>
  `
})
export class ReportePagosComponent implements OnInit {
  public donutOptions: any = {};

  constructor(private service: ReportesService) {}

  ngOnInit() {
    this.service.getMetodosPago(1).subscribe(data => {
      const labels = data.map(d => d.metodo.toUpperCase());
      const series = data.map(d => d.total);

      this.donutOptions = {
        series: series,
        labels: labels,
        chart: {
          type: 'donut',
          height: 320,
        },
        colors: ['#3B82F6', '#10B981', '#F59E0B', '#8B5CF6', '#EC4899'],
        legend: { position: 'bottom' },
        responsive: [
          {
            breakpoint: 480,
            options: {
              chart: { width: 280 },
              legend: { position: 'bottom' }
            }
          }
        ]
      };
    });
  }
}
```

---

### C. Rentabilidad de Propiedades (Horizontal Bar Chart)
Compara rápidamente los ingresos versus los gastos prorrateados por cada propiedad.

#### Configuración del Gráfico:
- **Tipo**: `'bar'` con propiedad `plotOptions: { bar: { horizontal: true } }`.
- **Eje Y**: Nombres de las propiedades.
- **Series**:
  - Serie 1 (Ingresos): Verde (`#10B981`)
  - Serie 2 (Gastos Asignados): Naranja (`#F59E0B`)
  - Serie 3 (Utilidad Neta): Azul (`#3B82F6`)

```typescript
this.chartOptions = {
  series: [
    { name: 'Ingresos', data: [9500, 6100] },
    { name: 'Gastos Pro-rata', data: [2500, 1700] },
    { name: 'Rentabilidad', data: [7000, 4400] }
  ],
  chart: {
    type: 'bar',
    height: 350
  },
  plotOptions: {
    bar: {
      horizontal: true,
      dataLabels: { position: 'top' }
    }
  },
  colors: ['#10B981', '#F59E0B', '#3B82F6'],
  xaxis: {
    categories: ['Edificio Los Portales', 'Residencial Primavera']
  }
};
```

---

### D. Tickets de Mantenimiento (Stacked Column Chart & Priority Pie)
Compara la carga de trabajo técnica del equipo de soporte. El gráfico de barras apiladas permite analizar las prioridades en base a los estados.

#### Estructura de Datos para Cargar:
* **Entrada JSON**:
  ```json
  {
    "total_tickets": 25,
    "por_estado": { "abierto": 5, "en_progreso": 8, "resuelto": 11, "anulado": 1 },
    "por_prioridad": { "baja": 4, "media": 15, "alta": 6 }
  }
  ```

* **Mapeo a Gráfico Circular de Prioridades**:
  - `series`: `[4, 15, 6]`
  - `labels`: `['Baja', 'Media', 'Alta']`
  - `colors`: `['#10B981', '#F59E0B', '#EF4444']` (Semáforo de prioridad).

---

## 💎 Tips de Diseño y Estética Premium en Angular

1. **Dark Mode Integrado**:
   ApexCharts tiene soporte nativo para temas oscuros. Puedes cambiarlo dinámicamente inyectando una propiedad de configuración:
   ```typescript
   theme: {
     mode: this.isDarkMode ? 'dark' : 'light',
     palette: 'palette1'
   }
   ```
2. **Glassmorphism en Contenedores**:
   En tu CSS de Angular (`styles.css`), aplica este estilo para los contenedores de los gráficos gerenciales para lograr un diseño moderno y premium:
   ```css
   .card-dashboard {
     background: rgba(255, 255, 255, 0.45);
     backdrop-filter: blur(12px);
     -webkit-backdrop-filter: blur(12px);
     border-radius: 16px;
     border: 1px solid rgba(255, 255, 255, 0.25);
     padding: 24px;
     box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.08);
   }
   ```
3. **Optimización de Tooltips**:
   Personaliza el tooltip para que muestre el símbolo monetario de forma elegante:
   ```typescript
   tooltip: {
     y: {
       formatter: function(val) {
         return "S/. " + val.toFixed(2);
       }
     }
   }
   ```
