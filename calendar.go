package main

import (
	"fmt"
	"strconv"
)

func main() {
	MuestraMeses([]string{"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"})
}

type MesDía map[string]int

var MesesDias = MesDía{
	"Enero":      31,
	"Febrero":    28,
	"Marzo":      31,
	"Abril":      30,
	"Mayo":       31,
	"Junio":      30,
	"Julio":      31,
	"Agosto":     31,
	"Septiembre": 30,
	"Octubre":    31,
	"Noviembre":  30,
	"Diciembre":  31,
}

func MuestraMeses(queMeses []string) {
	for i := 0; i < len(queMeses); i++ {
		MuestraMes(queMeses[i])
	}
}

func MuestraMes(queMes string) {
	fmt.Println("     " + queMes)
	for i := 1; i <= MesesDias[queMes]; i++ {
		if i < 10 {
			fmt.Print(" ")
		}
		fmt.Print(strconv.Itoa(i) + " ")
		if i%7 == 0 {
			fmt.Println("")
		}
	}
	fmt.Println("")
}
