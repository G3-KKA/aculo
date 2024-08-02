
# Aculo idea

The project was originally planned to accumulate logs with their analysis    
But after learning more about [ELK stack](https://gitinsky.com/elkstack) and comparing it with my idea -- purpose of project changed  
Now it is general purpose [ELT](https://habr.com/ru/articles/695546/) system, that may be used for logging, but not limited by them   

---

My alternative stack to ELK consists of:
- Kafka 
- ClickHouse 
- Grafana
- And a bit of intermediate sevices 

---

Tasks

- [x] RESTAPI interface
- [x] Swagger for REST interfaces
- [ ] Tests
- [ ] Grafana
- [ ] Direct sending logs and event into Kafka via
- [ ] Kafka Connect + Clickhouse , to eleminate /batch-inserter service 

---

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