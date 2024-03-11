# Fiber Optic Components

### Modulo para interação com um banco de dados SQLite de Components de Fibra Óptica 

Este modulo define diversas estruturas que representam diferentes componentes de fibra óptica, cada um deles implementa a 
interface `querys.go/FiberComponent` que permite que todos compartilhem as mesmas funções de interação com o banco de dados. Essa interface é batante
simples e apenas define uma função `GetId()` que retorno o id do objeto e `String()` que é útil para transcrever o objeot no arquivo final de instância
