BRANCH=${TRAVIS_BRANCH}
RELEASE=TRUE

if [ "${BRANCH}" == "master" ]; then
    BRANCH=latest
fi

if [ "${BRANCH}" == "dev" ]; then
    BRANCH=canary
    RELEASE=FALSE
fi

docker build . -t zekro/myrunes:${BRANCH} --build-arg RELEASE=${RELEASE}
docker login -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD}
docker push zekro/myrunes:${BRANCH}