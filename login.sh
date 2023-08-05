#!/bin/bash

counter=0

createObjects () {
  while sleep 1; do 
    counter=$(( $counter+1 )) ;
    echo "Here: $counter" ;
    curl localhost:4917/login > /dev/null 2>&1 ; 
    if [ $? != 0 ]
    then
      exit 0 ;
    fi
  done
}

createObjects
