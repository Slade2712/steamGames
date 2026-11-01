# Projet Steam Games

Mini-projet d'**API REST en Go** pour consulter et gérer un catalogue de jeux Steam, stocké dans **PostgreSQL**.

**Technologies :** Go · PostgreSQL · Docker Compose · Adminer · JWT · Swagger

## Roadmap

- [x] Setup : Docker et base de données PostgreSQL
- [x] Import des données Steam depuis un fichier JSON
- [x] `GET /health`, format des réponses et logs
- [x] `GET /games/{id}` et `GET /games`
- [x] Pagination et tri
- [x] Recherche et filtres
- [x] Utilisateurs, connexion (`POST /login`) et authentification JWT
- [x] Routes d'écriture (`POST`, `PUT`, `DELETE /games`) et validation
- [x] CORS et documentation Swagger / OpenAPI
- [x] Compléter les requêtes de test HTTP
- [x] Compléter le README

## Lancement

Créer le fichier `.env` à partir de `.env.example`, puis renseigner les variables de connexion à la base et `JWT_SECRET`.

```bash
cp .env.example .env
docker compose up -d
```

## Commandes utiles

**Importer les jeux depuis `uploads/steam.json` :**

```bash
docker compose run --rm api go run ./cmd/import ./uploads/steam.json
```

**Afficher les logs :**

```bash
docker compose logs -f api
```

**Arrêter le projet :**

```bash
docker compose down
```

## Logs

Les requêtes HTTP reçues par l'API sont journalisées dans le dossier `logs/` du projet. Les fichiers sont rangés par date et par créneau horaire, par exemple : `logs/2026/11/01/17-50/`.

## API

Les routes utilisent le préfixe `/api/v1` par défaut.

| Méthode     | Route                     | Accès  |
| ----------- | ------------------------- | ------ |
| GET         | `/health`                 | Public |
| GET         | `/games` et `/games/{id}` | Public |
| POST        | `/register` et `/login`   | Public |
| GET         | `/users/me`               | JWT    |
| GET         | `/users`                  | Admin  |
| POST        | `/games`                  | Admin  |
| PUT, DELETE | `/games/{id}`             | Admin  |

La documentation interactive est accessible sur `/swagger/`.

---

**NEDELEC Clément & KEUNEBROEK Lucas**
