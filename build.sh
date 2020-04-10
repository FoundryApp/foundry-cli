#!/bin/sh

flags=$@
build_name="foundry"

echo "Building with flags: '$@'...\n"

is_debug=0

for flag in `echo $flags`; do
  if [[ $flag == "debug"  ]]; then
    is_debug=1
  fi

  build_name=$build_name-$flag
done

if [[ $is_debug == 1 ]]; then
  # go build -race -tags "$flags" -o "./build/$build_name" .
  go build -tags "$flags" -o "./build/$build_name" .
else
  go build -tags "$flags" -o "./build/$build_name" .
fi

if [[ $? == 0 ]]; then
  echo "✅ SUCCESS"
else
  echo "\n❌ FAIL"
  exit 1
fi

echo  "\nMoving $build_name to /usr/local/bin..."
cp ./build/$build_name /usr/local/bin
echo "...done"
