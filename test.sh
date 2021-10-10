#!/bin/bash

RESULT=""

function run_testcase {
    echo
    echo "############################## run $TESTNAME test ################################"
    echo "go run . $OPTIONS $INPUT > $OUTPUT"
    go run . $OPTIONS $INPUT > $OUTPUT
    rc=$?
    if [ $rc -ne 0 ]; then
        echo "** $TESTNAME test Failed (exit $rc)**"
        RESULT="${RESULT}${TESTNAME} test failed (exit $rc)\n"
    else
        echo "------- diff ---------------------------"
        diff -urN $WANT $OUTPUT
        rc=$?
        echo "----------------------------------------"
        if [ $rc -ne 0 ]; then
            echo "** $TESTNAME Failed **"
            RESULT="${RESULT}${TESTNAME} test failed\n"
        else
            echo "** $TESTNAME test OK **"
        fi
    fi
}


# unit test
echo "############################## run unit test #################################"
go test -v || RESULT="Unit test failed\n"


# test Merge
TESTNAME="merge"
INPUT="testcase/Dockerfile.before-merge"
WANT="testcase/Dockerfile.after-merge"
OUTPUT="testcase/_tmp_Dockerfile.test-merge"
OPTIONS="-m"
run_testcase


# test Split
TESTNAME="split"
INPUT="testcase/Dockerfile.before-split"
WANT="testcase/Dockerfile.after-split"
OUTPUT="testcase/_tmp_Dockerfile.test-split"
OPTIONS="-s"
run_testcase


# test Include
TESTNAME="include"
INPUT="testcase/Dockerfile.before-include"
WANT="testcase/Dockerfile.after-include"
OUTPUT="testcase/_tmp_Dockerfile.test-include"
OPTIONS="-i"
export TEST_ENV=""
run_testcase


# test Include and Merge
TESTNAME="include and merge"
INPUT="testcase/Dockerfile.before-include-and-merge"
WANT="testcase/Dockerfile.after-include-and-merge"
OUTPUT="testcase/_tmp_Dockerfile.include-and-merge"
OPTIONS="-i -m"
run_testcase


echo
if [ "$RESULT" != "" ]; then
    echo -e "############## Results #################\n$RESULT"
    exit 1
else
    echo -e "############## Results #################"
    echo -e "All tests are successful"
fi

