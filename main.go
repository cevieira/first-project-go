package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramentos = 5
const delay = 5

func main() {
	exibeOpcoes()
	comando := leituraComando()

	switch comando {
	case 1:
		iniciarMonitoramento()
	case 2:
		fmt.Println("Exibindo os logs")
		imprimeLogs()
	case 3:
		fmt.Println("Saindo do programa")
		os.Exit(0)
	default:
		fmt.Println("Comando não reconhecido")
		os.Exit(-1)
	}

}

func exibeOpcoes() {
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir Logs")
	fmt.Println("0 - Sair do Programa")
}

func leituraComando() int {
	var comando int
	fmt.Scan(&comando)
	fmt.Println("O comando escolhido foi o: ", comando)

	return comando
}

func iniciarMonitoramento() {
	fmt.Println("Iniciando o monitoramento")
	sites := extraiSitesArquivo()

	for i := 0; i < monitoramentos; i++ {
		for i, site := range sites {
			fmt.Println("Executando monitoramento ", i, "site:", site)
			executaMonitoramento(site)
		}
		time.Sleep(delay * time.Minute)
	}

}

func executaMonitoramento(site string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("Ocorreu um erro: ", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println(time.Now().Format("02/01/2006 15:04:05")+" - "+"Site:", site, "foi carregado com sucesso")
		registraLog(site, true)
	} else {
		fmt.Println(time.Now().Format("02/01/2006 15:04:05")+" - "+"Site:", site, " não foi carregado com sucesso")
		registraLog(site, false)
	}
}

func extraiSitesArquivo() []string {
	var sites []string
	arquivo, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro: ", err)
	}

	leitor := bufio.NewReader(arquivo)

	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)
		sites = append(sites, linha)
		if err != nil {
			break
		}
	}

	arquivo.Close()
	return sites
}

func registraLog(site string, status bool) {
	arquivo, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Ocorreu um erro: ", err)
	}

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")

	arquivo.Close()

}

func imprimeLogs() {
	arquivo, err := ioutil.ReadFile("logs.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro: ", err)
	}

	fmt.Println(string(arquivo))
}
