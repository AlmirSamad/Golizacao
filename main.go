package main

import (
	"fmt"
	"golizacao/sgs"
	"golizacao/sidra"
)

func main() {

	codigos := []string{"25411", "25412", "24363", "24364", "25392", "25394"}
	dataInicial := "01/01/2026"
	dataFinal := "31/12/2026"

	for i := 0; i < len(codigos); i++ {

		codigo := codigos[i]
		fmt.Printf("Código atual buscado: %s", codigo)
		bcb, err := sgs.BuscarDadosSGS(codigo, dataInicial, dataFinal)

		if err != nil {
			fmt.Printf(" Falha ao buscar código %s: %v\n", codigos[i], err)
		}

		fmt.Printf("Total de registros encontrados: %d\n", len(bcb))

		limite := 7
		if len(bcb) < limite {
			limite = len(bcb)
		}

		for i := 0; i < limite; i++ {
			fmt.Printf("Data: %s | Valor: %s\n", bcb[i].Data, bcb[i].Valor)
		}
		if len(bcb) > limite {
			fmt.Println("...")
		}

		nomeArquivo := fmt.Sprintf("serie_%s_2026.csv", codigos[i])

		err = sgs.SalvarEmCSV(nomeArquivo, bcb)
		if err != nil {
			fmt.Printf("Erro ao salvar CSV: %v\n", err)
		} else {
			fmt.Printf("-> Sucesso! Arquivo '%s' gerado com %d registros.\n\n", nomeArquivo, len(bcb))
		}
	}

	sidra.BaixarDadosSidra()

}
