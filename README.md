# How to deploy and use

In development!

# Architecture 

In development! (image)

# Aculo idea

The project focused on Accumulating Logs with their analysis in the future    
It is an alternative to the [ELK stack](https://gitinsky.com/elkstack) 

It is an [ELT](https://habr.com/ru/articles/695546/) system, so feel free to throw unsorted mess into logs, we'll deal with that later in the process and not on the client side 

---

My alternative stack to ELK consists of:
- Kafka 
- ClickHouse 
- Grafana
- And a bit of intermediate sevices 

---

Project partially follows principles of TDD and Clean Architecture in places where i find it appropriate

---
Phase 0:
- [x] All parts have similar deploy and management parts (config//logging//docker images//structure)
- [x] Client API 
- [x] gRPC and HTTP master interface
- [x] Tests 
- [ ] Master--Stream Cluster logic on master side 
- [ ] Tests
- [ ] Swagger for REST interfaces
- [ ] Stream Accumulator Cluster
- [ ] Tests
- [ ] Ready to try 
- [ ] e2e tests
- [ ] Ready to use 

Phase 1:

- [ ] Internal logging into 
- [ ] Tests
- [ ] Kafka tweaks to increase throughput in default configuration
- [ ] Benchmarks
- [ ] Clickhouse tweaks  to increase throughput in default configuration
- [ ] Benchmarks
- [ ] Grafana
- [ ] send logs via gRPC stream
- [ ] send log files via FTP 