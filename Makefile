maria: 
	docker run -p 127.0.0.1:3306:3306 --name some-mariadb \
	-e MARIADB_ROOT_PASSWORD=my-secret-pw -e MARIADB_DATABASE=myapp -d mariadb:latest


image:
	docker build -t todo:test -f Dockerfile .

container:
	docker run -p 8081:8081 --env-file ./local.env --link some-mariadb:db \
	--name myapp todo:test
