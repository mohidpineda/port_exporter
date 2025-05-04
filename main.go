package portexporter

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus" // Per crear i registrar mètriques.
	"github.com/prometheus/client_golang/prometheus/promhttp" // Per exposar mètriques via HTTP.
)

var (
	// openPorts és una mèrica de tipus GaugeVec per indicar l'estat dels ports TCP.
	openPorts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "instance_open_ports",
			Help: "Indica si un port TCP està obert (1) o tancat (0) a la instància.",
		},
		[]string{"port"}, // La mèrica està etiquetada pel número de port.
	)
)

// checkPort intenta escoltar en un port TCP. Retorna true si el port està obert (no es pot escoltar).
func checkPort(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return true // El port està probablement en ús (obert).
	}
	ln.Close()
	return false // Hem pogut escoltar, el port està probablement tancat.
}

// collectOpenPorts itera per tots els ports TCP i actualitza la mèrica openPorts.
func collectOpenPorts() {
	for port := 1; port <= 65535; port++ {
		if checkPort(port) {
			openPorts.WithLabelValues(strconv.Itoa(port)).Set(1)
		}
		// Si el port està tancat, no establim explícitament el valor a 0 per defecte.
		// Prometheus mantindrà l'últim valor reportat per a aquest port.
	}
}

func main() {
	// Crea un nou registre de Prometheus.
	promRegistry := prometheus.NewRegistry()
	// Registra la nostra mèrica personalitzada.
	promRegistry.Register(openPorts)

	// Recull l'estat dels ports oberts abans de començar a servir les mètriques.
	collectOpenPorts()

	fmt.Println("Servidor de l'exporter: http://localhost:8000/")

	// Crea un handler per exposar les mètriques en el format d'Prometheus.
	handler := promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{})
	// Assigna el handler a la ruta arrel ("/").
	http.Handle("/", handler)
	// Inicia el servidor HTTP per servir les mètriques.
	http.ListenAndServe(":8000", nil)
}