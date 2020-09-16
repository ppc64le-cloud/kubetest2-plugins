module github.com/ppc64le-cloud/kubetest2-plugins

go 1.15

require (
	cloud.google.com/go v0.51.0 // indirect
	github.com/lucasjones/reggen v0.0.0-20200904144131-37ba4fa293bb
	github.com/pkg/errors v0.9.1
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546 // indirect
	github.com/spf13/pflag v1.0.5
	k8s.io/api v0.19.1 // indirect
	k8s.io/client-go v9.0.0+incompatible
	k8s.io/klog/v2 v2.2.0
	k8s.io/utils v0.0.0-20200729134348-d5654de09c73
	sigs.k8s.io/kubetest2 v0.0.0-20200910235614-8dd2cc76cff9
)

replace k8s.io/client-go => k8s.io/client-go v0.19.1
