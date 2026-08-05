import Factory
import UIKit

@MainActor
final class AppContainer {
    static let shared = AppContainer()

    let authenticationRepository: any AuthenticationRepository
    let persistenceController: PersistenceController

    private init(container: Container = .shared) {
        persistenceController = container.persistenceController()
        authenticationRepository = LiveAuthenticationRepository(
            client: container.httpClient(),
            credentialsStore: container.credentialsStore()
        )
    }

    func makeRootViewController() -> UIViewController {
        MainTabBarController()
    }
}

extension Container {
    var appConfiguration: Factory<AppConfiguration> {
        self { AppConfiguration.current }.singleton
    }

    var httpClient: Factory<any HTTPClient> {
        self { URLSessionHTTPClient(configuration: self.appConfiguration()) }.singleton
    }

    var credentialsStore: Factory<any CredentialsStoring> {
        self { KeychainCredentialsStore() }.singleton
    }

    var persistenceController: Factory<PersistenceController> {
        self { PersistenceController() }.singleton
    }
}
