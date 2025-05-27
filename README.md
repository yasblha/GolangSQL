# GoSQL - Un Moteur SQLite en Go 🚀

## Qu'est-ce que c'est ? 🤔

GoSQL est comme un petit assistant qui peut lire et comprendre les fichiers SQLite, comme si on lui donnait un livre et qu'il pouvait le lire et le comprendre ! 

## Comment ça marche ? 🎮

Imaginez que vous avez une boîte magique (notre base de données) qui contient des tiroirs (les tables) et dans chaque tiroir, il y a des fiches (les données). Notre programme peut :
1. Ouvrir la boîte
2. Lire ce qui est écrit sur les tiroirs
3. Ajouter ou chercher des fiches dans les tiroirs

```
┌─────────────────────────────────┐
│           GoSQL Engine          │
├─────────┬─────────┬────────────┤
│  Reader │ Parser  │ Executor   │
└─────────┴─────────┴────────────┘
```

## Les Parties Principales 🎯

### 1. Le Reader (db/) 📖
C'est comme quelqu'un qui sait lire le langage spécial de SQLite.

```
┌─────────────────────────────────┐
│           SQLite File           │
├─────────────────────────────────┤
│  ┌─────────┐  ┌─────────────┐  │
│  │ Header  │  │ Master Page │  │
│  └─────────┘  └─────────────┘  │
│  ┌─────────┐  ┌─────────────┐  │
│  │ Data    │  │ Indexes     │  │
│  └─────────┘  └─────────────┘  │
└─────────────────────────────────┘
```

#### Le Header (header.go)
- C'est comme la première page d'un livre
- Il nous dit :
  - La taille des pages (4096 bytes)
  - Le type d'encodage (UTF-8)
  - La version de SQLite

#### La Table Master (master.go)
- C'est comme la table des matières
- Elle nous dit quelles tables existent
- Elle nous donne la structure de chaque table

#### Les Indexes (index.go)
- C'est comme un index de livre
- Il nous aide à trouver les données plus vite
- Il utilise un arbre B (comme un arbre généalogique)

```
     [10]
    /    \
  [5]    [15]
 /   \   /   \
[1] [7] [12] [20]
```

### 2. Le Parser (parser/) 🔍
C'est comme un traducteur qui comprend le langage SQL.

#### Comment il fonctionne :
1. Il reçoit une commande SQL
2. Il la découpe en morceaux
3. Il comprend ce qu'on veut faire

```
┌─────────────────┐
│  SELECT * FROM  │
│  users WHERE    │
│  age > 18       │
└─────────────────┘
        ↓
┌─────────────────┐
│  - Type: SELECT │
│  - Table: users │
│  - Colonnes: *  │
│  - Condition:   │
│    age > 18     │
└─────────────────┘
```

#### Les Types de Requêtes :

1. SELECT (Chercher des données)
```sql
SELECT nom, age FROM utilisateurs WHERE age > 18
```

2. INSERT (Ajouter des données)
```sql
INSERT INTO utilisateurs (nom, age) VALUES ('Jean', 20)
```

3. CREATE TABLE (Créer une nouvelle table)
```sql
CREATE TABLE utilisateurs (
    id INTEGER PRIMARY KEY,
    nom TEXT,
    age INTEGER
)
```

### 3. L'Executor (à venir) ⚙️
C'est comme un robot qui exécute les ordres :
1. Il reçoit les instructions du parser
2. Il utilise le reader pour trouver les données
3. Il fait ce qu'on lui demande (chercher, ajouter, etc.)

## Comment Utiliser GoSQL ? 🛠️

### Installation
```bash
git clone https://github.com/votre-nom/gosql.git
cd gosql
go build
```

### Exemple d'Utilisation
```go
// Ouvrir une base de données
file, _ := os.Open("ma_base.db")

// Lire l'en-tête
info := db.ParseHeader(file)

// Lire les tables
tables := db.ReadMasterTable(file, info)

// Exécuter une requête
query := "SELECT * FROM users WHERE age > 18"
result := parser.ParseSelect(query)
```

## Les Fichiers Importants 📁

