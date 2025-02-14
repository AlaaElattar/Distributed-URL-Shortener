# Distributed URL Shortener with Rate Limiting & Analytics

A distributed URL shortener service with Redis storage, rate limiting, and analytics logging in
MongoDB.

## Table of Contents

- [🚀 Overview](#-overview)
- [🏗 System Architecture](#-system-architecture)

- [📌 Environment Configuration](#-environment-configuration)

- [🐳 Running the URL Shortener with Docker](#-running-the-url-shortener-with-docker)
  - [1️⃣ Ensure Docker is Installed](#1️⃣-ensure-docker-is-installed)
  - [2️⃣ Build and Start the Containers](#2️⃣-build-and-start-the-containers)
  - [3️⃣ Verify Running Containers](#3️⃣-verify-running-containers)
  - [4️⃣ Access the Application](#4️⃣-access-the-application)
  - [5️⃣ Stopping the Application](#5️⃣-stopping-the-application)
- [🚀 Running the URL Shortener with Kubernetes & Minikube](#-running-the-url-shortener-with-kubernetes--minikube)
  - [1️⃣ Start Minikube](#1️⃣-start-minikube)
  - [2️⃣ Create Kubernetes ConfigMap for Environment Variables](#2️⃣-create-kubernetes-configmap-for-environment-variables)
  - [3️⃣ Deploy to Kubernetes](#3️⃣-deploy-to-kubernetes)
  - [4️⃣ Verify Everything is Running](#4️⃣-verify-everything-is-running)
  - [5️⃣ Forward Ports for Local Testing](#5️⃣-forward-ports-for-local-testing)
  - [6️⃣ Stopping the Kubernetes Deployment](#6️⃣-stopping-the-kubernetes-deployment)
- [📚 Running Tests](#-running-tests)

## 🚀 Overview

This project is a high-performance, distributed URL shortener built with Go, designed for scalability, efficiency, and reliability. It provides:

- Fast URL shortening & redirection using Gin as the web framework.
- Efficient storage & caching with Redis, ensuring quick lookups.
- Real-time analytics tracking using MongoDB, logging each URL access.
- Rate limiting to prevent abuse, restricting users to 10 requests per minute.
- Worker goroutines to handle non-blocking analytics logging, improving API performance.
- Full containerized deployment with Docker & Kubernetes, enabling auto-scaling.
- Reverse proxy & load balancing with Nginx, ensuring high availability.

## 🏗 System Architecture

The following diagram illustrates the architecture of the Distributed URL Shortener:

```
   ┌────────────┐     ┌───────────┐     ┌───────────┐
   │    User    │ --> │   Nginx   │ --> │    API    │
   └────────────┘     └───────────┘     └───────────┘
                                     ↙        ↘
                          ┌──────────────┐ ┌───────────┐
                          │  Redis Cache │ │  MongoDB  │
                          └──────────────┘ └───────────┘
```

### 📌 Explanation

✔️ User → Sends a request to shorten a URL or retrieve the original.  
✔️ Nginx → Acts as a reverse proxy & load balancer.  
✔️ API (Gin Framework) → Handles business logic, rate limiting, and URL mapping.  
✔️ Redis Cache → Stores shortened URLs with a 30-day expiry.  
✔️ MongoDB → Stores analytics logs (shortID, timestamp, user IP).

## 📌 Environment Configuration

Before running the application, ensure you have a `.env` file configured with the necessary environment variables.

- **If running the app with Kubernetes:**

  ```bash
  REDIS_ADDRESS=redis-service:6379
  MONGO_URI=mongodb://mongo-service:27017
  SERVER_PORT=8080
  ```

- **If running the app with Docker Compose:**

  ```bash
  REDIS_ADDRESS=redis-db:6379
  MONGO_URI=mongodb://mongo-db:27017
  SERVER_PORT=8080
  ```

⚠️ Note: If you are running the application locally without Docker or Kubernetes, you should update all values accordingly to point to your local Redis and MongoDB instances.

## 🐳 Running the URL Shortener with Docker

### 1️⃣ Ensure Docker is Installed

Make sure you have Docker installed on your system.

### 2️⃣ Build and Start the Containers

Run the following command to build and start the application:

```bash
docker compose up --build
```

### 3️⃣ Verify Running Containers

Check if the containers are running using:

```bash
docker ps
```

### 4️⃣ Access the Application

The API should be accessible at:

```bash
http://localhost:8080
```

### 5️⃣ Stopping the Application

```bash
docker compose down
```

To remove all Docker volumes (including stored data):

```bash
docker compose down -v
```

## 🚀 Running the URL Shortener with Kubernetes & Minikube

### 1️⃣ Start Minikube

First, ensure Minikube is installed, then start it:

```bash
minikube start
```

2️⃣ Create Kubernetes ConfigMap for Environment Variables

The application requires environment variables, create a ConfigMap from your .env file:

```bash
kubectl create configmap url-shortener-env --from-env-file=.env
```

### 3️⃣ Deploy to Kubernetes

Apply all necessary configurations:

```bash
kubectl apply -f kubernetes/
```

This will create MongoDB, Redis, API, and Nginx deployments.

### 4️⃣ Verify Everything is Running

Check the running pods:

```bash
kubectl get pods
```

Check the services:

```bash
kubectl get services
```

### 5️⃣ Forward Ports for Local Testing

Since your API runs inside the Kubernetes cluster, you need to expose it to your local machine. Use the following command to forward the API service port:

```bash
kubectl port-forward service/url-shortener-service 8080:80
```

### 6️⃣ Stopping the Kubernetes Deployment

```bash
kubectl delete -f kubernetes/
```

To completely stop Minikube:

```bash
minikube stop
```

## 📚 Running Tests

To run tests, ensure you are inside the server/ directory where the Go modules are located.

### 1️⃣ Navigate to the server/ directory:

```bash
cd server
```

### 2️⃣ Run all tests inside the project:

```bash
go test ./...
```

### 3️⃣ Run tests inside a Docker container (if using Docker Compose):

```bash
docker compose exec api go test -v ./app
```
