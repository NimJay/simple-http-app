# Multi-cluster backend

To turn the backend service into a multi-cluster service in Google Kubernetes Engine (GKE) clusters, follow these steps:

1. Select your Google Cloud project.
    ```bash
    export PROJECT_ID=my-project-id
    ```

1. Enable the necessary APIs:
    ```bash
    gcloud --project="${PROJECT_ID}" services enable \
	    container.googleapis.com \
        dns.googleapis.com \
        gkeconnect.googleapis.com \
        gkehub.googleapis.com \
        cloudresourcemanager.googleapis.com \
        multiclusterservicediscovery.googleapis.com \
        trafficdirector.googleapis.com
    ```

1. Create 2 GKE Autopilot clusters:
    ```bash
    gcloud --project="${PROJECT_ID}" \
        beta container clusters create-auto my-cluster-1 \
        --region=us-central1 \
        --fleet-project="${PROJECT_ID}"
    gcloud --project="${PROJECT_ID}" \
        beta container clusters create-auto my-cluster-2 \
        --region=us-central1 \
        --fleet-project="${PROJECT_ID}"
    ```
    If you already have clusters created, you'll have to use
    [a different command to register those existing clusters into the Fleet](https://cloud.google.com/anthos/fleet-management/docs/register/gke#register_your_cluster).

1. Check that both clusters have been registered into the project's Fleet.
    ```bash
    gcloud container fleet memberships list --project="${PROJECT_ID}"
    ```

1. Enable multi-cluster services in the your Fleet.
    ```bash
    gcloud container fleet multi-cluster-services enable --project="${PROJECT_ID}"
    ```

1. Verify that multi-cluster services have been enabled.
    ```bash
    gcloud container fleet multi-cluster-services describe --project "${PROJECT_ID}"
    ```

1. Grant the required Identity and Access Management (IAM) permissions for MCS Importer:
    ```bash
    gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
        --member "serviceAccount:"${PROJECT_ID}".svc.id.goog[gke-mcs/gke-mcs-importer]" \
        --role "roles/compute.networkViewer"
    ```

1. Deploy the backend to cluster 1.
    ```bash
    kubectl --context="gke_${PROJECT_ID}_us-central1_my-cluster-1" \
        apply -f ./backend/kubernetes/
    ```

1. Deploy the `ServiceExport` to cluster 1, which turns the backend `Service` into a multi-cluster service.
    ```bash
    kubectl --context="gke_${PROJECT_ID}_us-central1_my-cluster-1" \
        apply -f ./multi-cluster-backend/service-export.yaml
    ```

1. Set the `BACKEND_URL` value inside the frontend Deployment to `http://my-backend-service.default.svc.clusterset.local`,
   so that the `frontend` refers to the multi-cluster backend Service. The starting `http://` is optional.

1. Replace `PROJECT_ID` inside both `deployment.yaml` Kubernetes YAML files.

1. Deploy the frontend service into both clusters.
    ```bash
    kubectl --context="gke_${PROJECT_ID}_us-central1_my-cluster-1" \
        apply -f ./frontend/kubernetes/
    kubectl --context="gke_${PROJECT_ID}_us-central1_my-cluster-2" \
        apply -f ./frontend/kubernetes/
    ```

1. Get the public IP addresses of the two frontend Services.
    ```bash
    kubectl --context="gke_${PROJECT_ID}_us-central1_my-cluster-1" \
        get service my-frontend-service | awk '{print $4}'
    kubectl --context="gke_${PROJECT_ID}_us-central1_my-cluster-2" \
        get service my-frontend-service | awk '{print $4}'
    ```
