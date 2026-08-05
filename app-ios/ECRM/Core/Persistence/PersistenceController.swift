import CoreData

final class PersistenceController: @unchecked Sendable {
    let container: NSPersistentContainer

    init(inMemory: Bool = false) {
        let model = NSManagedObjectModel()
        let cachedSession = NSEntityDescription()
        cachedSession.name = "CachedSession"
        cachedSession.managedObjectClassName = "NSManagedObject"
        cachedSession.properties = [
            attribute(name: "userID", type: .integer64AttributeType),
            attribute(name: "updatedAt", type: .dateAttributeType),
        ]
        model.entities = [cachedSession]

        container = NSPersistentContainer(name: "ECRM", managedObjectModel: model)
        if inMemory {
            container.persistentStoreDescriptions.first?.url = URL(fileURLWithPath: "/dev/null")
        }
        container.loadPersistentStores { _, error in
            if let error {
                AppLogger.app.error("Core Data 存储初始化失败：\(error.localizedDescription, privacy: .public)")
            }
        }
        container.viewContext.automaticallyMergesChangesFromParent = true
    }

    private func attribute(name: String, type: NSAttributeType) -> NSAttributeDescription {
        let attribute = NSAttributeDescription()
        attribute.name = name
        attribute.attributeType = type
        attribute.isOptional = false
        return attribute
    }
}
