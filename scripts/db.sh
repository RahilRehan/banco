#!/bin/sh

# shellcheck disable=SC2039
GREEN=$'\e[0;32m'
#RED=$'\e[0;31m'
NC=$'\e[0m'

# shellcheck disable=SC2039
source app.env

echo "Create root user for postgres: "; read -r DB_ROOT_USER; export DB_ROOT_USER
echo "Create root password for postgres: "; read -r DB_ROOT_PASSWORD; export DB_ROOT_PASSWORD
echo "Create password for normal user: "; read -r DB_PASSWORD; export DB_PASSWORD

echo "${GREEN}Brining up postgres docker container...${NC}"
docker run --name postgres-banco -p 5432:5432 -v banco-data:/var/lib/postgresql/data -e POSTGRES_USER="${DB_ROOT_USER}" -e POSTGRES_PASSWORD="${DB_ROOT_PASSWORD}" -d postgres

echo "${GREEN}Waiting for postgres container to boot...${NC}"
while ! curl "http://localhost:5432" 2>&1 | grep '52' >> /dev/null; do sleep 2; done
sleep 1

echo "${GREEN}Creating user and db...${NC}"
docker exec -it postgres-banco psql -U "${DB_ROOT_USER}" -c "create user ${DB_USER} password '${DB_PASSWORD}';"
docker exec -it postgres-banco psql -U "${DB_ROOT_USER}" -c "create database ${DB_NAME} owner=${DB_USER};"