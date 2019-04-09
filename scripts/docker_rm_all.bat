@echo off
FOR /F "usebackq" %%a IN (`docker ps -a -q`) DO (
	docker rm --force %%a
)
FOR /F "usebackq" %%a IN (`docker images -a -q`) DO (
	docker rmi --force %%a
)