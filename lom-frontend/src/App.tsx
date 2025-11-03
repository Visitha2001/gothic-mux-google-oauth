import GoogleLoginButton from './components/GoogleLoginButton';
import LogoutButton from './components/LogoutButton';

function App() {
  const isAuthenticated = false;

  return (
    <div className="p-8 text-center">
      <h1 className="text-4xl font-bold mb-8">Maldives Resort App</h1>
      
      {isAuthenticated ? (
        <div className="space-y-4">
          <p className="text-xl">Welcome back, User!</p>
          <LogoutButton />
        </div>
      ) : (
        <GoogleLoginButton />
      )}
    </div>
  );
}

export default App;