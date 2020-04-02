# Keycloak, Gatekeeper and Docker Compose sample

This is sample project intended to show how to configure Keycloak and Keycloak Gatekeeper to protect microservice 
using docker compose.

# Prerequisites

1. Install required software:<br/>
   Docker Engine v19.03.5 or greater<br/>
   Docker Compose v1.25.4 or greater

2. Add to /etc/hosts next 2 lines:

        127.0.0.1	keykloak
        127.0.0.1	demo-client.local

# FAQ

**Q:** How to start environment?<br/>
**A:** Run next command in terminal `docker-compose up`

**Q:** How to stop environment?<br/>
**A:** Run next command in terminal `docker-compose down --rmi all -v`

**Q**: How to test microservices?<br/>
**A**: Acquire access token and use it to run requests described in files `./server_gatekeeper_protected/test.http`
or `./server_middleware_protected/test.http`

**Q**: How to acquire access token?<br/>
**A**: Once environment described in `./docker-compose.yml` file is up and running, open next url in browser 
http://demo-client.local:5556/ then click on login link and enter credentials into keycloak login form.
Upon successful login you will see access token shown on screen in bold font. To get access token for different user
perform logout by clicking on Logout link.

**Q**: Which users is able to login via http://demo-client.local:5556/ ?
**A**: There are 2 configured users:

1. Regular user can access only `api/v1/info`

        username: test
        password: test
        
2. Admin user can access `api/v1/info` and `api/v1/info_admin`
 
        username: test_admin
        password: test_admin