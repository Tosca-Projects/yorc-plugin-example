#!/usr/bin/env bash 

set -euo pipefail
scriptDir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

function cleanup {
    echo  "removing Yorc container"
    docker stop yorc
    rm trace.log
}

function waitForYorcStartup {

    while true ; do 
        docker logs yorc  > trace.log  2>&1 || true
        if [[ "$(grep -c 'agent: Synced service "yorc"' trace.log)" != "0" ]]; then 
            break
        else 
            echo "waiting for Yorc startup..."
            sleep 2
        fi
    done
    rm -f trace.log
}

[[ -x "${scriptDir}/bin/my-plugin" ]] || { echo "Plugin binary missing, ensure that 'make build' completes properly" ; exit 1; }

echo "Running a Yorc container..."

docker pull ystia-docker.jfrog.io/ystia/yorc:latest

docker run -d --rm \
  -e 'YORC_LOCATIONS_FILE_PATH=/var/yorc/conf/locations.json' \
  -e 'YORC_LOG=1' \
  --mount "type=bind,src=${scriptDir}/conf,dst=/var/yorc/conf" \
	--mount "type=bind,src=${scriptDir}/bin,dst=/var/yorc/plugins" \
	--mount "type=bind,src=${scriptDir}/tosca,dst=/var/yorc/topology" \
    --name yorc \
	ystia-docker.jfrog.io/ystia/yorc:latest

trap cleanup EXIT
waitForYorcStartup

echo "Deploying application example"
docker exec -it yorc sh -c "yorc d deploy --id my-test-app /var/yorc/topology"

docker exec -it yorc sh -c "yorc d info --follow  my-test-app" 

echo "Checking deployment logs"
docker exec -it yorc sh -c "yorc d logs --no-stream  my-test-app" | tee trace.log > /dev/null

echo "Checking expected outputs..."

grep '**********Provisioning node "Compute" of type "mytosca.types.Compute"' trace.log || { echo "Missing Delegate executor log"; exit 1; }
grep '******Executing operation "standard.create" on node "Soft"' trace.log || { echo "Missing Operation executor log"; exit 1; }

echo "Test succeeded"

