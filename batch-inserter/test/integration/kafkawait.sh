#! /bin/bash

var="$(wait-for-port 29099 --timeout=1 --state=inuse)"

if [[ "$var" != "" ]]
then
    exit 1
else
    kafka-topics.sh --bootstrap-server 'localhost:29099' --create --topic 't_healthcheck' ;
    kafka-topics.sh --bootstrap-server 'localhost:29099' --describe --topic 't_healthcheck' ;
    exit $? ; 
fi