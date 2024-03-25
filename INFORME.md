# Parte 1: Introducción a Docker

### Ejercicio N°1:

Se modificó el archivo docker-compose-dev.yaml y se agregò un segundo cliente llamado client2 utilizando la misma estructura que tenia client1, cambiando su id

Para comprobar basta con ejecutar nuevamente

```bash
make docker-compose-up
make docker-compose-logs
```

Y se podrá observar en los logs el nuevo cliente y sus mensajes.

### Ejercicio N°1.1:

Se creó un nuevo archivo docker-compose-script-with-n-clients.sh para poder generar dinámicamente clientes basado en comandos por terminal; además se modificó el comando del makefile `docker-compose-up` para utilizar ahora dicho script para inicializar los contenedores y el comando `docker-compose-logs` para usar el nuevo docker compose.

Para usar y comprobar

```bash
make docker-compose-up <Number of clients>
make docker-compose-logs
```

### Ejercicio N°2:

Para este ejercicio se definieron volumenes en el server y client; donde se montan sus respectivos archivos de configuración. Para esto se removió del dockerfile el COPY del archivo de config y se agregaron al .dockerignore dichos archisvos para que no se detecten como cambios al hacer `docker compose up`.

Se agregó al script generador del docker compose el volumen con el archivo de configuración que correspondía según el caso.

Para probar se cambió el puerto del servidor en ambos archivos al 8080 y se observó que los cambios se reflejaron correctamente.

```bash
server | 2024-03-16 02:52:12 DEBUG action: config | result: success | port: 8080 | listen_backlog: 5 | logging_level: DEBUG
server | 2024-03-16 02:52:12 INFO action: accept_connections | result: in_progress
server | 2024-03-16 02:52:13 INFO action: accept_connections | result: success | ip: 172.25.125.4
server | 2024-03-16 02:52:13 INFO action: receive_message | result: success | ip: 172.25.125.4 | msg: [CLIENT 2] Message N°1
server | 2024-03-16 02:52:13 INFO action: accept_connections | result: in_progress
server | 2024-03-16 02:52:13 INFO action: accept_connections | result: success | ip: 172.25.125.3
server | 2024-03-16 02:52:13 INFO action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°1
server | 2024-03-16 02:52:13 INFO action: accept_connections | result: in_progress
client2 | time="2024-03-16 02:52:13" level=info msg="action: config | result: success | client_id: 2 | server_address: server:8080 | loop_lapse: 20s | loop_period: 5s | log_level: DEBUG"
```

Se comprobó que utilizó la cache al hacer `make docker-compose-down` seguido de `make docker-compose-up clients=2` queriendo decir que no se montó de nuevo toda la imagen por el cambio de configuración.

### Ejercicio N°3:

Se creó un nuevo contenedor `netcat-client` con su propio dockerfile en el que se instala netcat y se utiliza un script para comprobar que la respuesta al enviarle un mensaje al server sea la misma que dicho mensaje.

La ip y el puerto del serrver se pueden configurar en `netcat-client/config.env`

Para realizar el test basta con ejecutar `make run-netcat-test` y observar los resultados. En caso de éxito se verá por terminal

```bash
docker compose -f docker-compose-n-clients.yaml up -d --build --remove-orphans
[+] Running 1/0
 ✔ Container server  Running                                                                                                                                                         0.0s
make[1]: se sale del directorio '/home/erick/distribuidos/sistemas-distribuidos-tp0'
docker run --rm --network sistemas-distribuidos-tp0_testing_net --env-file ./netcat-client/config.env --name netcat-client netcat-client:latest
OK: Server response is correct: message: Test message to server, response: Test message to server
```

En caso de error; por ejemplo al cambiar el puerto del server a uno invalido:

```bash
docker compose -f docker-compose-n-clients.yaml up -d --build --remove-orphans
[+] Running 1/0
 ✔ Container server  Running                                                                                                                                                         0.0s
make[1]: se sale del directorio '/home/erick/distribuidos/sistemas-distribuidos-tp0'
docker run --rm --network sistemas-distribuidos-tp0_testing_net --env-file ./netcat-client/config.env --name netcat-client netcat-client:latest
ERROR: Server response is incorrect: message: Test message to server, response:
```