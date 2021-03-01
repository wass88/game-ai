ADDR:=game
DB_ADDR:=game-db

WEB_USER:=web
RSYNC:=rsync -av1
RSYNC_WEB:=rsync -av --chown=$(WEB_USER):$(WEB_USER) --rsync-path="sudo rsync"

deploy: build rsync restart-service

build:
	cd gameai-app; make
	cd backend; make bin

rsync:
	ssh $(ADDR) "\
		sudo mkdir -p ~$(WEB_USER)/game-ai/backend/.data/ai-docker &&\
		sudo chown -R $(WEB_USER) ~$(WEB_USER)/game-ai &&\
		sudo chgrp -R $(WEB_USER) ~$(WEB_USER)/game-ai "
	$(RSYNC_WEB) ./backend/bin ./backend/configs ./backend/migrations \
	             $(ADDR):~$(WEB_USER)/game-ai/backend
	$(RSYNC_WEB) ./gameai-app/build/ $(ADDR):~$(WEB_USER)/game-ai/frontend

restart-service:
	ssh $(ADDR) "\
		sudo systemctl restart game-ai-web &&\
		sudo systemctl restart game-ai-kick"

init-deploy: add-user rsync install-nginx install-docker install-service

add-user:
	ssh $(ADDR) "\
	(id -u ${WEB_USER} &>/dev/null || sudo useradd $(WEB_USER)) &&\
	sudo mkdir -p /home/$(WEB_USER) &&\
	sudo chown $(WEB_USER) ~$(WEB_USER) &&\
	sudo chgrp $(WEB_USER) ~$(WEB_USER) &&\
	sudo chmod 711 /home/${WEB_USER}"

install-nginx:
	ssh $(ADDR) "\
		sudo amazon-linux-extras install -y nginx1 &&\
		sudo cp ~$(WEB_USER)/game-ai/backend/configs/nginx.conf /etc/nginx/nginx.conf &&\
		sudo systemctl restart nginx"

install-service:
	ssh $(ADDR) "\
		sudo cp ~$(WEB_USER)/game-ai/backend/configs/{game-ai-web,game-ai-kick}.service /etc/systemd/system/ &&\
		sudo systemctl daemon-reload &&\
		sudo systemctl enable game-ai-web &&\
		sudo systemctl enable game-ai-kick"

start-service:
	ssh $(ADDR) "\
		sudo systemctl start game-ai-web &&\
		sudo systemctl start game-ai-kick"

install-docker:
	ssh $(ADDR) "\
		sudo amazon-linux-extras install docker &&\
		sudo gpasswd -a web docker &&\
		sudo systemctl start docker"

install-golang:
	ssh $(ADDR) "\
		cd /tmp &&\
		wget https://golang.org/dl/go1.15.3.linux-amd64.tar.gz &&\
		sudo tar -C /usr/local -xzf go1.15.3.linux-amd64.tar.gz"

init-db: add-user-db rsync-db install-mysql install-migrate

rsync-db:
	ssh $(DB_ADDR) "\
		sudo mkdir -p ~$(WEB_USER)/game-ai/backend/.data/ai-docker &&\
		sudo chown -R $(WEB_USER) ~$(WEB_USER)/game-ai &&\
		sudo chgrp -R $(WEB_USER) ~$(WEB_USER)/game-ai "
	$(RSYNC_WEB) ./backend/configs ./backend/migrations \
	             $(DB_ADDR):~$(WEB_USER)/game-ai/backend

add-user-db:
	ssh $(DB_ADDR) "\
	(id -u ${WEB_USER} &>/dev/null || sudo useradd $(WEB_USER)) &&\
	sudo mkdir -p /home/$(WEB_USER) &&\
	sudo chown $(WEB_USER) ~$(WEB_USER) &&\
	sudo chgrp $(WEB_USER) ~$(WEB_USER) &&\
	sudo chmod 711 /home/${WEB_USER}"

install-mysql:
	ssh $(DB_ADDR) "\
		cd /tmp &&\
		sudo yum localinstall -y https://dev.mysql.com/get/mysql80-community-release-el7-3.noarch.rpm &&\
		sudo yum install -y mysql-community-client &&\
		sudo yum install -y yum install mysql-server &&\
		sudo systemctl start mysqld"

DB_PASS_FILE=./mysql-root
NEW_DB_PASSWORD= $(shell cat $(DB_PASS_FILE))
DB_WEB_PASS_FILE=./mysql-web
DB_WEB_PASSWORD= $(shell cat $(DB_WEB_PASS_FILE))
init-database:
	ssh $(DB_ADDR) "\
		DB_PASSWORD=\$$(sudo grep 'A temporary password is generated' /var/log/mysqld.log | sed -e 's/.*root@localhost: //'); \
		mysql -uroot \"-p\$${DB_PASSWORD}\" --connect-expired-password -e \"ALTER USER 'root'@'localhost' IDENTIFIED BY '${NEW_DB_PASSWORD}'; flush privileges;\" "
	ssh $(DB_ADDR) " \
		echo \" \
			create database gameai; \
			create user gameai@'172.31.%.%' IDENTIFIED BY '${DB_WEB_PASSWORD}'; \
			grant all on gameai.* to gameai@'172.31.%.%'; \
		\" | mysql -uroot --connect-expired-password \"-p${NEW_DB_PASSWORD}\" "

install-migrate:
	ssh $(DB_ADDR) "\
		cd /tmp &&\
		curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-armv7.tar.gz | tar xvz &&\
		sudo mv ./migrate.linux-armv7 /usr/bin/migrate"

migrate:
	ssh $(DB_ADDR) "\
		sudo migrate -database 'mysql://gameai:${DB_WEB_PASSWORD}@tcp(127.0.0.1:3306)/gameai'\
		-path ~$(WEB_USER)/game-ai/backend/migrations -verbose up"

.PHONY: deploy
