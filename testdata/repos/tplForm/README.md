# Ritchie Templates

## 🇺🇸 How to add some language on Rit create formula

- Create a folder on languages folder
- The folders on language's folder will be presented to user on `rit create formula`
- Inside the created folder you should have:
  - `Makefile`, this file should do:
    - create `run.sh` and `run.bat`, these files should explain how to run the formula,
     `run.sh` is called on linux and mac system
     and `build.bat` will be called on windows system
    - copy files that `run.sh` and `run.bat` needs to run
    - remember `Makefile` can be called inside a docker using
     the dockerImageBuilder as the docker image name
  - `config.json`, with this file rit can ask the inputs
  and use dockerImageBuilder to build the formula
  - `/src` folder, on this folder you can create a simple formula
   using the language that you will add.
  - `Dockerfile`, this file will be use when --docker is pass to the formula,
  the objective is to run the formula inside the docker,
  so in this file you need to create a dockerfile
  that can run any formula of this language

* * *

## 🇧🇷 Como adicionar uma nova linguagem ao rit create formula

Para adicionar uma nova linguagem ao rit create fórmula,
você apenas precisa adicionar uma nova pasta no caminho
templates/create_formula/languages.
todos os nomes de pastas adicionadas nesse caminho serão listadas
para o usuário quando ele executar rit create formula.

Dentro dessa pasta você deve adicionar:

- /src                  (obrigatório)
- Makefile              (caso rode no linux)
- build.bat             (caso rode no windows)
- config.json           (obrigatório)
- Dockerfile            (caso rode com docker)
- README.md             (documentação github)
- set_umask.sh          (caso rode com docker)
- metadata.json         (documentação portal)

Vamos entender para que server cada componente:

### Source

- **/src** nessa pasta você vai colocar uma fórmula de exemplo
 na linguagem que você quer adicionar.

- **config.json** nesse arquivo você deve adicionar os inputs da fórmula
 que está na pasta src. alem disso você pode adicionar
 o campo *dockerImageBuilder* que é a imagem que deve ser utiliza para
 fazer o build do seu código. Lembrando que esse campo é opcional, caso ele
 não exista o ritchie sempre vai fazer o build local.

### Build Local

Todo build deve gerar uma pasta bin com no minimo os seguintes arquivos:

- run.sh
- run.bat

Esses arquivos são os **arquivos de execução** que serão executados ao
 chamar a fórmula, logo, esses arquivos devem saber como rodar o código.

- **Makefile** é responsável por fazer o build do código em máquinas linux.
- **build.bat** é responsável por fazer o build do código em máquinas windows.

Alem de gerar os **arquivos de execução** o build também deve copiar os
 arquivos necessários para a pasta bin para que os
 **arquivos de execução** consigam funcionar.

### Build com Docker

Caso você adicione o campo *dockerImageBuilder* no **config.json**, o
 ritchie vai tentar fazer o build utilizando o docker.
 Com isso, ele vai rodar o arquivo **Makefile** dentro de um docker com a
 imagem do campo **dockerImageBuilder**. O build com docker deve gerar
 a mesma pasta bin que um build local, sendo que a grande vantagem do build com
 docker é que o usuário não precisar ter as ferramentas necessárias para
 o build instaladas na máquina.

### Run com Docker

Caso você adicione o arquivo **Dockerfile**, o ritchie vai rodar o **run.sh**
 dentro de um docker. Para isso, ele vai utilizar os arquivos:

- **Dockerfile** utilizaremos esse arquivo para fazer o
 build do docker que vai rodar o **run.sh**
- **set_umask.sh** é o entrypoint do docker,
 normalmente esse arquivo utiliza o comando umask para que o volume
 dentro do docker tenha uma melhor compatibilidade.

**Lembrando que devemos copiar esses arquivos para a
 pasta bin ao fazer o build.**

### Documentação

- **metadata.json** arquivo utilizado pelo portal de fórmulas
 do ritchie para fazer a indexação.
- **README.md** arquivo para explicar como utilizar a fórmula.
 Quando alguém abrir a pasta da sua fórmula pelo github,
 ele vai ver o conteúdo desse arquivo.

### Pasta root

Quando o usuário cria uma fórmula pela primeira vez em um workspace,
ele copia os arquivos da pasta root para o workspace do usuário,
caso a linguagem tenha alguma regra nova de gitignore você pode
adicionar essa regra no arquivo .gitignore da pasta root.