### Dans le dossier `db/` :
- `header.go` : Lit l'en-tête du fichier SQLite
- `master.go` : Gère la table des matières
- `index.go` : Gère les index pour chercher vite
- `table.go` : Définit ce qu'est une table
- `page.go` : Lit les pages de données

### Dans le dossier `parser/` :
- `sql.go` : Comprend le langage SQL
- `parser_test.go` : Vérifie que tout fonctionne

## Les Fichiers Binaires 🔍

### Structure d'un Fichier SQLite
```
┌─────────────────────────────────────────────────┐
│                    SQLite File                  │
├─────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────┐   │
│  │              Header (100 bytes)          │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  │   │
│  │  │ Magic   │  │ Page    │  │ Version │  │   │
│  │  │ String  │  │ Size    │  │ Info    │  │   │
│  │  └─────────┘  └─────────┘  └─────────┘  │   │
│  └─────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────┐   │
│  │            Master Page (4096 bytes)      │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  │   │
│  │  │ Table   │  │ Index   │  │ Schema  │  │   │
│  │  │ Info    │  │ Info    │  │ Info    │  │   │
│  │  └─────────┘  └─────────┘  └─────────┘  │   │
│  └─────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────┐   │
│  │            Data Pages (4096 bytes)       │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  │   │
│  │  │ Row 1   │  │ Row 2   │  │ Row 3   │  │   │
│  │  │ Data    │  │ Data    │  │ Data    │  │   │
│  │  └─────────┘  └─────────┘  └─────────┘  │   │
│  └─────────────────────────────────────────┘   │
└─────────────────────────────────────────────────┘
```

### Format Binaire des Données
```
┌─────────┬─────────┬─────────┬─────────┐
│  Type   │  Size   │  Data   │  Next   │
└─────────┴─────────┴─────────┴─────────┘
   1 byte   2 bytes   N bytes   4 bytes
```

#### Types de Données Supportés
```
INTEGER:  ┌───┐┌─────────────┐
          │ 1 ││    Value    │
          └───┘└─────────────┘

TEXT:     ┌───┐┌───┐┌─────────────┐
          │ 2 ││ N ││    Text     │
          └───┘└───┘└─────────────┘

BLOB:     ┌───┐┌───┐┌─────────────┐
          │ 3 ││ N ││    Data     │
          └───┘└───┘└─────────────┘
```

## Exemples Détaillés de Requêtes 📝

### 1. Requêtes SELECT

#### Simple SELECT
```sql
SELECT * FROM users;
```
```
┌─────────┬─────────┬─────────┐
│   id    │  name   │  age    │
├─────────┼─────────┼─────────┤
│    1    │  John   │   25    │
│    2    │  Alice  │   30    │
└─────────┴─────────┴─────────┘
```

#### SELECT avec WHERE
```sql
SELECT name, age FROM users WHERE age > 25;
```
```
┌─────────┬─────────┐
│  name   │  age    │
├─────────┼─────────┤
│  Alice  │   30    │
└─────────┴─────────┘
```

#### SELECT avec Conditions Multiples
```sql
SELECT * FROM users 
WHERE age > 20 AND name LIKE 'J%';
```
```
┌─────────┬─────────┬─────────┐
│   id    │  name   │  age    │
├─────────┼─────────┼─────────┤
│    1    │  John   │   25    │
└─────────┴─────────┴─────────┘
```

### 2. Requêtes INSERT

#### Insertion Simple
```sql
INSERT INTO users (name, age) VALUES ('Bob', 35);
```
```
Avant:          Après:
┌─────────┐     ┌─────────┐
│  John   │     │  John   │
│  Alice  │     │  Alice  │
└─────────┘     │  Bob    │
                └─────────┘
```

#### Insertion Multiple
```sql
INSERT INTO users (name, age) VALUES 
    ('Eve', 28),
    ('Frank', 42);
```
```
┌─────────┬─────────┐
│  name   │  age    │
├─────────┼─────────┤
│  John   │   25    │
│  Alice  │   30    │
│  Eve    │   28    │
│  Frank  │   42    │
└─────────┴─────────┘
```

### 3. Création de Tables

