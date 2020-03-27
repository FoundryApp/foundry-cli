#!/bin/sh

flags=$@
build_name="foundry"

echo "Building with flags: '$@'...\n"

for flag in `echo $flags`; do
  build_name=$build_name-$flag
done

go build -tags "$flags" -o "./build/$build_name" .
if [ $? -eq 0 ]; then
  echo "✅ SUCCESS"
else
  echo "\n❌ FAIL"
fi
