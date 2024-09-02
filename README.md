
<a href="https://hauke.cloud" target="_blank"><img src="https://img.shields.io/badge/home-hauke.cloud-brightgreen" alt="hauke.cloud" style="display: block;" /></a>
<a href="https://github.com/hauke-cloud" target="_blank"><img src="https://img.shields.io/badge/github-hauke.cloud-blue" alt="hauke.cloud Github Organisation" style="display: block;" /></a>

# Fedora CoreOS Kubernetes Bootstrap

<img src="https://raw.githubusercontent.com/hauke-cloud/.github/main/resources/img/organisation-logo-small.png" alt="hauke.cloud logo" width="109" height="123" align="right">

This repository contains a collection of scripts to prepare a Fedora CoreOS for a Kubernetes installation and to join or initialize a cluster.

We are trying to achive following goals with this template:
- Standardization of Kubernetes cluster creation
- Configuration using Ignition/Cloud-Init
- Rolling updates of nodes after cluster initialization
- Providing the scripts using rpm so that they can be installed in Fedora CoreOS
- Cluster state determination based on API request and Kubernetes API server


## Table of Contents

- [Getting started](#-getting-started)
- [License](#license)
- [Contributing](#contributing)
- [Contact](#contact)

## 🚀 Getting started
To get started, you need to clone the repository containing this `README.md` file. Follow the steps below:

### 1. Clone the repository

Use the following command to clone the repository:

```bash
git clone https://github.com/hauke-cloud/fedora-coreos-kubernetes.git
```

### 2. Navigate to the repository directory

Once the repository is cloned, navigate to the directory:

```bash
cd fedora-coreos-kubernetes
```

### 3. Check the content

```bash
ls -la
```

This will display all the files and directories in the cloned repository.




## 📄 License

This Project is licensed under the GNU General Public License v3.0

- see the [LICENSE](LICENSE) file for details.

## :coffee: Contributing

To become a contributor, please check out the [CONTRIBUTING](CONTRIBUTING.md) file.
## :email: Contact

For any inquiries or support requests, please open an issue in this
repository or contact us at [contact@hauke.cloud](mailto:contact@hauke.cloud).
