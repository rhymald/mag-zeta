#!/bin/bash

counter=0

createObjects () {
  while sleep 0.618; do 
    counter=$(( $counter+1 )) ;
    echo "Here: $counter" ;
    curl localhost:4917/test/spawn > /dev/null 2>&1 ; 
    if [ $? != 0 ]
    then
      exit 0 ;
    fi
  done
}

createObjects
