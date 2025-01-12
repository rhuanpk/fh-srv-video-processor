# Video Processor

Para executar o MVP, tenha o processador de mídia [FFmpeg](https://www.ffmpeg.org/download.html) e a linguagem de programação [Golang](https://go.dev/dl/) instalados no sistema operacional (e disponíveis no PATH).

Pré-Execução:
1. Coloque videos de amostra dentro da pasta "video" (que está na raiz do projeto) e altere as constantes do programa
1. Se desejar alterar a quantidade de captura de _frames_ por segundo, troque o valor da constante `frameInterval` no arquivo **main.go**

Execução:
1. Abra o emulador de terminal do seu sistema operacional na pasta raiz do projeto
1. Prepare o programa (apenas na primeira execução):
	`go mod tidy`
1. Execute o programa:
	`go run ./`

ToDo:
- [x] Criar MVP do serviço
- [x] Realizar teste de carga com vídeos grandes
- [ ] Implementar algum padrão de arquitetura
- [X] Realizar mais de um processamento por vez
- [ ] Implementar integração com AWS
