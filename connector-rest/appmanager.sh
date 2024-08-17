#! /bin/bash













########################## Logic Starts At Line 158+ ##########################


















########################## Local Functions ##########################

#######################################
# If WORKSPACE is not set
# Then Set WORKSPACE to PWD and export it
# Globals:
#   WORKSPACE
#   PWD
# Arguments:
#   None
#######################################
function default_workspace() {
    if [[  -z "$WORKSPACE"  ]]; then 
        WORKSPACE=$(pwd)
        export WORKSPACE
    fi
}

#######################################
# Exporting env variables from .env
# Globals:
#   WORKSPACE
# Arguments:
#   None
# Exits (22) if:
#   .env file does not exist in WORKSPACE
#######################################
function must_set_env_from_workspace() {   
    if [[  -e "$WORKSPACE/.env"  ]]; then
        set -a;
        source ${WORKSPACE}/.env;
    else 
        exit 22
    fi
}

#######################################
# Validates pre-build requirements
# Globals:
#   WORKSPACE
#   APP
# Arguments:
#   None
# Exits (23,24) if:
#   - main.go file does not exist in WORKSPACE/cmd
#   - vendor directory does not exist in WORKSPACE
#######################################
function must_build_validate() {
    if [[  -e "${WORKSPACE}/bin/${APP}" ]]; then
        echo "[INFO] Overriding existing binary ..."
    fi
    if [[ ! -e "${WORKSPACE}/cmd/main.go" ]]; then
        echo "[FATAL] main.go file not found in WORKSPACE! "
        exit 23
    fi
    if [[ ! -e "${WORKSPACE}/vendor" ]]; then
        echo "[FATAL] vendor directory does not exist in WORKSPACE! "
        exit 24
    fi
}

#######################################
# Validates pre-run requirements
# Globals:
#   WORKSPACE
#   APP
# Arguments:
#   None
# Exits (25) if:
#   - binary not found in WORKSPACE/bin
#######################################
function must_run_validate() {
    if [[ ! -e "${WORKSPACE}/bin/${APP}" ]]; then
        echo "[FATAL] Executable binary not found in WORKSPACE! "
        exit 25
    fi
}

#######################################
# Validates pre-integration-test requirements
# Globals:
#   WORKSPACE
# Arguments:
#   None
# Exits (26,27,28) if:
#   - docker-compose.yaml not found in WORKSPACE
#   - .env file not found in WORKSPACE
#   - tmp directory not found in WORKSPACE and cant be created
#######################################
function must_integration_test_validate() {
    if [[ ! -e "${WORKSPACE}/docker-compose.yaml" ]]; then
        echo "[FATAL] docker-compose.yaml not found in WORKSPACE! "
        exit 27
    fi
    if [[ ! -e "${WORKSPACE}/.env" ]]; then
        echo "[FATAL] .env not found in WORKSPACE! "
        exit 28
    fi
    if [[ ! -e "${WORKSPACE}/tmp" ]]; then
        echo "[INFO] Trying to create tmp directory in WORKSPACE ..."
        mkdir ${WORKSPACE}/tmp
        if [[ $? -ne 0 ]]; then
            echo "[FATAL] Failed to create tmp directory in WORKSPACE! "
            exit 26
        fi
    fi

}

#######################################
# Checks that WORKSPACE != PWD
# Globals:
#   WORKSPACE
#   PWD
# Arguments:
#   None 
# Exits (29) if: 
#   WORKSPACE is not explicitly set
function must_workspace_pwd_different() {
    if [[ ${PWD} = ${WORKSPACE} ]]; then
        echo "[FATAL] $1 Requires explicit WORKSPACE, that is different from PWD"
        exit 29
    fi
}

########################### App Manager ############################

