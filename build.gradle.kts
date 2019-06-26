import com.pivotstir.gogradle.GoPluginExtension

plugins {
    id("com.pivotstir.gogradle") version "1.1.3"
}

group = "com.innobead"
version = "0.0.1"

repositories {
    jcenter()
    mavenLocal()
    mavenCentral()
}

extensions.getByType(GoPluginExtension::class).apply {
    env {
        useSandbox = true
    }

    build {
        packagePaths = listOf("./cmd/kubevent")
    }

    dep {
        thirdpartyIgnored = true
    }

    dependencies {
        //        build("github.com/golang/protobuf@v1.3.1")
    }
}