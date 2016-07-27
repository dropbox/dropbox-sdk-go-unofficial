#! /usr/bin/env bash
set -euo pipefail

if [[ $# -ne 0 ]]; then
    echo "$0: Not expecting any command-line arguments, got $#." 1>&2
    exit 1
fi

loc=$(realpath -e $0)
base_dir=$(dirname "$loc")
spec_dir="$base_dir/dropbox-api-spec"
gen_dir=$(dirname ${base_dir})/dropbox

python3 -m stone.cli -v -a :all go_types.stoneg.py "$gen_dir" "$spec_dir"/*.stone
python3 -m stone.cli -v -a :all go_client.stoneg.py "$gen_dir" "$spec_dir"/*.stone

goimports -l -w ${gen_dir}
