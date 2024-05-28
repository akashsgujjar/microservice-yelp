___
# Lab 0: Getting Started
Welcome to CSE453! This lab serves as a brief introduction to the following tools we will be using in cse453 and as a guide for setting up the lab infrastructure. 
We provide a brief overview of the technologies and end with a set of
non-graded assignments. Even though these assignments are not graded,
you still need to finish them in order to be prepared for the
remainder of the project.

First, a high-level overview of our toolchain:

- **Golang**: Efficient and concurrent programming language used for building high-performance applications, server-side software, and distributed systems.
- **Kubernetes**: Automated containerized application management platform used for deploying, scaling, and managing applications in a containerized environment, providing resilience, scalability, and declarative configuration.
- **protobuf**: Efficient and language-agnostic data serialization format used for structured data exchange between different systems, optimizing both message size and processing speed.
- **gRPC**: Fast and efficient microservices communication framework that uses protobuf as the interface definition language (IDL) and enables high-performance, language-independent communication between microservices in a distributed system.

## Assignment 1: Access your GCP instance
The first thing that we will need to do is set up a computing environment where we can eventually deploy our application. Since this is a cloud computing class, we will be using 🥁... the cloud!

### Google Cloud Platform
We will allocate a VM instance (type `e2-standard-8`) to each group, equipped with Ubuntu 20.04, 8 vCPUs, 4 cores, and 32 GB of memory. Please follow the steps below to access your instance.


1. **Upload your SSH keys**: To access your instance via a terminal or IDE, you will need to use either an existing SSH key pair or generate a new one. You can upload a public key to your VM instance two ways.

    - **Google Cloud UI**: Go to the instance details page (click on your VM name). Then click `Edit` towards the top of the page. Scroll down to the SSH Keys section and click `Add Item`. Enter your public key and then click `Save` at the bottom of the page.
    - **GCloud CLI**: Install the [gcloud CLI](https://cloud.google.com/sdk/docs/install#linux) and complete the setup steps from `gcloud init` if necessary. Running the following command should automatically generate and add an ssh key for you.
    
        ```bash
        gcloud compute ssh instance-<GROUP-ID>
        ```

2. **Find your instance's external IP address**: We created a VM instance for each group and the instance name will be `instance-<GROUP-ID>`. Navigate to Google Cloud Console dashboard to find your instance's external IP address.

3. **SSH into the virtual machine**: 
    ```bash
    # USERNAME: your username in the SSH key. 
    # EXTERNAL_IP: the external IP address of the VM.
    ssh -i PATH_TO_PRIVATE_KEY USERNAME@EXTERNAL_IP
    ```

Note: The instances will remain active for the duration of the class. To help us reduce costs, please suspend your instance when it is not in use and resume it as needed. While all your local files will remain intact when an instance is suspended, we recommend checking your files regularly to avoid any potential issues.


<!-- ### Option 2: Setup VM with Google Cloud Console
To set up a virtual machine in the Google Cloud Console, follow these steps:

1. Open the Google Cloud Console in your web browser.

2. Select the project name provided by your instructor from the project dropdown menu at the top of the page. 

3. In the Google Cloud Console sidebar, navigate to "Compute Engine" under the "Compute" section.

4. Click on "VM instances" in the sidebar to go to the VM instances page.

5. Click the "Create" button to create a new virtual machine instance.

6. On the "Create an instance" page, configure the following settings:

   - **Name**: Enter a name for your virtual machine. 
   - **Region/Zone**: Select the region and zone where you want your VM to be located. 
   - **Machine type**: Choose the desired machine type, CPU, and memory configuration for your VM.
   - **Boot disk**: Select a the Ubuntu 20.04 OS image for your VM. 
   - **Other options**: Select a boot disk size of 20GB.

7. Review the configuration and click the "Create" button to create the virtual machine.

8. Wait for the virtual machine to be provisioned. Once it's ready, you will see the VM instance details page with information about your new virtual machine. Under "VM Instances" you will see an "SSH" button to log into your VM from the cloud console.  -->


### Experiment Setup
1. **Clone Repository**:
Clone the class project repository. If prompted for a password, you will need to generate SSH keys and upload them to your GitLab profile.
    ```bash
    git clone git@gitlab.cs.washington.edu:syslab/cse453-welp.git
    cd cse453-welp
    ```
2. **Install Dependencies**: 

    8.1. ***Kubernetes***: The following script installs all the key Kubernetes tools that you will need including `kubectl`, `kubeadm`, and `kubelet`. It also initializes your virtual machine as a control plane node. 
    ```bash
    . ./scripts/k8s_setup.sh
    kubectl version # You should see the version of client and server if kubernetes is sucessfully installed
    ```
    Only one member of your group needs to run this script. Others can simply run the following:
    ```
    mkdir -p $HOME/.kube
    sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
    sudo chown $(id -u):$(id -g) $HOME/.kube/config
    ```

    8.2 ***Golang***:  The following script installs Golang.
    ```bash
    . ./scripts/golang_setup.sh
    go version  
    ```
    After one member of your group runs the script, the other group members just need to update their path via
    ```bash
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    source $HOME/.bashrc
    ```

    8.3 ***Conda and Python***: We will use [Python](https://python.org) for tasks such as populating our microservices application or checking expected whether actual application behavior matches expected behavior. [Miniconda](https://repo.anaconda.com/miniconda/) is a lightweight package manager that allows us to easily set up a Python environment. Python environments are essentially configurations of dependencies, Python libraries, and settings that provide isolated and self-contained spaces for running Python applications, ensuring consistency and avoiding conflicts between different projects or applications. The provided file `requirements.yaml` specifies all the needed dependencies and parameters for the cse453 environment. Each group member can run installation separately.
    ```
    wget https://repo.anaconda.com/miniconda/Miniconda3-py310_23.3.1-0-Linux-x86_64.sh -O Miniconda.sh
    bash Miniconda.sh

    conda env create -f requirements.yaml
    conda activate cse453
    ```
    If you get an error about `conda` not being found when trying to create the env, you may need to run `source $HOME/.bashrc`.

    8.4 `wrk2`: We will use an HTTP benchmarking tool called `wrk2` for generating load on our application and measuring performance. Only one group member needs to do `wrk2` setup.
    ```console
    cd wrk2
    sudo apt install luarocks
    sudo apt-get install libssl-dev
    sudo apt-get install libz-dev
    sudo luarocks install luasocket
    sudo apt install make
    make 
    ```
<!-- 9. **Generate Protocol Buffers**:
    ```bash 
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/detail/detail.proto

    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/detail/review.protom

    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/detail/reservation.proto
    ``` -->
Awesome! Now that we have finished setting up our VM, let's learn more about Go, Kubernetes, and gRPC!

## Assignment 2: Go, gRPC, and Kubernetes
Complete the following tutorials so that you are prepared for the remainder of the labs.

1. [A Tour of Go](https://go.dev/tour/list): A short, interactive tutorial that covers the basics of the Go programming language. You can complete this tutorial in the browser or on the virtual machine you set up in Assignment 1. 

2. [gRPC and Protobuf Tutorial](../tutorials/grpc.md): Complete this tutorial on the virtual machine you set up in Assignment 1. 

3. [Kubernetes Tutorial](https://kubernetes.io/docs/tutorials/kubernetes-basics/): Complete this tutorial on the virtual machine you set up in Assignment 1. 
    - You can skip the minikube setup since we'll be using kubeadm for this lab.
    - Only the "Deploy an app" modules is mandatory. However, you're encouraged to delve into other modules if you're interested.
