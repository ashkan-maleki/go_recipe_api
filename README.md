# go_recipe_api

commands:

docker compose -f docker-compose.yml -f docker-compose.override.yml up -d

docker-compose -f docker-compose.yml -f docker-compose.override.yml down

docker-compose -f docker-compose.yml -f docker-compose.override.yml up -d --build