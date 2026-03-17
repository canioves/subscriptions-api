## Тестовое задание Subscriptions API

### Запуск:
1. Скопировать репозиторий
```bash
git clone github.com/canioves/subscriptions-api
cd subscriptions-api
```
2. Добавить .env файл
``` env
GO_ENV="dev"              #для локального запуска

DB_NAME="db-name"         #имя БД
DB_USER="user"            #имя пользователя 
DB_PASSWORD="pass"        #пароль

DB_HOST_LOCAL="localhost" #локальный хост
DB_HOST_DOCKER="postgres" #хост в docker
DB_PORT="5432"            #порт БД

APP_PORT="8080"           #порт севрера
```
3. Запуск Docker
```bash
docker-compose build
docker-compose up
```

### Документация
Документация реализована через swagger и хранится по эндпоинту `swagger/index.html`
