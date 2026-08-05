import Foundation

enum HTTPMethod: String, Sendable {
    case get = "GET"
    case post = "POST"
    case put = "PUT"
    case delete = "DELETE"
}

struct RequestContext: Sendable {
    let accessToken: String?
    let storeAppID: String?

    static let anonymous = RequestContext(accessToken: nil, storeAppID: nil)
}

struct APIRequest: Sendable {
    let path: String
    let method: HTTPMethod
    let body: Data?

    init(path: String, method: HTTPMethod = .get, body: Data? = nil) {
        self.path = path
        self.method = method
        self.body = body
    }
}

protocol HTTPClient: Sendable {
    func send<Response: Decodable & Sendable>(
        _ request: APIRequest,
        context: RequestContext
    ) async throws -> Response
}

private struct APIEnvelope<Payload: Decodable & Sendable>: Decodable, Sendable {
    let status: Int
    let message: String?
    let data: Payload
}

enum APIError: Error, LocalizedError, Sendable {
    case invalidURL
    case invalidResponse
    case server(statusCode: Int, message: String?)
    case decoding

    var errorDescription: String? {
        switch self {
        case .invalidURL: return "接口地址无效"
        case .invalidResponse: return "接口响应无效"
        case let .server(statusCode, message): return message ?? "接口请求失败（\(statusCode)）"
        case .decoding: return "接口数据格式无效"
        }
    }
}

actor URLSessionHTTPClient: HTTPClient {
    private let session: URLSession
    private let configuration: AppConfiguration
    private let decoder: JSONDecoder

    init(configuration: AppConfiguration, session: URLSession = .shared) {
        self.configuration = configuration
        self.session = session
        decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
    }

    func send<Response: Decodable & Sendable>(
        _ request: APIRequest,
        context: RequestContext = .anonymous
    ) async throws -> Response {
        let normalizedPath = request.path.hasPrefix("/") ? String(request.path.dropFirst()) : request.path
        guard let url = URL(string: normalizedPath, relativeTo: configuration.apiBaseURL) else {
            throw APIError.invalidURL
        }

        var urlRequest = URLRequest(url: url)
        urlRequest.httpMethod = request.method.rawValue
        urlRequest.httpBody = request.body
        urlRequest.setValue("application/json", forHTTPHeaderField: "Accept")
        if request.body != nil {
            urlRequest.setValue("application/json", forHTTPHeaderField: "Content-Type")
        }
        if let accessToken = context.accessToken, !accessToken.isEmpty {
            // 服务端只接受此拼写，不能替换为标准 Authorization。
            urlRequest.setValue("Bearer \(accessToken)", forHTTPHeaderField: "Authori-zation")
        }
        if let storeAppID = context.storeAppID, !storeAppID.isEmpty {
            urlRequest.setValue(storeAppID, forHTTPHeaderField: "X-AppId")
        }

        if configuration.logsNetworkRequests {
            AppLogger.network.debug("请求 \(request.method.rawValue, privacy: .public) \(url.path, privacy: .public)")
        }

        let (data, response) = try await session.data(for: urlRequest)
        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }
        guard (200 ... 299).contains(httpResponse.statusCode) else {
            throw APIError.server(statusCode: httpResponse.statusCode, message: nil)
        }
        do {
            let envelope = try decoder.decode(APIEnvelope<Response>.self, from: data)
            guard envelope.status == 200 else {
                throw APIError.server(statusCode: envelope.status, message: envelope.message)
            }
            return envelope.data
        } catch {
            AppLogger.network.error("响应解码失败：\(error.localizedDescription, privacy: .public)")
            throw APIError.decoding
        }
    }
}
