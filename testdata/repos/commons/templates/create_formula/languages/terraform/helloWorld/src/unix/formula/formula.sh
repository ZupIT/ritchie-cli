#!/bin/sh
# shellcheck disable=SC2116
# shellcheck disable=SC2039
runFormula() {
  if [ -f /.dockerenv ] ; then
    cd /rit/tf-files || exit
    export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
  else
    cd tf-files || exit
  fi

  # Export TF VAR to Uppercase
  for e in $(printenv); do
    if [[ $e == "TF_VAR_"* ]]; then
      env_value=$(cut -d'=' -f2 <<< "$e")
      env_name=$(cut -d'=' -f1 <<< "$e")
      sub_env_name_lower=$(echo "$env_name" | cut -c8-| tr '[:upper:]' '[:lower:]')
      export "$(echo "TF_VAR_${sub_env_name_lower}=${env_value}")"
    fi
  done

  terraform version



  if [ "$BACKEND_S3" != "" ]; then
    {
    echo " "
    echo "terraform {"
    echo "  backend \"s3\" {}"
    echo "}"
    } >> main.tf


    terraform init \
    -backend-config="key=terraform.tfstate" \
    -backend-config="bucket=$BACKEND_S3" \
    -backend-config="region=$TF_VAR_BUCKET_REGION"
  else
    terraform init
  fi

  if [ "$ACTION" = "plan" ]; then
    terraform plan
  else
    terraform "$ACTION" -auto-approve
  fi
}
