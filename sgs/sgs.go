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

	// ---> ADICIONE ESTE BLOCO AQUI <---
	// Se o status não for 200 (OK), não tentamos ler como JSON.
	if resp.StatusCode != 200 {
		// Retornamos o erro mostrando o Status e um pedaço da mensagem HTML que o servidor enviou
		mensagemServidor := string(body)
		if len(mensagemServidor) > 100 {
			mensagemServidor = mensagemServidor[:100] + "..." // Corta para não poluir o terminal
		}
		return nil, fmt.Errorf("API retornou status %d. Detalhe: %s", resp.StatusCode, mensagemServidor)
	}
	// ----------------------------------

	var records []SGSRecord
	err = json.Unmarshal(body, &records)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON: %v", err)
	}

	return records, nil
}

func SalvarEmCSV(nomeArquivo string, registros []SGSRecord) error {
	// 1. Cria ou sobrescreve o arquivo no sistema
	file, err := os.Create(nomeArquivo)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %v", err)
	}
	defer file.Close()

	// 2. Cria um novo escritor CSV que usará o arquivo
	writer := csv.NewWriter(file)
	// Garante que todos os dados que estão na memória sejam gravados no arquivo antes da função terminar
	defer writer.Flush()

	// 3. Escreve o cabeçalho das colunas (opcional, mas recomendado)
	cabecalho := []string{"Data", "Valor"}
	if err := writer.Write(cabecalho); err != nil {
		return fmt.Errorf("erro ao escrever cabeçalho: %v", err)
	}

	// 4. Itera sobre os registros convertendo cada um para um slice de strings (linha)
	for _, record := range registros {
		linha := []string{record.Data, record.Valor}
		if err := writer.Write(linha); err != nil {
			return fmt.Errorf("erro ao escrever linha %s: %v", record.Data, err)
		}
	}

	return nil
}
