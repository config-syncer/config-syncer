#!/bin/bash
set -eou pipefail

echo "checking kubeconfig context"
kubectl config current-context || {
  echo "Set a context (kubectl use-context <context>) out of the following:"
  echo
  kubectl config get-contexts
  exit 1
}
echo ""

OS=""
ARCH=""
DOWNLOAD_URL=""
DOWNLOAD_DIR=""
TEMP_DIRS=()
ONESSL=""
ONESSL_VERSION=v0.13.1

# http://redsymbol.net/articles/bash-exit-traps/
function cleanup() {
  rm -rf ca.crt ca.key server.crt server.key
  # remove temporary directories
  for dir in "${TEMP_DIRS[@]}"; do
    rm -rf "${dir}"
  done
}

# detect operating system
# ref: https://raw.githubusercontent.com/helm/helm/master/scripts/get
function detectOS() {
  OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')

  case "$OS" in
    # Minimalist GNU for Windows
    cygwin* | mingw* | msys*) OS='windows';;
  esac
}

# detect machine architecture
function detectArch() {
  ARCH=$(uname -m)
  case $ARCH in
    armv7*) ARCH="arm";;
    aarch64) ARCH="arm64";;
    x86) ARCH="386";;
    x86_64) ARCH="amd64";;
    i686) ARCH="386";;
    i386) ARCH="386";;
  esac
}

detectOS
detectArch

# download file pointed by DOWNLOAD_URL variable
# store download file to the directory pointed by DOWNLOAD_DIR variable
# you have to sent the output file name as argument. i.e. downloadFile myfile.tar.gz
function downloadFile() {
  if curl --output /dev/null --silent --head --fail "$DOWNLOAD_URL"; then
    curl -fsSL ${DOWNLOAD_URL} -o $DOWNLOAD_DIR/$1
  else
    echo "File does not exist"
    exit 1
  fi
}

export APPSCODE_ENV=${APPSCODE_ENV:-prod}
trap cleanup EXIT

# ref: https://github.com/appscodelabs/libbuild/blob/master/common/lib.sh#L55
inside_git_repo() {
  git rev-parse --is-inside-work-tree >/dev/null 2>&1
  inside_git=$?
  if [ "$inside_git" -ne 0 ]; then
    echo "Not inside a git repository"
    exit 1
  fi
}

detect_tag() {
  inside_git_repo

  # http://stackoverflow.com/a/1404862/3476121
  git_tag=$(git describe --exact-match --abbrev=0 2>/dev/null || echo '')

  commit_hash=$(git rev-parse --verify HEAD)
  git_branch=$(git rev-parse --abbrev-ref HEAD)
  commit_timestamp=$(git show -s --format=%ct)

  if [ "$git_tag" != '' ]; then
    TAG=$git_tag
    TAG_STRATEGY='git_tag'
  elif [ "$git_branch" != 'master' ] && [ "$git_branch" != 'HEAD' ] && [[ "$git_branch" != release-* ]]; then
    TAG=$git_branch
    TAG_STRATEGY='git_branch'
  else
    hash_ver=$(git describe --tags --always --dirty)
    TAG="${hash_ver}"
    TAG_STRATEGY='commit_hash'
  fi

  export TAG
  export TAG_STRATEGY
  export git_tag
  export git_branch
  export commit_hash
  export commit_timestamp
}

onessl_found() {
  # https://stackoverflow.com/a/677212/244009
  if [ -x "$(command -v onessl)" ]; then
    onessl version --check=">=${ONESSL_VERSION}" >/dev/null 2>&1 || {
      # old version of onessl found
      echo "Found outdated onessl"
      return 1
    }
    export ONESSL=onessl
    return 0
  fi
  return 1
}

# download onessl if it does not exist
onessl_found || {
  echo "Downloading onessl ..."

  ARTIFACT="https://github.com/kubepack/onessl/releases/download/${ONESSL_VERSION}"
  ONESSL_BIN=onessl-${OS}-${ARCH}
  case "$OS" in
    cygwin* | mingw* | msys*)
      ONESSL_BIN=${ONESSL_BIN}.exe
    ;;
  esac

  DOWNLOAD_URL=${ARTIFACT}/${ONESSL_BIN}
  DOWNLOAD_DIR="$(mktemp -dt onessl-XXXXXX)"
  TEMP_DIRS+=($DOWNLOAD_DIR) # store DOWNLOAD_DIR to cleanup later

  downloadFile $ONESSL_BIN # downloaded file name will be saved as the value of ONESSL_BIN variable

  export ONESSL=${DOWNLOAD_DIR}/${ONESSL_BIN}
  chmod +x $ONESSL
}

