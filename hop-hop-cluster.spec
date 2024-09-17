Name:           hop-hop-cluster
Version:        0.0.3
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
mkdir -p %{buildroot}/usr/bin
cp -a test.sh %{buildroot}/usr/bin/
cp -a test2.sh %{buildroot}/usr/bin/

%files
/usr/bin/test.sh
/usr/bin/test2.sh

%changelog
* Tue Sep 17 2024 Hauke Mettendorf <hauke@mettendorf.it> - 0.0.3
- Added every binary to %files

* Thu Sep 13 2024 Hauke Mettendorf <hauke@mettendorf.it> - 0.0.2
- Changes binary path to /usr/bin since Fedora CoreOS doesn't use /usr/local/bin

* Thu Sep 02 2024 Hauke Mettendorf <hauke@mettendorf.it> - 0.0.1
- Initial version
