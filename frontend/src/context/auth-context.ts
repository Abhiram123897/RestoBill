import { createContext } from "react"
import type { User } from "@/types/auth"

export interface AuthContextValue {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<void>
  logout: () => void
}

// Undefined default lets useAuth detect usage outside the provider
export const AuthContext = createContext<AuthContextValue | undefined>(undefined)
