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