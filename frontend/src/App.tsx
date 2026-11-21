import { useAuth } from "@/hooks/useAuth"
import LoginPage from "@/pages/LoginPage"
import { Button } from "@/components/ui/button"

// Temporary switch: routing + real dashboard come in Story 2
function App() {
  const { isAuthenticated, user, logout } = useAuth()

  if (!isAuthenticated) {
    return <LoginPage />
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="bg-white p-8 rounded-lg shadow-md flex flex-col items-center gap-4">
        <h1 className="text-3xl font-bold text-blue-600">RestoBill</h1>
        <p className="text-gray-600">
          Logged in as <strong>{user?.name}</strong> ({user?.role})
        </p>
        <Button variant="outline" onClick={logout}>
          Log out
        </Button>
      </div>
    </div>
  )
}

export default App
