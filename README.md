# Kubernetes Controller for ClusterScan Custom Resource

🚀 This repository contains an implementation of a basic Kubernetes controller designed to reconcile a custom resource called ClusterScan. The primary objective of this controller is to facilitate the management of arbitrary jobs within a Kubernetes cluster and record their execution results.

## Tech Stack

- **Language:** Golang
- **Framework:** Kubebuilder

## Prerequisites

- Kubernetes Cluster 🌐
- kubectl ⌨️
- Golang 🐹
- Kubebuilder 🏗️
- Docker 🐳

## Features

**Custom Resource Definition (CRD):** Defines the structure of ClusterScan, allowing users to create instances of this resource.

**Reconciliation:** The controller continuously monitors the state of ClusterScans and ensures that the desired state matches the actual state by creating Jobs and/or CronJobs as necessary.

**Support for One-off and Recurring Executions:** ClusterScans can specify either one-off executions or recurring executions using CronJobs, providing flexibility for various use cases.

## Run Locally

1. **Clone the project**

    ```bash
    git clone https://github.com/SaharshPatel24/clusterscan-controller.git
    ```

2. **Go to the project directory**

    ```bash
    cd clusterscan-controller
    ```

3. **Install dependencies**

    ```bash
    go mod tidy
    ```
4. **Run Kubernetes cluster locally**

    ```bash
    minikube start
    ```

5. **Create Custom Resource**
    - Create a custom resource according to the custom resource definition available [here](https://github.com/SaharshPatel24/clusterscan-controller/blob/main/config/crd/bases/api.my.domain_clusterscans.yaml)..
    - For a sample custom resource, refer to [this file](https://github.com/SaharshPatel24/clusterscan-controller/blob/main/config/samples/api_v1alpha1_clusterscan.yaml).

6. **Run the Controller**
    ```bash
    make run
    ```

7. **Apply Custom Resource**

    In a second terminal, apply the custom resource to the controller using 
    ```bash
    kubectl apply -f /path/to/custom-resource
    ```
