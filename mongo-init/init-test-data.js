// MongoDB initialization script for test data
// This script creates 10 test users for development and testing

db = db.getSiblingDB('users_brm_dev');

// Drop existing collections to start fresh
db.users.drop();

// Create index on email field for uniqueness
db.users.createIndex({ "email": 1 }, { unique: true });

// Insert 10 test users
db.users.insertMany([
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Juan Carlos Pérez",
    "email": "juan.perez@example.com",
    "age": 28,
    "created_at": new Date("2024-01-15T10:30:00Z"),
    "updated_at": new Date("2024-01-15T10:30:00Z")
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "name": "María González López",
    "email": "maria.gonzalez@example.com",
    "age": 32,
    "created_at": new Date("2024-01-16T14:20:00Z"),
    "updated_at": new Date("2024-01-16T14:20:00Z")
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440003",
    "name": "Carlos Alberto Rodríguez",
    "email": "carlos.rodriguez@example.com",
    "age": 25,
    "created_at": new Date("2024-01-17T09:15:00Z"),
    "updated_at": new Date("2024-01-17T09:15:00Z")
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440004",
    "name": "Ana Sofía Martínez",
    "email": "ana.martinez@example.com",
    "age": 29,
    "created_at": new Date("2024-01-18T16:45:00Z"),
    "updated_at": new Date("2024-01-18T16:45:00Z")
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440005",
    "name": "Luis Fernando Herrera",
    "email": "luis.herrera@example.com",
    "age": 35,
    "created_at": new Date("2024-01-19T11:30:00Z"),
    "updated_at": new Date("2024-01-19T11:30:00Z")
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440006",
    "name": "Carmen Elena Vargas",
    "email": "carmen.vargas@example.com",
    "age": 27,
    "created_at": new Date("2024-01-20T13:20:00Z"),
    "updated_at": new Date("2024-01-20T13:20:00Z")
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440007",
    "name": "Roberto José Silva",
    "email": "roberto.silva@example.com",
    "age": 31,
    "created_at": new Date("2024-01-21T08:55:00Z"),
    "updated_at": new Date("2024-01-21T08:55:00Z")
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440008",
    "name": "Patricia Isabel Morales",
    "email": "patricia.morales@example.com",
    "age": 26,
    "created_at": new Date("2024-01-22T15:10:00Z"),
    "updated_at": new Date("2024-01-22T15:10:00Z")
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440009",
    "name": "Diego Alejandro Torres",
    "email": "diego.torres@example.com",
    "age": 33,
    "created_at": new Date("2024-01-23T12:40:00Z"),
    "updated_at": new Date("2024-01-23T12:40:00Z")
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440010",
    "name": "Valentina Andrea Jiménez",
    "email": "valentina.jimenez@example.com",
    "age": 24,
    "created_at": new Date("2024-01-24T17:25:00Z"),
    "updated_at": new Date("2024-01-24T17:25:00Z")
  }
]);

// Print confirmation
print("✅ Test data initialized successfully!");
print("📊 Created " + db.users.countDocuments() + " test users");
print("📋 Users created:");
db.users.find({}, {name: 1, email: 1, age: 1}).forEach(function(user) {
  print("   - " + user.name + " (" + user.email + ") - " + user.age + " años");
}); 