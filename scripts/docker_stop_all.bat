@echo off
FOR /F "usebackq" %%a IN (`docker ps -a -q`) DO (
	docker stop %%a
)
