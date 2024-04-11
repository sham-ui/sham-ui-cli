start_openobserv:
	mkdir -p ${PWD}/data/openobserve/
	docker run -v ${PWD}/data/openobserve:/data -e ZO_DATA_DIR="/data" -p 5080:5080 -e ZO_ROOT_USER_EMAIL="root@example.com" -e ZO_ROOT_USER_PASSWORD="pass" public.ecr.aws/zinclabs/openobserve:latest

open_openobserv:
	xdg-open http://localhost:5080

