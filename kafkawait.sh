#! /bin/bash

var="$(wait-for-port 29095 --timeout=1 --state=inuse)"

if [[ "$var" != "" ]]
then
    exit 1
else
    kafka-topics.sh --bootstrap-server '127.0.0.1:29095' --create --topic 't_healthcheck' ;
    kafka-topics.sh --bootstrap-server '127.0.0.1:29095' --describe --topic 't_healthcheck' ;
    exit $? ; 
fi