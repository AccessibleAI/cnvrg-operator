#!/usr/bin/env bash

WORK_DIR="$(pwd)/webhook_deployment"
COMMON_NAME=$1
EXPIRATION_DAYS=36500

echo "${WORK_DIR}"
echo "${COMMON_NAME}"

init () {
    if [ -d "$WORK_DIR" ]; then
        rm -fr "${WORK_DIR}"
    fi
    mkdir "${WORK_DIR}"
    cd "${WORK_DIR}" || exit
    cat << EOF > conf
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[v3_req]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth, serverAuth
EOF
}

create_ca () {
    openssl genrsa -out ca.key 2048
    openssl req -x509 -new -nodes -key ca.key -days ${EXPIRATION_DAYS} -out ca.crt -subj "/CN=admission_ca"
}

create_server_crts () {
    openssl genrsa -out server.key 2048
    openssl req -new -key server.key -out server.csr -subj "/CN=${COMMON_NAME}" -config conf
    openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days ${EXPIRATION_DAYS} -extensions v3_req -extfile conf
}

print_base64_certs (){
    echo -e "base64 encoded ca.crt\n"
    base64 -i "${WORK_DIR}"/ca.crt
    echo -e "\n"
    echo -e "base64 encoded server.crt\n"
    base64 -i "${WORK_DIR}"/server.crt
    echo -e "\n"
    echo -e "base64 encoded server.key\n"
    base64 -i "${WORK_DIR}"/server.key
    echo -e "\n"
}

print_k8s_webhook_def(){
export SERVICE_NAME=${COMMON_NAME}
export BASE64_CA_BUNDLE=$(base64 -i "${WORK_DIR}/ca.crt")
cat <<EOF > "${WORK_DIR}/adwebhook.yaml"
apiVersion: admissionregistration.k8s.io/v1beta1
kind: ValidatingWebhookConfiguration
metadata:
  name: uac
  labels:
    app: uac
webhooks:
  - name: ${SERVICE_NAME}
    clientConfig:
      url: https://${SERVICE_NAME}:8080/
      caBundle: ${BASE64_CA_BUNDLE}
    rules:
      - operations: [ "CREATE"]
        apiGroups: ["*"]
        apiVersions: ["*"]
        resources: ["oauthaccesstokens"]
    failurePolicy: Ignore
EOF
cat /tmp/adwebhook.yaml
echo "############### create webhook cmd ##################"
echo "# oc create -f ./webhook_deployment/adwebhook.yaml  #"
echo "#####################################################"
}

if [ "$#" -ne 1 ]; then
    echo "Missing certificate common name (CN). Example usage: ./create-certs.sh uac.bnhp-system.svc.cluster.local"
    exit 1
fi


init
create_ca
create_server_crts
print_base64_certs
print_k8s_webhook_def