// Trabalho 2 da disciplina de lab avanc, alunos Luciano uchoa e Gabriel Dias Goncalves
package main

// faz os imports necessarios
import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
)

// Cria os nossos arquivos locais
var (
	dataFile = "data.json"
	mu       sync.Mutex
)

// define o tipo de como vamos guardar os dados no caso key-value, isso e como se fosse o model da aplicaçao
type KeyValueStore map[string]string

// cria a pagina HTML usada no programa, define tb que as rotas de submit e query são chamadas no clique do botao, caso queira ver um erro de tipo de requisiçao basta trocar o tipo de metodo nas acoes dos botoes
const mainPage = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Chave-Conteúdo</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f0f0f0;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
        }
        .container {
            background-color: #fff;
            padding: 20px;
            box-shadow: 0 0 10px rgba(0,0,0,0.1);
            border-radius: 8px;
            text-align: center;
            width: 80%; /* Define a largura máxima do contêiner */
            max-width: 800px; /* Define a largura máxima do contêiner */
        }
        h1 {
            color: #333;
        }
        p {
            font-size: 18px;
            color: #666;
            word-wrap: break-word;
        }
        input[type="text"] {
            width: calc(100% - 22px);
            padding: 10px;
            margin: 8px 0;
            box-sizing: border-box;
            border: 2px solid #ccc;
            border-radius: 4px;
        }
        input[type="submit"] {
            background-color: #4CAF50;
            color: white;
            border: none;
            padding: 10px 20px;
            text-align: center;
            text-decoration: none;
            display: inline-block;
            font-size: 16px;
            margin: 4px 2px;
            cursor: pointer;
            border-radius: 4px;
        }
        input[type="submit"]:hover {
            background-color: #45a049;
        }
        form {
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Insira Chave e Conteúdo</h1>
        <form action="/submit" method="post">
            <label for="key">Chave:</label><br>
            <input type="text" id="key" name="key" required><br>
            <label for="value">Conteúdo:</label><br>
            <input type="text" id="value" name="value" required><br><br>
            <input type="submit" value="Salvar">
        </form>
        <h1>Consultar Conteúdo</h1>
        <form action="/query" method="get">
            <label for="key">Chave:</label><br>
            <input type="text" id="key" name="key" required><br><br>
            <input type="submit" value="Consultar">
        </form>
    </div>
</body>
</html>
`
const queryResultPage = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Resultado da Consulta</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f0f0f0;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
        }
        .container {
            background-color: #fff;
            padding: 20px;
            box-shadow: 0 0 10px rgba(0,0,0,0.1);
            border-radius: 8px;
            text-align: center;
            width: 80%;
            max-width: 800px;
        }
        h1 {
            color: #333;
        }
        p {
            font-size: 18px;
            color: #666;
        }
        a {
            text-decoration: none;
            color: #4CAF50;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Resultado da Consulta</h1>
        <p>{{.Result}}</p>
        <a href="/">Voltar à página inicial</a>
    </div>
</body>
</html> 
`

func loadStore() (KeyValueStore, error) {
	// Essa func serve pra carregar a store que no caso são nosso arquivos locais
	// tenta abrir os arquivos locais e nos retorna um arquivo vazio em caso de erro, esse arquivo tb e do tipo [string]string declarado na linha 19
	file, err := os.Open(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return make(KeyValueStore), nil
		}
		return nil, err
	}
	defer file.Close()
	// Uma variavel store vazia para receber os valores dos nossos arquivos locais
	store := make(KeyValueStore)
	// Fazemos o decode dos valores vindos dos arquivos locais
	if err := json.NewDecoder(file).Decode(&store); err != nil {
		return nil, err
	}
	// retorna os arquivos locais ja decodados para o formato json
	return store, nil
}

func saveStore(store KeyValueStore) error {
	// Essa func e responsavel por salvar um store(que possui key e value) nos nossos arquivos locais

	// Cria um file ou trunca um existente pra poder salvar os novos valores
	file, err := os.Create(dataFile)

	// trata o erro caso a operaçao de salvar de errado
	if err != nil {
		return err
	}
	defer file.Close()

	// faz o encode do arquivo em formato Json pra poder salvar, caso tenhamos erros para no if de cima
	return json.NewEncoder(file).Encode(store)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Essa func serve pra mostrar a pagina html principal criada
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(mainPage))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	// Essa func lida com o botao submit
	// Primeiro checamos se o metodo que tentamos chamar ao clicar submit e um POST
	if r.Method != http.MethodPost {
		http.Error(w, "Botao submit serve apenas pra posts", http.StatusMethodNotAllowed)
		return
	}
	// Guardamos os valores digitados nas variaveis
	key := r.FormValue("key")
	value := r.FormValue("value")

	mu.Lock()
	defer mu.Unlock()

	// cria instancia de store e carrega ela usando a func loadStore
	store, err := loadStore()
	if err != nil {
		http.Error(w, "Falha carregando arquivos locais", http.StatusInternalServerError)
		return
	}

	// atribui a key digitada a propriedade key da store que criamos
	store[key] = value

	// tentamos salvar a key nos arquivos locais e retornamos erro caso seja o caso
	// olhando com atençao percebe-se que a store que passamos por parametro pra func saveStore e a store que criamos na linha 156 e 157 e que possui a chave e o conteudo que temos que salvar
	if err := saveStore(store); err != nil {
		http.Error(w, "Falha ao salvar nos arquivos locais", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	// Essa func lida com a procura no arquivo local pela chave digitada

	// Primeiro checamos se o metodo que tentamos chamar ao clicar submit e um GET
	if r.Method != http.MethodGet {
		http.Error(w, "Botao submit serve apenas pra gets", http.StatusMethodNotAllowed)
		return
	}

	// Aqui iniciamos uma query(pesquisa ) pra entendermos se a chave digitada esta nos dados locais
	key := r.URL.Query().Get("key")

	mu.Lock()
	defer mu.Unlock()

	// cria instancia de store e carrega ela
	store, err := loadStore()
	if err != nil {
		http.Error(w, "Falha carregando arquivos locais", http.StatusInternalServerError)
		return
	}
	// Aqui checamos de fato se o valor da pesquisa esta nos arquivos locais, se estiver n estiver "ok", retornamos o erro chave n encontrada
	value, ok := store[key]
	if !ok {
		http.Error(w, "Chave nao encontrada", http.StatusNotFound)
		return
	}
	// aqui vamos criar um template pra apresentar uma pagina html com o resultado da query
	tmpl, err := template.New("queryResult").Parse(queryResultPage)
	// tratamos o erro caso n consigamos criar o template
	if err != nil {
		http.Error(w, "Falha ao carregar template", http.StatusInternalServerError)
		return
	}
	// definimos a variavel data que e o que sera mostrado no template
	data := struct {
		Result string
	}{
		Result: value,
	}
	// mostramos o a pagina html com a o resultado da query
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

func main() {
	// Aq definimos as rotas
	// Se der uma olhada na pagina html main page que criamos acima, vera que as açoes dos botoes estao linkadas com essa rotas, cada açao chama uma rota
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/query", queryHandler)
	// Inicia o Server e nos mostra um erro caso as coisas deem errado ao iniciar o server na porta 8080
	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Importante mutex, em ambas as funcs que mexem com store utilizamos o muLock esse bloco e importante pq faz com que essas funcoes rodem em um bloco que as auxilia a lidar com concorrencia, previnando as condiçoes de corrida vistas em sala de aula
