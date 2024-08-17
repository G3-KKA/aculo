#! /bin/bash













########################## Start At Line 199+ ##########################


















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
# Validates that .env file exists in WORKSPACE
# Globals:
#   WORKSPACE
# Arguments:
#   None
# Exits (22) if:
#   .env file does not exist in WORKSPACE
#######################################
function must_validate_env_exists_in_workspace() {
    if [[ !  -e "$WORKSPACE/.env"  ]]; then
        echo "[INFO] .env file not found in WORKSPACE! "
        exit 22
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
    must_validate_env_exists_in_workspace
    set -a;
    source ${WORKSPACE}/.env;
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
function must_validate_build() {
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
function must_validate_run() {
    if [[ ! -e "${WORKSPACE}/bin/${APP}" ]]; then
        echo "[FATAL] Executable binary not found in WORKSPACE! "
        exit 25
    fi
}
#######################################
# Validates that .env file exists in WORKSPACE
# Globals:
#   WORKSPACE
# Arguments:
#   None
# Exits (37,38) if:
#   - go.mod not found in APP_ROOT
#   - APP_ROOT is not set
#######################################
function must_validate_app_root() {
    if [[ -z "${APP_ROOT}" ]]; then
        echo "[FATAL] APP_ROOT is not set! \
        It should point to the directory containing bin/${APP}!"
        exit 38
    fi
    if [[ ! -e "${APP_ROOT}/go.mod" ]]; then 
        echo "[FATAL] go.mod not found in APP_ROOT! "
        exit 37
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
function must_validate_integration_test() {
    if [[ ! -e "${WORKSPACE}/docker-compose.yaml" ]]; then
        echo "[FATAL] docker-compose.yaml not found in WORKSPACE! "
        exit 27
    fi
    must_validate_env_exists_in_workspace
    if [[ ! -e "${WORKSPACE}/tmp" ]]; then
        echo "[INFO] Trying to create tmp directory in WORKSPACE ..."
        mkdir ${WORKSPACE}/tmp
        if [[ $? -ne 0 ]]; then
            echo "[FATAL] Failed to create tmp directory in WORKSPACE! "
            exit 26
        fi
    fi
    must_validate_app_root


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
#######################################
function must_workspace_pwd_different() {
    if [[ ${PWD} = ${WORKSPACE} ]]; then
        echo "[FATAL] $1 Requires explicit WORKSPACE, that is different from PWD"
        exit 29
    fi
}
#######################################
# Fails with message if previous command exited with non-zero exit code
# Globals:
#   saved_ex_code
#   $? of previous command ( exit code in bash )
# Arguments:
#   One or more messages -- then will be printed
#   OR
#   None -- then default message will be printed
# Usage:
#   Call right after command you want to check,
#   Nothing will happen if previous command was successful (exit code 0)
#   Otherwise it will exit with non-zero exit code and print messages
#######################################
saved_ex_code=0
function non_zero_exit_code(){
    ex_code=$?
    if [[ ${saved_ex_code} -ne 0 ]]; then
        ex_code=${saved_ex_code}
    fi
    if [[ ${ex_code} -ne 0 ]]; then
        if [[ $# -ne 0 ]]; then
            for msg in "$@"; do 
                echo $msg
            done
        else
            echo "[FATAL] Previous command exited with non-zero exit code but no message was given!"
        fi
        exit ${ex_code}
    fi
}
#next_exit_code = 40#
########################### Bool Flag ############################
# Bool flags!
# Flags are set to false by default
# Logic:
#   - Flag set to ANYTHING means it's true
#   - Flag not set means it's false
# Flags:
#   - INTERNAL_DO_NOT_CLEAR_INTEGRATION_TESTS

########################### Non-Bool Flag ############################
# Non-Bool flags!
# Flags have default values
# Flags:                       Default value:
#   - INTERNAL_DOCKER_WORKDIR  /app


########################### Aculo Manager ############################

# Checks that there is only one option
if [[ $# -ne 1 ]]; then
    echo "[INFO] Type 'help' for usage ..." 
    echo "[FATAL] Aculo Manager takes exactly one option"
    exit 30
fi
if [[ -z "$APP" ]]; then
    echo "[INFO] Type 'help' for usage ..."
    echo "[FATAL] APP is not set!"
    exit 39
fi

default_workspace

must_set_env_from_workspace

case "$1" in
build)
    echo "[INFO] Building started ..."

    must_validate_build


    echo "[INFO] Vendor mod build ..."
    go build                        \
        -C ${WORKSPACE}             \
        -mod=vendor                 \
        -o ${WORKSPACE}/bin/${APP}  \
        ${WORKSPACE}/cmd/main.go    \


    non_zero_exit_code "[FATAL] Build failed!"

    echo "[INFO] Building Successful."

    exit 0
    ;;
run)
    echo "[INFO] Running ..."


    must_validate_run


    ${WORKSPACE}/bin/${APP}

    non_zero_exit_code "[FATAL] Run failed!"

    exit 0
    ;;
test)
    echo "[INFO] Testing ..."
    echo "[INFO] Race detection enabled."

    root=$(pwd)
    cd ${WORKSPACE}

    go test               \
        -race             \
        -v                \
        ${WORKSPACE}/...  \


    non_zero_exit_code "[FATAL] Test failed!"

    cd $root

    exit 0
    ;;
integr-test)
    echo "[INFO] Integration testing started ..."


    must_workspace_pwd_different


    must_validate_integration_test


    echo "[INFO] Setting up environment for integration tests ..."
	docker compose                         \
        --project-directory ${WORKSPACE}   \
        up                                 \
        -d                                 \
        --wait kafka0 clickhouse kafka-ui  \


    non_zero_exit_code "[FATAL] Docker compose up before tests failed!"

    INTEGRATION_TEST="true";
    export INTEGRATION_TEST


    root=$(pwd)
    cd ${APP_ROOT}

    echo "[INFO] Compiling integration tests ..."
    go test                   \
        -v                    \
        -c                    \
        -o ${WORKSPACE}/tmp/  \
        ${APP_ROOT}/...       \


    non_zero_exit_code "[FATAL] Integration tests compilation failed, some compiled tests might be missing!"

    cd $root


    echo "[INFO] Running integration tests ..."
    for varb in $(ls ${WORKSPACE}/tmp/ | grep .test)
    do
        # Before execution test will be marked with EXEC_ prefix
        # This is used by the continue command
        # For ease of identification between processes
        mv ${WORKSPACE}/tmp/$varb ${WORKSPACE}/tmp/EXEC_$varb

        # TODO , ADD OPTIONAL TIMEOUT (${WORKSPACE}/tmp/EXEC_${varb} -test.timeout 20s)
        echo "[INFO] Running test ${varb} ..."
        # Actually executing tests
        (${WORKSPACE}/tmp/EXEC_${varb})

        saved_ex_code=$?

        if [[ -z ${INTERNAL_DO_NOT_CLEAR_INTEGRATION_TESTS} ]]; then
            rm ${WORKSPACE}/tmp/EXEC_$varb
        fi

        non_zero_exit_code "[FATAL] Integration test ${varb} failed, not all compiled test might be removed!"

    done


    echo "[INFO] Teardown environment for integration tests ..."
    docker compose                        \
        --project-directory ${WORKSPACE}  \
        down                              \


    non_zero_exit_code "[FATAL] Docker compose down after tests failed!"

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
build-image)
    echo "[INFO] Building docker image ..."

    build_date=$(date +"%Y_%m_%d_%H%M")
    echo "[INFO] Build date: ${build_date}"

    if [[ -z "$APP_ROOT" ]]; then
        APP_ROOT=${WORKSPACE}/${APP}
    fi

    must_validate_app_root
    dockerfile_name=""
    if [[ -e "${APP_ROOT}/dockerfile" ]]; then
        dockerfile_name="dockerfile"
    elif [[ -e "${APP_ROOT}/Dockerfile" ]]; then
        dockerfile_name="Dockerfile"
    else 
        echo "[FATAL] Dockerfile not found in WORKSPACE! "
        exit 34
    fi
    if [[ -z ${INTERNAL_DOCKER_WORKDIR} ]]; then
        INTERNAL_DOCKER_WORKDIR="/app"
    fi
    docker build                                                           \
        -t ${APP}:${build_date}                                            \
        -f ${APP_ROOT}/${dockerfile_name}                                  \
        --build-arg WORKDIR=${INTERNAL_DOCKER_WORKDIR}                      \
        --volume ./aculo-manager.sh:${INTERNAL_DOCKER_WORKDIR}/aculo-manager.sh  \
        .                                                                  \


    non_zero_exit_code "[FATAL] Docker build failed!"

    exit 0
    ;;
h|-h|--h|help|-help|--help)
    echo "CLI tool with simple commands and validation"
    echo ""
    echo "Usage:"
    echo ""
    echo "      aculo-manager.sh [OPTION]"
    echo "  Or"
    echo "      ENV=VALUE aculo-manager.sh [OPTION]"
    echo ""
    echo "Options:"
    echo ""
    echo "      build         build binary to /bin/ "
    echo "      run           execute binary from /bin/"
    echo "      test          unit-tests"
    echo "      integr-test   integration-tests"
    echo "      continue      continue, if test supports it"
    echo "      build-image   build docker image"
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
