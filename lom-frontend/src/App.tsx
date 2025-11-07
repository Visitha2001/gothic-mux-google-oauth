import GoogleLoginButton from './components/GoogleLoginButton';
import LogoutButton from './components/LogoutButton';

function App() {
  const isAuthenticated = false;

  return (
    <div className="flex items-center justify-center h-screen">
      <div className="container flex flex-col items-center bg-gray-800 p-8 rounded-lg shadow-lg">
        <h1 className="text-4xl font-bold mb-8 text-white text-center">Maldives Resort App</h1>
        {isAuthenticated ? (
          <div className="space-y-4 flex flex-col items-center">
            <p className="text-xl text-white">Welcome back, User!</p>
            <LogoutButton />
          </div>
        ) : (
          <GoogleLoginButton />
        )}
      </div>
    </div>
  );
}

export default App;