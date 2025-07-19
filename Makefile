
# start mongoDB container# start mongoDB container
mongodb:
	echo "Starting MongoDB container..."
	docker rm mongodb-test || true
	docker run --name mongodb-test -p 27017:27017 -v .data:/data/db mongodb/mongodb-community-server:latest

mongoDB-short:
	echo "Starting MongoDB container..."
	docker rm mongodb-test || true
	docker run --name mongodb-test -p 27017:27017 mongodb/mongodb-community-server:latest