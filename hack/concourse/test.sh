#!/bin/bash

set -x -e

# start docker and log-in to docker-hub
entrypoint.sh
docker login --username=$DOCKER_USER --password=$DOCKER_PASS
docker run hello-world

# install python pip
apt-get update > /dev/null
apt-get install -y python python-pip > /dev/null

# install kubectl
curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl &> /dev/null
chmod +x ./kubectl
mv ./kubectl /bin/kubectl

# install onessl
curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.3.0/onessl-linux-amd64 \
  && chmod +x onessl \
  && mv onessl /usr/local/bin/

# install pharmer
pushd /tmp
curl -LO https://cdn.appscode.com/binaries/pharmer/0.1.0-rc.4/pharmer-linux-amd64
chmod +x pharmer-linux-amd64
mv pharmer-linux-amd64 /bin/pharmer
popd

function cleanup_test_stuff {
    rm -rf $ONESSL ca.crt ca.key server.crt server.key || true

    # delete cluster on exit
    pharmer get cluster || true
    pharmer delete cluster $NAME || true
    pharmer get cluster || true
    sleep 300 || true
    pharmer apply $NAME || true
    pharmer get cluster || true

    # delete docker image on exit
    curl -LO https://raw.githubusercontent.com/appscodelabs/libbuild/master/docker.py || true
    chmod +x docker.py || true
    ./docker.py del_tag appscodeci kubed $KUBED_IMAGE_TAG || true
}
trap cleanup_test_stuff EXIT

# copy kubed to $GOPATH
mkdir -p $GOPATH/src/github.com/appscode
cp -r kubed $GOPATH/src/github.com/appscode
pushd $GOPATH/src/github.com/appscode/kubed

# name of the cluster
# nameing is based on repo+commit_hash
NAME=kubed-$(git rev-parse --short HEAD)

./hack/builddeps.sh
export APPSCODE_ENV=test-concourse
export DOCKER_REGISTRY=appscodeci
./hack/docker/setup.sh
./hack/docker/setup.sh push
popd

#create credential file for pharmer
cat > cred.json <<EOF
{
    "token" : "$TOKEN"
}
EOF

# create cluster using pharmer
# note: make sure the zone supports volumes, not all regions support that
pharmer create credential --from-file=cred.json --provider=DigitalOcean cred
pharmer create cluster $NAME --provider=digitalocean --zone=nyc1 --nodes=2gb=1 --credential-uid=cred --kubernetes-version=v1.10.0
pharmer apply $NAME
pharmer use cluster $NAME
#wait for cluster to be ready
sleep 300
kubectl get nodes

# create storageclass
cat > sc.yaml <<EOF
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: standard
parameters:
  zone: nyc1
provisioner: external/pharmer
EOF

# create storage-class
kubectl create -f sc.yaml
sleep 120
kubectl get storageclass

pushd $GOPATH/src/github.com/appscode/kubed

# run tests
source ./hack/deploy/kubed.sh --docker-registry=appscodeci
./hack/make.py test e2e --v=3 --kubeconfig=/root/.kube/config --selfhosted-operator=true
