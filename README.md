# Snapshot Creator Kubernetes Operator

A Kubernetes operator which creates VolumeSnapshot based on the inputs provided by the user in the CRD.

## How to run

- Run below command to generate Docker image:
    ```
    docker build -t <DOCKER_USERNAME>/snapshotter:1.0.0 .
    ```

- Push the generated image to Docker hub:
    ```
    docker push <DOCKER_USERNAME>/snapshotter:1.0.0
    ```

- Replace the image name in `manifest/deploy.yaml`.

- Run below command to deploy the controller on the Kubernetes:
    ```
    kubectl apply -f manifest/
    ```
