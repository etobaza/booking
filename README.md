## Требования

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Make](https://www.gnu.org/software/make/)

## Запуск проекта

1. **Клонируйте репозиторий:**

   ```bash
   git clone https://github.com/your-username/your-repo.git
   cd your-repo
   ```

2. **Запустите контейнеры:**

   ```bash
   make up
   ```

3. **Запуск тестов:**

   Чтобы запустить тесты, выполните следующую команду:

   ```bash
   make test
   ```

4. **Остановка сервисов:**

   Если вы хотите вручную остановить сервисы, напишите:

   ```bash
   make down
   ```

## Команды Makefile

- `make up`: Поднимает все сервисы в фоновом режиме.
- `make test`: Поднимает сервисы, запускает тесты и завершает работу сервисов.
- `make down`: Завершает работу всех сервисов.

## Конфигурация базы данных

База данных PostgreSQL будет автоматически настроена и доступна по следующим параметрам:

- **Хост:** `localhost`
- **Порт:** `5432`
- **Имя базы данных:** `booking`
- **Пользователь:** `user`
- **Пароль:** `password`

Эти параметры можно изменить в файле `docker-compose.yml`.
