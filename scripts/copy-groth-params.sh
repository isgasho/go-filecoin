#!/usr/bin/env bash
copy_groth_params() {
  local params_dir="${FILECOIN_PARAMETER_CACHE:-"/tmp/filecoin-proof-parameters/"}"
  local params_file=v9-zigzag-proof-of-replication-52431242c129794fe51d373ae29953f2ff52abd94c78756e318ce45f3e4946d8
  mkdir -p ${params_dir}
  pushd ${params_dir}
  if [[ ! -e ${params_file} ]]; then
    echo "groth parameter ${params_file} was not found"
    echo "try rerunning 'go run ./build/*.go deps'"
    exit 1
  fi
  ln -sf ${params_file} params.out
  popd
}

copy_groth_params
