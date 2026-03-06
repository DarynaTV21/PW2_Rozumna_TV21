package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/calc", calcPage)
	http.HandleFunc("/calc/result", calcResult)

	fmt.Println("Сервер запущено на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "home.html", nil)
}

func calcPage(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "calc.html", nil)
}

func round(val float64, decimals int) float64 {
	format := "%." + strconv.Itoa(decimals) + "f"
	v, _ := strconv.ParseFloat(fmt.Sprintf(format, val), 64)
	return v
}

func calcResult(w http.ResponseWriter, r *http.Request) {

	masaVug, _ := strconv.ParseFloat(r.FormValue("masaVug"), 64)
	masaMaz, _ := strconv.ParseFloat(r.FormValue("masaMaz"), 64)
	masaGas, _ := strconv.ParseFloat(r.FormValue("masaGas"), 64)

	lowHeatVug := 20.47
	lowHeatMaz := 39.48
	lowHeatGas := 33.08

	masaContVug := 25.2
	masaContMaz := 0.15

	masaContSubstancesVug := 1.5
	masaContSubstancesMaz := 0.0

	efficiencyCleaning := 0.985

	fractionVug := 0.8
	fractionMaz := 1.0

	// Вугілля
	emvug := (math.Pow(10, 6) / lowHeatVug) *
		fractionVug *
		(masaContVug / (100 - masaContSubstancesVug)) *
		(1 - efficiencyCleaning)

	vukvug := math.Pow(10, -6) * emvug * lowHeatVug * masaVug

	// Мазут
	emmaz := (math.Pow(10, 6) / lowHeatMaz) *
		fractionMaz *
		(masaContMaz / (100 - masaContSubstancesMaz)) *
		(1 - efficiencyCleaning)

	vukmaz := math.Pow(10, -6) * emmaz * lowHeatMaz * masaMaz

	// Газ
	emgaz := 0.0
	vukgaz := math.Pow(10, -6) * emgaz * lowHeatGas * masaGas

	result := fmt.Sprintf(`
Показник емісії твердих частинок при спалюванні вугілля становитиме: 
%.3f г/ГДж
Валовий викид при спалюванні вугілля становитиме: 
%.3f т

Показник емісії твердих частинок при спалюванні мазуту становитиме: 
%.3f г/ГДж
Валовий викид при спалюванні мазуту становитиме: 
%.3f т

Показник емісії твердих частинок при спалюванні природного газу становитиме:
%.3f г/ГДж
Валовий викид при спалюванні природного газу становитиме: 
%.3f т
`,
		round(emvug, 3),
		round(vukvug, 3),
		round(emmaz, 3),
		round(vukmaz, 3),
		round(emgaz, 3),
		round(vukgaz, 3),
	)

	templates.ExecuteTemplate(w, "result.html", result)
}
