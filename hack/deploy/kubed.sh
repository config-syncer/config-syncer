#!/bin/bash
set -eou pipefail

echo "checking kubeconfig context"
kubectl config current-context || { echo "Set a context (kubectl use-context <context>) out of the following:"; echo; kubectl config get-contexts; exit 1; }
echo ""

# http://redsymbol.net/articles/bash-exit-traps/
function cleanup {
	rm -rf $ONESSL ca.crt ca.key server.crt server.key
}

export APPSCODE_ENV=${APPSCODE_ENV:-prod}
if [ "$APPSCODE_ENV" != "test-concourse" ]; then
    trap cleanup EXIT
fi

# ref: https://github.com/appscodelabs/libbuild/blob/master/common/lib.sh#L55
inside_git_repo() {
    git rev-parse --is-inside-work-tree > /dev/null 2>&1
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

# https://stackoverflow.com/a/677212/244009
if [ -x "$(command -v onessl)" ]; then
    export ONESSL=onessl
else
    # ref: https://stackoverflow.com/a/27776822/244009
    case "$(uname -s)" in
        Darwin)
            curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.3.0/onessl-darwin-amd64
            chmod +x onessl
            export ONESSL=./onessl
            ;;

        Linux)
            curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.3.0/onessl-linux-amd64
            chmod +x onessl
            export ONESSL=./onessl
            ;;

        CYGWIN*|MINGW32*|MSYS*)
            curl -fsSL -o onessl.exe https://github.com/kubepack/onessl/releases/download/0.3.0/onessl-windows-amd64.exe
            chmod +x onessl.exe
            export ONESSL=./onessl.exe
            ;;
        *)
            echo 'other OS'
            ;;
    esac
fi

# ref: https://stackoverflow.com/a/7069755/244009
# ref: https://jonalmeida.com/posts/2013/05/26/different-ways-to-implement-flags-in-bash/
# ref: http://tldp.org/LDP/abs/html/comparison-ops.html

export KUBED_NAMESPACE=kube-system
export KUBED_SERVICE_ACCOUNT=kubed-operator
export KUBED_ENABLE_RBAC=true
export KUBED_RUN_ON_MASTER=0
export KUBED_DOCKER_REGISTRY=${DOCKER_REGISTRY:-appscode}
export KUBED_IMAGE_TAG=0.7.0-rc.1
export KUBED_IMAGE_PULL_SECRET=
export KUBED_IMAGE_PULL_POLICY=IfNotPresent
export KUBED_ENABLE_ANALYTICS=true
export KUBED_UNINSTALL=0

export SCRIPT_LOCATION="curl -fsSL https://raw.githubusercontent.com/appscode/kubed/0.7.0-rc.1/"
if [[ "$APPSCODE_ENV" = "dev" || "$APPSCODE_ENV" = "test-concourse" ]]; then
    detect_tag
    export SCRIPT_LOCATION="cat "
    export KUBED_IMAGE_TAG=$TAG
    export KUBED_IMAGE_PULL_POLICY=Always
fi

show_help() {
    echo "kubed.sh - install Kubernetes cluster daemon"
    echo " "
    echo "kubed.sh [options]"
    echo " "
    echo "options:"
    echo "-h, --help                         show brief help"
    echo "-n, --namespace=NAMESPACE          specify namespace (default: kube-system)"
    echo "    --rbac                         create RBAC roles and bindings (default: true)"
    echo "    --docker-registry              docker registry used to pull kubed images (default: appscode)"
    echo "    --image-pull-secret            name of secret used to pull kubed operator images"
    echo "    --run-on-master                run kubed operator on master"
    echo "    --enable-analytics             send usage events to Google Analytics (default: true)"
    echo "    --uninstall                    uninstall kubed"
}

while test $# -gt 0; do
    case "$1" in
        -h|--help)
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
            export KUBED_NAMESPACE=`echo $1 | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        --docker-registry*)
            export KUBED_DOCKER_REGISTRY=`echo $1 | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        --image-pull-secret*)
            secret=`echo $1 | sed -e 's/^[^=]*=//g'`
            export KUBED_IMAGE_PULL_SECRET="name: '$secret'"
            shift
            ;;
        --enable-analytics*)
            val=`echo $1 | sed -e 's/^[^=]*=//g'`
            if [ "$val" = "false" ]; then
                export SEARCHLIGHT_ENABLE_ANALYTICS=false
            fi
            shift
            ;;
        --rbac*)
            val=`echo $1 | sed -e 's/^[^=]*=//g'`
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
    kubectl get clusterrole appscode:kubed:view

    exit 0
fi

echo "checking whether extended apiserver feature is enabled"
$ONESSL has-keys configmap --namespace=kube-system --keys=requestheader-client-ca-file extension-apiserver-authentication || { echo "Set --requestheader-client-ca-file flag on Kubernetes apiserver"; exit 1; }
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
kubectl get secret kubed-config -n $KUBED_NAMESPACE > /dev/null 2>&1 || CONFIG_FOUND=0
if [ $CONFIG_FOUND -eq 0 ]; then
    config=`${SCRIPT_LOCATION}hack/deploy/config.yaml`
    kubectl create secret generic kubed-config -n $KUBED_NAMESPACE \
        --from-literal=config.yaml="${config}"
fi
kubectl label secret kubed-config app=kubed -n $KUBED_NAMESPACE --overwrite

${SCRIPT_LOCATION}hack/deploy/operator.yaml | $ONESSL envsubst | kubectl apply -f -

if [ "$KUBED_ENABLE_RBAC" = true ]; then
    kubectl create serviceaccount $KUBED_SERVICE_ACCOUNT --namespace $KUBED_NAMESPACE
    kubectl label serviceaccount $KUBED_SERVICE_ACCOUNT app=kubed --namespace $KUBED_NAMESPACE
    ${SCRIPT_LOCATION}hack/deploy/rbac-list.yaml | $ONESSL envsubst | kubectl auth reconcile -f -
    ${SCRIPT_LOCATION}hack/deploy/user-roles.yaml | $ONESSL envsubst | kubectl auth reconcile -f -
fi

if [ "$KUBED_RUN_ON_MASTER" -eq 1 ]; then
    kubectl patch deploy kubed-operator -n $KUBED_NAMESPACE \
      --patch="$(${SCRIPT_LOCATION}hack/deploy/run-on-master.yaml)"
fi

echo "waiting until kubed deployment is ready"
$ONESSL wait-until-ready deployment kubed-operator --namespace $KUBED_NAMESPACE || { echo "Kubed deployment failed to be ready"; exit 1; }

echo "waiting until kubed apiservice is available"
$ONESSL wait-until-ready apiservice v1alpha1.kubed.appscode.com || { echo "Kubed apiservice failed to be ready"; exit 1; }

echo
echo "Successfully installed Kubed cluster daemon in $KUBED_NAMESPACE namespace!"
