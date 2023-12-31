build-all:
	cd checkout && GOOS=linux GOARCH=amd64 make build && echo "checkout done!"
	cd loms && GOOS=linux GOARCH=amd64 make build && echo "loms done!"
	cd notifications && GOOS=linux GOARCH=amd64 make build && echo "notifications done!"

run-all: build-all
	# sudo docker compose up --force-recreate --build
	mkdir -p checkout_pgdata loms_pgdata notifications_pgdata
	docker volume create checkout_pgdata
	docker volume create loms_pgdata
	docker volume create notifications_pgdata
	docker-compose up --force-recreate --build

down-all:
	docker-compose down

precommit:
	cd checkout && make precommit
	cd loms && make precommit
	cd notifications && make precommit