#### Table Simple
```sql
CREATE TABLE products (
    id INTEGER PRIMARY KEY,
    name TEXT,
    price REAL
);
```
```
┌─────────────────────────┐
│      products Table     │
├─────────┬───────────────┤
│   id    │ INTEGER (PK)  │
│  name   │     TEXT      │
│  price  │     REAL      │
└─────────┴───────────────┘
```

#### Table avec Index
```sql
CREATE TABLE orders (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    date DATETIME,
    INDEX idx_user (user_id)
);
```
```
┌─────────────────────────┐
│       orders Table      │
├─────────┬───────────────┤
│   id    │ INTEGER (PK)  │
│ user_id │   INTEGER     │
│  date   │   DATETIME    │
└─────────┴───────────────┘
     │
     ▼
┌─────────────┐
│  B-Tree     │
│  Index      │
└─────────────┘
```

## Processus d'Exécution des Requêtes 🔄

### 1. Parsing
```
Requête SQL → Tokens → Arbre Syntaxique
```

### 2. Validation
```
Arbre Syntaxique → Vérification → Plan d'Exécution
```

### 3. Exécution
```
Plan d'Exécution → Lecture des Données → Résultats
```

## Gestion de la Mémoire et du Cache 💾

### 1. Ce qui est en RAM 🚀

```
┌─────────────────────────────────┐
│           En Mémoire            │
├─────────────────────────────────┤
│  • En-tête du fichier           │
│  • Table des matières           │
│  • Index principaux             │
│  • Cache des pages récentes     │
└─────────────────────────────────┘
```

#### Exemple avec un fichier de 1 To
```
Fichier DB (1 To)
    │
    ├── En RAM (quelques Mo)
    │   ├── En-tête (100 bytes)
    │   ├── Table des matières (4 Ko)
    │   └── Cache (100 Mo max)
    │
    └── Sur Disque (1 To)
        ├── Données
        └── Index secondaires
```

### 2. Comment ça marche ? 🔄

#### Lecture d'une Donnée
```
1. Vérifier le cache en RAM
   ┌─────────┐
   │  Cache  │ → Si trouvé, retourner
   └─────────┘

2. Si pas en cache
   ┌─────────┐    ┌─────────┐
   │  Disque │ → │  Cache  │
   └─────────┘    └─────────┘
```

#### Exemple Concret
```
Requête: SELECT * FROM users WHERE id = 1000

1. Vérifie l'index en RAM
   ┌─────────┐
   │  Index  │ → Page 42
   └─────────┘

2. Vérifie le cache
   ┌─────────┐
   │  Cache  │ → Page 42 non trouvée
   └─────────┘

3. Lit du disque
   ┌─────────┐    ┌─────────┐
   │  Disque │ → │  Cache  │
   └─────────┘    └─────────┘
```

### 3. Gestion du Cache 🎯

#### Taille du Cache
```
┌─────────────────────────────────┐
│         Taille du Cache         │
├─────────────────────────────────┤
│  • Par défaut: 100 Mo           │
│  • Configurable                 │
│  • Maximum: 1 Go                │
└─────────────────────────────────┘
```

#### Politique de Remplacement
```
┌─────────────────────────────────┐
│     Pages les moins utilisées   │
│     sont supprimées du cache    │
└─────────────────────────────────┘
```

### 4. Optimisations 🚀

#### Index en RAM
```
┌─────────────────────────────────┐
│           Index en RAM          │
├─────────────────────────────────┤
│  • Clés primaires               │
│  • Index fréquemment utilisés   │
└─────────────────────────────────┘
```

#### Cache Intelligent
```
┌─────────────────────────────────┐
│         Cache Intelligent       │
├─────────────────────────────────┤
│  • Garde les pages fréquentes   │
│  • Libère les pages rares       │
└─────────────────────────────────┘
```

### 5. Exemple avec un Gros Fichier 📊

```
Fichier DB (1 To)
    │
    ├── En RAM (100 Mo)
    │   ├── En-tête (100 bytes)
    │   ├── Table des matières (4 Ko)
    │   └── Cache (99.99 Mo)
    │       ├── Pages récentes
    │       └── Index actifs
    │
    └── Sur Disque (1 To)
        ├── Données (999.99 Go)
        └── Index secondaires (1 Go)
```

