#! /bin/bash





if [[ $# -ne 1 && $# -ne 2 ]]; then
    echo "[INFO] Type 'help' for usage ..." 
    echo "[FATAL] Aculo takes one or two arguments"
    exit 20
fi
if [[ ! -e aculo-manager.sh ]]; then
    echo "[FATAL] INTERNAL ERROR !!!" 
    echo "[FATAL] aculo-manager.sh not found in PWD!"
    exit 22
fi


case "$1" in
batch-inserter)
    ;;
connector-rest)
    ;;
frontend-rest)
    ;;
start)
    echo "[INFO] Starting app ..."
    echo "NOT IMPLEMENTED YET"
    exit 1
    ;;
h|-h|--h|help|-help|--help)
    echo "User-Friendly CLI interface built on top of aculo-manager.sh"
    echo ""
    echo "Usage:"
    echo ""
    echo "      aculo.sh [COMMAND]"
    echo ""
    echo "      aculo.sh [APP] [OPTION]"
    echo ""
    echo "Command:"
    echo ""
    echo "      help          show this help"
    echo "      start         start whole app"
    echo ""
    echo "App:"
    echo "      batch-inserter"
    echo "      connector-rest"
    echo "      frontend-rest"
    echo ""
    echo "Options:"
    echo ""
    echo "      go            build and run app"
    echo "      build         build binary to /bin/ "
    echo "      run           execute binary from /bin/"
    echo "      build-image   build docker image"
    echo "      test          unit-tests"
    echo "      integr-test   integration-tests"
    echo "      continue      continue, if test supports it"
    echo ""
    exit 0
    ;;
*)
    echo "[INFO] Type 'help' for usage ..."  
    echo "[FATAL] Unknown option or command: $1"

    exit 35
    ;;
esac 
if [[ $# -ne 2 ]]; then
    echo "[INFO] Type 'help' for usage ..." 
    echo "[FATAL] App is specifiend but command is not"
    exit 21
#WORKSPACE={$APP_ROOT}/test/integration
fi

APP=$1
export APP

case "$2" in
go)
    APP_ROOT=$(pwd)/${APP}
    WORKSPACE=${APP_ROOT}

    export WORKSPACE
    export APP_ROOT

    ./aculo-manager.sh build
    ./aculo-manager.sh run

    exit $?

;;
build)
    APP_ROOT=$(pwd)/${APP}
    WORKSPACE=${APP_ROOT}

    export WORKSPACE
    export APP_ROOT

    ./aculo-manager.sh build

    exit $?

;;
run)

    APP_ROOT=$(pwd)/${APP}
    WORKSPACE=${APP_ROOT}

    export WORKSPACE
    export APP_ROOT

    ./aculo-manager.sh run

    exit $?

;;
build-image)
    APP_ROOT=$(pwd)/${APP}
    WORKSPACE=${APP_ROOT}

    export WORKSPACE
    export APP_ROOT

    ./aculo-manager.sh build-image

    exit $?

;;
test)
;;
integr-test)
;;

continue)
;;
*)
    echo "[INFO] Type 'help' for usage ..."
    echo "[FATAL] Unknown option: $2"
    exit 36
esac