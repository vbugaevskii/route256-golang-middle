build-all:
	cd checkout && GOOS=linux GOARCH=amd64 make build && echo "checkout done!"
	cd loms && GOOS=linux GOARCH=amd64 make build && echo "loms done!"
	cd notifications && GOOS=linux GOARCH=amd64 make build && echo "notifications done!"

run-all: build-all
	# sudo docker compose up --force-recreate --build
	mkdir -p checkout_pgdata loms_pgdata
	docker-compose up --force-recreate --build

precommit:
	cd checkout && make precommit
	cd loms && make precommit
	cd notifications && make precommit