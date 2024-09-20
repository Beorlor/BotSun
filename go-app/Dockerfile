# go-app/Dockerfile

FROM golang:1.20-alpine

# Installer les dépendances nécessaires
RUN apk update && apk add --no-cache git

# Installer Air pour le rechargement à chaud
RUN go install github.com/cosmtrek/air@latest

WORKDIR /app

# Copier les fichiers de configuration
COPY go.mod go.sum ./

# Télécharger les dépendances
RUN go mod download

# Exposer le port de l'application
EXPOSE 8080

# Commande pour démarrer l'application avec Air
CMD ["air"]