# Checks that there is only one option
if [[ $# -ne 1 ]]; then
    echo "[INFO] Type 'help' for usage ..." 
    echo "[FATAL] App Manager takes exactly one option"
    exit 30
fi

APP=connector-rest

default_workspace

must_set_env_from_workspace

case "$1" in
build)
    echo "[INFO] Building started ..."

    must_build_validate


    echo "[INFO] Vendor mod build ..."
    go build                       \
        -mod=vendor                \
        -o ${WORKSPACE}/bin/${APP} \
    ${WORKSPACE}/cmd/main.go 

    echo "[INFO] Building Successful."
    exit 0
    ;;
run)
    echo "[INFO] Running ..."


    must_run_validate


    ${WORKSPACE}/bin/${APP}

    exit 0
    ;;
test)
    echo "[INFO] Testing ..."
    echo "[INFO] Race detection enabled."
    go test   \
        -race \
        -v    \
    ./...

    exit $?
    ;;
integr-test)
    echo "[INFO] Integration testing started ..."


    must_workspace_pwd_different


    must_integration_test_validate


    echo "[INFO] Setting up environment for integration tests ..."
	docker compose                       \
        --project-directory ${WORKSPACE} \
        up                               \
        -d                               \
    --wait kafka0 clickhouse kafka-ui


    INTEGRATION_TEST="true";
    export INTEGRATION_TEST


    echo "[INFO] Compiling integration tests ..."
    go test -v -c -o ${WORKSPACE}/tmp/ ./...


    echo "[INFO] Running integration tests ..."
    for varb in $(ls ${WORKSPACE}/tmp/ | grep .test)
    do
        # Before execution test will be marked with EXEC_ prefix
        # This is used by the continue command
        # For ease of identification between processes
        mv ${WORKSPACE}/tmp/$varb ${WORKSPACE}/tmp/EXEC_$varb

        # TODO , ADD OPTIONAL TIMEOUT (${WORKSPACE}/tmp/EXEC_${varb} -test.timeout 20s)
        # Actually executing tests
        echo "[INFO] Running test ${varb} ..."
        (${WORKSPACE}/tmp/EXEC_${varb})
    done


    echo "[INFO] Teardown environment for integration tests ..."
    docker compose                       \
        --project-directory ${WORKSPACE} \
    down


    exit 0
    ;;

continue)
    echo "[INFO] Trying to continue tests ..."

    must_workspace_pwd_different

    testRegex=${WORKSPACE}/tmp/EXEC_.*test\$

    pname=$( ps -ef | grep "${testRegex}" -o )
    if [[ -z ${pname} ]]; then
        echo "[FATAL] No running test found!"
        exit 31
    fi


    pid=$(pidof ${pname})
    if [[ -z ${pid} ]]; then
        echo "[FATAL] PID of ${pname} not found!"
        exit 32
    fi

    if [[ ! -z  ${DO_NOT_WAIT_FOR_CONTINUE}  ]]; then
        echo "[INFO] Killing ${pname} with PID ${pid}, \
              .continue flag may not be set yet!"

        
        kill -s 10 "${pid}"
        exit $?
    fi

    if [[ ! -e "${WORKSPACE}/tmp/.continue" ]]; then
        echo "[FATAL] .continue flag not found"
        exit 33
    fi
    

    kill -s 10 "${pid}"


    exit $?
    ;;
h|-h|--h|help|-help|--help)
    echo "CLI tool with simple commands and validation"
    echo ""
    echo "Usage:"
    echo ""
    echo "      appmanager.sh [OPTION]"
    echo "  Or"
    echo "      ENV=VALUE appmanager.sh [OPTION]"
    echo ""
    echo "Options:"
    echo ""
    echo "      build         build binary to /bin/ "
    echo "      run           execute binary from /bin/"
    echo "      test          unit-tests"
    echo "      integr-test   integration-tests"
    echo "      continue      continue, if test supports it"
    echo ""
    echo "      help          show this help"
    echo ""
    exit 0
    ;;
*)
    echo "[INFO] Type 'help' for usage ..."  
    echo "[FATAL] Unknown option: $1"

    exit 35
    ;;
esac
