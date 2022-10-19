# wheeler-trader

The end users of the Wheeler Trader application are the workers in car showrooms. Our application does not yet have a UI (it currently exposes REST APIs), but here one may find schematic diagrams.

What can be done through the application is to display all vehicles, vehicles for sale only, for rent only, as well as vehicles that are neither for sale nor for rent. From the sections "All Vehicles" and "Unassigned Vehicles" we can add a new vehicle, and from "Sale" and "Rent" we can set a vehicle for sale or rent. Each section supports entity filtering, updating and deleting entities.


## How to run

Docker Desktop is required to run the application. Execute the following commands:
- docker pull postgres:14.3			(download prebuilt postgres image)
- docker run -d -e POSTGRES_PASSWORD=postgres --name wheeler-trader_postgres postgres:14.3	(start a container using the postgres image)
- docker inspect wheeler-trader_postgres	(get the IP of the container "NetworkSettings" -> "IPAddress" -> "172.17.0.2" (example value))
- cd server/migrations
- docker build . -t wheeler-trader_migrator	(build db schema migrator image)
- docker run -e APP_POSTGRES_URI="postgres://postgres:postgres@172.17.0.2:5432/postgres?sslmode=disable" --name wheeler-trader_migrator wheeler-trader_migrator:latest
		(start a container with the migrator image that will create the database tables)
- cd ..
- docker build . -t wheeler-trader_app		(build application image)
- docker run -d -p 8080:8080 -e APP_POSTGRES_URI="postgres://postgres:postgres@172.17.0.2:5432/postgres" --name wheeler-trader_app wheeler-trader_app:latest
		(start a container with the application, the application could be requested on localhost 8080)


## Cleanup

- docker container rm wheeler-trader_app --force
- docker container rm wheeler-trader_migrator
- docker container rm wheeler-trader_postgres --force
- docker image rm wheeler-trader_app
- docker image rm wheeler-trader_migrator

## UI mockups

![image](https://user-images.githubusercontent.com/36930531/172374318-21d8bf13-1702-400a-8d0c-1c7175fb2ef2.png)
<br/><br/>
![image](https://user-images.githubusercontent.com/36930531/172374388-44c109c3-bb75-4c48-8ed6-cd1365b14350.png)
<br/><br/>
![image](https://user-images.githubusercontent.com/36930531/172374416-56b141d7-0651-42b9-9534-7f60d265182a.png)
<br/><br/>
![image](https://user-images.githubusercontent.com/36930531/172374436-d543a8b7-5320-424c-9389-85079d9e5c9a.png)
<br/><br/>
![image](https://user-images.githubusercontent.com/36930531/172374478-5af6fe95-b13b-46f2-a682-c3c3375326fc.png)
<br/><br/>
![image](https://user-images.githubusercontent.com/36930531/172374493-5c6b6862-93e3-4dc7-a357-5123308ee449.png)
<br/><br/>
![image](https://user-images.githubusercontent.com/36930531/172374507-c1fdd31f-0c0d-4c90-aaf7-85134db5ccf8.png)
<br/><br/>
![image](https://user-images.githubusercontent.com/36930531/172374517-30ee5aa9-54cc-4ea2-a908-837ff8f2711a.png)
<br/><br/>
![image](https://user-images.githubusercontent.com/36930531/172374534-fc1fe380-aee2-4528-911b-51ca91bc7855.png)
<br/><br/>
