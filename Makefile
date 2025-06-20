dev:
	docker-compose -f ./docker-compose-dev.yaml down && docker-compose -f ./docker-compose-dev.yaml build && docker-compose -f ./docker-compose-dev.yaml up

prod:
	docker-compose down && docker-compose build && docker-compose up -d
