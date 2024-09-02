Name:           hop-hop-cluster
Version:        0.0.1
Release:        1%{?dist}
Summary:        A package which provides kubeadm bootstraps script

License:        GPL-3.0
Source0:        %{name}-%{version}.tar.gz

%global debug_package %{nil}

%description
The scripts included in this make it possible to initialize or join a Kubernetes cluster.
The necessary tokens, configurations and certificates can be provided using Ignition/Cloud-Init.

%prep
%setup -q

%build

%install
mkdir -p %{buildroot}/usr/local/bin
cp -a test.sh %{buildroot}/usr/local/bin/

%files
/usr/local/bin/

%changelog
* Thu Sep 02 2024 Hauke Mettendorf <hauke@mettendorf.it> - 0.0.1
- Initial version
