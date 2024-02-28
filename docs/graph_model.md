# Graph Model

Modelo de grafo para interação da API.
Este modulo define a classe `CSV_Graph` que representa um grafo bidirecional com nodos `Nodes` dentro de um limite `Limiter` com origem de linha `Olt`

Csv no nome se da ao fato que estes grafos são instanciados na função `graph_csv_parser.go/New_csvg()`que le um arquivo csv e gera um grafo valido a partir disto, em geral esses grafos
são armazenados com sessões (ver `sessions/sessions.go`), para que possam ser usados mais de uma vez pelo mesmo usuário

Outras funções de importância são `graph_extended_methods.go/LimitedBranchingFrom()` que recebe um nodo e um limite maximo em metros do nodo original, e calcula o subgrafo de caminhos possíves
com soma da distância dos arcos menor que o limite. E a função `graph_extended_methods.go/ClosestNode()` que recebe uma coordenada e responde com o nodo mais próximo

`graph_string_conversion.go` tem algumas funções uteis para converter a informação contida no grafo para um formato mais amigáel para usar na geração da instância, e `graph_csvg_adapter.gp`
contém algumas funções uteis para tranformar o grafo em caminho de coordenadas por exemplo, de forma mais conveniente para desenhar em um mapa.
