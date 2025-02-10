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
- [x] Implementar algum padrão de arquitetura
- [X] Realizar mais de um processamento por vez
- [x] Implementar integração com AWS
- [x] Realizar do serviço de status
- [ ] Implementar limitador de goroutines
- [x] Add GitHub Actions

# Workflow de CI/CD

Este projeto utiliza o GitHub Actions para automação do pipeline de CI/CD. O workflow realiza as seguintes etapas:

1. Configuração do ambiente Docker: Habilita o suporte ao QEMU e configura o Buildx para builds otimizados.
2. Autenticação no Docker Hub: Realiza login no Docker Hub utilizando credenciais seguras.
3. Criação e Envio da Imagem Docker: O código é empacotado em uma imagem Docker e enviado para o Docker Hub. Para mais detalhes, consulte o arquivo .github/workflows/deploy.yml
