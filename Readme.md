## First draft for the whole project

```
version: '3.8'

services:
  timescaledb:
    image: timescale/timescaledb:latest-pg14
    container_name: timescaledb
    environment:
      POSTGRES_USER: your_user
      POSTGRES_PASSWORD: your_password
      POSTGRES_DB: your_db
    ports:
      - "5432:5432"
    volumes:
      - timescaledb-data:/var/lib/postgresql/data
    networks:
      - backend

  js-surveillance:
    build: ./js-surveillance
    container_name: js-surveillance
    depends_on:
      - go-app
    volumes:
      - ./js-surveillance:/app
    working_dir: /app
    command: "node surveillance.js"
    networks:
      - backend

  js-trade:
    build: ./js-trade
    container_name: js-trade
    depends_on:
      - go-app
    volumes:
      - ./js-trade:/app
    working_dir: /app
    command: "node trade.js"
    networks:
      - backend

  go-app:
    build: ./go-app
    container_name: go-app
    depends_on:
      - timescaledb
    volumes:
      - ./go-app:/app
    working_dir: /app
    command: "./go-app"
    networks:
      - backend

  go-api:
    build: ./go-api
    container_name: go-api
    depends_on:
      - go-app
    volumes:
      - ./go-api:/app
    working_dir: /app
    command: "./go-api"
    networks:
      - backend
      - frontend
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.go-api.rule=Host(`api.yourdomain.com`)"
      - "traefik.http.routers.go-api.entrypoints=websecure"
      - "traefik.http.routers.go-api.tls.certresolver=myresolver"
      # Middleware for authentication (optional)
      - "traefik.http.middlewares.go-api-auth.basicauth.users=user:$$apr1$$H6uskkkW$$IgXLP6EWZwLWziLQGoT8M/"
      - "traefik.http.routers.go-api.middlewares=go-api-auth"

  pgadmin:
    image: dpage/pgadmin4
    container_name: pgadmin
    environment:
      PGADMIN_DEFAULT_EMAIL: admin@example.com
      PGADMIN_DEFAULT_PASSWORD: your_secure_password
    volumes:
      - pgadmin-data:/var/lib/pgadmin
    depends_on:
      - timescaledb
    networks:
      - frontend
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.pgadmin.rule=Host(`pgadmin.yourdomain.com`)"
      - "traefik.http.routers.pgadmin.entrypoints=websecure"
      - "traefik.http.routers.pgadmin.tls.certresolver=myresolver"
      # Middleware for authentication (optional)
      - "traefik.http.middlewares.pgadmin-auth.basicauth.users=admin:$$apr1$$nUuR6gHf$$r4EiaPRPzQZ2h8GDz6.7Y."
      - "traefik.http.routers.pgadmin.middlewares=pgadmin-auth"

  traefik:
    image: traefik:v2.9
    container_name: traefik
    command:
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
      - "--certificatesresolvers.myresolver.acme.httpchallenge=true"
      - "--certificatesresolvers.myresolver.acme.httpchallenge.entrypoint=web"
      - "--certificatesresolvers.myresolver.acme.email=your-email@example.com"
      - "--certificatesresolvers.myresolver.acme.storage=/letsencrypt/acme.json"
      - "--api.dashboard=true"
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080"  # Traefik dashboard
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./letsencrypt:/letsencrypt
    networks:
      - frontend
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.traefik.rule=Host(`traefik.yourdomain.com`)"
      - "traefik.http.routers.traefik.entrypoints=websecure"
      - "traefik.http.routers.traefik.tls.certresolver=myresolver"
      - "traefik.http.routers.traefik.service=api@internal"
      # Middleware for authentication (optional)
      - "traefik.http.middlewares.traefik-auth.basicauth.users=admin:$$apr1$$k.yDy8xO$$bwjQp6oxpOyW1qJKGmP7h/"
      - "traefik.http.routers.traefik.middlewares=traefik-auth"

volumes:
  timescaledb-data:
  pgadmin-data:

networks:
  backend:
    driver: bridge
  frontend:
    driver: bridge
```
