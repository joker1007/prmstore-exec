VERSION=$(shell gobump show -r)
COMMIT=$(shell git rev-parse HEAD)

build:
	go build -ldflags "-X main.Gitcommit=${COMMIT}" ./cmd/prmstore-exec

crossbuild: pkg

pkg: cmd/prmstore-exec/*.go
	gox -os="linux darwin" -arch="amd64 arm64" -output="pkg/{{.OS}}_{{.Arch}}/{{.Dir}}" -ldflags="-w -s" ./cmd/prmstore-exec

archive: crossbuild archive_linux_amd64 archive_linux_arm64 archive_darwin_amd64 archive_darwin_arm64

archive_linux_amd64:
	cp README.md LICENSE pkg/linux_amd64/ && tar cvzf prmstore-exec_${VERSION}_linux_amd64.tar.gz pkg/linux_amd64 && mkdir -p releases && mv prmstore-exec_${VERSION}_linux_amd64.tar.gz releases

archive_linux_arm64:
	cp README.md LICENSE pkg/linux_arm64/ && tar cvzf prmstore-exec_${VERSION}_linux_arm64.tar.gz pkg/linux_arm64 && mkdir -p releases && mv prmstore-exec_${VERSION}_linux_arm64.tar.gz releases

archive_darwin_amd64:
	cp README.md LICENSE pkg/darwin_amd64/ && tar cvzf prmstore-exec_${VERSION}_darwin_amd64.tar.gz pkg/darwin_amd64 && mkdir -p releases && mv prmstore-exec_${VERSION}_darwin_amd64.tar.gz releases

archive_darwin_arm64:
	cp README.md LICENSE pkg/darwin_arm64/ && tar cvzf prmstore-exec_${VERSION}_darwin_arm64.tar.gz pkg/darwin_arm64 && mkdir -p releases && mv prmstore-exec_${VERSION}_darwin_arm64.tar.gz releases

clean:
	rm -rf pkg
	rm -rf releases
