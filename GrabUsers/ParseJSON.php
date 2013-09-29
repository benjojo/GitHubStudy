<?php

$f = fopen( 'php://stdin', 'r' );

while( $line = fgets( $f ) ) {
  $a = json_decode($line);
  if(isset($a->actor)) echo($a->actor."\n");
}