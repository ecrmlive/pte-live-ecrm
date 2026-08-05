import Foundation

private struct LoginPayload: Encodable {
    let account: String
    let password: String
    let channel = "ios"
    let captchaToken: String
}

private struct LoginResponse: Decodable, Sendable {
    let token: TokenPair
    let user: UserProfile
}

private struct StoreContextResponse: Decodable, Sendable {
    let token: TokenPair
    let context: StoreContext
}

final class LiveAuthenticationRepository: AuthenticationRepository, @unchecked Sendable {
    private let client: any HTTPClient
    private let credentialsStore: any CredentialsStoring

    init(client: any HTTPClient, credentialsStore: any CredentialsStoring) {
        self.client = client
        self.credentialsStore = credentialsStore
    }

    func logIn(account: String, password: String, captchaToken: String) async throws -> UserSession {
        let body = try JSONEncoder().encode(LoginPayload(
            account: account,
            password: password,
            captchaToken: captchaToken
        ))
        let response: LoginResponse = try await client.send(
            APIRequest(path: "api/app/v1/auth/login", method: .post, body: body),
            context: .anonymous
        )
        let session = UserSession(token: response.token, profile: response.user)
        try await credentialsStore.save(session)
        return session
    }

    func currentSession() async throws -> UserSession? {
        try await credentialsStore.load()
    }

    func issueStoreContext(appID: String) async throws -> UserSession {
        guard var session = try await credentialsStore.load() else {
            throw APIError.server(statusCode: 401, message: "请先登录")
        }
        let response: StoreContextResponse = try await client.send(
            APIRequest(path: "api/app/v1/auth/store-context", method: .post),
            context: RequestContext(accessToken: session.token.accessToken, storeAppID: appID)
        )
        session = UserSession(
            token: response.token,
            profile: session.profile,
            merchantAppID: response.context.merchantAppID
        )
        try await credentialsStore.save(session)
        return session
    }

    func signOut() async throws {
        try await credentialsStore.clear()
    }
}
