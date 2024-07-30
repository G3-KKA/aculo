app1{
    log.Info(Event happened){=>

        (l *logger).Info(data){
        data = structurise(data)

        loggger.Kafka1
        .ToTopik("structured logs")
        .Send()
        }
    }
}
kafka1{
    topics[
        "structured"{

        }
    ]    
    select {
        event := <-structuredEvents:
            topics("structured").Append(event)
    }
}
topics_parser_service1{
    lis := listen.Kafka.Topic("structured")
    select{
        event := <-lis:
            event = parse(event)
        .Kafka2
        .ToTopik("clickhouse_ready")
        .Send()
    }

}
kafka2{
    topics[
        "clickhouse_ready"{

        }
    ]    
    select {
        event := <-fromParserEvents:
            topics("clickhouse_ready").Append(event)
    }
    IF topic("clickhouse_ready").Fullnes > 80% {
        buffer := topic("clickhouse_ready")
        chConnect.FlushBuffer(buffer.All())
        
    }
}