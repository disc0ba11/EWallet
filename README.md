# Запуск
Скачайте последнюю версию программы в разделе [Releases](https://github.com/disc0ba11/EWallet/releases).

Установите пакет [curl](https://github.com/curl/curl):

	Arch: sudo pacman -S curl
	Ubuntu: sudo apt install curl
	Debian: sudo apt-get install curl
<!-- Программа написана на языке golang, поэтому для её запуска Для её компиляции у вас должны быть установлены:
1. Пакет go: https://go.dev/doc/install
2. Библиотека go-sqlite3: https://github.com/mattn/go-sqlite3 -->

Прежде, чем запустить программу, следует выполнить ряд действий:
1. Разместить программу в отдельной папке.
2. Убедиться, что у вас открыт и не занят порт 8080.
3. Проверить, какой у вас браузер (например, браузер на основе Chromium не распарсит данные в JSON-формате, а Mozilla Firefox — распарсит).

Для запуска в ОС Windows дважды щёлкните по исполняемому файлу.

Для запуска в Linux воспользуйтесь командой:

	./EWallet
Пример:

При первом запуске, если у вас в папке нет БД с названием 'wallets.db', программа автоматически её создаст и заполнит десятью кошельками со случайными адресами. После запуска программы вы увидите список кошельков в БД: их ID в базе, адреса и балансы. Если в базе данных на момент запуска было больше десяти кошельков, то она автоматически уменьшит их количество до десяти, если меньше — увеличит до десяти. Отсутствие сообщения об ошибке означает, что программа запущена и работает корректно, можно переходить к эксплуатации.

# Эксплуатация
Программа предполагает три метода: getBalance, send и getLast.
## GetBalance
В этом методе реализована функция просмотра баланса кошелька. Эндпоинт: `GET /api/wallet/адрес_кошелька/balance`. Для просмотра баланса кошелька запустите программу и в адресной строке браузера введите следующий запрос:

	http://localhost:8080/api/wallet/адрес_кошелька/balance
Пример:

## Send
Данный метод позволяет перевести деньги с одного счёта на другой. Эндпоинт: `POST /api/send`. Для перевода денег с одного счёта на другой необходимо отправить POST-запрос при помощи curl или другого HTTP-клиента. Тело запроса обладает **строгим** шаблоном: `{"From": "адрес_первого_кошелька", "To": "адрес_второго_кошелька", "Amount": "количество"}`. Данный формат является JSON-объектом, и именно его метод Send принимает на вход.

	curl -d `{"From": "адрес_первого_кошелька", "To": "адрес_второго_кошелька", "Amount": "количество"}` http://localhost:8080/api/send

Пример:

## GetLast
Последний метод позволяет при помощи браузера посмотреть последние N запросов в формате JSON. Эндпоинт: `GET /api/transactions?count=N`. В адресной строке браузера необходимо ввести следующий запрос:

	http://localhost:8080/api/transactions?count=количество_запросов
Пример:

# Компиляция
Для компиляции программы вам будет необходим пакет [Go](https://github.com/golang/go).

	Arch: sudo pacman -S go
	Ubuntu: sudo apt install golang
Для установки Go в дистрибутиве Debian воспользуйтесь руководством в разделе Linux на [официальном сайте](https://go.dev/doc/install).

Затем вам понадобится библиотека [go-sqlite3](https://github.com/mattn/go-sqlite3). Для её установки сначала установите [gcc](https://gcc.gnu.org/):

	Arch: sudo pacman -S gcc
	Ubuntu: sudo apt -y install build-essential
	Debian: sudo apt install build-essential
Для Debian так же можно установить man-страницы, включающие документацию для GNU/Linux:

	sudo apt-get install manpages-dev
Затем скачайте библиотеку go-sqlite3:

	go get github.com/mattn/go-sqlite3
И установите:

	go install github.com/mattn/go-sqlite3
После этого можно будет скомпилировать программу при помощи команды:

	go build -o EWallet