# ref: https://stackoverflow.com/a/7069755/244009
# ref: https://jonalmeida.com/posts/2013/05/26/different-ways-to-implement-flags-in-bash/
# ref: http://tldp.org/LDP/abs/html/comparison-ops.html

export KUBED_NAMESPACE=kube-system
export KUBED_SERVICE_ACCOUNT=kubed-operator
export KUBED_ENABLE_RBAC=true
export KUBED_RUN_ON_MASTER=0
export KUBED_DOCKER_REGISTRY=${DOCKER_REGISTRY:-appscode}
export KUBED_IMAGE_TAG=v0.11.0
export KUBED_IMAGE_PULL_SECRET=
export KUBED_IMAGE_PULL_POLICY=IfNotPresent
export KUBED_USE_KUBEAPISERVER_FQDN_FOR_AKS=true
export KUBED_ENABLE_ANALYTICS=true
export KUBED_UNINSTALL=0

export KUBED_CONFIG_CLUSTER_NAME=unicorn
export KUBED_CONFIG_ENABLE_APISERVER=false
export KUBED_PRIORITY_CLASS=system-cluster-critical

export APPSCODE_ENV=${APPSCODE_ENV:-prod}
export SCRIPT_LOCATION="curl -fsSL https://raw.githubusercontent.com/appscode/kubed/v0.11.0/"
if [[ "$APPSCODE_ENV" == "dev" ]]; then
  detect_tag
  export SCRIPT_LOCATION="cat "
  export KUBED_IMAGE_TAG=$TAG
  export KUBED_IMAGE_PULL_POLICY=Always
fi

KUBE_APISERVER_VERSION=$(kubectl version -o=json | $ONESSL jsonpath '{.serverVersion.gitVersion}')
$ONESSL semver --check='<1.9.0' $KUBE_APISERVER_VERSION || { export KUBED_CONFIG_ENABLE_APISERVER=true; }

show_help() {
  echo "kubed.sh - install Kubernetes cluster daemon"
  echo " "
  echo "kubed.sh [options]"
  echo " "
  echo "options:"
  echo "-h, --help                             show brief help"
  echo "-n, --namespace=NAMESPACE              specify namespace (default: kube-system)"
  echo "    --rbac                             create RBAC roles and bindings (default: true)"
  echo "    --docker-registry                  docker registry used to pull kubed images (default: appscode)"
  echo "    --image-pull-secret                name of secret used to pull kubed operator images"
  echo "    --run-on-master                    run kubed operator on master"
  echo "    --cluster-name                     name of cluster (default: unicorn)"
  echo "    --enable-apiserver                 enable/disable kubed apiserver"
  echo "    --use-kubeapiserver-fqdn-for-aks   if true, uses kube-apiserver FQDN for AKS cluster to workaround https://github.com/Azure/AKS/issues/522 (default true)"
  echo "    --enable-analytics                 send usage events to Google Analytics (default: true)"
  echo "    --uninstall                        uninstall kubed"
}

while test $# -gt 0; do
  case "$1" in
    -h | --help)
      show_help
      exit 0
      ;;
    -n)
      shift
      if test $# -gt 0; then
        export KUBED_NAMESPACE=$1
      else
        echo "no namespace specified"
        exit 1
      fi
      shift
      ;;
    --namespace*)
      export KUBED_NAMESPACE=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --docker-registry*)
      export KUBED_DOCKER_REGISTRY=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --image-pull-secret*)
      secret=$(echo $1 | sed -e 's/^[^=]*=//g')
      export KUBED_IMAGE_PULL_SECRET="name: '$secret'"
      shift
      ;;
    --use-kubeapiserver-fqdn-for-aks*)
      val=$(echo $1 | sed -e 's/^[^=]*=//g')
      if [ "$val" = "false" ]; then
        export KUBED_USE_KUBEAPISERVER_FQDN_FOR_AKS=false
      else
        export KUBED_USE_KUBEAPISERVER_FQDN_FOR_AKS=true
      fi
      shift
      ;;
    --enable-analytics*)
      val=$(echo $1 | sed -e 's/^[^=]*=//g')
      if [ "$val" = "false" ]; then
        export KUBED_ENABLE_ANALYTICS=false
      fi
      shift
      ;;
    --rbac*)
      val=$(echo $1 | sed -e 's/^[^=]*=//g')
      if [ "$val" = "false" ]; then
        export KUBED_SERVICE_ACCOUNT=default
        export KUBED_ENABLE_RBAC=false
      fi
      shift
      ;;
    --run-on-master)
      export KUBED_RUN_ON_MASTER=1
      shift
      ;;
    --cluster-name*)
      export KUBED_CONFIG_CLUSTER_NAME=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --enable-apiserver*)
      val=$(echo $1 | sed -e 's/^[^=]*=//g')
      if [ "$val" = "false" ]; then
        export KUBED_CONFIG_ENABLE_APISERVER=false
      else
        export KUBED_CONFIG_ENABLE_APISERVER=true
      fi
      shift
      ;;
    --uninstall)
      export KUBED_UNINSTALL=1
      shift
      ;;
    *)
      show_help
      exit 1
      ;;
  esac
