ADDR:=game
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
	$(RSYNC_WEB) ./backend/bin/ $(ADDR):~$(WEB_USER)/game-ai/backend/bin
	$(RSYNC_WEB) ./backend/configs/ $(ADDR):~$(WEB_USER)/game-ai/backend/configs
	$(RSYNC_WEB) ./backend/migrations/ $(ADDR):~$(WEB_USER)/game-ai/backend/migrations
	$(RSYNC_WEB) ./gameai-app/build/ $(ADDR):~$(WEB_USER)/game-ai/frontend

restart-service:
	ssh $(ADDR) "\
		systemctl restart game-ai-web &&\
		systemctl restart game-ai-kick"

init-deploy: add-user rsync install-mysql install-nginx install-docker install-migrate install-service

add-user:
	ssh $(ADDR) "\
	(sudo useradd $(WEB_USER) || true) &&\
	sudo mkdir -p /home/$(WEB_USER) &&\
	sudo chown $(WEB_USER) ~$(WEB_USER) &&\
	sudo chgrp $(WEB_USER) ~$(WEB_USER)"

install-mysql:
	ssh $(ADDR) "\
		cd /tmp &&\
		sudo yum localinstall -y https://dev.mysql.com/get/mysql80-community-release-el7-3.noarch.rpm &&\
		sudo yum install -y mysql-community-client"

init-database:
	ssh $(ADDR) "\
		echo \"create database gameai\" | sudo mysql -paaaaaaaa&&\
		echo \"create user gameai@localhost IDENTIFIED BY 'goodpassXYZ'\" | sudo mysql -paaaaaaaa&&\
		echo \"grant all on gameai.* to gameai@localhost\" | sudo mysql -paaaaaaaa"

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
		sudo systemctl start game-ai-web &&\
		sudo systemctl enable game-ai-kick &&\
		sudo systemctl start game-ai-kick"

install-docker:
	ssh $(ADDR) "\
		sudo amazon-linux-extras install docker &&\
		sudo gpasswd -a web docker"

install-golang:
	ssh $(ADDR) "\
		cd /tmp &&\
		wget https://golang.org/dl/go1.15.3.linux-amd64.tar.gz &&\
		sudo tar -C /usr/local -xzf go1.15.3.linux-amd64.tar.gz"

install-migrate:
	ssh $(ADDR) "\
		cd /tmp &&\
		curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-armv7.tar.gz | tar xvz &&\
		sudo mv ./migrate.linux-armv7 /usr/bin/migrate"

migrate:
	ssh $(ADDR) "\
		sudo migrate -database 'mysql://gameai:goodpassXYZ@tcp(127.0.0.1:3306)/gameai'\
		-path ~$(WEB_USER)/game-ai/backend/migrations -verbose up"

.PHONY: deploy
