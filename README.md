# courseProjectDB 

Command to pull PosrgreSQL server Image:
 
```
docker pull postgres
```
Enter the command "docker run," which is used to launch container: 
```
docker run --rm   --name pg-docker -e POSTGRES_PASSWORD=password -d -p 5431:5432 -v $HOME/docker/volumes/postgres:/var/lib/postgresql/data  postgres
 
```
Start PostgreSQL server:
```
psql -h 0.0.0.0 -U postgres -d postgres --port=5431

```
It's time to create all tables and triggers (init.sql).

Pull Redis server Image and launch container:
```
docker pull redis
sudo docker run --restart=always -d --name redis_1    -v /opt/redis/etc/redis.conf:/usr/local/etc/redis/redis.conf    -v /opt/redis/data:/data    -p 127.0.0.1:6380:6379 redis redis-server /usr/local/etc/redis/redis.conf 
 
```
Run:
```
go run main.go
```
The answer should be: "Successfully connected!"
