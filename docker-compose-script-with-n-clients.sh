MAX_CLIENTS=10

usage() {
    echo "Usage: $0 <Number of clients>" >&2
    exit 1
}

# Check if the number of arguments is correct
if [ "$#" -ne 1 ]; then
    usage
fi

# Check if the argument is a number
if ! [[ $1 =~ ^[0-9]+$ ]]; then
    echo "Error: Argument must be a number" >&2
    usage
fi

if [ $1 -gt $MAX_CLIENTS ]; then
    echo "Error: Maximum number of clients is $MAX_CLIENTS" >&2
    usage
fi

# Number of clients to add
num_clients=$1

# Generate Docker Compose configuration
echo "
version: '3.9'
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - LOGGING_LEVEL=DEBUG
    networks:
      - testing_net
    volumes:
      - ./server/config.ini:/config.ini
"

NAMES=("Erick" "Juan" "Pedro" "Sofia" "Laura")
LASTNAMES=("Martinez" "Gomez" "Perez" "Rodriguez" "Gonzalez")

# Add clients to the configuration
for ((i = 1; i <= num_clients; i++)); do
echo "  
  client$i:
    container_name: client$i
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=$i
      - CLI_LOG_LEVEL=DEBUG
      - CLI_NAME=${NAMES[$((i % 5))]}
      - CLI_LASTNAME=${LASTNAMES[$((i % 5))]}
      - CLI_DOCUMENT=1234567$i
      - CLI_BIRTHDATE=1990-01-0$i
      - CLI_NUMBER=757$i
    networks:
      - testing_net
    depends_on:
      - server
    volumes:
      - ./client/config.yaml:/config.yaml
"
done

# Add network to the configuration
echo "
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
"