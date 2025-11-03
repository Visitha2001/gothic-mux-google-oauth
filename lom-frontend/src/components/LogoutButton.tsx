const GO_BACKEND_URL = "http://localhost:8082";

function LogoutButton() {
  const logoutUrl = `${GO_BACKEND_URL}/auth/google/logout`;

  return (
    <a 
      href={logoutUrl}
      className="bg-red-600 hover:bg-red-700 text-white font-bold py-2 px-4 rounded transition duration-300"
    >
      Logout
    </a>
  );
}

export default LogoutButton;