default: clean build

.PHONY: clean
clean:
	./gradlew goclean

.PHONY: build
build: clean
	./gradlew gobuild

.PHONY: publish
publish: publish
	echo "TODO: Release to github release page"

# TODO Supported in gogradle
#code-generate:
#	go get -u k8s.io/code-generator/...
#    cd $GOPATH/src/k8s.io/code-generator
#    ./generate-groups.sh all "github.com/innobead/kubevent/pkg/client" "github.com/innobead/kubevent/pkg/api" "kubevent:v1alpha1"
