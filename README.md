# Simple HTTP app

This is a simple demo app consisting of a frontend service and a backend service — both written in Golang.
This app is primarily used for Kubernetes-related testing.

<img src="/diagram.png" alt="Two arrows, the first arrow pointing from a user to the frontend, the arrow second pointing from the frontend to the backend." height="100" />

### Build and push images

To build and push the images into a Google Cloud [Artifact Registry](https://cloud.google.com/artifact-registry/docs/overview) repository, follow these steps:

1. Select your Google Cloud project.
    ```bash
    export PROJECT_ID=my-project-id
    ```

1. Enable Artifact Registry.
    ```bash
    gcloud --project="$PROJECT_ID" \
        services enable artifactregistry.googleapis.com
    ```

1. Create an Artifact Registry Docker repository.
    ```bash
    gcloud --project="$PROJECT_ID" \
        artifacts repositories create simple-http-app --repository-format=docker \
        --location=us --description="Docker repository used for https://github.com/NimJay/simple-http-app"
    ```

1. Build the images.
    ```bash
    cd frontend
    docker build --tag "us-docker.pkg.dev/${PROJECT_ID}/simple-http-app/frontend:0.0.5" .
    cd ..
    cd backend
    docker build --tag "us-docker.pkg.dev/${PROJECT_ID}/simple-http-app/backend:0.0.5" .
    cd ..
    ```

1. Push the images
    ```bash
    docker push "us-docker.pkg.dev/${PROJECT_ID}/simple-http-app/frontend:0.0.5"
    docker push "us-docker.pkg.dev/${PROJECT_ID}/simple-http-app/backend:0.0.5"
    ```

### Run locally using Docker

Assuming you have already built the images locally (see `docker build` commands above), you can run them locally:
```bash
docker run -p 8081:8080 -d "us-docker.pkg.dev/${PROJECT_ID}/simple-http-app/frontend:0.0.5"
docker run -p 8082:8080 -d "us-docker.pkg.dev/${PROJECT_ID}/simple-http-app/backend:0.0.5"
```

You can then access the endpoints at [localhost:8081](http://localhost:8081) and [localhost:8082](http://localhost:8082).

If you have the backend running at some public IP address like `123.456.789`, you can use use:
```bash
docker run -p 8081:8080 --env BACKEND_URL=http://123.456.789 \
    -d "us-docker.pkg.dev/${PROJECT_ID}/simple-http-app/frontend:0.0.5"
```

### Deploy app into Kubernetes cluster

1. Replace `PROJECT_ID` inside both `deployment.yaml` Kubernetes YAML files.

1. Assuming your current `kubectl` context is set to your desired Kubernetes cluster, deploy the Kubernetes YAML files:
    ```bash
    kubectl apply -f ./frontend/kubernetes/
    kubectl apply -f ./backend/kubernetes/
    ```

1. Get the public IP address of the frontend Service:
    ```bash
    kubectl get service my-frontend-service | awk '{print $4}'
    ```
