const GO_BACKEND_URL = "http://localhost:8082";

function GoogleLoginButton() {
  const loginUrl = `${GO_BACKEND_URL}/auth/google`;

  return (
    <a href={loginUrl} className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded-xl">
      <span className="text-white">Login with Google</span>
    </a>
  );
}

export default GoogleLoginButton;