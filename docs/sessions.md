# Sessions

### Modulo simples para gerenciar sessões persistentes no servidro

Defines duas estruturas principais `cookies.go/ServerCookie` armazeda um identificador único que é gerado assim que o usuário gera seu grafo, e um mapa 
chamado `sessions/Sessions` que liga um dado cookie a um Grafo. 

Uma função interessante de notar é `sessions/CleanExpiredSessions()` que é chamada de 15 em 15 minutos de maneira assincrona em `main.go`
