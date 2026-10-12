## Projet Steam Games

- Go
- Postgresql
- Adminer

Ordre possible des étapes :

X - Setup (Docker + BDD)
X - Import des données Steam (et conversion si PostgreSQL)
O - GET /health + format d’erreur + logs
O - GET /games/:id puis GET /games
O - Pagination et tri
O - Recherche / filtres
O - Users admin + POST /login + middleware JWT
O - Routes d’écriture (POST, PUT, DELETE) + validation
O - CORS, doc Swagger/OpenAI
O - Compléter test.http
O - Compléter README.md

# Commande utiles :

## Lancer le seeding de la bdd avec le fichier steam.json

```cmd
docker compose run --rm api go run ./cmd/import ./uploads/steam.json
```

**NEDELEC Clément & KEUNEBROEK Lucas**
