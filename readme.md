# Servidor Web em Go

Este projeto é uma aplicação web simples implementada em Go. Ele fornece um servidor web que atende na porta 8080 e tem as seguintes funcões:

1. **Registrar um objeto do tipo CHAVE : CONTEÚDO**
2. **Salvar esse objeto em um arquivo local que consegue ser encontrado na raiz desse diretorio**
3. **Consulta esse arquivo local buscando por uma chave especifica que o usuario digitar**

## Informações do Projeto

- **Equipe**: Luciano Uchoa & Gabriel Dias
- **Tecnologia Utilizada**: Go (Golang)

## Como usar:
  # Registrar chave
   - Digite um valor e uma chave nos campos solicitados
   - Clique no botao salvar.
   - Apos salvar veja o arquivo json criado na diretriz deste projeto, la sera possivel ver as chaves registradas
  # Consultar chave
   - Digite o valor da chave que deseja consultar
   - Clique em consultar
   - Caso a chave ja tenha sido cadastrada voce ira para uma pagina com o valor do conteudo daquela chave, caso contrario apresentaremos um erro


- **`servido.go`**: Código fonte da aplicação Go.
- **`readme.md`**: Este arquivo de documentação.

## Como Executar o Projeto

1. **Instale o Go**:
   - Certifique-se de ter o Go instalado em sua máquina. Você pode baixá-lo em [golang.org](https://golang.org/dl/).

2. **Clone o Repositório**:
   ```bash
   git clone [URL_DO_REPOSITORIO]
   cd Lab
   ```

3. **Compile e Execute o Servidor:**:
    ```bash
    go run servidor.go
    ```

4. **Acesse o Servidor:**:
    ```
    Abra seu navegador e acesse http://localhost:8080 para ver a página.
    ```
