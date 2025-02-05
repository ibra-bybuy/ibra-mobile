// swift-tools-version: 5.7

import PackageDescription

let package = Package(
  name: "XPN",
  platforms: [.iOS(.v12)],
  products: [
    .library(name: "XPN", targets: ["XPN"])
  ],
  targets: [
    .binaryTarget(
      name: "XPN",
      url: "github.com/ibra-bybuy/xray-mobile/releases/download/1.8.1/XPN.xcframework.zip",
      checksum: "803a4561f614971744b044fe2943710025297cb6064f78824f55f7f9f1f46fb0"
    )
  ]
)
