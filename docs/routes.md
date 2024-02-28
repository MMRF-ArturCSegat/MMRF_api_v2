# Routes

Esse modulo representa a API geral do servidor e suas rotas, por meio dele que o usuário interage com qualquer outro modulo

`CRUD.go` contem as rotas para criar o grafo (`grap_model/CSV_graph`), visualiz o grafo como caminho desnhavel em mapa, e interagir com suas sessões (`sessions/Session`). Criado ambos grafo
e sessão por meio de `CRUD.go/parse_csv_to_obj()` pode se realizar operações mais complexas no grafo por meio das rotas contidas em `graph_interaction_endpoint.go`

cada um dos componentes de fibra definidos em `fiber_optic_components/` tem seu respectivo arquivo contendo rotas para interagir com sua respectiva tabela no banco de dados que contem informação
sobre os componentes de fibra.

`instance_generation_endpoint.go` contém a rota para gerar o arquivo de texto com a instância do grafo armazenado na sessão do usuário que realizou o pedido na rota. 
