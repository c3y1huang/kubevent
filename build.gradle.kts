import com.pivotstir.gogradle.GoPluginExtension

plugins {
    id("com.pivotstir.gogradle") version "1.1.5"
}

apply(plugin = "com.pivotstir.gogradle")

group = "com.innobead"
version = "0.0.1"

repositories {
    jcenter()
    mavenLocal()
    mavenCentral()
}

extensions.getByType(GoPluginExtension::class).apply {
    pluginConfig.modulePath = "github.com/innobead/kubevent"

    env {
        useSandbox = false
    }

    build {
        packagePaths = listOf("./cmd/kubevent")
    }

    dep {
        thirdpartyIgnored = true
    }

    dependencies {
        build("k8s.io/client-go@master") // waiting v12.0.0 release supporting module

        // workaround: request to use 1.14 api related libraries to avoid incompatible usage when using controller-runtime
        build("k8s.io/api@kubernetes-1.14.0")
        build("k8s.io/apiextensions-apiserver@kubernetes-1.14.0")
        build("k8s.io/apimachinery@kubernetes-1.14.0")
        build("sigs.k8s.io/controller-runtime@master") // controller-runtime master supports module which decides what version of client-go should be used.

        build("github.com/spf13/cobra/cobra@v0.0.5")
        build("github.com/sirupsen/logrus@v1.4.2")
        build("github.com/thoas/go-funk@v0.4.0")
        build("github.com/spf13/viper@v1.4.0")
        build("github.com/mitchellh/mapstructure@v1.1.2")
    }
}
