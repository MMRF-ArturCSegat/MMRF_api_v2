# Instance Generation

Este módulo lida com os objetos `Instance` e `Solution` que representam respectivamente uma instância, e subsequentemente sua solução, as partes do modulo que lidam com as soluções (`solution.go`)
não foram realmente testadas, pois não tive um par de instância gerada pelo Routelib com sua respectiva solução, mas o código dentro de `solution.go` deve ser capaz de ler e interpretar um arquivo de solução,
enquanto o código dentro de `instance.go` deve ser capaz de criar um arquivo de instância com base em nos dados contidos em um objeto `Instance` que é criado em `routes/instance_generation_endpoint.go`

### Instance

A classe `Instance` é composta da seguinte maneira

```
type Instance struct {
    Paths               [][]gm.GraphPath     `json:"paths"`
    Clients             []util.Coord         `json:"clients"`
    Cable_id            uint32               `json:"Cable_id"`
    Splicebox_id        uint32               `json:"Splicebox_id"`
    Uspliters_id        []uint32             `json:"Uspliters_id"`
    Bspliters_id        []uint32             `json:"Bspliters_id"`
}
```
*Paths*: Resultado de chamadas ao metodo `graph_model/graph_extended_methods.go/limitedBranchingFrom()` cada `[]gm.GraphPath` representa um conjunto de caminhos possíveis de um dado cliente, isso geralmente é calculado
da seguitne maneira:

1 Interface recebe coordenada de um dado cliente

2 Calculase por meio de `graph_model/graph_extended_methods.go/closestNode()` qual o nodo mais próximo deste cliente

3 Usando o nodo calculado na última etapa se usa `graph_model/graph_extended_methods.go/limitedBranchingFrom()` para calcular todos os outos nodos relavantes àquele cliente

*Clients*: As coordenadas dos clientes específicados

*Cable_id*: O id do tipo de cabo selecionado para esta instância

*Splicebox_id*: O id do tipo de splicebox selecionado para esta instância

*Uspliters_id*: Vetor contendo os ids dos dipos de splites não balanceado (_Unbalanced_) para esta instância 

*Bspliters_id*: Vetor contendo os ids dos dipos de splites balanceado (_Balanced_) para esta instância 

Essa Classe tem também alguns metodos auxiliares para converter os campos com id's em campos com objetos com tipos complexos como `GetCable()` que retorna o objeto de cabo armazenado no banco de dados com o
id = `Cable_id`.

Por último a principal função da classe `Instance` é o metodo `GenerateSubGraphOptimizationFile()` que recebe a propria classe assim como o grafo sobre a qual foi gerada e gera um arquivo no formato de otimização
que deve ser respondido para o cliente, um exemplo disso pode ser visto em `routes/instance_generation_endpoint.go`

### Solution

`solution.go` define 3 estruturas importantes

`Solution`
```
type Solution struct {
    Path                     [][2]util.Coord
    NodeIdToBspliter         map[uint32]foc.FiberBalancedSpliter
    NodeIdToUspliter         map[uint32]foc.FiberUnbalancedSpliter
    NodeIdToCable            map[uint32]foc.FiberCable
    SpliceboxNodesId         []uint32
}
```
Define a solução em si:
Path: lista de pares de coordenadas, cada par representa uma linha para ser desenhada ou percorrida, em geral os pares são as coordenadas de 2 nodos com arcos entre si.

NodeIdToBspliter, NodeIdToUsplite e NodeIdToCable: são mapas que ligam ID de um nodo a um componente de fibra que foi aplicado nesse nodo.

SpliceboxNodesID: como cada solução tem apenas um tipo de splicebox, cada id neste vetor é o id de um nodo onde esta splicebox foi aplicada.

`VirualNetwork`
```
type VirtualNetwork struct {
    Path                [][2]util.Coord
    BspliterMap         map[uint32]foc.FiberBalancedSpliter
}
```
A rede virtual é basicamente uma sub solução contendo 2 dos campos da solução, esses campos cumprem a mesma função aqui.

`solutionOrVirtnet`
```
type solutionOrVirtNet interface {
    addPath([2]util.Coord)
}
```

Esta interface define um método comum `addPath()` que adiciona um par de coordenadas coonectadas ao caminho do objeto, serve para q funções como `loadInfrastructure()`
possam interagir tanto com obejetos de tipo `Solution` quando com `VirtualNetwork` sem ter que definir duas funções diferentes, um exemplo é a função `loadInfrastructure()`
que recebe um _buffer_ contendo as linhas definindo arcos e preenche os vetores do campo _Path_ de acordo.

A função mais importante é `ParseSolutionFile()` que recebe um arquivo de solução, a `Instance` da qual ele foi gerado e o `CSV_Graph` no qual esta `Instance` foi gerada.
Essa função então le o arquivo e separa as linhas dele em diferentes _buffers_, de acordo com o que cada linha representa. Esses _buffers_ são então lidos separadamente para
preencher objetos de `Solution` e `VirtualNetwork`.






