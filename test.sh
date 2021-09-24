#!/bin/bash

RESULT=""

# unit test

echo "############################## run unit test #################################"
go test -v || RESULT="Unit test failed\n"

# test Merge
echo
echo "############################## run merge test ################################"
INPUT="testcase/Dockerfile.before-merge"
WANT="testcase/Dockerfile.after-merge"
OUTPUT="testcase/_tmp_Dockerfile.test-merge"

go run . -m $INPUT > $OUTPUT
echo "------- diff ---------------------------"
diff -urN $WANT $OUTPUT
rc=$?
echo "----------------------------------------"
if [ $rc -ne 0 ]; then
    echo "** merge test Failed **"
    RESULT="${RESULT}Merge test failed\n"
else
    echo "** merge test OK **"
fi


# test Split
echo
echo "############################## run split test ################################"
INPUT="testcase/Dockerfile.before-split"
WANT="testcase/Dockerfile.after-split"
OUTPUT="testcase/_tmp_Dockerfile.test-split"

go run . -s $INPUT > $OUTPUT
echo "------- diff ---------------------------"
diff -urN $WANT $OUTPUT
rc=$?
echo "----------------------------------------"
if [ $rc -ne 0 ]; then
    echo "** split test Failed **"
    RESULT="${RESULT}Split test failed\n"
else
    echo "** split test OK **"
fi

# test Include
echo
echo "############################## run include test ################################"
INPUT="testcase/Dockerfile.before-include"
WANT="testcase/Dockerfile.after-include"
OUTPUT="testcase/_tmp_Dockerfile.test-include"

go run . -i $INPUT > $OUTPUT
echo "------- diff ---------------------------"
diff -urN $WANT $OUTPUT
rc=$?
echo "----------------------------------------"
if [ $rc -ne 0 ]; then
    echo "** include test Failed **"
    RESULT="${RESULT}Include test failed\n"
else
    echo "** include test OK **"
fi

# test Include and Merge
echo
echo "############################## run include and merge test ################################"
INPUT="testcase/Dockerfile.before-include-and-merge"
WANT="testcase/Dockerfile.after-include-and-merge"
OUTPUT="testcase/_tmp_Dockerfile.include-and-merge"

go run . -i -m $INPUT > $OUTPUT
echo "------- diff ---------------------------"
diff -urN $WANT $OUTPUT
rc=$?
echo "----------------------------------------"
if [ $rc -ne 0 ]; then
    echo "** include and merge test Failed **"
    RESULT="${RESULT}Include and Merge test failed\n"
else
    echo "** include and merge test OK **"
fi

echo
echo -e "############## Failed tests is below #################\n$RESULT"

