// Types matching the backend auth DTOs (internal/dto/auth_dto.go)

export type Role = "admin" | "cashier" | "waiter"

export interface User {
  id: string
  name: string
  email: string
  role: Role
  is_active: boolean
}

export interface LoginRequest {
  email: string
  password: string
}

export interface AuthResponse {
  token: string
  user: User
}