done

if [ "$KUBED_NAMESPACE" != "kube-system" ]; then
    export KUBED_PRIORITY_CLASS=""
fi

if [ "$KUBED_UNINSTALL" -eq 1 ]; then
  # delete webhooks and apiservices
  kubectl delete validatingwebhookconfiguration -l app=kubed || true
  kubectl delete mutatingwebhookconfiguration -l app=kubed || true
  kubectl delete apiservice -l app=kubed
  # delete kubed operator
  kubectl delete deployment -l app=kubed --namespace $KUBED_NAMESPACE
  kubectl delete service -l app=kubed --namespace $KUBED_NAMESPACE
  kubectl delete secret -l app=kubed --namespace $KUBED_NAMESPACE
  # delete RBAC objects, if --rbac flag was used.
  kubectl delete serviceaccount -l app=kubed --namespace $KUBED_NAMESPACE
  kubectl delete clusterrolebindings -l app=kubed
  kubectl delete clusterrole -l app=kubed
  kubectl delete rolebindings -l app=kubed --namespace $KUBED_NAMESPACE
  kubectl delete role -l app=kubed --namespace $KUBED_NAMESPACE
  # delete user roles
  kubectl delete clusterrole appscode:kubed:view

  exit 0
fi

echo "checking whether extended apiserver feature is enabled"
$ONESSL has-keys configmap --namespace=kube-system --keys=requestheader-client-ca-file extension-apiserver-authentication || {
  echo "Set --requestheader-client-ca-file flag on Kubernetes apiserver"
  exit 1
}
echo ""

env | sort | grep KUBED*
echo ""

# create necessary TLS certificates:
# - a local CA key and cert
# - a webhook server key and cert signed by the local CA
$ONESSL create ca-cert
$ONESSL create server-cert server --domains=kubed-operator.$KUBED_NAMESPACE.svc
export SERVICE_SERVING_CERT_CA=$(cat ca.crt | $ONESSL base64)
export TLS_SERVING_CERT=$(cat server.crt | $ONESSL base64)
export TLS_SERVING_KEY=$(cat server.key | $ONESSL base64)

CONFIG_FOUND=1
kubectl get secret kubed-config -n $KUBED_NAMESPACE >/dev/null 2>&1 || CONFIG_FOUND=0
if [ $CONFIG_FOUND -eq 0 ]; then
  config=$(${SCRIPT_LOCATION}hack/deploy/config.yaml | $ONESSL envsubst)
  kubectl create secret generic kubed-config -n $KUBED_NAMESPACE \
    --from-literal=config.yaml="${config}"
fi
kubectl label secret kubed-config app=kubed -n $KUBED_NAMESPACE --overwrite

${SCRIPT_LOCATION}hack/deploy/operator.yaml | $ONESSL envsubst | kubectl apply -f -

if [ "$KUBED_ENABLE_RBAC" = true ]; then
  ${SCRIPT_LOCATION}hack/deploy/service-account.yaml | $ONESSL envsubst | kubectl apply -f -
  ${SCRIPT_LOCATION}hack/deploy/rbac-list.yaml | $ONESSL envsubst | kubectl auth reconcile -f -
  ${SCRIPT_LOCATION}hack/deploy/user-roles.yaml | $ONESSL envsubst | kubectl auth reconcile -f -
fi

if [ "$KUBED_RUN_ON_MASTER" -eq 1 ]; then
  kubectl patch deploy kubed-operator -n $KUBED_NAMESPACE \
    --patch="$(${SCRIPT_LOCATION}hack/deploy/run-on-master.yaml)"
fi

if [ "$KUBED_CONFIG_ENABLE_APISERVER" = true ]; then
  ${SCRIPT_LOCATION}hack/deploy/apiservices.yaml | $ONESSL envsubst | kubectl apply -f -
fi

echo "waiting until kubed deployment is ready"
$ONESSL wait-until-ready deployment kubed-operator --namespace $KUBED_NAMESPACE || {
  echo "Kubed deployment failed to be ready"
  exit 1
}

if [ "$KUBED_CONFIG_ENABLE_APISERVER" = true ]; then
  echo "waiting until kubed apiservice is available"
  $ONESSL wait-until-ready apiservice v1alpha1.kubed.appscode.com || {
    echo "Kubed apiservice failed to be ready"
    exit 1
  }
fi

echo
echo "Successfully installed Kubed cluster daemon in $KUBED_NAMESPACE namespace!"
