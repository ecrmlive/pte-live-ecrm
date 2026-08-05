import Foundation

struct UserProfile: Codable, Sendable, Equatable {
    let id: Int64
    let account: String
    let nickname: String?
    let avatar: String?
}

struct TokenPair: Codable, Sendable, Equatable {
    let accessToken: String
    let refreshToken: String
    let expiresIn: Int64
}

struct UserSession: Codable, Sendable, Equatable {
    let token: TokenPair
    let profile: UserProfile
    let merchantAppID: String?

    init(token: TokenPair, profile: UserProfile, merchantAppID: String? = nil) {
        self.token = token
        self.profile = profile
        self.merchantAppID = merchantAppID
    }
}

struct StoreContext: Codable, Sendable, Equatable {
    let merchantID: Int64
    let storeID: Int64
    let merchantAppID: String
    let imSDKAppID: String?
}

protocol AuthenticationRepository: Sendable {
    func logIn(account: String, password: String, captchaToken: String) async throws -> UserSession
    func currentSession() async throws -> UserSession?
    func issueStoreContext(appID: String) async throws -> UserSession
    func signOut() async throws
}
