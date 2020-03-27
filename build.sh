#!/bin/sh

debug_flag="debug"

build_flag=$1 # "debug" or empty (= release)

if [ "$build_flag" = "$debug_flag" ]
then
  echo "Debug build...\n"
  build_name="foundry-debug"
elif [ -z "$build_flag" ] # empty build_flag is a release version
then
  echo "Release build..."
  build_name="foundry"
else
  echo "ERROR: Unknown build flag '$build_flag'"
  exit 1
fi

go build -tags "$build_flag" -o "./build/$build_name" .
if [ $? -eq 0 ]; then
  echo "✅ SUCCESS"
else
  echo "\n❌ FAIL"
fi