## Organisation du Code 📁

### Structure des Fonctions

```
mydb/
├── db/
│   ├── header.go
│   │   ├── ParseHeader()      // Lit l'en-tête SQLite
│   │   └── ValidateHeader()   // Vérifie la validité
│   │
│   ├── master.go
│   │   ├── ReadMasterTable()  // Lit la table des matières
│   │   └── ParseTableInfo()   // Analyse les infos des tables
│   │
│   ├── index.go
│   │   ├── ReadIndex()        // Lit un index
│   │   ├── SearchIndex()      // Recherche dans l'index
│   │   └── UpdateIndex()      // Met à jour l'index
│   │
│   └── page.go
│       ├── ReadPage()         // Lit une page
│       ├── ParseCells()       // Analyse les cellules
│       └── WritePage()        // Écrit une page
│
└── parser/
    ├── sql.go
    │   ├── ParseSelect()      // Analyse les requêtes SELECT
    │   ├── ParseInsert()      // Analyse les requêtes INSERT
    │   └── ParseCreate()      // Analyse les requêtes CREATE
    │
    └── parser_test.go
        ├── TestParseSelect()  // Tests des SELECT
        └── TestParseInsert()  // Tests des INSERT
```

### Détails des Fonctions Principales

#### 1. Lecture du Fichier (db/)
```
┌─────────────────────────────────┐
│           header.go             │
├─────────────────────────────────┤
│ • ParseHeader()                 │
│   - Lit les 100 premiers bytes  │
│   - Vérifie le magic number     │
│   - Extrait la taille des pages │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│           master.go             │
├─────────────────────────────────┤
│ • ReadMasterTable()             │
│   - Lit la page 1               │
│   - Extrait les infos tables    │
│   - Construit la structure      │
└─────────────────────────────────┘
```

#### 2. Gestion des Index (db/)
```
┌─────────────────────────────────┐
│           index.go              │
├─────────────────────────────────┤
│ • ReadIndex()                   │
│   - Lit la structure B-tree     │
│   - Charge les pages nécessaires│
│                                 │
│ • SearchIndex()                 │
│   - Recherche binaire           │
│   - Navigation dans l'arbre     │
└─────────────────────────────────┘
```

#### 3. Parser SQL (parser/)
```
┌─────────────────────────────────┐
│           sql.go                │
├─────────────────────────────────┤
│ • ParseSelect()                 │
│   - Tokenize la requête         │
│   - Construit l'arbre syntaxique│
│   - Valide la structure         │
│                                 │
│ • ParseInsert()                 │
│   - Extrait les valeurs         │
│   - Vérifie les types           │
└─────────────────────────────────┘
```

### Flux d'Exécution

```
Requête SQL
    │
    ▼
Parser (parser/sql.go)
    │
    ▼
Validation
    │
    ▼
Lecture (db/header.go, master.go)
    │
    ▼
Index (db/index.go)
    │
    ▼
Données (db/page.go)
```

### Tests et Validation

```
┌─────────────────────────────────┐
│         parser_test.go          │
├─────────────────────────────────┤
│ • TestParseSelect()             │
│   - Requêtes simples            │
│   - Conditions WHERE            │
│   - Erreurs de syntaxe          │
│                                 │
│ • TestParseInsert()             │
│   - Insertions simples          │
│   - Valeurs multiples           │
│   - Types de données            │
└─────────────────────────────────┘
```

## Comment Contribuer ? 🤝

1. Fork le projet
2. Créez une branche (`git checkout -b feature/AmazingFeature`)
3. Committez vos changements (`git commit -m 'Add some AmazingFeature'`)
4. Push sur la branche (`git push origin feature/AmazingFeature`)
5. Ouvrez une Pull Request

## Licence 📝
Ce projet est sous licence MIT. Voir le fichier `LICENSE` pour plus de détails.

## Contact 📧
Votre Nom - [@votre_twitter](https://twitter.com/votre_twitter)

Lien du projet : [https://github.com/votre-nom/gosql](https://github.com/votre-nom/gosql)
