import OSLog

enum AppLogger {
    static let app = Logger(subsystem: Bundle.main.bundleIdentifier ?? "com.qixi.ecrm.ios", category: "app")
    static let network = Logger(subsystem: Bundle.main.bundleIdentifier ?? "com.qixi.ecrm.ios", category: "network")
}
