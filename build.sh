#!/bin/bash
docker-compose -f ./geoservice/docker-compose.yml up --force-recreate --build  -d # Запуск сервиса geo
docker-compose -f ./authservice/docker-compose.yml up --force-recreate --build  -d # Запуск сервиса auth
docker-compose -f ./userservice/docker-compose.yml up --force-recreate --build  -d # Запуск сервиса user
docker-compose -f ./proxy/docker-compose.yml up --force-recreate --build  -d # Запуск сервиса proxy
docker network connect mynetwork2 proxy
docker network connect mynetwork2 geoservice
docker network connect mynetwork2 authservice
docker network connect mynetwork2 userservice