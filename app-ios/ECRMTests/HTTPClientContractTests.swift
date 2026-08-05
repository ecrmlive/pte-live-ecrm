import XCTest
@testable import ECRM

final class HTTPClientContractTests: XCTestCase {
    func testAnonymousContextHasNoCredentials() {
        XCTAssertNil(RequestContext.anonymous.accessToken)
        XCTAssertNil(RequestContext.anonymous.storeAppID)
    }

    func testStoreContextKeepsTokenAndStoreAppIDSeparate() {
        let context = RequestContext(accessToken: "access-token", storeAppID: "store-app-id")
        XCTAssertEqual(context.accessToken, "access-token")
        XCTAssertEqual(context.storeAppID, "store-app-id")
    }
}
