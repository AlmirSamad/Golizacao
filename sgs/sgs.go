package sgs

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type SGSRecord struct {
	Data  string `json:"data"`
	Valor string `json:"valor"`
}

// Nova função orquestradora que será chamada pelo app.go
func BaixarDadosSGS() {
	codigos := []string{"25411", "25412", "24363", "24364", "25392", "25394"}
	dataInicial := "01/01/2026"
	dataFinal := "31/12/2026"

	for _, codigo := range codigos {
		fmt.Printf("Código atual buscado: %s\n", codigo)

		// Chama a função BuscarDadosSGS diretamente (sem 'sgs.')
		bcb, err := BuscarDadosSGS(codigo, dataInicial, dataFinal)

		if err != nil {
			fmt.Printf(" Falha ao buscar código %s: %v\n", codigo, err)
			continue // Usa continue para não parar o programa se um código falhar
		}

		fmt.Printf("Total de registros encontrados: %d\n", len(bcb))

		limite := 7
		if len(bcb) < limite {
			limite = len(bcb)
		}

		// Loop de impressão
		for j := 0; j < limite; j++ {
			fmt.Printf("Data: %s | Valor: %s\n", bcb[j].Data, bcb[j].Valor)
		}
		if len(bcb) > limite {
			fmt.Println("...")
		}

		nomeArquivo := fmt.Sprintf("serie_%s_2026.csv", codigo)

		// Chama a função SalvarEmCSV diretamente
		err = SalvarEmCSV(nomeArquivo, bcb)
		if err != nil {
			fmt.Printf("Erro ao salvar CSV: %v\n", err)
		} else {
			fmt.Printf("-> Sucesso! Arquivo '%s' gerado com %d registros.\n\n", nomeArquivo, len(bcb))
		}
	}
}

func BuscarDadosSGS(codigo, dataInicial, dataFinal string) ([]SGSRecord, error) {
	url := fmt.Sprintf("https://api.bcb.gov.br/dados/serie/bcdata.sgs.%s/dados?formato=json&dataInicial=%s&dataFinal=%s",
		codigo, dataInicial, dataFinal)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição HTTP: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o corpo da resposta: %v", err)
	}

	if resp.StatusCode != 200 {
		mensagemServidor := string(body)
		if len(mensagemServidor) > 100 {
			mensagemServidor = mensagemServidor[:100] + "..."
		}
		return nil, fmt.Errorf("API retornou status %d. Detalhe: %s", resp.StatusCode, mensagemServidor)
	}

	var records []SGSRecord
	err = json.Unmarshal(body, &records)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON: %v", err)
	}

	return records, nil
}

func SalvarEmCSV(nomeArquivo string, registros []SGSRecord) error {
	file, err := os.Create(nomeArquivo)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	cabecalho := []string{"Data", "Valor"}
	if err := writer.Write(cabecalho); err != nil {
		return fmt.Errorf("erro ao escrever cabeçalho: %v", err)
	}

	for _, record := range registros {
		linha := []string{record.Data, record.Valor}
		if err := writer.Write(linha); err != nil {
			return fmt.Errorf("erro ao escrever linha %s: %v", record.Data, err)
		}
	}

	return nil
}
