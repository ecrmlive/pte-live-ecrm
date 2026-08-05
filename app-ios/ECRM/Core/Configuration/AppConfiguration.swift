import Foundation

struct AppConfiguration: Sendable {
    let apiBaseURL: URL
    let logsNetworkRequests: Bool

    static let current: AppConfiguration = {
        let value = Bundle.main.object(forInfoDictionaryKey: "ECRMAPIBaseURL") as? String
        let apiBaseURL = URL(string: value ?? "") ?? URL(string: "https://api.example.invalid/")!
        let logsNetworkRequests = Bundle.main.object(forInfoDictionaryKey: "ECRMLogNetwork") as? Bool ?? false
        return AppConfiguration(apiBaseURL: apiBaseURL, logsNetworkRequests: logsNetworkRequests)
    }()
}
