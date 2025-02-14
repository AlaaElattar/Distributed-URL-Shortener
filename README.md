# Distributed URL Shortener with Rate Limiting & Analytics

A distributed URL shortener service with Redis storage, rate limiting, and analytics logging in
MongoDB.

## Table of Contents




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
