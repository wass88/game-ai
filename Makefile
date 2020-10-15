ADDR:=game
USER:=root
ADDR:=$(USER)@$(ADDR)
WEB_USER:=web
RSYNC_WEB:=rsync -av --chown=$(WEB_USER):$(WEB_USER)

deploy: build rsync restart-service

build:
	cd gameai-app; make
	cd backend; make bin

rsync:
	ssh $(ADDR) "\
		mkdir -p ~$(WEB_USER)/game-ai/backend &&\
		chown -R $(WEB_USER) ~$(WEB_USER)/game-ai &&\
		chgrp -R $(WEB_USER) ~$(WEB_USER)/game-ai"
	$(RSYNC_WEB) ./backend/bin/ $(ADDR):~$(WEB_USER)/game-ai/backend/bin
	$(RSYNC_WEB) ./backend/configs/ $(ADDR):~$(WEB_USER)/game-ai/backend/configs
	$(RSYNC_WEB) ./backend/migrations/ $(ADDR):~$(WEB_USER)/game-ai/backend/migrations
	$(RSYNC_WEB) ./gameai-app/build/ $(ADDR):~$(WEB_USER)/game-ai/frontend

restart-service:
	ssh $(ADDR) "\
		systemctl restart game-ai-web &&\
		systemctl restart game-ai-kick"

install-service:
	ssh $(ADDR) "\
		cp ~$(WEB_USER)/game-ai/backend/configs/{game-ai-web,game-ai-kick}.service /etc/systemd/system/ &&\
		systemctl daemon-reload &&\
		systemctl enable game-ai-web &&\
		systemctl start game-ai-web &&\
		systemctl enable game-ai-kick &&\3306
		systemctl start game-ai-kick"

install-mysql:
	ssh $(ADDR) "\
		cd /tmp
		wget https://dev.mysql.com/get/mysql-apt-config_0.8.15-1_all.deb &&\
		DEBIAN_FRONTEND=noninteractive dpkg -i mysql-apt-config_0.8.15-1_all.deb &&\
		apt update -y &&\
		DEBIAN_FRONTEND=noninteractive apt install -y mysql-server mysql-client &&\
		echo \"create database gameai\" | mysql &&\
		echo \"create user gameai@localhost IDENTIFIED BY 'goodpassXYZ'\" | mysql &&\
		echo \"grant all on gameai.* to gameai@localhost\" | mysql"

install-nginx:
	ssh $(ADDR) "\
		cp ~$(WEB_USER)/game-ai/backend/configs/nginx.conf /etc/nginx/nginx.conf &&\
		systemctl restart nginx"

add-user:
	ssh $(ADDR) "\
	useradd $(WEB_USER) &&\
	mkdir -p /home/$(WEB_USER) &&\
	chown $(WEB_USER) ~$(WEB_USER) &&\
	chgrp $(WEB_USER) ~$(WEB_USER)"

install-required:
	ssh $(ADDR) "\
		apt install -y rsync nginx"

install-docker:
	ssh $(ADDR) "\
		apt install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common &&\
		curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add - &&\
		add-apt-repository \"deb [arch=amd64] https://download.docker.com/linux/debian $$(lsb_release -cs) stable\" &&\
		apt update &&\
		apt install -y docker-ce docker-ce-cli containerd.io &&\
		gpasswd -a web docker"

install-golang:
	ssh $(ADDR) "\
		cd /tmp &&\
		wget https://golang.org/dl/go1.15.3.linux-amd64.tar.gz &&\
		tar -C /usr/local -xzf go1.15.3.linux-amd64.tar.gz"

install-migrate:
	ssh $(ADDR) "\
		cd /tmp &&\
		curl -L https://github.com/golang-migrate/migrate/releases/download/v4.11.0/migrate.linux-amd64.tar.gz | tar xvz &&\
		mv ./migrate.linux-amd64 /usr/bin/migrate"

migrate:
	ssh $(ADDR) "\
		migrate -database 'mysql://gameai:goodpassXYZ@tcp(127.0.0.1:3306)/gameai'\
		-path ~$(WEB_USER)/game-ai/backend/migrations -verbose up"

.PHONY: deploy
