ADDR:=game
USER:=root
ADDR:=$(USER)@$(ADDR)
WEB_USER:=web
RSYNC_WEB:=rsync -av --chown=$(WEB_USER):$(WEB_USER)

build:
	cd gameai-app; make
	cd backend; make bin

deploy:
	ssh $(ADDR) "\
		mkdir -p ~$(WEB_USER)/game-ai/backend &&\
		chown -R $(WEB_USER) ~$(WEB_USER)/game-ai &&\
		chgrp -R $(WEB_USER) ~$(WEB_USER)/game-ai"
	$(RSYNC_WEB) ./backend/bin/ $(ADDR):~$(WEB_USER)/game-ai/backend/bin
	$(RSYNC_WEB) ./backend/configs/ $(ADDR):~$(WEB_USER)/game-ai/backend/configs
	$(RSYNC_WEB) ./backend/migrations/ $(ADDR):~$(WEB_USER)/game-ai/backend/migrations
	$(RSYNC_WEB) ./gameai-app/build/ $(ADDR):~$(WEB_USER)/game-ai/frontend

install-service:
	ssh $(ADDR) "\
		cp ~$(WEB_USER)/game-ai/backend/configs/{game-ai-web,game-ai-kick}.service /etc/systemd/system/ &&\
		systemctl daemon-reload &&\
		systemctl enable game-ai-web &&\
		systemctl start game-ai-web &&\
		systemctl enable game-ai-kick &&\
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

.PHONY: deploy