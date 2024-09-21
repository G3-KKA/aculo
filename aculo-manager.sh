#! /bin/bash













########################## Start At Line 257+ ##########################


















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
        echo "[INFO] WORKSPACE is not set, setting it to PWD: ${PWD}"
        WORKSPACE=$(pwd)
        export WORKSPACE
    fi
}


#######################################
# Exits if given file not exist
#
# Arguments:
#   One -- existance of this file will be checked, if it is not exist default message will be printed
#   OR
#   Two or more -- first will be checked as file or, if it is not exist all other arguments will be printed as error mesages before exiting (77)
#   OR
#   Zero - will result in early `return 1`
#     without checking actual file existence
#     message of missused function will be printed
# Usage:
#   required_to_exist ${WORKSPACE}/mario
#   required_to_exist $(pwd)../luigi [FATAL] luigi should be set otherwise mario will be upset!
# Exits (77) if:
#   - ARG_1 does not exist
#######################################
function required_to_exist(){
    if [[ $# -eq 0 ]]; then
        echo "[INFO] required_to_exist called without arguments!"
        return 1 
    fi
    if [[ ! -e $1 ]]; then
        if [[ $# -gt 1 ]]; then
            shift 
            for message in "$@" ; do
                echo $message
            done
            exit 77
        fi
        echo "[FATAL] $1 does not exist!"
        exit 77
    fi

    return 0

}

#######################################
# Asserts that command line args 
# Arguments:
#   - 2 integers, first is min, second is max, min <= max required
#   OR
#   - any other variant will result in early `return 1`
#     without checking actual range of command line args
#     message of missused function will be printed
# Globals:
#   COMMAND_LINE_ARGS_LENGTH
# Usage:
#   required_command_line_args_in_rage 0 2 => zero, one, two command line args are allowed    
#   required_command_line_args_in_rage 1 1 => strictly one command line arg is allowed
#  
# Exits (69) if:
#   - Number of command line args is not in range
#######################################
function required_command_line_args_in_rage() {
    # validate that function call is valid
    if [[ $# -eq 0 ]]; then
        echo "[INFO] required_command_line_args_in_rage called without arguments!"
        return 1
    fi
    if [[ $# -ne 2 ]]; then
        echo "[INFO] required_command_line_args_in_rage called with wrong number of arguments, expected 2 got $#!"
        return 1
    fi
    if [[ ! $1 -le $2 ]]; then
        echo "[INFO] required_command_line_args_in_rage called with wrong arguments, expected ARG_1 <= ARG_2 got $1 and $2!"
        return 1
    fi
    # compare that command line args are in given range
    if [[  ! $1 -le ${COMMAND_LINE_ARGS_LENGTH} ]]; then
        echo "[INFO] Type 'help' for usage ..."
        echo "[FATAL] Number of command line args is not in range! Expected in range of [$1,$2] got ${COMMAND_LINE_ARGS_LENGTH}!"
        exit 69
    fi
    if [[  ! $2 -ge ${COMMAND_LINE_ARGS_LENGTH} ]]; then
        echo "[INFO] Type 'help' for usage ..."
        echo "[FATAL] Number of command line args is not in range! Expected in range of [$1,$2] got ${COMMAND_LINE_ARGS_LENGTH}!"
        exit 69
    fi
}

#######################################
# Check that argument(s) are set  
# Arguments:
#   One or more ENV variables -- all will be checked
#   OR
#   Zero -- will result in early `return 1`
#     without checking actual range of command line args
#     message of missused function will be printed
# Usage:
#    
#   required_to_be_set REDIS_ADDRESS
#
# Description:
#   Call once on every required ENV variable
#   This way you do not need to write multiple if [[]] fi on every ENV
#   ! Downside: prints defult message if ENV is not set
#   Which is not enough verbose in some cases 
#
# Exits (45) if:
#   One of ENV variables is not set
#######################################
function required_to_be_set() {
    if [[ $# -eq 0 ]]; then
        echo "[INFO] required_to_be_set called without arguments-env to check!"
        return 1 
    fi

    if [[  -z "${!1}"  ]]; then
        if [[ $# -gt 1 ]]; then
            shift 
            for message in "$@" ; do
                echo $message
            done
            exit 45
        fi
        echo "[FATAL] $1 not set!"
        exit 45
    fi
}
#######################################
# Fails with message if previous command exited with non-zero exit code
# Globals:
#   global_saved_ex_code
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
global_saved_ex_code=0
function on_non_zero_exit_code(){
    __ex_code=$?
    if [[ ${global_saved_ex_code} -ne 0 ]]; then
        __ex_code=${global_saved_ex_code}
    fi
    if [[ ${__ex_code} -ne 0 ]]; then
        if [[ $# -ne 0 ]]; then
            for msg in "$@"; do 
                echo $msg
            done
        else
            echo "[FATAL] Previous command exited with non-zero exit code but no message was given!"
        fi
        exit ${__ex_code}
    fi
}
#######################################
# Checks that two ENV variables have different values
# Globals:
# Arguments:
#   Two ENV variables -- will be checked to be different
#   OR
#   Three or more arguments -- first 2 will be checked as env, others will be printed on error 
#   OR
#   [0,1] arguments -- considered as error
#   Will print error msg of missung function but will not exit
# Usage:
#   required_to_be_different PWD WORKSPACE
#   required_to_be_different PWD WORKSPACE [FATAL] luigi should be set otherwise mario will be upset!
# Exits (29) if:
#   Two ENV variables have same value
#######################################
function required_to_be_different() {   
    if [[ $# -lt 2 ]]; then
        echo "[INFO] required_to_be_different called without arguments-env to check!"
        return 1 
    fi
    if [[ ${!1} = ${!2} ]]; then
        if [[ $# -gt 2 ]]; then
            shift 2
            for message in "$@" ; do
                echo $message
            done
            exit 29
        fi 
        echo "[FATAL] $1 Required to be different from $2"
        exit 29
    fi
}

########################### Bool Flag ############################
# Bool flags!
# Flags are set to false by default
# Logic:
#   - Flag set to ANYTHING means it's true
#   - Flag not set means it's false
# Flags:
#   - INTERNAL_DO_NOT_CLEAR_INTEGRATION_TESTS
#   - INTERNAL_DO_NOT_WAIT_FOR_CONTINUE

########################### Non-Bool Flag ############################
# Non-Bool flags!
# Flags have default values
# Flags:                       Default value:
#   - INTERNAL_DOCKER_WORKDIR  /app


########################### Aculo Manager ############################

COMMAND_LINE_ARGS=$@
COMMAND_LINE_ARGS_LENGTH=$#

if [[ $# -eq 0  ]]; then
    echo "[INFO] Type 'help' for usage ..." 
    exit 30
fi

current_datetime=$(date +"%Y_%m_%d_%H%M")

if [[ $1 != "help" ]]; then
    default_workspace
fi

case "$1" in
##########################
# Required ENV:
#
########################## 
h|-h|--h|help|-help|--help)

    echo "CLI tool with simple commands and validation                                             "
    echo "                                                                                         "
    echo "Usage:                                                                                   "
    echo "                                                                                         "
    echo "      aculo-manager.sh [OPTION]                                                          "
    echo "  Or                                                                                     "
    echo "      ENV=VALUE aculo-manager.sh [OPTION]                                                "
    echo "                                                                                         "
    echo "Options:                                                                                 "
    echo "                                                                                         "
    echo "      build          build binary to /bin/                                               "
    echo "      run            execute binary from /bin/                                           "
    echo "      test           unit-tests                                                          "
    echo "      integr-test    integration-tests                                                   "
    echo "      continue       continue, if test supports it                                       "
    echo "      build-image    build docker image                                                  "
    echo "      run-image      run docker image                                                    "
    echo "      self-populate  copies the RX_ONLY aculo-manager.sh to all directories in arguments "
    echo "      preserve-logs  copies the logs to backup directory                                 "
    echo "                                                                                         "
    echo "      help            show this help                                                     "
    echo "                                                                                         "
    
    exit 0
    ;;
##########################
# Required ENV:
#   - WORKSPACE
#   - APP
##########################
build)
    echo "[INFO] Building started ..."

    required_command_line_args_in_rage 1 1

    required_to_be_set WORKSPACE 
    required_to_be_set APP

    required_to_exist ${WORKSPACE}/cmd/main.go "[FATAL] main.go file not found in WORKSPACE : ${WORKSPACE}!"
    required_to_exist ${WORKSPACE}/vendor "[FATAL] vendor directory not found in WORKSPACE : ${WORKSPACE}!"


    if [[  -e "${WORKSPACE}/bin/${APP}" ]]; then
        echo "[INFO] Overriding existing binary ..."
    fi


    echo "[INFO] Vendor mod build ..."
    go build                        \
        -C ${WORKSPACE}             \
        -mod=vendor                 \
        -o ${WORKSPACE}/bin/${APP}  \
        ${WORKSPACE}/cmd/main.go    \


    on_non_zero_exit_code "[FATAL] Build failed!"

    echo "[INFO] Building Successful."

    exit 0
    ;;
##########################
# Required ENV:
#   - WORKSPACE
#   - APP
##########################
run)
    echo "[INFO] Running ..."

    required_command_line_args_in_rage 1 1 

    required_to_be_set WORKSPACE 
    required_to_be_set APP

    required_to_exist ${WORKSPACE}/.env "[FATAL] .env file not found in WORKSPACE : ${WORKSPACE}!"
    required_to_exist ${WORKSPACE}/bin/${APP} "[FATAL] ${APP} binary not found in WORKSPACE : ${WORKSPACE}!"
   
    set -a;
    source ${WORKSPACE}/.env;

    ${WORKSPACE}/bin/${APP}

    on_non_zero_exit_code "[FATAL] Run failed!"

    exit 0
    ;;
##########################
# Required ENV:
#   - WORKSPACE
##########################
test)
    echo "[INFO] Testing ..."
    required_command_line_args_in_rage 1 1
    required_to_be_set WORKSPACE

    required_to_exist ${WORKSPACE}/.env "[FATAL] .env file for testing not found in WORKSPACE : ${WORKSPACE}!"
    set -a;
    source ${WORKSPACE}/.env;

    echo "[INFO] Race detection enabled."

    root=$(pwd)
    cd ${WORKSPACE}

    go test               \
        -race             \
        -v                \
        ${WORKSPACE}/...  \


    on_non_zero_exit_code "[FATAL] Test failed!"

    cd $root

    exit 0
    ;;
##########################
# Required ENV:
#   - WORKSPACE
#   - APP_ROOT
#   - APP
# Optional ENV:
#   - INTERNAL_DO_NOT_CLEAR_INTEGRATION_TESTS
##########################
integr-test)
    echo "[INFO] Integration testing started ..."
    required_command_line_args_in_rage 1 1 

    required_to_be_set WORKSPACE 
    required_to_be_set APP
    required_to_be_set APP_ROOT "[FATAL] APP_ROOT is not set! It should point to the directory containing /bin/${APP}!"
   
    required_to_be_different PWD WORKSPACE

    required_to_exist ${WORKSPACE}/.env "[FATAL] .env file not found in WORKSPACE : ${WORKSPACE}!"
    required_to_exist ${WORKSPACE}/docker-compose.yaml "[FATAL] docker-compose.yaml not found in WORKSPACE : ${WORKSPACE}!"
    required_to_exist ${APP_ROOT}/go.mod
    
    set -a;
    source ${WORKSPACE}/.env;

    if [[ ! -e "${WORKSPACE}/tmp" ]]; then
        echo "[INFO] Trying to create /tmp directory in WORKSPACE ..."
        mkdir ${WORKSPACE}/tmp
        required_to_exist ${WORKSPACE}/tmp "[FATAL] Failed to create /tmp directory in WORKSPACE : ${WORKSPACE}! "
    fi

    ## HARDCODED WAIT FOR KAFKA CLICKHOUSE AND KAFKA UI ##

    echo "[INFO] Setting up environment for integration tests ..."
	docker compose                         \
        --project-directory ${WORKSPACE}   \
        up                                 \
        -d                                 \
        --wait kafka0 clickhouse kafka-ui  \


    on_non_zero_exit_code "[FATAL] Docker compose up before tests failed!"

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


    on_non_zero_exit_code "[FATAL] Integration tests compilation failed, some compiled tests might be missing!"

    cd $root


    echo "[INFO] Running integration tests ..."
    for tesibinary in $(ls ${WORKSPACE}/tmp/ | grep .test)
    do
        # Before execution test will be marked with EXEC_ prefix
        # This is used by the continue command
        # For ease of identification between processes
        mv ${WORKSPACE}/tmp/$tesibinary ${WORKSPACE}/tmp/EXEC_$tesibinary

        # TODO , ADD OPTIONAL TIMEOUT (${WORKSPACE}/tmp/EXEC_${tesibinary} -test.timeout 20s)
        echo "[INFO] Running test ${tesibinary} ..."
        # Actually executing tests
        (${WORKSPACE}/tmp/EXEC_${tesibinary})

        global_saved_ex_code=$?

        if [[ -z ${INTERNAL_DO_NOT_CLEAR_INTEGRATION_TESTS} ]]; then
            rm ${WORKSPACE}/tmp/EXEC_$tesibinary
        fi

        on_non_zero_exit_code "[FATAL] Integration test ${tesibinary} failed, not all compiled test might be removed!"

    done


    echo "[INFO] Teardown environment for integration tests ..."
    docker compose                        \
        --project-directory ${WORKSPACE}  \
        down                              \


    on_non_zero_exit_code "[FATAL] Docker compose down after tests failed!"

    exit 0
    ;;

##########################
# Required ENV:
#   - WORKSPACE
# Optional ENV:
#   - INTERNAL_DO_NOT_WAIT_FOR_CONTINUE
##########################
continue)
    echo "[INFO] Trying to continue tests, if that is possible ..."
    
    required_command_line_args_in_rage 1 1 

    required_to_be_set WORKSPACE

    required_to_be_different PWD WORKSPACE "[FATAL] WORKSPACE should point to the directory with test config and /tmp/.test files"

    test_regex=${WORKSPACE}/tmp/EXEC_.*test\$

    process_name=$( ps -ef | grep "${test_regex}" -o )
    if [[ -z ${process_name} ]]; then
        echo "[FATAL] No running test found!"
        exit 31
    fi


    pid=$(pidof ${process_name})
    if [[ -z ${pid} ]]; then
        echo "[FATAL] PID of ${process_name} not found!"
        exit 32
    fi

    if [[ ! -z  ${INTERNAL_DO_NOT_WAIT_FOR_CONTINUE}  ]]; then
        echo "[INFO] Killing ${process_name} with PID: ${pid}, \
              .continue flag may not be set yet!"

        
        kill -s 10 "${pid}"
        exit $?
    fi

    required_to_exist ${WORKSPACE}/tmp/.continue "[FATAL] .continue flag not found in WORKSPACE : ${WORKSPACE}!"

    kill -s 10 "${pid}"

    exit $?
    ;;
##########################
# Required ENV:
#   - WORKSPACE
#   - APP
# Optional ENV:
#   - APP_ROOT
#   - INTERNAL_DOCKER_WORKDIR
##########################
build-image)
    echo "[INFO] Building docker image ..."

    required_command_line_args_in_rage 1 1
    required_to_be_set WORKSPACE 
    required_to_be_set APP

    if [[ -z "$APP_ROOT" ]]; then
        APP_ROOT=${WORKSPACE}/${APP}
    fi

    required_to_exist ${APP_ROOT}/go.mod "[FATAL] APP_ROOT should point to actual root directory of a go project, which is checked by existing of go.mod! It HAS default value, this ENV should be set explicitly if you see this error, APP_ROOT: ${APP_ROOT}"

    echo "[INFO] Build date: ${current_datetime}"


    dockerfile_name=""
    if [[ -e "${APP_ROOT}/dockerfile" ]]; then
        dockerfile_name="dockerfile"
    elif [[ -e "${APP_ROOT}/Dockerfile" ]]; then
        dockerfile_name="Dockerfile"
    else 
        echo "[FATAL] Dockerfile not found in WORKSPACE : ${WORKSPACE}! "
        exit 34
    fi

    if [[ -z ${INTERNAL_DOCKER_WORKDIR} ]]; then
        INTERNAL_DOCKER_WORKDIR="/app"
    fi

    cp ./aculo-manager.sh ${APP_ROOT}/aculo-manager.sh
    root=$(pwd)
    cd ${WORKSPACE}

    docker build                                                           \
        -t ${APP}:${current_datetime}                                          \
        -f $(pwd)/${dockerfile_name}                                       \
        --build-arg WORKDIR=${INTERNAL_DOCKER_WORKDIR}                     \
        --build-arg APP=${APP}                                             \
        .                                                                  \

    global_saved_ex_code=$?
    cd ${root}

    on_non_zero_exit_code "[FATAL] Docker build failed!"

    exit 0
    ;;
##########################
# Required ENV:
#   - WORKSPACE
#   - APP
# Optional ENV:
#   - INTERNAL_DOCKER_WORKDIR
##########################
run-image)
    echo "[INFO] Running docker image ..."

    required_command_line_args_in_rage 1 1
    required_to_be_set WORKSPACE 
    required_to_be_set APP

    required_to_exist ${WORKSPACE}/.env 
    required_to_exist ${WORKSPACE}/config.yaml

    # sort by 2nd column (TAG) , return in reversed order
    # grep our app 
    # get only first line  
    image=$(docker images                                               \
            --format "{{.Repository}} {{.Tag}} {{.ID}}"                 \
            | sort -rk 2                                                \
            | grep "^${APP} 2"                                          \
            | head -n 1                                                 \
        )                                                               \

    echo "[INFO] Found Image: ${image}"

    id=$(echo ${image} | awk 'NR==1{print $2}')


    if [[ -z ${id} ]]; then
        echo "[FATAL] Docker image not found!"
        exit 35
    fi
    if [[ -z ${INTERNAL_DOCKER_WORKDIR} ]]; then
        INTERNAL_DOCKER_WORKDIR="/app"
    fi
    docker run                                                              \
        -v ${WORKSPACE}/.env:${INTERNAL_DOCKER_WORKDIR}/.env                \
        -v ${WORKSPACE}/config.yaml:${INTERNAL_DOCKER_WORKDIR}/config.yaml  \
        ${APP}:${id}                                                        \

    on_non_zero_exit_code "[FATAL] Docker run failed!"

    exit 0

    ;;
##########################
# Required ENV:
#   - LOG_DIR
# Optional ENV:
#   - INTERNAL_DO_NOT_TRY_TO_CREATE_COPY
#
##########################
preserve-logs) 
    echo "[INFO] Preserving logs ..."
    required_command_line_args_in_rage 1 2 
    required_to_be_set LOG_DIR

    if [[ ! -e ${LOG_DIR} ]]; then
        echo "[INFO] ${LOG_DIR} not found. Creating ..."
        mkdir ${LOG_DIR}
        on_non_zero_exit_code "[FATAL] Failed to create ${LOG_DIR}!"
    fi
    backup_dir=${LOG_DIR}/log_backup_${current_datetime}

    mkdir ${backup_dir}
    ex_code=$?
    global_saved_ex_code=${ex_code}

    if [[ $ex_code -ne 0 ]]; then
        if [[ -z ${INTERNAL_DO_NOT_TRY_TO_CREATE_COPY} ]]; then

            echo "[INFO] Failed to create ${backup_dir} !"
            echo "[INFO] Trying to create copy ..."

            backup_dir=${LOG_DIR}/log_backup_copy_${current_datetime}

            mkdir ${backup_dir}
            global_saved_ex_code=$?
        fi
    fi
    on_non_zero_exit_code "[FATAL] Failed to create ${backup_dir}!"
    found="false"
    recreate="false"
    if [[ $2 == "recreate" ]]; then
        recreate="true"
    fi

    for logfile in $(ls ${LOG_DIR} | grep .log)
    do
        found="true"

        mv ${LOG_DIR}/$logfile ${backup_dir}/$logfile
        global_saved_ex_code=$?
        if [[ ${recreate} == "true" ]]; then
            touch ${LOG_DIR}/$logfile
            chmod a=rwx ${LOG_DIR}/$logfile
        fi

        on_non_zero_exit_code "[FATAL] Failed to move $logfile to ${backup_dir}!"

    done

    if [[ ${found} == "false" ]]; then
        echo "[FATAL] No log files found in ${LOG_DIR}!"
        rmdir ${backup_dir}
        exit 37
    fi

    exit 0
    ;;
##########################
# Required ENV:
#   - WORKSPACE
#
##########################
self-populate)
    echo "[INFO] Self-populating ..."
    required_to_be_set WORKSPACE 
    
    shift
    if [[ $# -eq 0 ]]; then
        echo "[FATAL] self-populate called without directories to populate!"
        exit 89
    fi
    for relative_dir in $@ ; do
        dir=${WORKSPACE}/${relative_dir}
        required_to_exist ${dir}

        # Populating
        if [[ -e "${dir}/aculo-manager.sh" ]]; then
            echo "[INFO] Overriding existing ${dir}/aculo-manager.sh"
            chmod a=rwx ${dir}/aculo-manager.sh
        fi
        # chmod here is a simple trick to enfore changes to be made only to main script and not its own copies

        cp ./aculo-manager.sh ${dir}/aculo-manager.sh

        on_non_zero_exit_code "[FATAL] Failed to populate ${dir}/aculo-manager.sh"

        chmod a=rx ${dir}/aculo-manager.sh

        on_non_zero_exit_code "[FATAL] Failed to chmod ${dir}/aculo-manager.sh"
    done

    exit 0
    ;;
*)
    echo "[INFO] Type 'help' for usage ..."  
    echo "[FATAL] Unknown option: $1"

    exit 35
    ;;
esac
