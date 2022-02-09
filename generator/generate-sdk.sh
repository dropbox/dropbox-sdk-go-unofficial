#! /usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
    echo "$0: Expecting exactly one command-line argument, got $#." 1>&2
    exit 1
fi

version=$(echo $1 | cut -f1 -d'.')
loc=$(realpath -e $0)
base_dir=$(dirname "$loc")
spec_dir="$base_dir/dropbox-api-spec"
gen_dir=$(dirname ${base_dir})/v$version/dropbox

stone -v -a :all go_types.stoneg.py "$gen_dir" "$spec_dir"/*.stone
stone -v -a :all go_client.stoneg.py "$gen_dir" "$spec_dir"/*.stone

# Update SDK and API spec versions
sdk_version=${1}
pushd ${spec_dir}
spec_version=$(git rev-parse --short HEAD)
popd

sed -i.bak -e "s/UNKNOWN SDK VERSION/${sdk_version}/" \
    -e "s/UNKNOWN SPEC VERSION/${spec_version}/" ${gen_dir}/sdk.go
rm ${gen_dir}/sdk.go.bak
pushd ${gen_dir}
goimports -l -w ${gen_dir}
popd
