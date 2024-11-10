# Servidor de Anúncios Rotativos com Filtro de Preferências

## Descrição

Este projeto é um sistema de exibição de anúncios personalizados com base nas preferências dos usuários, processadas em tempo real e sem a complexidade de leilões e bidding. O servidor recebe uma solicitação de um cliente, filtra anúncios com base no perfil do usuário e escolhe um anúncio relevante para exibição. Esse tipo de sistema é ideal para pequenas redes de publicidade ou sites que desejam mostrar anúncios relevantes sem a necessidade de um sistema de leilões complexo.

## Funcionalidades do Projeto

### 1. Receber Solicitações de Anúncios

- Um cliente (site ou app) envia uma requisição ao servidor, incluindo o ID do usuário e informações adicionais (ex: categoria de interesse, localização).

### 2. Filtro de Anúncios

- O servidor mantém um banco de dados de anúncios classificados em categorias, como "tecnologia", "moda", "esportes" etc.
- Com base nas preferências do usuário ou nas informações enviadas, o servidor seleciona um anúncio da categoria correspondente.

### 3. Rotação de Anúncios

- Para evitar repetição excessiva de anúncios, o servidor implementa um sistema de rotação, exibindo um anúncio diferente a cada nova solicitação.

### 4. Respostas em Tempo Real

- O servidor responde rapidamente ao cliente, enviando o anúncio mais relevante para exibição na interface do usuário.

### 5. Relatórios Básicos

- Implementação de um sistema básico de contagem de impressões e cliques, para monitorar quantas vezes um anúncio foi exibido e quantos cliques ele recebeu.

## Ferramentas e Tecnologias

### Backend em Go

- O servidor de anúncios é construído em Go, utilizando frameworks como Gin ou Echo para o servidor HTTP.

### Banco de Dados (SQLite ou Redis)

- SQLite é utilizado para armazenar dados dos anúncios e preferências dos usuários devido ao pequeno volume de dados.
- Redis pode ser uma opção adicional para cache rápido das preferências e anúncios disponíveis.

### Simulação de Dados

- Dados fictícios de usuários e anúncios categorizados são gerados para teste, sem a necessidade de integrações com DSPs ou ad-exchanges.

## Estrutura da Aplicação

### 1. Tabela de Anúncios

- Armazena anúncios em uma tabela com campos como `id`, `categoria`, `conteúdo` e `impressões`, para monitorar o número de exibições.

### 2. Tabela de Preferências de Usuários

- Armazena preferências de usuários com campos como `id`, `nome`, `categoria_preferida`, para filtrar anúncios de acordo com os dados de cada usuário.

### 3. Lógica de Seleção

- Com base nas preferências, o servidor escolhe um anúncio da categoria correspondente. Em caso de ausência de preferência, a seleção pode ser feita de forma aleatória ou por rotação.


### 4. Sistema de Relatórios

- Implementação de contagem simples de quantas vezes cada anúncio foi exibido e quantos cliques foram registrados.

## Como Executar

1. Clone o repositório.
2. Instale as dependências do Go.
3. Configure o banco de dados (SQLite ou Redis).
4. Inicie o servidor e faça solicitações para o endpoint de anúncios.

## Contribuições

Contribuições são bem-vindas! Sinta-se à vontade para enviar pull requests com melhorias, correções de bugs ou novas funcionalidades.

## Licença

Este projeto é distribuído sob a licença MIT.
