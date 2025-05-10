# Exporter de Ports TCP per a Prometheus

Aquest és un exporter de Prometheus que monitoritza l'estat dels ports TCP en la instància on s'executa.

## Funcionament

L'exporter funciona intentant obrir una connexió TCP a cada port possible (de l'1 al 65535). Si no es pot establir la connexió (perquè un altre servei ja està escoltant en aquest port), es considera que el port està **obert**. Si la connexió s'estableix amb èxit i es tanca immediatament, es considera que el port està **tancat**.

L'estat de cada port obert es registra com una mèrica de tipus `Gauge` anomenada `instance_open_ports`. Aquesta mèrica té una etiqueta `port` que indica el número del port. El valor de la mèrica serà `1` si el port està obert. Els ports tancats no generen una mèrica activa per defecte.

L'exporter exposa les mètriques al punt final `/` a través d'HTTP, escoltant al port `8000` per defecte.

## Ús

1. **Clonar aquest repositori (si escau).**
2. **Construir l'exporter:**

   ```bash
   go build -o port_exporter main.go
   ```

3. **Executar l'exporter:**

   ```bash
   ./port_exporter
   ```

   L'exporter començarà a servir les mètriques a `http://localhost:8000/`.

4. **Configurar Prometheus:**

   Afegeix la següent configuració al fitxer de configuració de Prometheus (`prometheus.yml`) per començar a recollir les mètriques:

   ```yaml
   scrape_configs:
     - job_name: 'port_exporter'
       static_configs:
         - targets: ['localhost:8000']
   ```

## Consideracions sobre el Rendiment

Aquest exporter, per defecte, comprova **tots** els ports TCP (de l'1 al 65535) en cada cicle de recollida de mètriques. Això pot generar una càrrega significativa en el sistema. Per a un ús més eficient, es recomana **modificar la funció `collectOpenPorts` per comprovar només un conjunt específic de ports** que siguin rellevants per a la vostra monitorització.

## Dependències

Aquest projecte utilitza la **llibreria client oficial de Prometheus per a Go**, que permet crear i exposar mètriques de manera senzilla.

Les principals dependències són:

- [`github.com/prometheus/client_golang/prometheus`](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus): Per definir i registrar mètriques.
- [`github.com/prometheus/client_golang/prometheus/promhttp`](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus/promhttp): Per exposar les mètriques a través d'HTTP.

### Instal·lació de dependències

Des de la carpeta del projecte, simplement executa:

```bash
go mod tidy
```

Aquest comandament crearà automàticament el fitxer `go.mod` (si no existeix) i descarregarà totes les dependències necessàries.
