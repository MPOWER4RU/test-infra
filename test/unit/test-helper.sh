#!/bin/bash

# Copyright 2018 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# Useful exit code
readonly SUCCESS=0
readonly FAILURE=1

# Call a function and verify its return value and output.
# Parameters: $1 - expected return code.
#             $2 - expected output ("" if no output is expected)
#             $3..$n - function to call and its parameters.
function test_function() {
  local expected_retcode=$1
  local expected_string=$2
  local output="$(mktemp)"
  local output_code="$(mktemp)"
  shift 2
  echo -n "$(trap '{ echo $? > ${output_code}; }' EXIT ; "$@")" &> ${output}
  local retcode=$(cat ${output_code})
  if [[ ${retcode} -ne ${expected_retcode} ]]; then
    cat ${output}
    echo "Return code ${retcode} doesn't match expected return code ${expected_retcode}"
    return 1
  fi
  if [[ -n "${expected_string}" ]]; then
    local found=1
    grep "${expected_string}" ${output} > /dev/null || found=0
    if (( ! found )); then
      cat ${output}
      echo "String '${expected_string}' not found"
      return 1
    fi
  else
    if [[ -s ${output} ]]; then
      ls ${output}
      cat ${output}
      echo "Unexpected output"
      return 1
    fi
  fi
  echo "'$@' returns code ${expected_retcode} and displays '${expected_string}'"
}

# Run the function with gcloud mocked (does nothing and outputs nothing).
# Parameters: $1..$n - parameters passed to the function.
function mock_gcloud_function() {
  set -e
  function gcloud() {
    echo ""
  }
  "$@" 2>&1
}
