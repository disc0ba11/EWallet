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

Для запуска в ОС Windows скачайте EWallet.exe и дважды щёлкните по исполняемому файлу.

Для запуска в Linux скачайте EWallet и воспользуйтесь командой:

	./EWallet
Пример:

![launch_example](https://user-images.githubusercontent.com/114620640/208523135-f1246ed4-6a8a-4325-819e-62d1bc44f045.png)

При первом запуске, если у вас в папке нет БД с названием 'wallets.db', программа автоматически её создаст и заполнит десятью кошельками со случайными адресами. После запуска программы вы увидите список кошельков в БД: их ID в базе, адреса и балансы. Если в базе данных на момент запуска было больше десяти кошельков, то она автоматически уменьшит их количество до десяти, если меньше — увеличит до десяти. Отсутствие сообщения об ошибке означает, что программа запущена и работает корректно, можно переходить к эксплуатации.

# Эксплуатация
Программа предполагает три метода: `getBalance`, `send` и `getLast`.

## GetBalance
В этом методе реализована функция просмотра баланса кошелька. Эндпоинт: `GET /api/wallet/адрес_кошелька/balance`. Для просмотра баланса кошелька запустите программу и в адресной строке браузера введите следующий запрос:

	http://localhost:8080/api/wallet/адрес_кошелька/balance
Пример:

![getBalance_example](https://user-images.githubusercontent.com/114620640/208523228-eafb0570-56bd-4ffe-b3dc-50efd8096a7a.png)


## Send
Данный метод позволяет перевести деньги с одного счёта на другой. Эндпоинт: `POST /api/send`. Для перевода денег с одного счёта на другой необходимо отправить POST-запрос при помощи curl или другого HTTP-клиента. Тело запроса обладает **строгим** шаблоном: `{"From": "адрес_первого_кошелька", "To": "адрес_второго_кошелька", "Amount": "количество"}`. Данный формат является JSON-объектом, и именно его метод Send принимает на вход.

	curl -d `{"From": "адрес_первого_кошелька", "To": "адрес_второго_кошелька", "Amount": "количество"}` http://localhost:8080/api/send

Пример:

![send_example](https://user-images.githubusercontent.com/114620640/208523247-af785941-ae7e-4f6b-99d2-087feae40595.png)


## GetLast
Последний метод позволяет при помощи браузера посмотреть последние N транзакций в формате JSON. Эндпоинт: `GET /api/transactions?count=N`. В адресной строке браузера необходимо ввести следующий запрос:

	http://localhost:8080/api/transactions?count=количество_запросов
Пример:

![image](https://user-images.githubusercontent.com/114620640/208523407-1d77b732-a9e8-4b82-90bb-68c52eacec97.png)


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
Затем инициализируйте go.mod в склонированном репозитории:

	go mod init
	go mod tidy
И установите зависимость:

	go install
После этого можно будет скомпилировать программу при помощи команды:

	go build -o EWallet

# Docker
Вы можете скомпилировать и запустить программу в Docker-контейнере. Для этого склонируйте репозиторий и запустите в нём следующие команды:

	sudo docker build -t ewallet-build -f Dockerfile.build .
	sudo docker run ewallet-build > EWallet.tar.gz
Будет создан архив `EWallet.tar.gz`, содержащий скомпилированную программу. На этом этапе её можно разархивировать и запустить, если в запуске из Docker-контейнера нет необходимости. Если необходимость есть, воспользуйтесь следующими командами:

	sudo docker build -t ewallet .
	sudo docker run --publish 80:8080 --name ewallet --rm ewallet
Программа будет запущена из Docker-контейнера. Флаг `--publish` пробрасывает порт 8080 из контейнера на 80 в хост-системе. Поэтому в дальнейшем для эксплуатации программы следует использовать порт 80, а не 8080.
