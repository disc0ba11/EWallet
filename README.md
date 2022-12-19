EWallet

Программа написана на языке golang
Для её компиляции у вас должны быть установлены:
1. Пакет go: https://go.dev/doc/install
2. Библиотека go-sqlite3: https://github.com/mattn/go-sqlite3

Прежде, чем запустить программу, следует выполнить ряд действий:
1. Разместить программу в отдельной папке
2. Убедиться, что у вас открыт и не занят порт 8080
3. Проверить, какой у вас браузер (например, браузер на основе Chromium не может отображать на странице данные в JSON-формате, а Mozilla Firefox — может)

При первом запуске, если у вас в папке нет БД с названием 'wallets.db', программа автоматически её создаст и заполнит десятью кошельками со случайными адресами
После запуска программы вы увидите список кошельков в БД: их ID в базе, адреса и балансы
Если в базе данных на момент запуска было больше десяти кошельков, то она автоматически уменьшит их количество до десяти, если меньше — увеличит до десяти
Отсутствие сообщения об ошибке означает, что программа запущена и работает корректно, можно переходить к эксплуатации

Программа предполагает три метода:
1. getBalance, просмотр баланса кошелька (эндпоинт: GET /api/wallet/адрес_кошелька/balance)
2. send, транзакция средств между двумя кошельками (эндпоинт: POST /api/send)
		Метод send принимает POST-запрос с JSON-объектом, пример использования:
		curl -d '{"From": "29b0223beea5f4f74391f445d15afd4294040374f6924b98cbf8713f8d962d7c", "To": "8d019192c24224e2cafccae3a61fb586b14323a6bc8f9e7df1d929333ff99393", "Amount": "0.50"}' http://localhost:8080/api/send
3. getLast, просмотр последних N транзакций (эндпоинт: GET /api/transactions?count=N)