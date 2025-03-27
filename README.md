# Chat App with WebSocket, RabbitMQ, Redis, MongoDB, and Nginx
This project implements a scalable chat application using WebSocket for real-time messaging, RabbitMQ for message queuing, Redis for caching, MongoDB for data storage, and Nginx as a reverse proxy for WebSocket connections.

The project is configured to run with Docker and Docker Compose, providing a seamless setup for the entire application stack.

# Table of Contents
Prerequisites
Project Structure
Setup and Configuration
Running the Application
Accessing the Services
Stopping the Services
Troubleshooting
License
Prerequisites
Before running this project, ensure you have the following tools installed:

Docker — To containerize the application.
Docker Compose — To orchestrate all the services.
You can verify if Docker and Docker Compose are installed with the following commands:

bash
Copy
docker --version
docker-compose --version
Project Structure
graphql
Copy
chat-app/
│
├── nginx/                    # Nginx configuration files
│   ├── Dockerfile             # Dockerfile for Nginx
│   └── nginx.conf             # Nginx config for WebSocket reverse proxy
│
├── internal/
│   ├── queue/                 # RabbitMQ queue client code
│   ├── websocket/             # WebSocket client code
│   └── message/               # Logic for consuming and sending messages
│
├── Dockerfile                 # Dockerfile for the Go WebSocket server
├── docker-compose.yml         # Docker Compose file to run all services
└── README.md                  # This README file
Setup and Configuration
1. Clone the Repository
Start by cloning the repository to your local machine:

bash
Copy
git clone https://github.com/your-username/chat-app.git
cd chat-app
2. Docker Configuration
The project uses Docker Compose to set up multiple services. The docker-compose.yml file defines the following services:

WebSocket Server: A Go WebSocket server running on port 8080.
Nginx: A reverse proxy server that handles WebSocket connections on port 80.
RabbitMQ: Message queuing service running on port 5672 and the management UI on port 15672.
Redis: In-memory cache service running on port 6379.
MongoDB: Data storage service running on port 27017.
3. Nginx Configuration
The nginx/nginx.conf file configures Nginx to forward WebSocket requests to the WebSocket server. The WebSocket connections are handled at the /ws endpoint.
4. MongoDB, Redis, and RabbitMQ
MongoDB stores chat messages.
Redis is used for caching WebSocket connections.
RabbitMQ is used to queue chat messages and distribute them to WebSocket clients.
Running the Application
1. Build and Start the Services
To set up the application and start the services, run the following command in the root directory where docker-compose.yml is located:

bash
Copy
docker-compose up --build
This will:

Build the WebSocket server Docker image.
Set up Nginx to proxy WebSocket requests.
Start RabbitMQ, Redis, and MongoDB containers.
Docker Compose will download the necessary images and build the WebSocket service. It might take a few minutes to complete.

2. Verify the Running Services
Once the services are started, Docker Compose will display the logs for each service. You can also view logs with:

bash
Copy
docker-compose logs <service-name>
For example, to view the logs for RabbitMQ:

bash
Copy
docker-compose logs rabbitmq
Accessing the Services
Once the application is up and running, you can access the services:

1. Nginx WebSocket Proxy
WebSocket Proxy: The WebSocket proxy is handled by Nginx, which listens on port 80 and forwards WebSocket connections to the Go WebSocket server.
URL: http://localhost/
WebSocket URL: ws://localhost/ws
2. RabbitMQ Management UI
RabbitMQ has a management UI available on port 15672. You can access it via:

URL: http://localhost:15672/
Username: guest
Password: guest
3. Redis
Redis is running on port 6379. You can access it through any Redis client or from your Go application.

URL: redis://localhost:6379
4. MongoDB
MongoDB is available on port 27017. You can connect to it using any MongoDB client or through your application.

URL: mongodb://localhost:27017
Stopping the Services
To stop and remove all running containers, run:

bash
Copy
docker-compose down
If you want to stop the services but keep the containers running (e.g., for restarting later), use:

bash
Copy
docker-compose stop
To remove all containers, networks, and volumes, run:

bash
Copy
docker-compose down --volumes
Troubleshooting
1. Docker Compose Build Fails
If the build process fails, make sure the following:

Docker and Docker Compose are installed correctly.
All necessary files (such as Dockerfile, nginx.conf, etc.) are in the correct paths.
Check the logs with docker-compose logs <service-name> to debug any issues.
2. Port Conflicts
If the ports 80, 5672, 6379, or 27017 are already in use on your machine, you can modify the docker-compose.yml to map the services to different local ports. For example:

yaml
Copy
rabbitmq:
  ports:
    - "5673:5672"
3. Service Connectivity
If your services are unable to connect to each other, ensure that they are all on the same network. Docker Compose sets up a default network for all services, but if necessary, check the networks section in docker-compose.yml.

License
This project is licensed under the MIT License - see the LICENSE file for details.

