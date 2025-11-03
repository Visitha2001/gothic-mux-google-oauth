const GO_BACKEND_URL = "http://localhost:8082";

function GoogleLoginButton() {
  const loginUrl = `${GO_BACKEND_URL}/auth/google`;

  return (
    <a href={loginUrl}>
      Login with Google
    </a>
  );
}

export default GoogleLoginButton;