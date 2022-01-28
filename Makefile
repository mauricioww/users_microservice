USER_PATH = './user_srv/userpb/'
DETAILS_PATH = './user_details_srv/detailspb/'
SPECS_PATH = 'server/swagger/specs/'

build:
	docker-compose up --build

clean_docker: 
	docker rm $(shell docker ps -a -q) && \
	docker image rm http_server:2.0 user_server:2.0

update_user:
	cd $(USER_PATH) && \
		buf mod update

run: user
	cd server && go run main.go

user: rm_user
	cd $(USER_PATH) && \
		buf generate
		
rm_user:
	cd $(SPECS_PATH) && \
		rm user.swagger.json