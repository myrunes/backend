BRANCH=${TRAVIS_BRANCH}

if [ "${BRANCH}" == "master" ]; then
    BRANCH=latest
fi

if [ "${BRANCH}" == "dev" ]; then
    BRANCH=canary
fi

docker build . -t zekro/vctr:${BRANCH}
docker login -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD}
docker push zekro/vctr:${BRANCH}