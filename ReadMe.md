* Цель удобное поздравление сотрудников
* Получение списка сотрудников любым способом(api/ad ldap/прямая регистрация)
* Авторизация (jwt)
* Возможность подписаться/отписаться от оповещения о дне рождения
* Оповещение о ДР того, на кого подписан
* Внешнее взаимодействие (тг бот https://t.me/bday_notifications_bot)

Телеграм бот отправляет сообщения о Днях Рождения каждый день в 10:00 по МСК.
**Перед запуском необходимо создать базу данных самостоятельно**. В файле [.env](env/.env) можно изменить порт и конфигурацию базы данных.


### Примеры запросов
Регистрация
```shell
 curl -X POST http://localhost:8080/sign_in -H "Content-Type: application/json" -d '{"username": "user1", "password": "password1"}' -i
```

Войти в аккаунт (необходимо после регистрации и для всех других запросов)
```shell
curl -X POST http://localhost:8080/sign_in -H "Content-Type: application/json" -d '{"username": "user1", "password": "password1"}' -i 
```

Добавить информацию о себе
```shell
curl -X POST http://localhost:8080/api/info \
  -H "Content-Type: application/json" \
  -H "Authorization: YOUR_TOKEN" \
  -d '{
        "userId": 5,
        "name": "John",
        "surname": "Doe",
        "birth":  "2024-06-27T00:00:00Z"
      }'
```

Получить список работников
```shell
curl -X GET http://localhost:8080/api/employees -H "Authorization: YOUR_TOKEN" -i
```

Подписаться
```shell
curl -X POST http://localhost:8080/api/subscribe \
  -H "Content-Type: application/json" \
  -H "Authorization: YOUR_TOKEN" \
  -d '{"id": 1, "subscribeTo": [2, 3, 4]}'
```

Отписаться
```shell
curl -X POST http://localhost:8080/api/unsubscribe \
  -H "Content-Type: application/json" \
  -H "Authorization: YOUR_TOKEN" \
  -d '{"id": 1, "subscribeTo": [2, 3, 4]}'
```

