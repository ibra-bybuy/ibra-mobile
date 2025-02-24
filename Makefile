GOMOBILE=gomobile
GOBIND=$(GOMOBILE) bind
BUILDDIR=$(shell pwd)/build
IOS_ARTIFACT=$(BUILDDIR)/Ibra.xcframework
ANDROID_ARTIFACT=$(BUILDDIR)/xray.aar
IOS_TARGET=ios
IOS_VERSION=15.6
MACOS_TARGET=macos
ANDROID_TARGET=android
# LDFLAGS='-s -w -X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn'
LDFLAGS='-s -w -extldflags -lresolv'
IMPORT_PATH=github.com/ibra-bybuy/ibra-mobile

BUILD_IOS="cd $(BUILDDIR) && $(GOBIND) -a -ldflags $(LDFLAGS) -target=$(IOS_TARGET) -iosversion=$(IOS_VERSION) -o $(IOS_ARTIFACT) $(IMPORT_PATH)"
BUILD_MACOS="cd $(BUILDDIR) && $(GOBIND) -a -ldflags $(LDFLAGS) -target=$(MACOS_TARGET) -o $(IOS_ARTIFACT) $(IMPORT_PATH)"
BUILD_ANDROID="cd $(BUILDDIR) && $(GOBIND) -a -ldflags $(LDFLAGS) -target=$(ANDROID_TARGET) -tags=gomobile -o $(ANDROID_ARTIFACT) $(IMPORT_PATH)"

all: ios android

ios:
	mkdir -p $(BUILDDIR)
	eval $(BUILD_IOS)
macos:
	mkdir -p $(BUILDDIR)
	eval $(BUILD_MACOS)
android:
	rm -rf $(BUILDDIR) 2>/dev/null
	mkdir -p $(BUILDDIR)
	eval $(BUILD_ANDROID)

clean:
	rm -rf $(BUILDDIR)
