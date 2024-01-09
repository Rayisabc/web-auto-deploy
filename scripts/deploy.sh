#!/bin/bash

##
#
# Author:Rui
# Date: 2024/01/03
# Desc: Docker deploy shell
#
##

# Debug model
# set -x

##
#
# @Param current_version
# current_version++
# @Return new_version
##
function version_rise() {
    local current_version=$1
    version_without_v="${current_version#v}"
    IFS='.' read -ra version_numbers <<<"$version_without_v"
    major=${version_numbers[0]}
    minor=${version_numbers[1]}
    patch=${version_numbers[2]}

    # auto increment patch|minor
    if [ $patch -eq 9 ]; then
        ((minor++))
        patch=0
        if [ $minor -eq 9 ]; then
            ((major++))
            minor=0
            patch=0
        fi
    else
        ((patch++))
    fi

    new_version="v$major.$minor.$patch"
    # new_version
    echo "$new_version"
}
# @see config.toml:{ServiceDir}
current_dir=test
# yours service name
service_name={service-name}
# yours docker repository
image_name={yours_repository}/$service_name

echo '### start check image version'
echo

# get the newest version
max_version=$(docker images | grep "$image_name" | awk '{print $2}' | sort -r | head -n 1)
if [[ -z $max_version ]]; then
    # new image
    max_version='v1.0.0'
fi

echo "current newest version is:$max_version"
echo
# version++
new_version=$(version_rise $max_version)
echo "system auto-assigned new version:$new_version"
image_tag=$new_version

# building
echo
echo "### start building{$service_name}image, version : $image_tag,repository:$image_name"
echo

cd /usr/local/{your_workspace}/$current_dir

docker build -t $image_name:$image_tag . -f {Dockerfile_name}

if [ $? -eq 0 ]; then
    echo
    echo '### image built success,pushing image...'

    echo
    docker push $image_name:$image_tag
    echo

    if [ $? -eq 0 ]; then
        echo '### image pushed success'
    else
        echo '### oh no image push failed'
        exit 1
    fi
else
    echo
    echo '### build image failed'
    exit 1
fi
