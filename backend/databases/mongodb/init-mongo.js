db.auth('admin-user', 'admin-password')
// db = db.getSiblingDB('drone_delivery')

db.createUser(
    {
        user: "drone-user",
        pwd: "drone-pwd",
        roles: [
            {
                role: "readWrite",
                db: "drone_delivery"
            }
        ]
    }
)

db.getCollection("warehouse")
db.getCollection("telemetry")
db.getCollection("drone")
db.getCollection("parcel")
db.warehouse.insertOne(
    { id: 1,  location: { latitude: 48.080922, longitude: 20.766208} }
)
