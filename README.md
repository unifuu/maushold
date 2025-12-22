# Maushold

Tryna make a microservices demo.

## Project Strucrture

```
├── services/
│   ├── player-service/
│   ├── monster-service/
│   ├── battle-service/
│   └── ranking-service/
├── frontend/
|   ├── maushold-react/
├── docker-compose.yml
├── Makefile
```

## Build

``` bash
make setup
make dev
make clear  # Clear all PostgresSQL data
make seed   # Generate seed data
```