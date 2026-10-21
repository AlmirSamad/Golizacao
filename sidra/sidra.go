package sidra

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type SidraRecord struct {
	Codigo string
	Link   string
}

func BaixarDadosSidra() {

	arquivo, err := os.Open("sidralinks.txt")
	if err != nil {
		fmt.Printf("Erro ao abrir o arquivo: %v\n", err)
		return
	}
	defer arquivo.Close()

	var fontes []SidraRecord

	scanner := bufio.NewScanner(arquivo)

	for scanner.Scan() {
		linha := scanner.Text()

		if strings.TrimSpace(linha) == "" {
			continue
		}

		partes := strings.SplitN(linha, " : ", 2)

		if len(partes) == 2 {
			codigo := strings.TrimSpace(partes[0])
			link := strings.TrimSpace(partes[1])

			novoItem := SidraRecord{
				Codigo: codigo,
				Link:   link,
			}
			fontes = append(fontes, novoItem)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Erro durante a leitura do arquivo: %v\n", err)
		return
	}

	fmt.Printf("Foram encontradas %d tabelas:\n\n", len(fontes))

	for i, fonte := range fontes {
		// ---------------------------------------------------------
		// A MELHORIA: O timer agora fica no TOPO do loop.
		// "Se o índice for maior que 0 (ou seja, se NÃO for a primeira tabela), espere."
		if i > 0 {
			fmt.Println("Aguardando 60 segundos antes de tentar a próxima tabela...")
			time.Sleep(60 * time.Second)
			fmt.Println()
		}
		// ---------------------------------------------------------

		fmt.Printf("A variável Código é: %s\n", fonte.Codigo)
		fmt.Printf("A variável Link é: %s\n", fonte.Link)
		fmt.Println("--------------------------------------------------")

		req, err := http.NewRequest("GET", fonte.Link, nil)
		if err != nil {
			fmt.Printf("Erro ao montar a requisição: %v\n", err)
			continue
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
		req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Referer", "https://sidra.ibge.gov.br/")
		req.Header.Set("Upgrade-Insecure-Requests", "1")

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			fmt.Printf("Erro ao fazer requisição: %v\n", err)
			continue
		}

		if resp.StatusCode != 200 {
			fmt.Printf("Erro: o servidor retornou o status %d\n", resp.StatusCode)
			resp.Body.Close()
			continue
		}

		nomeArquivo := fmt.Sprintf("sidra_%s.xlsx", fonte.Codigo)

		out, err := os.Create(nomeArquivo)
		if err != nil {
			fmt.Printf("Erro ao criar o arquivo local: %v\n", err)
			resp.Body.Close()
			continue
		}

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			fmt.Printf("Erro ao salvar os dados no arquivo: %v\n", err)
		} else {
			fmt.Printf("Sucesso! Arquivo salvo como '%s'\n", nomeArquivo)
		}

		out.Close()
		resp.Body.Close()

		// (O timer do final foi removido daqui)
	}
}
