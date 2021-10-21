docker pull mysql
docker-compose -f docker-compose.mysql.yml up -d
sleep 25s
docker exec -it mysqlTestDB bash -c "mysql -uroot -ppa55w0rd test_db < /test_db/test_db.sql"
echo "done"
