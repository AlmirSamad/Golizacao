package app

import (
	"fmt"
	"golizacao/sgs"
	"golizacao/sidra"

	"github.com/urfave/cli"
)

// Retorna a aplicação de linha de comnado pronta para ser executada
func Gerar() *cli.App {
	app := cli.NewApp()
	app.Name = "Aplicação de linha de comando"
	app.Usage = "Atualiza dados do Amazonas"

	app.Commands = []cli.Command{
		{
			Name:  "sidra",
			Usage: "Atualizar todos os dados do Sidra",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "tempo", // Nome da flag (--tempo)
					Value: 60,      // Valor padrão caso o usuário não digite nada
					Usage: "Tempo de espera em segundos entre as requisições",
				},
				cli.StringFlag{
					Name:  "arquivo",        // Nome da flag (--arquivo)
					Value: "sidralinks.txt", // Valor padrão
					Usage: "Caminho do arquivo txt contendo os links",
				},
			},
			Action: atualizarDadosDoSidra,
		},

		{

			Name:   "sgs",
			Usage:  "Atualizar dados do Banco Central (SGS)",
			Action: atualizarDadosDoSGS,
		},
	}

	return app
}

func atualizarDadosDoSidra(c *cli.Context) {
	// Resgata os valores das flags pelos nomes definidos acima
	tempo := c.Int("tempo")
	arquivo := c.String("arquivo")

	fmt.Printf("Iniciando rotina do Sidra usando o arquivo '%s' com pausa de %d segundos.\n", arquivo, tempo)

	// Agora você pode passar esses valores para a sua função principal!
	// *Lembre-se de atualizar a função BaixarDadosSidra para aceitar esses parâmetros
	sidra.BaixarDadosSidra(arquivo, tempo)
}

func atualizarDadosDoSGS(c *cli.Context) {
	fmt.Println("Iniciando rotina de extração do SGS (Banco Central)...")
	sgs.BaixarDadosSGS()
}
