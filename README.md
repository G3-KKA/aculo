
# Aculo idea

The project was originally planned to accumulate logs with their analysis in the future    
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
- [x] Tests
- [ ] Grafana
- [ ] Kafka Connect + Clickhouse , as an alternative to /batch-inserter service 
- [ ] gRPC-connector
- [ ] send filebatches via FTP 